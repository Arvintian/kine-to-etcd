# kine-to-etcd


Migrate k3s's [kine](https://github.com/k3s-io/kine) data to etcd.

## Usage

```
./kine-to-etcd --endpoint sqlite://./tests/state.db --etcd-endpoint http://localhost:2379
```