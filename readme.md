copy the proto file and add/remove the services and messages you want

protoc -I . ./environment.proto --go_out=plugins=grpc:environment

