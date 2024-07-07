package batches

import (
	"encoding/json"
    "net/http"
    "io"

	"github.com/ephex2/go-gpt-cli/api"
)

const BaseBatchesRoute string = "/v1/batches"

func CreateBatch(fileid string, apiEndpoint string) (resp Batch, err error) {
	p, err := getDefaultProfile()
	if err != nil {
		return
	}

    p.CreateBatchBody.FileId = fileid
    p.CreateBatchBody.Endpoint = apiEndpoint

    err = p.CreateBatchBody.Validate()
    if err != nil {
        return
    }

    marshal, err := json.Marshal(p.CreateBatchBody)
    if err != nil {
        return
    }

    buf, err := api.GenericRequest(nil, marshal, BaseBatchesRoute, "POST", p.OverrideUrl())
    if err != nil {
        return
    }

	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return
	}

	return
}

func CancelBatch(batchId string) (b Batch, err error) {
    p, err := getDefaultProfile()
    if err != nil {
        return
    }

    route := BaseBatchesRoute + "/" + batchId + "/cancel"
    buf, err := api.GenericRequest(nil, nil, route, "POST", p.OverrideUrl())
	if err != nil {
		return
	}

    err = json.Unmarshal(buf, &b)
    if err != nil {
        return
    }

	return
}

func GetBatch(batchId string) (b Batch, err error) {
    p, err := getDefaultProfile()
    if err != nil {
        return
    }

	route := BaseBatchesRoute + "/" + batchId
    buf, err := api.GenericRequest(nil, nil, route, "GET", p.OverrideUrl())
	if err != nil {
		return
	}

    err = json.Unmarshal(buf, &b)
    if err != nil {
        return
    }

	return
}

func ListBatches() (batches BatchList, err error) {
    p, err := getDefaultProfile()
    if err != nil {
        return
    }

	buf, err := api.GenericRequest(nil, nil, BaseBatchesRoute, "GET", p.OverrideUrl())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &batches)
	if err != nil {
		return
	}

	return
}

type listBatchesPaginator struct {
	listBatchesReturn []BatchList
}

// Based on current Open AI API spec, must use query parameters "after" and "limit" to control pagination.
func (paginator *listBatchesPaginator) Continue(req *http.Request, res *http.Response) (more bool, nextReq *http.Request, err error) {
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var list BatchList
	err = json.Unmarshal(buf, &list)
	if err != nil {
		return
	}

    nextReq = req
	nextReq.Header.Set("after", list.LastId)
	paginator.listBatchesReturn = append(paginator.listBatchesReturn, list)

    more = list.HasMore

    return
}

