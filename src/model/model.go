package model

const ModelRoute string = "/v1/models"

/*
Describes an OpenAI model offering that can be used with the API.

id: string
The model identifier, which can be referenced in the API endpoints.

created: integer
The Unix timestamp (in seconds) when the model was created.

object: string
The object type, which is always "model".

owned_by: string
The organization that owns the model.
*/
type Model struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type ListModelResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

type DeleteModelResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}
