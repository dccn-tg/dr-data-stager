package middleware

import (
	"time"

	"github.com/Donders-Institute/dr-data-stager/pkg/tasks"
	"github.com/hibiken/asynq"
)

type DataNotification struct {
	ID           string
	State        asynq.TaskState
	StagerUser   string
	DrUser       string
	SrcURL       string
	DstURL       string
	CreatedAt    time.Time
	CompletedAt  time.Time
	LastFailedAt time.Time
	LastErr      string
	Result       tasks.StagerTaskResult
}

const templateNotificationCompleted string = `
<html>
<style>
	div { width: 100%; padding-top: 10px; padding-bottom: 10px;}
	table { width: 95%; border-collapse: collapse; }
	th { width: 20%; border: 1px solid #ddd; background-color: #f5f5f5; text-align: left; padding: 10px; }
	td { width: 80%; border: 1px solid #ddd; text-align: left; padding: 10px; }
</style>
<body>
  <b>Please be informed by the following completed stager job:</b>
  <div style="width: 100%; padding-top: 10px; padding-bottom: 10px;">
		<table style="width: 95%; border-collapse: collapse;">
			<tr>
				<th style="width: 20%; border: 1px solid #ddd; background-color: #f5f5f5; text-align: left; padding: 10px;">id</th>
				<td style="width: 80%; border: 1px solid #ddd; text-align: left; padding: 10px;">{{ .ID }}</td>
			</tr>
			<tr>
				<th>state</th>
				<td>{{ .State }}</td>
			</tr>
			<tr>
				<th>owner</th>
				<td>{{ .StagerUser }}</td>
			</tr>
			<tr>
				<th>repository user</th>
				<td>{{ .DrUser }}</td>
			</tr>
			<tr>
				<th>submitted at</th>
				<td>{{ .CreatedAt }}</td>
			</tr>
			<tr>
				<th>complete at</th>
				<td>{{ .CompletedAt }}</td>
			</tr>
			<tr>
				<th>source</th>
				<td>{{ .SrcURL }}</td>
			</tr>
			<tr>
				<th>destination</th>
				<td>{{ .DstURL }}</td>
			</tr>
		</table>
	</div>
</html>`

const templateNotificationFailed string = `
<html>
<style>
	div { width: 100%; padding-top: 10px; padding-bottom: 10px;}
	table { width: 95%; border-collapse: collapse; }
	th { width: 20%; border: 1px solid #ddd; background-color: #f5f5f5; text-align: left; padding: 10px; }
	td { width: 80%; border: 1px solid #ddd; text-align: left; padding: 10px; }
</style>
<body>
  <b>Please be informed by the following completed stager job:</b>
  <div style="width: 100%; padding-top: 10px; padding-bottom: 10px;">
		<table style="width: 95%; border-collapse: collapse;">
			<tr>
				<th style="width: 20%; border: 1px solid #ddd; background-color: #f5f5f5; text-align: left; padding: 10px;">id</th>
				<td style="width: 80%; border: 1px solid #ddd; text-align: left; padding: 10px;">{{ .ID }}</td>
			</tr>
			<tr>
				<th>state</th>
				<td>{{ .State }}</td>
			</tr>
			<tr>
				<th>owner</th>
				<td>{{ .StagerUser }}</td>
			</tr>
			<tr>
				<th>repository user</th>
				<td>{{ .DrUser }}</td>
			</tr>
			<tr>
				<th>submitted at</th>
				<td>{{ .CreatedAt }}</td>
			</tr>
			<tr>
				<th>last failed at</th>
				<td>{{ .LastFailedAt }}</td>
			</tr>
			<tr>
				<th>source</th>
				<td>{{ .SrcURL }}</td>
			</tr>
			<tr>
				<th>destination</th>
				<td>{{ .DstURL }}</td>
			</tr>
			<tr>
				<th>progress</th>
				<td>{{ .Result.Progress.Processed }} / {{ .Result.Progress.Total }}</td>
			</tr>
			<tr>
				<th>last error</th>
				<td>{{ .LastErr }}</td>
			</tr>
		</table>
	</div>
</html>`
