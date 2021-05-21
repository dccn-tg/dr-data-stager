package job

// Stager defines the data required for performing the data transfer task.
type Stager struct {

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

	// short description about the job
	Title string `json:"title"`
}

// Progress defines data structure of the progress information.
type Progress struct {
	Total     int64 `json:"total"`
	Processed int64 `json:"processes"`
}
