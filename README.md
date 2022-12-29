# kine-to-etcd

Migrate [kine](https://github.com/k3s-io/kine) 's data to etcd.

## Build

```
make build
```

## Usage

```
k2e --endpoint sqlite://./tests/state.db --etcd-endpoint http://localhost:2379
```