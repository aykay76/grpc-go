package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/aykay76/grpc-go/environment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var ()

func newClient() {
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	// *** Now to do the actual gRPC ***
	client := pb.NewEnvironmentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.GetEnvironmentVariables(ctx, nil)
	if err != nil {
		log.Fatalf("%v.GetEnvironmentVariables(_) = _, %v", client, err)
	}
	for {
		kvp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetEnvironmentVariables(_) = _, %v", client, err)
		}
		log.Println(kvp)
	}
}
