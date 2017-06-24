package chat

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"io/ioutil"

	"google.golang.org/grpc"
)

// A Server enables people to chat with eachother.
type Server struct {
	name        string
	addr        string
	logger      *log.Logger
	mu          sync.Mutex
	connections map[ChatService_CommunicateServer]string
}

// NewServer returns a new server.
func NewServer(addr string, logger *log.Logger) *Server {
	if logger == nil {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}
	return &Server{
		name:        "Server",
		addr:        addr,
		logger:      logger,
		connections: make(map[ChatService_CommunicateServer]string),
	}
}

// Run regsisters the chat server and starts listening.
func (s *Server) Run() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	RegisterChatServiceServer(srv, s)
	s.logger.Printf("GRPC server running at: %s\n", s.addr)
	return srv.Serve(l)
}

// Communicate implements the grpc chat server stub.
func (s *Server) Communicate(stream ChatService_CommunicateServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}

		if msg.Register {
			if err := s.register(msg.User, stream); err != nil {
				return err
			}
		}

		if msg.Data != "" {
			if err := s.broadcast(msg); err != nil {
				return err
			}
		}

		if msg.Close {
			if err := s.close(msg.User, stream); err != nil {
				return err
			}
		}
	}
}

func (s *Server) register(user string, stream ChatService_CommunicateServer) error {
	// Simple username validation.
	if user == "" || user == s.name || len(user) > 20 {
		return errors.New("server: invalid username")
	}

	if err := s.addUser(stream, user); err != nil {
		return err
	}

	m := &Message{User: s.name, Data: fmt.Sprintf("%s has joined the chat.\n", user), Register: true}
	return s.broadcast(m)
}

func (s *Server) addUser(stream ChatService_CommunicateServer, user string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, name := range s.connections {
		if name == user {
			return errors.New("server: user with that name already registered")
		}
	}

	s.connections[stream] = user
	return nil
}

func (s *Server) broadcast(msg *Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for stream := range s.connections {
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) close(user string, ownStream ChatService_CommunicateServer) error {
	if err := s.removeUser(ownStream); err != nil {
		return err
	}

	m := &Message{User: s.name, Data: fmt.Sprintf("%s has left the chat.\n", user), Register: true}
	return s.broadcast(m)
}

func (s *Server) removeUser(ownStream ChatService_CommunicateServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.connections[ownStream]; !ok {
		return errors.New("server: error closing, stream not in connection list")
	}

	delete(s.connections, ownStream)
	return nil
}
