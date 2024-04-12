package finetuning

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/ephex2/go-gpt-cli/api"
)

const BaseFineTuningRoute string = "/v1/fine_tuning"

type listJobsPaginator struct {
	listJobsReturn []JobList
}

type listJobEventsPaginator struct {
	listJobEventsReturn []JobEventList
}

// Based on current Open AI API spec, must use query parameters "after" and "limit" to control pagination.
func (paginator *listJobsPaginator) Continue(req *http.Request, res *http.Response) (more bool, nextReq *http.Request, err error) {
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var jobList JobList
	err = json.Unmarshal(buf, &jobList)
	if err != nil {
		return
	}

	paginator.listJobsReturn = append(paginator.listJobsReturn, jobList)

	if jobList.HasMore {
		var after int
		var limit int

		after, err = strconv.Atoi(req.Header.Get("after"))
		if err != nil {
			return
		}

		limit, err = strconv.Atoi(req.Header.Get("limit"))
		if err != nil {
			return
		}

		newAfter := strconv.Itoa(after + limit)
		req.Header.Set("after", newAfter)

		http.DefaultClient.Do(req)
	}

	return
}

func (paginator *listJobEventsPaginator) Continue(req *http.Request, res *http.Response) (more bool, nextReq *http.Request, err error) {
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var eventList JobEventList
	err = json.Unmarshal(buf, &eventList)
	if err != nil {
		return
	}

	paginator.listJobEventsReturn = append(paginator.listJobEventsReturn, eventList)

	if eventList.HasMore {
		var after int
		var limit int

		after, err = strconv.Atoi(req.Header.Get("after"))
		if err != nil {
			return
		}

		limit, err = strconv.Atoi(req.Header.Get("limit"))
		if err != nil {
			return
		}

		newAfter := strconv.Itoa(after + limit)
		req.Header.Set("after", newAfter)

		http.DefaultClient.Do(req)
	}

	return
}

func CancelJob(id string) (resp Job, err error) {
    p, err := getDefaultProfile()
    if err != nil {
        return
    }

	route := BaseFineTuningRoute + "/jobs/" + id + "/cancel"
	buf, err := api.GenericRequest(nil, nil, route, "POST", p.OverrideUrl())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return
	}

	return
}

func GetJob(id string) (resp Job, err error) {
    p, err := getDefaultProfile()
    if err != nil {
        return
    }

	route := BaseFineTuningRoute + "/jobs/" + id
	buf, err := api.GenericRequest(nil, nil, route, "GET", p.OverrideUrl())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return
	}

	return
}

// List all fine-tune jobs on the API.
func ListJobs() (resp []Job, err error) {
    p, err := getDefaultProfile()
    if err != nil {
        return
    }

	paginator := listJobsPaginator{}
	route := BaseFineTuningRoute + "/jobs"

	err = api.GenericPaginatedRequest(&paginator, defaultPaginationQueryParameters, nil, route, "GET", p.OverrideUrl())
	if err != nil {
		return
	}

	for _, jobList := range paginator.listJobsReturn {
		resp = append(resp, jobList.Data...)
	}

	return
}

func ListEvents(id string) (resp []JobEvent, err error) {
    p, err := getDefaultProfile()
    if err != nil {
        return
    }

	paginator := listJobEventsPaginator{}
	route := BaseFineTuningRoute + "/jobs/" + id + "/events"

	err = api.GenericPaginatedRequest(&paginator, defaultPaginationQueryParameters, nil, route, "GET", p.OverrideUrl())
	if err != nil {
		return
	}

	for _, eventList := range paginator.listJobEventsReturn {
		resp = append(resp, eventList.Data...)
	}

	return
}

func CreateJob(id string) (resp Job, err error) {
	p, err := getDefaultProfile()
	if err != nil {
		return
	}

	if id == "" && p.CreateFineTuneBody.TrainingFile == "" {
		err = errors.New("no training_file id provided in default profile, and no id provided to the function")
		return
	} else if id != "" {
		// case where we override profile id with runtime id
		p.CreateFineTuneBody.TrainingFile = id
	}

	reqBuf, err := json.Marshal(p.CreateFineTuneBody)
	if err != nil {
		return
	}

	route := BaseFineTuningRoute + "/jobs"
	buf, err := api.GenericRequest(nil, reqBuf, route, "POST", p.OverrideUrl())
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return
	}

	return
}
