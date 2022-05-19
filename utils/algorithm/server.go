package algorithm

import (
	"context"
	pb "ess/gRPC"
	"ess/utils/setting"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	connection *grpc.ClientConn
	client     pb.AlgorithmClient
}

func (s *server) Setup() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	s.connection, err = grpc.DialContext(ctx, setting.GRPCSetting.Host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect rpc server: %v", err)
	}
	s.client = pb.NewAlgorithmClient(s.connection)
}

func (s *server) Ping(req *pb.PingRequest) (*pb.PingReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	r, err := s.client.Ping(ctx, req)
	return r, err
}

func (s *server) Schedule(req *pb.ScheduleRequest) (*pb.ScheduleReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	r, err := s.client.Schedule(ctx, req)
	return r, err
}
