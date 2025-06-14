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
	"github.com/dev-shimada/grpc-federation-playground/user/ent"
	"github.com/dev-shimada/grpc-federation-playground/user/ent/user"

	userv1 "github.com/dev-shimada/grpc-federation-playground/user/gen/user/v1"

	"github.com/dev-shimada/grpc-federation-playground/user/gen/user/v1/userv1connect"
	"github.com/google/uuid"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	Host = "0.0.0.0"
	Port = 8080
)

type server struct {
	*http.Server
	client *ent.Client
}

// domain model
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt string
	UpdatedAt string
}

// repository interface
type UserRepository interface {
	Save(ctx context.Context, msg User) error
	Get(ctx context.Context, id string) (*User, error)
}

// ent-based repository implementation
type EntRepository struct {
	client *ent.Client
}

func NewEntRepository(client *ent.Client) UserRepository {
	return &EntRepository{client: client}
}

func (r *EntRepository) Save(ctx context.Context, user User) error {
	ID, err := uuid.Parse(user.ID)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	_, err = r.client.User.
		Create().
		SetID(uuid.MustParse(ID.String())).
		SetEmail(user.Email).
		SetName(user.Name).
		Save(ctx)

	return err
}

func (r *EntRepository) Get(ctx context.Context, id string) (*User, error) {
	ID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id format: %w", err)
	}

	entMsg, err := r.client.User.
		Query().
		Where(user.ID(ID)).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return &User{
		ID:        entMsg.ID.String(),
		Email:     entMsg.Email,
		Name:      entMsg.Name,
		CreatedAt: entMsg.CreatedAt.Format(time.RFC3339),
		UpdatedAt: entMsg.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *server) Post(ctx context.Context, req *connect.Request[userv1.PostRequest]) (*connect.Response[userv1.PostResponse], error) {
	slog.Info("Received Post request", "email", req.Msg.Email, "name", req.Msg.Name)
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("Failed to generate UUID", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generate UUID: %w", err))
	}

	u := User{
		ID:    id.String(),
		Email: req.Msg.Email,
		Name:  req.Msg.Name,
	}

	repository := NewEntRepository(s.client)
	if err := repository.Save(ctx, u); err != nil {
		slog.Error("Failed to save user", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to save user: %w", err))
	}

	return &connect.Response[userv1.PostResponse]{
		Msg: &userv1.PostResponse{
			Id:        id.String(),
			Email:     u.Email,
			Name:      u.Name,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
	}, nil
}

func (s *server) Get(ctx context.Context, req *connect.Request[userv1.GetRequest]) (*connect.Response[userv1.GetResponse], error) {
	slog.Info("Received Get request", "id", req.Msg.Id)

	repository := NewEntRepository(s.client)
	u, err := repository.Get(ctx, req.Msg.Id)
	if err != nil {
		slog.Error("Failed to get user", "error", err)
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found: %w", err))
	}

	return &connect.Response[userv1.GetResponse]{
		Msg: &userv1.GetResponse{
			Id:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
	}, nil
}

func (s *server) PingPong(ctx context.Context, req *connect.Request[userv1.PingPongRequest]) (*connect.Response[userv1.PingPongResponse], error) {
	res := connect.NewResponse(&userv1.PingPongResponse{
		Email: req.Msg.Email,
		Name:  req.Msg.Name,
	})
	res.Header().Set("User-Version", "v1")
	return res, nil
}

func main() {
	// json logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(log.Writer(), nil)))

	// データベース接続の初期化
	client, err := ent.Open("sqlite3", "./user.db?_fk=1")
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
		userv1connect.UserServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// health check
	checker := grpchealth.NewStaticChecker(
		userv1connect.UserServiceName,
	)
	mux.Handle(grpchealth.NewHandler(checker))

	svc := &server{
		Server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", Host, Port),
			Handler: h2c.NewHandler(mux, &http2.Server{}),
		},
		client: client,
	}

	// user
	path, handler := userv1connect.NewUserServiceHandler(svc)
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
