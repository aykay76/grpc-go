package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/aykay76/grpc-go/environment"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	server = flag.Bool("server", false, "Set to true if acting as a server, false if acting as a client")
	tls    = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")

	// server variables
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 10000, "The server port")

	// client variables
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func main() {
	flag.Parse()
	var opts []grpc.ServerOption

	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server1.pem")
		}
		if *keyFile == "" {
			*keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	if *server {
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		defer lis.Close()

		grpcServer := grpc.NewServer(opts...)
		pb.RegisterEnvironmentServiceServer(grpcServer, &EnvironmentServer{})

		fmt.Println("Serving /metrics for Prometheus...")
		// run the http server for Prometheus scraping in the background
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			// TODO: make the port configurable
			http.ListenAndServe(":9093", nil)
			if err != nil {
				fmt.Println("Could not initiate metrics")
			}
		}()

		fmt.Printf("Serving GRPC listener on port %d...\n", *port)
		grpcServer.Serve(lis)
	} else {
		NewClient()
	}
}
