# kine-to-etcd


Migrate k3s's [kine](https://github.com/k3s-io/kine) data to etcd.

## Usage

```
./k2e --endpoint sqlite://./tests/state.db --etcd-endpoint http://localhost:2379
```