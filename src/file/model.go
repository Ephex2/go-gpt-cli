package file

var AllowedFilePurposes = struct {
	Assistants       string
	AssistantsOutput string
	FineTune         string
	FineTuneResults  string
}{
	Assistants:       "assistants",
	AssistantsOutput: "AssistantsOutput",
	FineTune:         "fine-tune",
	FineTuneResults:  "fine-tune-results",
}

var CreateFileBody = map[string]string{
	"purpose": "DecidedAtRuntime",
}

type FileList struct {
	Object string `json:"object"`
	Data   []File `json:"data"`
}

type File struct {
	ID        string `json:"id"`
	Bytes     int    `json:"bytes"`
	CreatedAt int    `json:"created_at"`
	FileName  string `json:"filename"`
	Object    string `json:"object"`
	Purpose   string `json:"purpose"`
}

type DeleteStatus struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

func GetDefaultBody() map[string]string {
	return CreateFileBody
}
