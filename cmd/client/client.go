package main

import (
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/craigivy/grem/pkg/common"
	"io"
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
		log.Fatalf("error client: %v, err: %v", client, err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a reminder : %v", err)
			}
			log.Printf("Got reminder %s for (%s)", in.Note, in.NodeID)
		}
	}()

	for _, reminder := range reminders {
		if err := stream.Send(reminder); err != nil {
			log.Fatalf("Failed to send a reminder: %v", err)
		}
		log.Printf("reminder sent: %v", reminder)
	}

	forever := make(chan struct{})
	<-forever

	stream.CloseSend()
	<-waitc
}
