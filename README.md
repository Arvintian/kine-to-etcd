# kine-to-etcd


Migrate [k3s](https://github.com/k3s-io/k3s) or [kine] 's (https://github.com/k3s-io/kine) data to etcd.

## install

```
go install github.com/Arvintian/kine-to-etcd
```

## Usage

```
kine-to-etcd --endpoint sqlite://./tests/state.db --etcd-endpoint http://localhost:2379
```