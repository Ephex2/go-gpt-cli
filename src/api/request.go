package api

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/ephex2/go-gpt-cli/config"
	"github.com/ephex2/go-gpt-cli/log"
)

var allowedMethods = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"PATCH":   true,
	"OPTIONS": true,
	"HEAD":    true,
}

func isValidHTTPMethod(method string) bool {
	_, ok := allowedMethods[strings.ToUpper(method)]
	return ok
}

func GenericRequest(queryParameters map[string]string, body []byte, route string, method string, overrideUrl string) (buf []byte, err error) {
	log.Debug("Body is : %s\n", string(body))

    var rawUrl string
    if overrideUrl != "" {
        rawUrl = overrideUrl + route
    }else {
    	rawUrl = config.BaseUrl() + route
    }

	method = strings.ToUpper(method)

	req, err := http.NewRequest(method, rawUrl, bytes.NewBuffer(body))
	if err != nil {
		err = errors.New("Error while initializing http request, error is: %s" + err.Error())
		return
	}

	err = defaultHeaders(req)
	if err != nil {
		return
	}

	for k, v := range queryParameters {
		req.Header.Add(k, v)
	}

	log.Debug("Request is : %s\n", req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		errBuf, _ := io.ReadAll(res.Body)

		if len(errBuf) == 0 {
			errBuf = []byte("EMPTY")
		}

		err = errors.New("Response from API does not indicate success: " + res.Status + "\nBody of response: " + string(errBuf))
		return
	}

	buf, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}

	return
}

func GenericPaginatedRequest(paginator Paginator, queryParameters map[string]string, body []byte, route string, method string, overrideUrl string) (err error) {
	log.Debug("Body is : %s\n", string(body))

    var rawUrl string
    if overrideUrl != "" {
        rawUrl = overrideUrl + route
    } else {
    	rawUrl = config.BaseUrl() + route
    }

	method = strings.ToUpper(method)

	req, err := http.NewRequest(method, rawUrl, bytes.NewBuffer(body))
	if err != nil {
		err = errors.New("Error while initializing http request, error is: %s" + err.Error())
		return
	}

	err = defaultHeaders(req)
	if err != nil {
		return
	}

	for k, v := range queryParameters {
		req.Header.Add(k, v)
	}

	log.Debug("Request is : %s\n", req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		errBuf, _ := io.ReadAll(res.Body)

		if len(errBuf) == 0 {
			errBuf = []byte("EMPTY")
		}

		err = errors.New("Response from API does not indicate success: " + res.Status + "\nBody of response: " + string(errBuf))
		return
	}

	more := true
	var nextReq *http.Request
	var nextRes *http.Response

	// Pagination loop. response behavior should be covered in paginator's Continue() function.
	for more == true {
		if nextReq != nil {
			nextRes, err = http.DefaultClient.Do(nextReq)
			if err != nil {
				return
			}

			more, nextReq, err = paginator.Continue(nextReq, nextRes)
		} else {
			more, nextReq, err = paginator.Continue(req, res)
			if err != nil {
				return
			}
		}
	}

	return
}

func MultiPartFormRequest(fileDetails []FileUploadDetails, fields map[string]string, route string, method string, overrideUrl string) (outputBuf []byte, err error) {
	if !isValidHTTPMethod(method) {
		err = errors.New("Method provided to api.MultiPartFormRequest() is not a valid http method: " + method)
		return
	}

	buf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(buf)

	// Add fields to MultiPartFormRequest
	for key, val := range fields {
		var part io.Writer
		part, err = bodyWriter.CreateFormField(key)
		if err != nil {
			return
		}

		_, err = part.Write([]byte(val))
		if err != nil {
			return
		}
	}

	// Add file(s) to MultiPartFormRequest
	for _, details := range fileDetails {
		f := details.File

		_, name := filepath.Split(f.Name())
		if name == "" {
			err = errors.New("error while parsing file name, we were unable to determine the base name. Source: " + f.Name())
			return
		}

		var fileWriter io.Writer
		fileWriter, err = bodyWriter.CreateFormFile(details.UploadFormFieldName, name)
		if err != nil {
			return
		}

		_, err = io.Copy(fileWriter, f)
		if err != nil {
			return
		}
	}

	err = bodyWriter.Close()
	if err != nil {
		return
	}

    var rawUrl string
    if overrideUrl != "" {
        rawUrl = overrideUrl + route
    } else {
    	rawUrl = config.BaseUrl() + route
    }

	// Setup request, using buf generated for multi part form fields
	req, err := http.NewRequest(method, rawUrl, buf)
	if err != nil {
		return
	}

	err = defaultHeaders(req)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		errBuf, _ := io.ReadAll(res.Body)

		if len(errBuf) == 0 {
			errBuf = []byte("EMPTY")
		}

		err = errors.New("Response from API does not indicate success: " + res.Status + "\nBody of response: " + string(errBuf))
		return
	}

	outputBuf, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}

	return
}

func defaultHeaders(req *http.Request) (err error) {
	apiKey, err := config.GetApiKey()
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	return
}
