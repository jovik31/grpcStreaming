gen:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

clean:
	rm proto/*.pb.go
