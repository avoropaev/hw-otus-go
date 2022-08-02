package cmd

import (
	"context"
	"encoding/json"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/cmd/config"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/rmq"
	rmq_models "github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/rmq/models"
)

func schedulerCommand(ctx context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "scheduler",
		Short: "scheduler",
		RunE:  schedulerCommandRunE(ctx),
	}

	command.Flags().StringVar(&cfgFile, "config", "", "Path to configuration file")

	err := command.MarkFlagRequired("config")
	if err != nil {
		return nil
	}

	return command
}

func schedulerCommandRunE(ctx context.Context) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		configFile := cmd.Flag("config").Value.String()

		cfg, err := config.ParseSchedulerConfig(configFile)
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

		producer := rmq.NewProducer(cfg.Producer.URI, cfg.Producer.Queue)
		defer func() {
			if err := producer.Disconnect(); err != nil {
				log.Error().Err(err).Send()
			}
		}()

		go notifications(ctx, application, producer)
		go cleaner(ctx, application)

		<-ctx.Done()

		return nil
	}
}

func notifications(ctx context.Context, app app.Application, producer *rmq.Producer) {
	log.Info().Msg("notifications producer started")
	defer log.Info().Msg("notifications producer stopped")

	for {
	next:
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(time.Second * 10):
			events, err := app.GetEventForNotify(ctx)
			if err != nil {
				log.Error().Err(err).Send()

				break next
			}

			for _, event := range events {
				message := rmq_models.Notification{
					EventGUID:     event.GUID,
					EventTitle:    event.Title,
					EventStartAt:  event.StartAt,
					EventUserGUID: event.UserGUID,
				}

				body, err := json.Marshal(message)
				if err != nil {
					log.Error().Err(err).Send()

					continue
				}

				err = producer.Publish(ctx, body)
				if err != nil {
					log.Error().Err(err).Send()
				}

				l := log.With().Fields(map[string]interface{}{
					"event_guid": event.GUID.String(),
				}).Logger()

				l.Info().Msg("notification sent to queue")
			}
		}
	}
}

func cleaner(ctx context.Context, app app.Application) {
	log.Info().Msg("cleaner started")
	defer log.Info().Msg("cleaner stopped")

	for {
	next:
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(time.Second * 10):
			rowsAffected, err := app.RemoveOldEvents(ctx)
			if err != nil {
				log.Error().Err(err).Send()

				break next
			}

			l := log.With().Fields(map[string]interface{}{
				"rows_affected": rowsAffected,
			}).Logger()

			l.Info().Msg("old events deleted")
		}
	}
}
