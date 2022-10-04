// TODO add copyright

// Package main implements a client for Metric service.
package main

import (
	"context"
	"log"
	"time"

	config "github.com/kuhuAvtx/twamp/conf"
	pb "github.com/kuhuAvtx/twamp/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var conf = config.ReadConfig()
	// Set up a connection to the server.
	conn, err := grpc.Dial(conf.GrpcServer.GrpcHost+":"+conf.GrpcServer.GrpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTwampMetricsServiceClient(conn)

	// Contact the server and print out its response.
	//TODO do indefinitely for timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	for {
		r, err := c.GetMetrics(ctx, &pb.TwampMetricsRequest{Name: "latency"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		recvd, recvErr := r.Recv()
		if recvErr != nil {
			log.Fatalf("could not Recv: %v", recvErr)
		}
		log.Printf("Latency: %g", recvd.Latency)
		time.Sleep(1 * time.Second)
	}
}
