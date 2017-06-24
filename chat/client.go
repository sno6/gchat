package chat

import (
	"context"

	"google.golang.org/grpc"
)

// Client represents a user chatting via grpc.
type Client struct {
	Name   string
	stream ChatService_CommunicateClient
	conn   *grpc.ClientConn
}

// A Notifier notifies the user interface with a message.
type Notifier interface {
	Notify(*Message)
}

func (c *Client) register() error {
	return c.stream.Send(&Message{User: c.Name, Register: true})
}

// Send sends a message via the server stream.
func (c *Client) Send(data string) error {
	return c.stream.Send(&Message{User: c.Name, Data: data})
}

// Close terminates a grpc connection.
func (c *Client) Close() error {
	if err := c.stream.Send(&Message{User: c.Name, Close: true}); err != nil {
		return err
	}
	return c.conn.Close()
}

// ReadPump reads from the server and if it receives anything it notifies the notifier
// which in this case is the CUI.
func (c *Client) ReadPump(n Notifier) error {
	for {
		msg, err := c.stream.Recv()
		if err != nil {
			return err
		}
		n.Notify(msg)
	}
}

// NewClient creates a new grpc client.
func NewClient(addr string, name string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := NewChatServiceClient(conn)
	stream, err := client.Communicate(context.Background())
	if err != nil {
		return nil, err
	}
	c := &Client{Name: name, conn: conn, stream: stream}
	if err := c.register(); err != nil {
		return nil, err
	}
	return c, nil
}
