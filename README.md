
## 开始准备
```shell
go1.14 ，ETCD，redis-cluster，kafka-cluster

ETCD配置
/micro/config/kafka-cluster
{"addr":["192.168.1.113:9191","192.168.1.113:9192","192.168.1.113.9193"]}

/micro/config/redis-cluster
{"addr":["192.168.1.38:9001","192.168.1.38:9002","192.168.1.38:9003","192.168.1.38:9004","192.168.1.38:9005","192.168.1.38:9006"],"password":""}

/micro/config/room-service
{"topic":"im-room-dev"}
```

## frontier-service 编译&启动
```shell
编译
go build ./cmd/frontier/main.go

启动
./cmd/frontier/main --registry=etcd --registry_address=:{ETCD端口} --server_address=:1234 --frontier_id={边界机ID}
```

## room-service 编译&启动
```shell
编译
go build ./cmd/room-service/main.go

启动
./cmd/room-service/main --registry=etcd --registry_address=:{ETCD端口} --server_address=:7880
```

## tenant-web-service 编译&启动
```shell
编译
go build ./cmd/tenant-web-service/main.go

启动
./cmd/tenant-web-service/main --registry=etcd --registry_address=:{ETCD端口} --server_address=:8090
```

## 参数描述
--registry_address ETCD地址

--server_address 当前服务的监听 地址:端口
