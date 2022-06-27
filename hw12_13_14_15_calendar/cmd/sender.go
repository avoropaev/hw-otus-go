package cmd

import (
	"context"
	"encoding/json"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/cmd/config"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/rmq"
	rmq_models "github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/rmq/models"
)

func senderCommand(ctx context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "sender",
		Short: "sender",
		RunE:  senderCommandRunE(ctx),
	}

	command.Flags().StringVar(&cfgFile, "config", "", "Path to configuration file")

	err := command.MarkFlagRequired("config")
	if err != nil {
		return nil
	}

	return command
}

func senderCommandRunE(ctx context.Context) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		configFile := cmd.Flag("config").Value.String()

		cfg, err := config.ParseSenderConfig(configFile)
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

		go func() {
			<-ctx.Done()

			log.Info().Msg("stopping an sender...")
		}()

		consumer := rmq.NewConsumer(
			cfg.Consumer.ConsumerTag,
			cfg.Consumer.URI,
			cfg.Consumer.ExchangeName,
			cfg.Consumer.ExchangeType,
			cfg.Consumer.Queue,
			cfg.Consumer.BindingKey,
		)

		store, err := getStore(ctx, cfg.DB.Type, cfg.DB.PSQL.URL, ctx.Done())
		if err != nil {
			return err
		}

		application := app.New(store)

		log.Info().Msg("starting an sender...")

		err = consumer.Handle(ctx, getWorker(application), cfg.Consumer.Threads)
		if err != nil {
			cancel()

			return err
		}

		return nil
	}
}

func getWorker(app app.Application) rmq.Worker {
	return func(ctx context.Context, messages <-chan amqp.Delivery) {
		log.Info().Msg("worker started")
		defer log.Info().Msg("worker stopped")

		var noti rmq_models.Notification

		for {
		next:
			select {
			case message := <-messages:
				if len(message.Body) == 0 {
					break next
				}

				err := json.Unmarshal(message.Body, &noti)
				if err != nil {
					log.Error().Err(err).Send()

					messageNackWithLog(message)

					break next
				}

				err = app.SendNotificationAndMarkAsNotifier(ctx, noti.EventGUID, noti.EventTitle, noti.EventStartAt, noti.EventUserGUID)
				if err != nil {
					log.Error().Err(err).Send()

					messageNackWithLog(message)

					break next
				}

				log.Info().Msg("send ack")
				err = message.Ack(false)
				if err != nil {
					log.Error().Err(err).Send()

					return
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

func messageNackWithLog(message amqp.Delivery) {
	log.Info().Msg("send nack")

	err := message.Nack(false, false)
	if err != nil {
		log.Error().Err(err).Send()
	}
}
