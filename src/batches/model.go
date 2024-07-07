package batches

import (
    "errors"
)

var allowedBatchApiEndpoints = []string{
    "/v1/completions",
    "/v1/chat/completions",
    "/v1/embeddings",
}

var allowedCompletionWindows = []string{
    "24h",
}

func AllowedBatchApiEndpoint(in string) error {
    err := stringAllowed(in, allowedBatchApiEndpoints)
    if err != nil {
        return err
    }

    return nil
}

func AllowedCompletionWindow(in string) error {
    err := stringAllowed(in, allowedCompletionWindows)
    if err != nil {
        return err
    }

    return nil
}

func stringAllowed(in string, allowedValues []string) (err error) {
    for _, endpointName := range allowedValues {
       if endpointName == in {
          return
       }
    }

    var allowedValuesString string
    for i, str := range allowedValues {
        if i == 0 {
            allowedValuesString += str
            continue
        }

        allowedValuesString += ",\n" + str
    }

    errorString := "value provided is not in the list of allowed api values.Value was: " + in + "\nAllowed values are:" + allowedValuesString
    err = errors.New(errorString)

    return
}

type CreateBatchBody struct {
    FileId string `json:"input_file_id"`
    Endpoint string `json:"endpoint"`
    CompletionWindow string `json:"completion_window"`
    Metadata map[string]string `json:"metadata,omitempty"`
}

func (cbb CreateBatchBody) Validate() error {
    var err error
    err = AllowedBatchApiEndpoint(cbb.Endpoint)
    err = AllowedCompletionWindow(cbb.CompletionWindow)

    return err
}

// Response objects //


type BatchList struct {
    Object string `json:"object"`
    Data []Batch `json:"data"`
    FirstId string `json:"first_id,omitempty"`
    LastId string `json:"last_id,omitempty"`
    HasMore bool `json:"has_more"`
}

type Batch struct {
	ID              string    `json:"id"`
	Object          string    `json:"object"`
	Endpoint        string    `json:"endpoint"`
	Errors          BatchError     `json:"errors,omitempty"`
	InputFileID     string    `json:"input_file_id"`
	CompletionWindow string   `json:"completion_window"`
	Status          string    `json:"status"`
	OutputFileID    string    `json:"output_file_id,omitempty"`
	ErrorFileID     string    `json:"error_file_id,omitempty"`
	CreatedAt       int64     `json:"created_at,omitempty"`
	InProgressAt    int64     `json:"in_progress_at,omitempty"`
	ExpiresAt       int64     `json:"expires_at,omitempty"`
	FinalizingAt    int64     `json:"finalizing_at,omitempty"`
	CompletedAt     int64     `json:"completed_at,omitempty"`
	FailedAt        int64     `json:"failed_at,omitempty"`
	ExpiredAt       int64     `json:"expired_at,omitempty"`
	CancellingAt    int64     `json:"cancelling_at,omitempty"`
	CancelledAt     int64     `json:"cancelled_at,omitempty"`
	RequestCounts   struct {
		Total     int `json:"total,omitempty"`
		Completed int `json:"completed,omitempty"`
		Failed    int `json:"failed,omitempty"`
	} `json:"request_counts"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type BatchError struct {
    Object string `json:"object"`
    Data []BatchErrorData `json:"data"`
}

type BatchErrorData struct {
    Code string `json:"code"`
    Message string `json:"message"`
    Param string `json:"param,omitempty"`
    Line int `json:"line,omitempty"`
}


// -------------------------- //


func GetDefaultBody() CreateBatchBody {
    var completionStr string
    if len(allowedCompletionWindows) > 0 {
        completionStr = allowedCompletionWindows[0]
    }

    defaultMetaData := make(map[string]string)
    defaultMetaData["started_by"] = "go-gpt-cli"

    // Do not validate ?
	return CreateBatchBody{
        FileId: "setAtRunTime",
        Endpoint: "setAtRunTime",
        CompletionWindow: completionStr,
        Metadata: defaultMetaData,
    }
}

func GetAllowedBatchApiEndpoints() []string {
    var out []string
    for _, str := range allowedBatchApiEndpoints {
        out = append(out, str)
    }

    return out
}

func GetAllowedCompletionWindows() []string {
    var out []string
    for _, str := range allowedCompletionWindows {
        out = append(out, str)
    }

    return out
}

