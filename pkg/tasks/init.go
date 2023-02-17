package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

// A list of task types.
const (
	TypeStager  = "stager"
	TypeEmailer = "email"
)

// StagerPayload defines the data structure of the stager file transfer payload.
type StagerPayload struct {
	// short description about the job
	Title string `json:"title"`

	// username of the DR account
	DrUser string `json:"drUser"`

	// path or DR namespace (prefixed with irods:) of the destination endpoint
	DstURL string `json:"dstURL"`

	// path or DR namespace (prefixed with irods:) of the source endpoint
	SrcURL string `json:"srcURL"`

	// username of stager's local account
	StagerUser string `json:"stagerUser"`

	// allowed duration in seconds for entire transfer job (0 for no timeout)
	Timeout int64 `json:"timeout,omitempty"`

	// allowed duration in seconds for no further transfer progress (0 for no timeout)
	TimeoutNoprogress int64 `json:"timeout_noprogress,omitempty"`
}

// EmailerPayload defines the data structure of the emailer payload.
type EmailerPayload struct {
	Recipients []string
	Message    string
}

func NewEmailNotifyTask(Recipients []string, Message string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailerPayload{
		Recipients: Recipients,
		Message:    Message,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailer, payload), nil
}

func NewStagerTask(Title, DrUser, DstURL, SrcURL, StagerUser string, Timeout, TimeoutNoprogress int64) (*asynq.Task, error) {
	payload, err := json.Marshal(StagerPayload{
		Title:             Title,
		DrUser:            DrUser,
		DstURL:            DstURL,
		SrcURL:            SrcURL,
		StagerUser:        StagerUser,
		Timeout:           Timeout,
		TimeoutNoprogress: TimeoutNoprogress,
	})
	if err != nil {
		return nil, err
	}
	// task options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(
		TypeStager,
		payload,
		asynq.MaxRetry(5),
		asynq.Timeout(time.Duration(Timeout)*time.Second),
	), nil
}

// Emailer implements asynq.Handler interface.
type Emailer struct {
}

func (emailer *Emailer) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p EmailerPayload

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("cannot unmarshal emailer payload: %v: %w", err, asynq.SkipRetry)
	}
	log.Debugf("[%s] emailer payload data: %+v", t.ResultWriter().TaskID(), p)
	return nil
}

func NewEmailer() *Emailer {
	// load mail server configuration
	return &Emailer{}
}

// Stager implements asynq.Handler interface.
type Stager struct {
}

func (stager *Stager) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p StagerPayload

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("cannot unmarshal stager payload: %v: %w", err, asynq.SkipRetry)
	}
	log.Debugf("[%s] payload: %+v", t.ResultWriter().TaskID(), p)

	ticker := time.NewTicker(1 * time.Second)
	stopMonitor := make(chan error, 1)

	total := make(chan int64, 1)
	processed := make(chan int64, 1)

	// progress monitor routine
	go func() {
		rslt := new(StagerTaskResult)
		for {
			select {
			case total := <-total:
				// update total steps
				rslt.Progress.Total = total
				if d, err := json.Marshal(rslt); err == nil {
					t.ResultWriter().Write(d)
				}
			case processed := <-processed:
				// update processed steps
				rslt.Progress.Processed = int64(processed)
				if d, err := json.Marshal(rslt); err == nil {
					t.ResultWriter().Write(d)
				}
			case err := <-stopMonitor:
				// stop monitor and write error to the task result
				if err != nil {
					rslt.Error = err.Error()
					if d, err := json.Marshal(rslt); err == nil {
						t.ResultWriter().Write(d)
					}
				}
				return
			case <-ticker.C:
				// one second ticker
			}
		}
	}()

	// task processing routine
	done := make(chan error, 1)
	go func() {

		log.Debugf("[%s] started", t.ResultWriter().TaskID())

		time.Sleep(30 * time.Second)
		total <- int64(100) // total steps are resolved

		i := 0
		for {

			if i == 100 { // all steps are processed
				break
			}

			// emulating step processing
			i++
			time.Sleep(2 * time.Second)
			processed <- int64(i)
		}

		log.Debugf("[%s] finished", t.ResultWriter().TaskID())

		// stop ticker
		ticker.Stop()
		done <- nil
	}()

	select {
	case <-ctx.Done():
		stopMonitor <- fmt.Errorf("task cancelled")
		return ctx.Err()
	case err := <-done:
		stopMonitor <- err
		return err
	}
}

func NewStager() *Stager {
	return &Stager{}
}

// StagerTaskResult
type StagerTaskResult struct {
	Progress struct {
		Total     int64 `json:"total"`
		Processed int64 `json:"processes"`
	} `json:"progress"`
	Error string `json:"error"`
}
