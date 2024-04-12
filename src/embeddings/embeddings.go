package embeddings

import (
	"encoding/json"

	"github.com/ephex2/go-gpt-cli/api"
)

const BaseEmbeddingsRoute string = "/v1/embeddings"

// Should the whole response be returned or just the embeddings themselves?
func CreateEmbeddings(input []string) (ceResp CreateEmbeddingResponse, err error) {
	embeddingsP, err := getDefaultProfile()
	if err != nil {
		return
	}

	embeddingsP.CreateEmbeddingBody.Input = input

	bodyBuf, err := json.Marshal(embeddingsP.CreateEmbeddingBody)
	if err != nil {
		return
	}

	buf, err := api.GenericRequest(nil, bodyBuf, BaseEmbeddingsRoute, "POST", embeddingsP.OverrideUrl())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &ceResp)
	if err != nil {
		return
	}

	return
}
