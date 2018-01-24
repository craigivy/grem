
package main

import (
	"log"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/craigivy/grem/pkg/common"
	"io"
	"time"
)

const (
	port = ":50051"
)

//var endpointRegistry = map[string]common.Endpoint{}

// server is used to implement helloworld.GreeterServer.
type server struct{
	endpointRegistry map[string]common.Endpoint
}

func (s *server) Remind(stream common.ReminderService_RemindServer) error {
	for {

		in, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Node connected: %v", in)
		_, ok := s.endpointRegistry["1"]
		if !ok  {
			s.endpointRegistry["1"] = common.NewServerEndpoint(stream)
		}

//		s.startEndpointReminder(in)
	}
}

//func (s *server) startEndpointReminder(in *common.Reminder) {
//
//	go func() {
//		for {
//			time.Sleep(5 * time.Second)
//			log.Printf("sending reminder %v", in)
//			s.ServerEndpoint.Send(in)
//		}
//	}()
//}


func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	epr := map[string]common.Endpoint{}
	s2 := &server{endpointRegistry: epr}
	common.RegisterReminderServiceServer(s, s2)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	go func(){
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	}()

	for {
		time.Sleep(5 * time.Second)
		if epr["1"] != nil {
			r := common.Reminder{"2", "what time is it",  "server"}
			epr["1"].Send(&r)
		}
	}


}
