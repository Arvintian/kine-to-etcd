module github.com/Arvintian/kine-to-etcd

go 1.16

replace (
	github.com/rancher/wrangler => github.com/rancher/wrangler v0.8.11-0.20220211163748-d5a8ee98be5f
	go.etcd.io/etcd/client/v3 => github.com/k3s-io/etcd/client/v3 v3.5.1-k3s1
)

require (
	github.com/k3s-io/kine v0.8.1
	github.com/rancher/dynamiclistener v0.3.1
	github.com/rancher/wrangler v0.8.11
	github.com/sirupsen/logrus v1.7.0
	github.com/urfave/cli v1.21.0
	go.etcd.io/etcd/client/v3 v3.5.0
)
