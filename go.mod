module github.com/rhizomata/bridge-chain-etcd

go 1.13

require (
	github.com/coreos/etcd v3.3.18+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/ethereum/go-ethereum v1.9.9
	github.com/gin-gonic/gin v1.5.0
	github.com/google/uuid v1.1.1
	github.com/rhizomata/js v0.0.0-20191231120211-6f5be962b23e
	go.etcd.io/etcd v3.3.18+incompatible
	go.uber.org/zap v1.13.0 // indirect
	google.golang.org/grpc v1.26.0 // indirect
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
