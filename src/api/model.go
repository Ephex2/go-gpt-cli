package api

import (
	"net/http"
	"os"
)

type FileUploadDetails struct {
	File                *os.File
	UploadFormFieldName string
}

// Paginators receive the request sent and the response it generated from the api.
// They must handle the response for their calling functions and then return:
// - more: whether another request will be required
// - nextReq: the next *http.Request to be executed
// - err: error or nil
type Paginator interface {
	Continue(req *http.Request, res *http.Response) (more bool, nextReq *http.Request, err error)
}
