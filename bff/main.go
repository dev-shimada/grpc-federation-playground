package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	bffv1 "github.com/dev-shimada/grpc-federation-playground/bff/gen/bff/v1"
	messagev1 "github.com/dev-shimada/grpc-federation-playground/bff/gen/message/v1"
	userv1 "github.com/dev-shimada/grpc-federation-playground/bff/gen/user/v1"
)

const (
	Host               = "0.0.0.0"
	Port               = 8080
	UserServiceAddr    = "user:8082"
	MessageServiceAddr = "message:8081"
)

// BffServiceClientFactory はgrpc-federationで必要なクライアントファクトリを実装
type BffServiceClientFactory struct {
	userConn    *grpc.ClientConn
	messageConn *grpc.ClientConn
}

// User_V1_UserServiceClient はUserServiceクライアントを作成
func (f *BffServiceClientFactory) User_V1_UserServiceClient(cfg bffv1.BffServiceClientConfig) (userv1.UserServiceClient, error) {
	return userv1.NewUserServiceClient(f.userConn), nil
}

// Message_V1_MessageServiceClient はMessageServiceクライアントを作成
func (f *BffServiceClientFactory) Message_V1_MessageServiceClient(cfg bffv1.BffServiceClientConfig) (messagev1.MessageServiceClient, error) {
	return messagev1.NewMessageServiceClient(f.messageConn), nil
}

type server struct {
	grpcServer    *grpc.Server
	bffService    *bffv1.BffService
	clientFactory *BffServiceClientFactory
}

func newServer() (*server, error) {
	logger := slog.Default()

	// UserServiceへの接続
	userConn, err := grpc.NewClient(
		UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	// MessageServiceへの接続
	messageConn, err := grpc.NewClient(
		MessageServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		userConn.Close()
		return nil, fmt.Errorf("failed to connect to message service: %w", err)
	}

	// ClientFactoryの作成
	clientFactory := &BffServiceClientFactory{
		userConn:    userConn,
		messageConn: messageConn,
	}

	// BffServiceの設定
	config := bffv1.BffServiceConfig{
		Client: clientFactory,
		Logger: logger,
	}

	// BffServiceの作成
	bffService, err := bffv1.NewBffService(config)
	if err != nil {
		userConn.Close()
		messageConn.Close()
		return nil, fmt.Errorf("failed to create BFF service: %w", err)
	}

	// gRPCサーバーの作成
	grpcServer := grpc.NewServer()

	// BffServiceをgRPCサーバーに登録
	bffv1.RegisterBffServiceServer(grpcServer, bffService)

	// リフレクションサービスを有効化（開発用）
	reflection.Register(grpcServer)

	return &server{
		grpcServer:    grpcServer,
		bffService:    bffService,
		clientFactory: clientFactory,
	}, nil
}

func (s *server) start() error {
	// リスナーの作成
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", Host, Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	slog.Info("Starting BFF server", "address", fmt.Sprintf("%s:%d", Host, Port))

	// gRPCサーバーの開始
	if err := s.grpcServer.Serve(lis); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *server) shutdown(ctx context.Context) error {
	slog.Info("Shutting down BFF server...")

	// grpc-federationのクリーンアップ
	if s.bffService != nil {
		bffv1.CleanupBffService(ctx, s.bffService)
	}

	// gRPCサーバーのグレースフルシャットダウン
	if s.grpcServer != nil {
		done := make(chan struct{})
		go func() {
			s.grpcServer.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			slog.Info("gRPC server stopped gracefully")
		case <-ctx.Done():
			slog.Warn("Force stopping gRPC server due to timeout")
			s.grpcServer.Stop()
		}
	}

	// 外部サービスへの接続を閉じる
	if s.clientFactory != nil {
		if s.clientFactory.userConn != nil {
			s.clientFactory.userConn.Close()
		}
		if s.clientFactory.messageConn != nil {
			s.clientFactory.messageConn.Close()
		}
	}

	return nil
}

func main() {
	// JSONロガーの設定
	slog.SetDefault(slog.New(slog.NewJSONHandler(log.Writer(), nil)))

	// サーバーの作成
	srv, err := newServer()
	if err != nil {
		slog.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	// グレースフルシャットダウンのためのシグナルハンドリング
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		slog.Info("Received shutdown signal")

		// シャットダウンタイムアウトの設定
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := srv.shutdown(shutdownCtx); err != nil {
			slog.Error("Failed to shutdown server gracefully", "error", err)
			os.Exit(1)
		}

		os.Exit(0)
	}()

	// サーバーの開始
	if err := srv.start(); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}

	slog.Info("BFF server stopped")
}
