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
	"github.com/dev-shimada/grpc-federation-playground/ent"
	"github.com/dev-shimada/grpc-federation-playground/ent/message"
	messagev1 "github.com/dev-shimada/grpc-federation-playground/gen/message/v1"
	"github.com/dev-shimada/grpc-federation-playground/gen/message/v1/messagev1connect"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	Host = "0.0.0.0"
	Port = 8081
)

type server struct {
	*http.Server
	client *ent.Client
}

// domain model
type Message struct {
	ID     string
	UserID string
	Text   string
}

// repository interface
type MessageRepository interface {
	Save(ctx context.Context, msg Message) error
	Get(ctx context.Context, id string) (*Message, error)
}

// ent-based repository implementation
type EntRepository struct {
	client *ent.Client
}

func NewEntRepository(client *ent.Client) MessageRepository {
	return &EntRepository{client: client}
}

func (r *EntRepository) Save(ctx context.Context, msg Message) error {
	userID, err := uuid.Parse(msg.UserID)
	if err != nil {
		return fmt.Errorf("invalid user_id format: %w", err)
	}

	_, err = r.client.Message.
		Create().
		SetID(uuid.MustParse(msg.ID)).
		SetUserID(userID.String()).
		SetText(msg.Text).
		Save(ctx)

	return err
}

func (r *EntRepository) Get(ctx context.Context, id string) (*Message, error) {
	messageID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid message id format: %w", err)
	}

	entMsg, err := r.client.Message.
		Query().
		Where(message.ID(messageID)).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return &Message{
		ID:     entMsg.ID.String(),
		UserID: entMsg.UserID,
		Text:   entMsg.Text,
	}, nil
}

func (s *server) Post(ctx context.Context, req *connect.Request[messagev1.PostRequest]) (*connect.Response[messagev1.PostResponse], error) {
	slog.Info("Received Post request", "user", req.Msg.UserId, "text", req.Msg.Text)
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("Failed to generate UUID", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generate UUID: %w", err))
	}

	msg := Message{
		ID:     id.String(),
		UserID: req.Msg.UserId,
		Text:   req.Msg.Text,
	}

	repository := NewEntRepository(s.client)
	if err := repository.Save(ctx, msg); err != nil {
		slog.Error("Failed to save message", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to save message: %w", err))
	}

	return &connect.Response[messagev1.PostResponse]{
		Msg: &messagev1.PostResponse{
			Id: id.String(),
		},
	}, nil
}

func (s *server) Get(ctx context.Context, req *connect.Request[messagev1.GetRequest]) (*connect.Response[messagev1.GetResponse], error) {
	slog.Info("Received Get request", "id", req.Msg.Id)

	repository := NewEntRepository(s.client)
	msg, err := repository.Get(ctx, req.Msg.Id)
	if err != nil {
		slog.Error("Failed to get message", "error", err)
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("message not found: %w", err))
	}

	return &connect.Response[messagev1.GetResponse]{
		Msg: &messagev1.GetResponse{
			UserId: msg.UserID,
			Text:   msg.Text,
		},
	}, nil
}

func (s *server) PingPong(ctx context.Context, req *connect.Request[messagev1.PingPongRequest]) (*connect.Response[messagev1.PingPongResponse], error) {
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

	// データベース接続の初期化
	client, err := ent.Open("sqlite3", "./message.db?_fk=1")
	if err != nil {
		slog.Error(fmt.Sprintf("failed opening connection to sqlite: %v", err))
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error(fmt.Sprintf("Failed to close database connection: %v", err))
		}
	}()

	// マイグレーション実行
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		slog.Error(fmt.Sprintf("failed creating schema resources: %v", err))
	}

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
		Server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", Host, Port),
			Handler: h2c.NewHandler(mux, &http2.Server{}),
		},
		client: client,
	}

	// message
	path, handler := messagev1connect.NewMessageServiceHandler(svc)
	mux.Handle(path, handler)

	// start server
	signalCtx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
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
	<-signalCtx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	slog.Info("Shutting down server...")
	if err := svc.Shutdown(shutdownCtx); err != nil {
		slog.Error(fmt.Sprintf("Failed to shutdown server: %v", err))
	}
}
