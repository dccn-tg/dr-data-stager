package middleware

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/dccn-tg/dr-data-stager/internal/worker/config"
	"github.com/dccn-tg/dr-data-stager/pkg/tasks"
	"github.com/hibiken/asynq"

	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

func init() {
	// setup logger
	log.NewLogger(
		log.Configuration{
			EnableConsole:     true,
			ConsoleJSONFormat: false,
			ConsoleLevel:      log.Debug,
		},
		log.InstanceLogrusLogger,
	)
}

func TestLoadConfig(t *testing.T) {

	cfg, err := config.LoadConfig(os.Getenv("TEST_CONFIG_FILE"))

	if err != nil {
		t.Fatalf("%s\n", err)
	}

	// SMTP mailer
	client := stagerMailer{
		config: cfg.Mailer,
	}

	payload, _ := json.Marshal(tasks.StagerPayload{
		CreatedAt:         time.Now().Add(-3 * time.Hour).Unix(),
		Title:             "test email notification",
		DrUser:            "u1234567@ru.nl",
		DrPass:            "xxxx",
		DstURL:            "/project/3010000.01/a/b/c/台灣.txt",
		SrcURL:            "irods:/nl.ru.donders/di/dccn/DAC_3010000.01_173/台灣.txt",
		StagerUser:        "piepuk",
		StagerUserEmail:   "h.lee@donders.ru.nl",
		Timeout:           86400,
		TimeoutNoprogress: 3600,
	})

	rslt := new(tasks.StagerTaskResult)
	rslt.Progress.Total = 100
	rslt.Progress.Processed = 90
	rslt.Progress.Failed = 1

	drslt, _ := json.Marshal(rslt)

	tinfo := asynq.TaskInfo{
		ID:           "123@abc",
		Queue:        "default",
		Type:         tasks.TypeStager,
		Payload:      payload,
		State:        asynq.TaskStateCompleted,
		MaxRetry:     3,
		Retried:      3,
		LastErr:      "no data has been uploaded",
		LastFailedAt: time.Now().Add(-2 * time.Hour),
		CompletedAt:  time.Now().Add(-1 * time.Minute),
		Result:       drslt,
	}

	sendEmailNotification(&client, &tinfo, nFailed, cfg.Admins...)

}
