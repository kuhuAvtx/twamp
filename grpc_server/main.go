/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Metric service.
package main

import (
	"log"
	"net"

	config "github.com/kuhuAvtx/twamp/conf"
	newclient "github.com/kuhuAvtx/twamp/newclient"
	pb "github.com/kuhuAvtx/twamp/proto"
	"google.golang.org/grpc"
)

// server is used to implement Server.
type server struct {
	pb.UnimplementedTwampMetricsServiceServer
}

func NewServer() *server {
	return &server{}
}

// GetMetrics that will fetch the latency metrics
func (s *server) GetMetrics(req *pb.TwampMetricsRequest, server pb.TwampMetricsService_GetMetricsServer) error {
	log.Printf("KUHU in GetMetrics")
	latency := newclient.GetLatency()
	sendErr := server.Send(&pb.TwampMetricsReply{Latency: latency})
	if sendErr != nil {
		panic(sendErr)
	}
	return nil
}

func main() {
	var conf = config.ReadConfig()
	lis, err := net.Listen("tcp", conf.GrpcServer.GrpcHost+":"+conf.GrpcServer.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTwampMetricsServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
