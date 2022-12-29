package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/k3s-io/kine/pkg/client"
	"github.com/k3s-io/kine/pkg/endpoint"
	"github.com/k3s-io/kine/pkg/version"
	certutil "github.com/rancher/dynamiclistener/cert"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConfig struct {
	Endpoints      string
	ClientETCDCert string
	ClientETCDKey  string
	ETCDServerCA   string
}

var (
	config     endpoint.Config
	etcdConfig EtcdConfig
)

const (
	defaultDialTimeout      = 2 * time.Second
	defaultKeepAliveTime    = 30 * time.Second
	defaultKeepAliveTimeout = 10 * time.Second
)

func main() {
	app := cli.NewApp()
	app.Name = "k2e"
	app.Usage = "Migrate kine's data to etcd."
	app.Version = fmt.Sprintf("%s (%s)", version.Version, version.GitCommit)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "endpoint",
			Usage:       "Storage endpoint (default is sqlite)",
			Destination: &config.Endpoint,
		},
		cli.StringFlag{
			Name:        "etcd-endpoint",
			Usage:       "Endpoints for target etcd",
			Destination: &etcdConfig.Endpoints,
		},
		cli.StringFlag{
			Name:        "etcd-ca-file",
			Usage:       "Client ca file for target etcd",
			Destination: &etcdConfig.ETCDServerCA,
		},
		cli.StringFlag{
			Name:        "etcd-cert-file",
			Usage:       "Client cert file for target etcd",
			Destination: &etcdConfig.ClientETCDCert,
		},
		cli.StringFlag{
			Name:        "etcd-key-file",
			Usage:       "Client key file for target etcd",
			Destination: &etcdConfig.ClientETCDKey,
		},
		cli.StringFlag{
			Name:        "ca-file",
			Usage:       "CA cert for DB connection",
			Destination: &config.BackendTLSConfig.CAFile,
		},
		cli.StringFlag{
			Name:        "cert-file",
			Usage:       "Certificate for DB connection",
			Destination: &config.BackendTLSConfig.CertFile,
		},
		cli.StringFlag{
			Name:        "key-file",
			Usage:       "Key file for DB connection",
			Destination: &config.BackendTLSConfig.KeyFile,
		},
		cli.BoolFlag{Name: "debug"},
	}
	app.Action = migrate

	if err := app.Run(os.Args); err != nil {
		if !errors.Is(err, context.Canceled) {
			logrus.Fatal(err)
		}
	}
}

func migrate(c *cli.Context) error {
	if c.Bool("debug") {
		logrus.SetLevel(logrus.TraceLevel)
	}

	ctx := signals.SetupSignalContext()
	kineEndpoint, err := endpoint.Listen(ctx, config)
	if err != nil {
		return err
	}

	kineClient, err := client.New(kineEndpoint)
	if err != nil {
		return err
	}

	defer kineClient.Close()

	etcdClient, err := getEtcdClient(ctx, etcdConfig)
	if err != nil {
		return err
	}
	defer etcdClient.Close()

	values, err := kineClient.List(ctx, "/registry/", 0)
	if err != nil {
		return err
	}

	for _, value := range values {
		logrus.Infof("Migrating etcd key %s", value.Key)
		if !c.Bool("debug") {
			_, err := etcdClient.Put(ctx, string(value.Key), string(value.Data))
			if err != nil {
				return err
			}
		}
	}

	logrus.Info("Migrating successed")

	return nil
}

func getEtcdClient(ctx context.Context, cfg EtcdConfig) (*clientv3.Client, error) {

	endpoints := strings.Split(cfg.Endpoints, ",")

	config := clientv3.Config{
		Endpoints:            endpoints,
		Context:              ctx,
		DialTimeout:          defaultDialTimeout,
		DialKeepAliveTime:    defaultKeepAliveTime,
		DialKeepAliveTimeout: defaultKeepAliveTimeout,
	}

	var err error
	if strings.HasPrefix(endpoints[0], "https://") {
		config.TLS, err = toTLSConfig(cfg)
		if err != nil {
			return nil, err
		}
	}

	return clientv3.New(config)
}

func toTLSConfig(cfg EtcdConfig) (*tls.Config, error) {
	if cfg.ClientETCDCert == "" || cfg.ClientETCDKey == "" || cfg.ETCDServerCA == "" {
		return nil, errors.New("runtime is not ready yet")
	}

	clientCert, err := tls.LoadX509KeyPair(cfg.ClientETCDCert, cfg.ClientETCDKey)
	if err != nil {
		return nil, err
	}

	pool, err := certutil.NewPool(cfg.ETCDServerCA)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{clientCert},
	}, nil
}
