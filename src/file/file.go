package file

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/ephex2/go-gpt-cli/api"
)

const BaseFileRoute string = "/v1/files"

// Should the whole response be returned or just the embeddings themselves?
func CreateFile(purpose string, filePath string) (resp File, err error) {
	p, err := getDefaultProfile()
	if err != nil {
		return
	}

	// Check is not dynamic at the moment
	if purpose != AllowedFilePurposes.Assistants && purpose != AllowedFilePurposes.FineTune {
		err = errors.New("purpose for file creation is not currently supported. purpose provided is: " + purpose + "Allowed purposes are: fine-tune, assistants")
		return
	}

	// Fine-tuning only supports .jsonl files in the OpenAI API at the moment.
	// i=f this check blocks alternative vendors / local code implementing similar APIs, consider unblocking.
	if purpose == AllowedFilePurposes.FineTune {
		ext := filepath.Ext(filePath)
		if ext != ".jsonl" {
			err = errors.New("only .jsonl files are supported. File provided has extension: " + ext)
			return
		}
	}

	p.CreateFileBody["purpose"] = purpose

	f, err := os.Open(filePath)
	if err != nil {
		return
	}

	fileDetails := []api.FileUploadDetails{
		{
			File:                f,
			UploadFormFieldName: "file",
		},
	}

	buf, err := api.MultiPartFormRequest(fileDetails, p.CreateFileBody, BaseFileRoute, "POST", p.OverrideUrl())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return
	}

	return
}

func DeleteFile(fileId string) (status DeleteStatus, err error) {
    p, err := getDefaultProfile()
    if err != nil {
        return
    }

	route := BaseFileRoute + "/" + fileId
	buf, err := api.GenericRequest(nil, nil, route, "DELETE", p.OverrideUrl())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &status)
	if err != nil {
		return
	}

	return
}

func GetFile(fileId string) (buf []byte, err error) {
    p, err := getDefaultProfile()
    if err != nil {
        return
    }

	route := BaseFileRoute + "/" + fileId + "/content"
	buf, err = api.GenericRequest(nil, nil, route, "GET", p.OverrideUrl())
	if err != nil {
		return
	}

	return
}

func StatFile(fileId string) (file File, err error) {
    p, err := getDefaultProfile()
    if err != nil {
        return
    }

	route := BaseFileRoute + "/" + fileId
	buf, err := api.GenericRequest(nil, nil, route, "GET", p.OverrideUrl())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &file)
	if err != nil {
		return
	}

	return
}

func ListFiles() (files FileList, err error) {
    p, err := getDefaultProfile()
    if err != nil {
        return
    }

	buf, err := api.GenericRequest(nil, nil, BaseFileRoute, "GET", p.OverrideUrl())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &files)
	if err != nil {
		return
	}

	return
}
