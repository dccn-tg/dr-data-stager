package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	hsrv "github.com/Donders-Institute/dr-data-stager/internal/api-server/handler"
	"github.com/Donders-Institute/dr-data-stager/internal/job"
	ppath "github.com/Donders-Institute/dr-data-stager/pkg/path"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
	"github.com/thoas/bokchoy"
)

// TaskResults defines the output structure of the task
type TaskResults struct {
	Error error  `json:"errors"`
	Info  string `json:"info"`
}

// StagerJobRunner implements `bokchoy.Handler` for applying update on project resource.
type StagerJobRunner struct {
	// Configuration file for the worker
	ConfigFile string

	// bokchoy Queue for updating task data during task runtime
	Queue *bokchoy.Queue
}

// Handle performs stager job based on the payload.
func (h *StagerJobRunner) Handle(r *bokchoy.Request) error {

	t := r.Task

	nattempt := hsrv.MaxRetries - t.MaxRetries + 1

	res, err := json.Marshal(t.Payload)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	var data job.Stager
	err = json.Unmarshal(res, &data)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}
	log.Debugf("payload data: %+v", data)

	// internal context that is closed upon timeout
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(data.Timeout)*time.Second)
	defer cancel()

	// logic of performing data transfer.
	srcPathInfo, err := ppath.GetPathInfo(ctx, data.SrcURL)
	if err != nil {
		return fmt.Errorf("invalid source url: %s", err)
	}

	nf, err := ppath.GetNumberOfFiles(ctx, srcPathInfo)
	if err != nil {
		return err
	}

	if nf == 0 {
		log.Warnf("%s is empty. Nothing to sync.", data.SrcURL)
		t.Result = job.Progress{
			Total:     0,
			Processed: 0,
		}
		return nil
	}

	// update result with total files to be processed
	t.Result = job.Progress{
		Total:     int64(nf),
		Processed: 0,
	}

	log.Infof("%+v", t.Result)
	h.Queue.Save(r.Context(), t)

	dstPathInfo, _ := ppath.GetPathInfo(ctx, data.DstURL)

	success, failure := ppath.ScanAndSync(ctx, srcPathInfo, dstPathInfo, 4)

	// timer for checking timeout noprogress
	timer := time.NewTimer(time.Duration(data.TimeoutNoprogress) * time.Second)

	// ticker to update queue data every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
		ticker.Stop()
	}()

	c := 0
	for {
		select {
		case _, ok := <-success:
			if !ok {
				success = nil
			} else {

				// new file being processed, stop the timer.
				if !timer.Stop() {
					<-timer.C
				}

				// increase the counter by 1, and update the queue data
				c++
				t.Result = job.Progress{
					Total:     int64(nf),
					Processed: int64(c),
				}

				// reset timer
				timer.Reset(time.Duration(data.TimeoutNoprogress) * time.Second)
			}
		case e, ok := <-failure:
			if !ok {
				failure = nil
			} else {
				// append error of this attempt to the errors from previous attempts
				return fmt.Errorf("[attempt %d] %s, path: %s", nattempt, e.Error, e.File)
			}
		case <-ticker.C:
			log.Infof("save to queue")
			// update queue with progress data
			h.Queue.Save(r.Context(), t)
		case <-timer.C:
			log.Infof("no progress")
			// reach the progress check timer
			return fmt.Errorf("[attempt %d] timeout after %d seconds without progress", nattempt, data.TimeoutNoprogress)
		case <-ctx.Done():
			// fail job if context is closed due to timeout.
			return fmt.Errorf("[attempt %d] timeout after %d seconds", nattempt, data.Timeout)
		}

		if success == nil && failure == nil {
			log.Debugf("finish")
			break
		}
	}

	return nil
}
