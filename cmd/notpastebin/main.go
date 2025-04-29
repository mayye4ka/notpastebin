package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/mayye4ka/notpastebin/internal/config"
	"github.com/mayye4ka/notpastebin/internal/db"
	"github.com/mayye4ka/notpastebin/internal/server"
	"github.com/mayye4ka/notpastebin/internal/service"
	api "github.com/mayye4ka/notpastebin/pkg/api/go"
	"github.com/rs/zerolog"
	migrate "github.com/rubenv/sql-migrate"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func applyMigrations(cfg *config.Config) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.DbUser, cfg.DbPassword, cfg.DbAddr, cfg.DbName)

	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("can't connect: %w", err)
	}

	_, err = migrate.Exec(conn, "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("can't migrate: %w", err)
	}
	return nil
}

func main() {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zerolog.New(os.Stdout)
	config, err := config.ReadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("can't read config")
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?pool_max_conns=10",
		config.DbUser, config.DbPassword, config.DbAddr, config.DbName)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		logger.Fatal().Err(err).Msg("can't read config")
	}

	err = applyMigrations(config)
	if err != nil {
		logger.Fatal().Err(err).Msg("can't apply migrations")
	}

	db := db.New(pool, &logger)
	service := service.New(db)
	server := server.New(service)

	eg, ectx := errgroup.WithContext(ctx)

	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	api.RegisterNotPasteBinServer(grpcServer, server)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcPort))
	if err != nil {
		logger.Fatal().Err(err).Msg("can't create listener")
	}

	eg.Go(func() error {

		err = grpcServer.Serve(lis)
		if err != nil {
			return fmt.Errorf("serve grpc: %w", err)
		}
		return nil
	})

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = api.RegisterNotPasteBinHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", config.GrpcPort), opts)
	if err != nil {
		logger.Fatal().Err(err).Msg("register auth handler from endpoint")
	}
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpPort),
		Handler: mux,
	}

	eg.Go(func() error {
		err = httpServer.ListenAndServe()
		if err != nil {
			return fmt.Errorf("serve http: %w", err)
		}
		return nil
	})

	shutdown := func() {
		shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), time.Minute)
		defer shutdownCancel()
		httpServer.Shutdown(shutdownContext)
		grpcServer.GracefulStop()
		cancel()
		pool.Close()
	}

	shutdownDone := make(chan struct{})
	go func() {
		select {
		case <-termChan:
		case <-ectx.Done():
		}
		shutdown()
		close(shutdownDone)
	}()

	err = eg.Wait()
	if err != nil {
		logger.Error().Err(err).Msg("eg wait error")
	}

	<-shutdownDone
}
