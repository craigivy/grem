package main

import (
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/craigivy/grem/pkg/common"
	"math/rand"
	"time"
	"strconv"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := common.NewReminderServiceClient(conn)
	remind(c)

	forever := make(chan struct{})
	<-forever

}

func remind(client common.ReminderServiceClient) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	nodeID := strconv.Itoa(r1.Intn(100))
	log.Printf("I'm node %s", nodeID)

	reminders := []*common.Reminder{
		{"1", "First reminder", nodeID},
		{"2", "Second reminder",  nodeID},
		{"3", "Third reminder",  nodeID},
	}

	stream, err := client.Remind(context.Background())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	cep := common.NewClientEndpoint(stream)

	for {
		err := cep.Send(reminders[0])
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		time.Sleep(5 * time.Second)
	}

}
