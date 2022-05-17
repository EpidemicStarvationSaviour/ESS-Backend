package algorithm

import (
	"ess/utils/logging"
	"ess/utils/setting"
	"log"
	"time"

	pb "ess/gRPC"
)

var s server
var enable bool
var timeout time.Duration

func Setup() {
	enable = setting.GRPCSetting.Enable
	if enable {
		timeout = time.Millisecond * time.Duration(setting.GRPCSetting.Timeout)
		s.Setup()

		_, err := s.Ping(&pb.PingRequest{Message: "ping"})
		if err != nil {
			log.Fatalf("failed to ping rpc server: %v", err)
		}
	}
}

func Schedule(gid int64) error {
	if enable { // TODO(TO/GA)
		request := pb.ScheduleRequest{
			Request: &pb.ItemList{
				Items: map[uint32]float64{1: 5, 2: 1.5},
			},
			NumDeliverer: 1,
			Itemlists: []*pb.ItemList{
				{
					Items: map[uint32]float64{1: 5, 2: 1.5},
				},
			},
			Distance: []uint64{1, 2, 3},
		}
		r, err := s.Schedule(&request)
		if err != nil {
			logging.ErrorF("could not schedule: %v", err)
		}
		log.Printf("%v", r)
	}

	return nil
}
