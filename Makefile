.PHONY: clean proto

all: clean proto

clean:
	find ./ -type f -name "*.pb.go" | xargs rm
	find ./ -type f -name "*_grpc.pb.go" | xargs rm

proto:
	# /usr/local/include 包含 wellknown protobuf 文件.
	protoc -I=. -I=/usr/local/include --go_out=:. proto/*/*.proto
	protoc -I=. -I=/usr/local/include --go-grpc_out=:. proto/*/*.proto
	protoc -I=. -I=/usr/local/include --go_out=:. proto/*/*/*.proto
	protoc -I=. -I=/usr/local/include --go-grpc_out=:. proto/*/*/*.proto
	cp -r wchat.im/* .
	rm -rf wchat.im
