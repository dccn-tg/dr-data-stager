package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/dccn-tg/dr-data-stager/internal/worker/config"
	"github.com/dccn-tg/dr-data-stager/pkg/tasks"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"

	"github.com/hibiken/asynq"
)

// internal data structure for a cancelled task
type ct struct {
	queue string
	id    string
}

// notification mode
type nmode int

const (
	nFailed nmode = iota
	nCompleted
)

func (m nmode) String() string {
	switch m {
	case nFailed:
		return "failed"
	case nCompleted:
		return "completed"
	}
	return "unknown"
}

// Notifier is a asynq middleware to handle cancelled jobs and sending out email notification.
func Notifier(inspector *asynq.Inspector, cfg config.Configuration) func(asynq.Handler) asynq.Handler {

	// SMTP mailer
	client := stagerMailer{
		config: cfg.Mailer,
	}

	cct := make(chan ct)

	// internal go routine to archive cancelled task
	go func() {

		log.Debugf("task archiver started\n")

		defer log.Debugf("task archiver stopped\n")

		var tinfo *asynq.TaskInfo
		var err error

	loop:
		for t := range cct {

			tinfo, err = inspector.GetTaskInfo(t.queue, t.id)
			if err != nil {
				log.Errorf("cannot get task info of task %s in queue %s: %s\n", t.id, t.queue, err)
				continue loop
			}

			for {
				if tinfo.State != asynq.TaskStateActive {
					break
				}

				log.Debugf("task %s status %s\n", t.id, tinfo.State)

				tinfo, err = inspector.GetTaskInfo(t.queue, t.id)
				if err != nil {
					log.Errorf("cannot get task info of task %s in queue %s: %s\n", t.id, t.queue, err)
					continue loop
				}
			}

			err = inspector.ArchiveTask(t.queue, t.id)
			if err != nil {
				log.Errorf("cannot archive cancelled task %s: %s\n", t.id, err)
			}
		}
	}()

	return func(next asynq.Handler) asynq.Handler {

		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {

			err := next.ProcessTask(ctx, t)

			switch {
			case err == nil:
				log.Debugf("job completed, notifying job owner: %+v\n", client)
				id, _ := asynq.GetTaskID(ctx)
				q, _ := asynq.GetQueueName(ctx)
				if tinfo, err := inspector.GetTaskInfo(q, id); err != nil {
					log.Errorf("cannot get task %s: %s\n", id, err)
					break
				} else {
					sendEmailNotification(&client, tinfo, nCompleted)
				}
			case err == asynq.SkipRetry:
				log.Debugf("job retry skipped")
			case errors.Is(err, context.Canceled):
				log.Debugf("job canceled")

				id, _ := asynq.GetTaskID(ctx)
				q, ok := asynq.GetQueueName(ctx)
				if !ok {
					log.Errorf("cannot find queue name for cancelled job: %s\n", id)
					break
				}

				cct <- ct{
					queue: q,
					id:    id,
				}
			case errors.Is(err, context.DeadlineExceeded):
				log.Debugf("job exceeded deadline")
			default:
				retried, _ := asynq.GetRetryCount(ctx)
				maxRetry, _ := asynq.GetMaxRetry(ctx)
				if retried >= maxRetry {
					log.Debugf("job failued after retries, notifying job ower and admin: %+v\n", client)
					id, _ := asynq.GetTaskID(ctx)
					q, _ := asynq.GetQueueName(ctx)
					if tinfo, err := inspector.GetTaskInfo(q, id); err != nil {
						log.Errorf("cannot get task %s: %s\n", id, err)
						break
					} else {
						sendEmailNotification(&client, tinfo, nFailed, cfg.Admins...)
					}
				}
			}

			return err
		})
	}
}

func sendEmailNotification(client *stagerMailer, tinfo *asynq.TaskInfo, nt nmode, cc ...string) {

	var p tasks.StagerPayload
	if err := json.Unmarshal(tinfo.Payload, &p); err != nil {
		log.Errorf("fail to unmarshal task payload %s: %s\n", tinfo.ID, err)
		return
	}

	idparts := strings.Split(tinfo.ID, ".")

	recipients := []string{}
	if p.StagerUserEmail != "" {
		recipients = append(recipients, p.StagerUserEmail)
	}

	var subject string
	switch nt {
	case nFailed:
		subject = fmt.Sprintf("[ALERT] stager job %s failed", idparts[len(idparts)-1])
		// in case there is no recipient, take the first one in `cc` list if it is possible
		if len(recipients) == 0 && len(cc) > 0 {
			recipients = append(recipients, cc[0])
			cc = cc[1:]
		}
	case nCompleted:
		subject = fmt.Sprintf("[OK] stager job %s completed", idparts[len(idparts)-1])
	}

	if len(recipients) == 0 {
		log.Warnf("no notification recipient for task %s\n", tinfo.ID)
		return
	}

	body := composeMailBody(tinfo, p, nt)

	err := client.SendHtmlMail(
		"datasupport@donders.ru.nl",
		subject,
		body,
		recipients,
		cc...,
	)

	if err != nil {
		log.Errorf("fail to send out notification: %s\n", err)
	}
}

func composeMailBody(tinfo *asynq.TaskInfo, payload tasks.StagerPayload, mode nmode) string {

	var tmsg, ntStr string

	tcompleted := tinfo.CompletedAt.Truncate(time.Second)
	tlastfailed := tinfo.LastFailedAt.Truncate(time.Second)

	idparts := strings.Split(tinfo.ID, ".")

	switch mode {
	case nFailed:
		tmsg = templateNotificationFailed
		ntStr = "failed"
		tlastfailed = time.Now().Truncate(time.Second)
	case nCompleted:
		tmsg = templateNotificationCompleted
		ntStr = "completed"
		tcompleted = time.Now().Truncate(time.Second)
	}

	t, err := template.New("msg").Parse(tmsg)
	if err != nil {
		log.Errorf("cannot parse email template for %s job %s: %s\n", ntStr, tinfo.ID, err)
		return fmt.Sprintf("stager job %s %s", tinfo.ID, ntStr)
	}

	var rslt tasks.StagerTaskResult
	if err := json.Unmarshal(tinfo.Result, &rslt); err != nil {
		log.Errorf("fail to unmarshal task result %s: %s\n", tinfo.ID, err)
		//return fmt.Sprintf("stager job %s %s", tinfo.ID, ntStr)
	}

	buf := new(bytes.Buffer)

	err = t.Execute(buf, DataNotification{
		ID:           idparts[len(idparts)-1],
		State:        mode,
		StagerUser:   payload.StagerUser,
		DrUser:       payload.DrUser,
		SrcURL:       payload.SrcURL,
		DstURL:       payload.DstURL,
		CreatedAt:    time.Unix(payload.CreatedAt, 0),
		CompletedAt:  tcompleted,
		LastFailedAt: tlastfailed,
		LastErr:      tinfo.LastErr,
		Result:       rslt,
	})
	if err != nil {
		log.Errorf("cannot compose email for %s job %s: %s\n", ntStr, tinfo.ID, err)
		return fmt.Sprintf("stager job %s %s", tinfo.ID, ntStr)
	}

	return buf.String()
}
