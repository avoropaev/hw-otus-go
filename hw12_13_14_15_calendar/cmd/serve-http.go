package cmd

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	pgx "github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/cmd/config"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/app"
	internalgrpc "github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	psqlstorage "github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

func serveHTTPCommand(ctx context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "serve-http",
		Short: "serves http api",
		RunE:  serveHTTPCommandRunE(ctx),
	}

	command.Flags().StringVar(&cfgFile, "config", "", "Path to configuration file")

	err := command.MarkFlagRequired("config")
	if err != nil {
		return nil
	}

	return command
}

func serveHTTPCommandRunE(ctx context.Context) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		configFile := cmd.Flag("config").Value.String()

		cfg, err := config.ParseConfig(configFile)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse config")

			return err
		}

		logLevel, err := zerolog.ParseLevel(cfg.Logger.Level)
		if err != nil {
			log.Error().Err(err).Msg("failed to install log level")

			return err
		}

		zerolog.SetGlobalLevel(logLevel)

		ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		store, err := getStore(ctx, cfg.DB.Type, cfg.DB.PSQL.URL, ctx.Done())
		if err != nil {
			return err
		}

		application := app.New(store)
		grpcServer := internalgrpc.NewServer(cfg.GRPC.Host, cfg.GRPC.Port, application)
		httpServer := internalhttp.NewServer(cfg.HTTP.Host, cfg.HTTP.Port, cfg.GRPC.Host+":"+strconv.Itoa(cfg.GRPC.Port), application)

		go func() {
			<-ctx.Done()

			log.Info().Msg("stopping an http server...")

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			if err := httpServer.Stop(ctx); err != nil {
				log.Error().Err(err).Msg("failed to stop http server")
			}

			log.Info().Msg("stopping an grpc server...")

			ctx, cancel = context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			if err := grpcServer.Stop(ctx); err != nil {
				log.Error().Err(err).Msg("failed to stop grpc server")
			}
		}()

		log.Info().Msg("calendar is running...")

		go func() {
			if err := grpcServer.Start(ctx); err != nil {
				cancel()

				if !errors.Is(err, http.ErrServerClosed) {
					log.Error().Err(err).Msg("failed to start grpc server")
				}
			}
		}()

		if err := httpServer.Start(ctx); err != nil {
			cancel()

			if !errors.Is(err, http.ErrServerClosed) {
				log.Error().Err(err).Msg("failed to start http server")

				return err
			}
		}

		return nil
	}
}

func getStore(ctx context.Context, dbType string, psqlURL string, done <-chan struct{}) (app.Storage, error) {
	var store app.Storage

	switch dbType {
	case "psql":
		conn, err := pgx.Connect(ctx, psqlURL)
		if err != nil {
			log.Error().Err(err).Msg("unable to connect to database")

			return nil, err
		}

		go func() {
			<-done

			err := conn.Close(ctx)
			if err != nil {
				log.Error().Err(err).Msg("unable to close connect to database")
			}
		}()

		err = conn.Ping(ctx)
		if err != nil {
			log.Error().Err(err).Msg("unable to connect to database")

			return nil, err
		}

		store = psqlstorage.New(conn)
	case "memory":
		store = memorystorage.New()
	default:
		err := errors.New("unknown db type")
		log.Error().Err(err).Send()

		return nil, err
	}

	return store, nil
}
