copy the proto file and add/remove the services and messages you want

protoc -I . ./environment.proto --go_out=plugins=grpc:environment

Ensure Dockerfile contains parameters in CMD entry to run as server.

Make sure server is listening on 0.0.0.0 NOT localhost so that it works in container also.

## To build the image

`docker build . --tag localhost:5000/grpc-go`
`docker push localhost:5000/grpc-go`

## To run the image

`docker run --rm --name grpc-go --publish 10000:10000 localhost:5000/grpc-go`