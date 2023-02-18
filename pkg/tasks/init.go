package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"

	ppath "github.com/Donders-Institute/dr-data-stager/pkg/path"
	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

// A list of task types.
const (
	TypeStager  = "stager"
	TypeEmailer = "email"
)

// Queues for different task types, with their associated task priority
var (
	StagerQueues = map[string]int{
		"critical": 6,
		"default":  3,
		"low":      1,
	}
	EmailerQueues = map[string]int{
		"default": 3,
	}
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

	updateRslt := func(rslt *StagerTaskResult) {
		if d, err := json.Marshal(rslt); err == nil {
			t.ResultWriter().Write(d)
		}
	}

	var p StagerPayload

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("cannot unmarshal stager payload: %v: %w", err, asynq.SkipRetry)
	}
	log.Debugf("[%s] payload: %+v", t.ResultWriter().TaskID(), p)

	// setup child context to handle `timeoutNoprogress`
	ctxnp, cancel := context.WithCancel(ctx)
	defer cancel()

	timer := time.NewTimer(time.Duration(p.TimeoutNoprogress) * time.Second)

	// logic of performing data transfer.
	srcPathInfo, err := ppath.GetPathInfo(ctxnp, p.SrcURL)
	if err != nil {
		log.Errorf("[%s] %s", t.ResultWriter().TaskID(), err)
		return fmt.Errorf("invalid source url: %s", err)
	}

	nf, err := ppath.GetNumberOfFiles(ctxnp, srcPathInfo)
	if err != nil {
		return err
	}

	rslt := new(StagerTaskResult)
	if nf == 0 {
		log.Warnf("[%s] %s is empty. Nothing to sync.", t.ResultWriter().TaskID(), p.SrcURL)
		updateRslt(rslt)
		return nil
	}

	// update total steps
	rslt.Progress.Total = int64(nf)
	updateRslt(rslt)

	dstPathInfo, _ := ppath.GetPathInfo(ctxnp, p.DstURL)

	success, failure := ppath.ScanAndSync(ctxnp, srcPathInfo, dstPathInfo, 4)

	var c int64 = 0
	var e ppath.SyncError
	ok1, ok2 := false, false
	for {
		ok1, ok2 = true, true
		// handle a successful transfer
		select {
		case _, ok1 = <-success:
			if ok1 {
				// stop timer
				if !timer.Stop() {
					<-timer.C
				}

				// increase the counter by 1, and update the queue data
				c++
				rslt.Progress.Processed = c
				updateRslt(rslt)

				// reset timer
				timer.Reset(time.Duration(p.TimeoutNoprogress) * time.Second)
			}
		default:
		}

		// handle a failed transfer
		select {
		case e, ok2 = <-failure:
			if ok2 {
				// append error of this attempt to the errors from previous attempts
				return fmt.Errorf("%s: %s", e.File, e.Error.Error())
			}
		default:
		}

		// both `success` and `failure` channels are closed, indicating the transfer is
		// completed.
		if !ok1 && !ok2 {
			log.Debugf("[%s] finished", t.ResultWriter().TaskID())
			return nil
		}

		// handle `timeoutNoprogress` and abort signal from parent context
		select {
		case <-timer.C:
			// reach the timer limit for the period of timeoutNoProgress
			log.Errorf("[%s] no progress timeout", t.ResultWriter().TaskID())
			return fmt.Errorf("no progress more than %d seconds", p.TimeoutNoprogress)
		case <-ctx.Done():
			// receive abort signal from parent context
			log.Debugf("[%s] aborted", t.ResultWriter().TaskID())
			return fmt.Errorf("received aborting signal: %w", asynq.SkipRetry)
		default:
		}
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
}
