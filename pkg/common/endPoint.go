package common

import (
	"io"
	"log"
)

type Endpoint interface {
	Send(reminder *Reminder) error
}

type ServerEndpoint struct {
	stream ReminderService_RemindServer
}

func NewServerEndpoint(stream ReminderService_RemindServer) *ServerEndpoint {

	se := ServerEndpoint{}
	se.stream = stream
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				//				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a reminder : %v", err)
			}
			log.Printf("Got reminder %s for (%s)", in.Note, in.NodeID)
		}
	}()
	return &se

}

func (se ServerEndpoint) Send(reminder *Reminder) error {
	return se.stream.Send(reminder)
}



type ClientEndpoint struct {
	stream ReminderService_RemindClient
}

func NewClientEndpoint(stream ReminderService_RemindClient) ClientEndpoint {

	ce := ClientEndpoint{}
	ce.stream = stream
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
//				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a reminder : %v", err)
			}
			log.Printf("Got reminder %s for (%s)", in.Note, in.NodeID)
		}
	}()
	return ce


}

func (ce ClientEndpoint) Send(reminder *Reminder) error {
	return ce.stream.Send(reminder)
}




