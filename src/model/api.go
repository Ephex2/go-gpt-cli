package model

import (
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/api"
)

func ListModels() (models []Model, err error) {
	res, err := api.GenericRequest(nil, nil, ModelRoute, "GET", "")
	if err != nil {
		return
	}

	var resp ListModelResponse
	err = json.Unmarshal(res, &resp)
	if err != nil {
		return
	}

	models = resp.Data
	return
}

func RetrieveModel(name string) (model Model, err error) {
	buf, err := api.GenericRequest(nil, nil, ModelRoute+"/"+name, "GET", "")
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &model)
	if err != nil {
		return
	}

	return
}

func DeleteModel(name string) (dres DeleteModelResponse, err error) {
	route := ModelRoute + "/" + name
	buf, err := api.GenericRequest(nil, nil, route, "DELETE", "")
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &dres)
	if err != nil {
		return
	}

	return
}
