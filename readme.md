copy the proto file and add/remove the services and messages you want

protoc -I . ./environment.proto --go_out=plugins=grpc:environment

Ensure Dockerfile contains parameters in CMD entry to run as server.

Make sure server is listening on 0.0.0.0 NOT localhost so that it works in container also.