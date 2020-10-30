<<<<<<< HEAD
# wechat

## 环境准备

```shell
go install github.com/bazelbuild/bazelisk
alias bazel=bazelisk
```

## 编译

所有构建物编译

```shell
cd $workspace
bazel build //...
```

frontier 编译
```shell
bazel build //cmd/frontier:frontier
```

room-node 编译
```shell
bazel build //cmd/room-node:room-node
```

frontier 镜像编译
```shell
bazel build //cmd/frontier:frontier_image
```

room-node 镜像编译
```shell
bazel build //cmd/room-node:room-node_image
```

## Bazel 构建文件

添加构建文件
```shell
bazel run //:gazelle
```

## 添加依赖

```shell
bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=repositories.bzl%go_repositories
```

## 构建物

### frontier

可执行文件
bazel-bin/cmd/frontier/frontier_/frontier

docker 镜像
bazel-bin/cmd/frontier/frontier_image-layer.tar

### room-node

可执行文件
bazel-bin/cmd/room-node/room-node_/room-node

docker 镜像
bazel-bin/cmd/room-node/room-node_image-layer.tar

## 部署

### frontier

frontier 依赖 room-node 和 token-service 必要启动参数为
room-node 和 token-service 地址

./frontier -room-node-addr {IP}:{PORT} -token-service-addr {IP}:{PORT} -log_dir {LOG_DIR} -v 1

更多参数以及具体含义通过 ./frontier -h 查看.

frontier 可以部署多实例支持更多连接数.

frontier 主要任务是承接长连接，因此 frontier 所在机器需要做文件具柄和 TCP 优化.

[调优参考](https://www.jianshu.com/p/e0b52dc702d6)

### room-node

room-node 当前不依赖其他服务

./room-node -log_dir {LOG_DIR} -v 1
=======

## 开始准备
```shell
go1.14 ，ETCD，redis-cluster，kafka-cluster

ETCD配置
/micro/config/frontier
{"logLevel":4,"heartbeatTimeout":90,"writerBufferSize":1024,"readerBufferSize":1024}
logLevel 0 关闭日志,1 错误日志,2 警告日志,3 信息日志,4 全部日志
值越大，打印的信息越全面

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
./cmd/frontier/main --registry=etcd --registry_address=:{ETCD端口} --server_address=:1234 --prometheus_address=:9090
```

## room-service 编译&启动
```shell
编译
go build ./cmd/room-service/main.go

启动
./cmd/room-service/main --registry=etcd --registry_address=:{ETCD端口} --server_address=:7880 --prometheus_address=:9090
```

## tenant-web-service 编译&启动
```shell
编译
go build ./cmd/tenant-web-service/main.go

启动
./cmd/tenant-web-service/main --registry=etcd --registry_address=:{ETCD端口} --server_address=:8090 --prometheus_address=:9090
```

## 参数描述
--registry_address ETCD地址

--server_address 当前服务的监听 地址:端口
>>>>>>> master
