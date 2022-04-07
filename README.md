# kine-to-etcd


Migrate [k3s](https://github.com/k3s-io/k3s) or [kine](https://github.com/k3s-io/kine) 's data to etcd.

## Install

```
go install github.com/Arvintian/kine-to-etcd
```

## Usage

```
kine-to-etcd --endpoint sqlite://./tests/state.db --etcd-endpoint http://localhost:2379
```

## Issue

This tool can help you migrate data from kine to etcd, but can not guarantee all workload works fine, because kine is not etcd, migrated key/value not have reversion data. Maybe your controller list and watch will fail because kine's data not keep the reversion data.

But this tool help you migrate workload's spec from kine to etcd, you can just trigger recreate workload by kubectl, but it is no pain. Maybe You should shutdown you workload when you migrate cluster's data.