package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	pb "github.com/aykay76/grpc-go/environment"
	empty "github.com/golang/protobuf/ptypes/empty"
)

// EnvironmentServer : server for the environment service
type EnvironmentServer struct {
	pb.UnimplementedEnvironmentServiceServer
}

// GetEnvironmentVariable : allow clients to get specified environment variable
func (server *EnvironmentServer) GetEnvironmentVariable(ctx context.Context, kvp *pb.KeyValuePair) (*pb.KeyValuePair, error) {
	var result pb.KeyValuePair
	result.Key = kvp.Key
	result.Value = os.Getenv(kvp.Key)

	return &result, nil
}

// GetEnvironmentVariables : allows clients to get all environment variables on a stream
func (server *EnvironmentServer) GetEnvironmentVariables(req *empty.Empty, stream pb.EnvironmentService_GetEnvironmentVariablesServer) error {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		fmt.Println(pair[0])
		var entry pb.KeyValuePair
		entry.Key = pair[0]
		entry.Value = pair[1]
		err := stream.Send(&entry)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

// SetEnvironmentVariable : allows clients to set environment variables
func (server *EnvironmentServer) SetEnvironmentVariable(ctx context.Context, kvp *pb.KeyValuePair) (*empty.Empty, error) {
	os.Setenv(kvp.Key, kvp.Value)
	return nil, nil
}
