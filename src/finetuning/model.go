package finetuning

var defaultPaginationQueryParameters = map[string]string{
	"after": "0",
	"limit": "20",
}

type CreateFineTuneRequest struct {
	Model           string          `json:"model"`
	TrainingFile    string          `json:"training_file"`
	HyperParameters HyperParameters `json:"hyperparameters"`
	Suffix          *string         `json:"suffix,omitempty"`
	ValidationFile  *string         `json:"validation_file,omitempty"`
}

type HyperParameters struct {
	BatchSize              string `json:"batch_size,omitempty"`               // string or integer
	LearningRateMultiplier string `json:"learning_rate_multiplier,omitempty"` // string or number
	NEpochs                string `json:"n_epochs,omitempty"`                 // string or integer
}

type Job struct {
	ID             string    `json:"id"`
	CreatedAt      int       `json:"created_at"`
	Error          *JobError `json:"error,omitempty"`
	FineTunedModel *string   `json:"fine_tuned_model,omitempty"`
	FinishedAt     *int      `json:"finished_at,omitempty"`
	Model          string    `json:"model"`
	Object         string    `json:"object"`
	OrganizationID string    `json:"organization_id"`
	ResultFiles    []string  `json:"result_files"`
	Status         string    `json:"status"`
	TrainedTokens  *int      `json:"trained_tokens,omitempty"`
	TrainingFile   string    `json:"training_file"`
	ValidationFile *string   `json:"validation_file,omitempty"`
}

type JobList struct {
	Object  string `json:"object"`
	Data    []Job  `json:"data"`
	HasMore bool   `json:"has_more"`
}

type JobError struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Param   *string `json:"param"`
}

type JobEvent struct {
	ID        string `json:"id"`
	CreatedAt int    `json:"created_at"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Object    string `json:"object"`
}

type JobEventList struct {
	Object  string     `json:"object"`
	Data    []JobEvent `json:"data"`
	HasMore bool       `json:"has_more"`
}

func DefaultCreateFineTuneRequest() CreateFineTuneRequest {
	return CreateFineTuneRequest{
		Model:           "gpt-3.5-turbo",
		TrainingFile:    "", // empty training file IDs can be overwritten at runtime
		HyperParameters: defaultHyperParameters(),
		Suffix:          nil,
		ValidationFile:  nil,
	}
}

func defaultHyperParameters() HyperParameters {
	return HyperParameters{
		BatchSize:              "auto",
		LearningRateMultiplier: "auto",
		NEpochs:                "auto",
	}
}
