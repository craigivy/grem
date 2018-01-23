package main

import (
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/craigivy/grem/pkg/common"
	"io"
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

	reminders := []*common.Reminder{
		{"1", "First reminder"},
		{"2", "Second reminder"},
		{"3", "Third reminder"},
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
			log.Printf("Got reminder %s with ID(%s)", in.Note, in.ID)
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
