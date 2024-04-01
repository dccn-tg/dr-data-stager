package middleware

import (
	"context"
	"errors"

	"github.com/Donders-Institute/dr-data-stager/internal/worker/config"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"

	"github.com/dccn-tg/tg-toolset-golang/pkg/mailer"
	"github.com/hibiken/asynq"
)

func Notifier(cfg config.Configuration) func(asynq.Handler) asynq.Handler {

	// SMTP mailer
	client := mailer.New(cfg.Mailer)

	return func(next asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {

			err := next.ProcessTask(ctx, t)

			switch {
			case err == nil:
				log.Debugf("job completed, notifying job owner: %+v\n", client)
			case err == asynq.SkipRetry:
				log.Debugf("job retry skipped")
			case errors.Is(err, context.Canceled):
				log.Debugf("job canceled")
			case errors.Is(err, context.DeadlineExceeded):
				log.Debugf("job exceeded deadline")
			default:
				retried, _ := asynq.GetRetryCount(ctx)
				maxRetry, _ := asynq.GetMaxRetry(ctx)
				if retried >= maxRetry {
					log.Debugf("job failued after retries, notifying job ower and admin: %+v\n", client)
				}
			}

			return err
		})
	}
}
