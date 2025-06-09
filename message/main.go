package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	messagev1 "github.com/dev-shimada/grpc-federation-playground/gen/message/v1"
	"github.com/dev-shimada/grpc-federation-playground/gen/message/v1/messagev1connect"
	"github.com/google/uuid"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	Host = "0.0.0.0"
	Port = 8081
)

type server struct {
	*http.Server
}

// domain model
type Message struct {
	ID   string
	Text string
}

// repository interface
type MessageRepository interface {
	Save(Message) error
	Get(id string) (*Message, error)
}

// infrastructure implementation
type Repository struct {
	ID   string
	Text string
}

var repo []Message = make([]Message, 0)

func (mr Repository) Save(msg Message) error {
	repo = append(repo, msg)
	return nil
}
func (mr Repository) Get(id string) (*Message, error) {
	for _, msg := range repo {
		if msg.ID == id {
			return &msg, nil
		}
	}
	return nil, fmt.Errorf("message with id %s not found", id)
}

func (s server) Post(ctx context.Context, req *connect.Request[messagev1.PostRequest]) (*connect.Response[messagev1.PostResponse], error) {
	slog.Info("Received Post request", "user", req.Msg.UserId, "text", req.Msg.Text)
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("Failed to generate UUID", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generate UUID: %w", err))
	}
	msg := Message{
		ID:   id.String(),
		Text: req.Msg.Text,
	}
	repository := Repository{
		ID:   id.String(),
		Text: req.Msg.Text,
	}
	if err := repository.Save(msg); err != nil {
		slog.Error("Failed to save message", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to save message: %w", err))
	}
	return &connect.Response[messagev1.PostResponse]{
		Msg: &messagev1.PostResponse{
			Id: &messagev1.UUID{Value: []byte(id.String())},
		},
	}, nil
}

func (s server) Get(ctx context.Context, req *connect.Request[messagev1.GetRequest]) (*connect.Response[messagev1.GetResponse], error) {
	slog.Info("Received Get request", "id", req.Msg.Id)
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("Failed to generate UUID", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generate UUID: %w", err))
	}

	return &connect.Response[messagev1.GetResponse]{
		Msg: &messagev1.GetResponse{
			UserId: &messagev1.UUID{Value: []byte(id.String())},
			Text:   "This is a sample message",
		},
	}, nil
}

func (s server) PingPong(ctx context.Context, req *connect.Request[messagev1.PingPongRequest]) (*connect.Response[messagev1.PingPongResponse], error) {
	res := connect.NewResponse(&messagev1.PingPongResponse{
		UserId: req.Msg.UserId,
		Text:   req.Msg.Text,
	})
	res.Header().Set("Message-Version", "v1")
	return res, nil
}

func main() {
	// json logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(log.Writer(), nil)))

	mux := http.NewServeMux()

	// reflection
	reflector := grpcreflect.NewStaticReflector(
		messagev1connect.MessageServicePingPongProcedure,
		messagev1connect.MessageServicePostProcedure,
		messagev1connect.MessageServiceGetProcedure,
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// health check
	checker := grpchealth.NewStaticChecker(
		messagev1connect.MessageServicePingPongProcedure,
		messagev1connect.MessageServicePostProcedure,
		messagev1connect.MessageServiceGetProcedure,
	)
	mux.Handle(grpchealth.NewHandler(checker))

	svc := &server{
		&http.Server{
			Addr:    fmt.Sprintf("%s:%d", Host, Port),
			Handler: h2c.NewHandler(mux, &http2.Server{}),
		},
	}

	// message
	path, handler := messagev1connect.NewMessageServiceHandler(svc)
	mux.Handle(path, handler)

	// start server
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	slog.Info(fmt.Sprintf("Server is running at %s:%d", Host, Port))
	go func() {
		if err := svc.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				slog.Info("Server closed.")
			} else {
				slog.Error(fmt.Sprintf("Failed to serve: %v", err))
			}
		}
	}()
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	slog.Info("Shutting down server...")
	if err := svc.Shutdown(shutdownCtx); err != nil {
		slog.Error(fmt.Sprintf("Failed to shutdown server: %v", err))
	}
}
