package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	pb "github.com/aykay76/grpc-go/environment"
	empty "github.com/golang/protobuf/ptypes/empty"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
)

type environmentServer struct {
	pb.UnimplementedEnvironmentServiceServer
}

func (server *environmentServer) GetEnvironmentVariable(ctx context.Context, kvp *pb.KeyValuePair) (*pb.KeyValuePair, error) {
	var result pb.KeyValuePair
	result.Key = kvp.Key
	result.Value = os.Getenv(kvp.Key)

	return &result, nil
}

func (server *environmentServer) GetEnvironmentVariables(req *empty.Empty, stream pb.EnvironmentService_GetEnvironmentVariablesServer) error {
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

func (server *environmentServer) SetEnvironmentVariable(ctx context.Context, kvp pb.KeyValuePair) {
	os.Setenv(kvp.Key, kvp.Value)
}

func newServer() *environmentServer {
	s := &environmentServer{}
	return s
}
