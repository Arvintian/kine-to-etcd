package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/k3s-io/kine/pkg/client"
	"github.com/k3s-io/kine/pkg/endpoint"
	"github.com/k3s-io/kine/pkg/version"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	config   endpoint.Config
	toConfig endpoint.Config
)

func main() {
	app := cli.NewApp()
	app.Name = "k2e"
	app.Usage = "Migrate kine's data."
	app.Version = fmt.Sprintf("%s (%s)", version.Version, version.GitCommit)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "endpoint",
			Usage:       "Storage endpoint (default is sqlite)",
			Destination: &config.Endpoint,
		},
		cli.StringFlag{
			Name:        "to-endpoint",
			Usage:       "Endpoints for target storage",
			Destination: &toConfig.Endpoint,
		},
		cli.StringFlag{
			Name:        "to-ca-file",
			Usage:       "Client ca file for target storage",
			Destination: &toConfig.BackendTLSConfig.CAFile,
		},
		cli.StringFlag{
			Name:        "to-cert-file",
			Usage:       "Client cert file for target storage",
			Destination: &toConfig.BackendTLSConfig.CertFile,
		},
		cli.StringFlag{
			Name:        "to-key-file",
			Usage:       "Client key file for target storage",
			Destination: &toConfig.BackendTLSConfig.KeyFile,
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

	toEndpoint, err := endpoint.Listen(ctx, toConfig)
	if err != nil {
		return err
	}

	toClient, err := client.New(toEndpoint)
	if err != nil {
		return err
	}

	defer toClient.Close()

	values, err := kineClient.List(ctx, "/registry/", 0)
	if err != nil {
		return err
	}

	for _, value := range values {
		logrus.Infof("Migrating etcd key %s", value.Key)
		if !c.Bool("debug") {
			val, err := toClient.Get(ctx, string(value.Key))
			if err != nil {
				err = toClient.Create(ctx, string(value.Key), value.Data)
			} else {
				err = toClient.Update(ctx, string(val.Key), val.Modified, val.Data)
			}
			if err != nil {
				return err
			}
		}
	}

	logrus.Info("Migrating successed")

	return nil
}
