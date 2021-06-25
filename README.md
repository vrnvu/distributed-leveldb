# distributed-db

- [Reference LevelDB vs RockDB](https://medium.com/walmartglobaltech/https-medium-com-kharekartik-rocksdb-and-embedded-databases-1a0f8e6ea74f)

- [Why we built CockroachDB on top of RocksDB](https://www.cockroachlabs.com/blog/cockroachdb-on-rocksd/)

- [Netflix, applcation data caching using ssds](https://netflixtechblog.com/application-data-caching-using-ssds-5bf25df851ef)

# To RUN

```
$  go run main.go -id node0 ./node0
$  go run main.go -id node1 -haddr 127.0.0.1:11001 -raddr 127.0.0.1:12001 -join :11000 ./node1
$  go run main.go -id node2 -haddr 127.0.0.1:11002 -raddr 127.0.0.1:12002 -join :11000 ./node2
$  curl -XPOST localhost:11000/key -d '{"user1": "batman"}'
$  curl -XGET localhost:11000/key/user1
```

Makefile, only tests

```
make test
```