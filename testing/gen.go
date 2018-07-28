package testing

//go:generate protoc -I ./ -I ../vendor ./test.proto --go_out=plugins=grpc:.
