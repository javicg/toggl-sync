package api

import (
	"encoding/json"
	"fmt"
	"github.com/javicg/toggl-sync/config"
	"net/http"
	"strconv"
	"time"
)

type TogglApi interface {
	GetMe() (*Me, error)
	GetTimeEntries(startDate time.Time, endDate time.Time) ([]TimeEntry, error)
	GetProjectById(id int) (*Project, error)
}

type TogglApiHttpClient struct {
	baseUrl string
	client  *http.Client
}

func NewTogglApi() TogglApi {
	api := &TogglApiHttpClient{}
	api.baseUrl = config.GetTogglServerUrl()
	api.client = &http.Client{}
	return api
}

type Me struct {
	Data PersonalInfo
}

type PersonalInfo struct {
	Email    string
	Fullname string
}

func (toggl *TogglApiHttpClient) GetMe() (*Me, error) {
	resp, err := toggl.getAuthenticated("/me")
	if err != nil {
		return nil, fmt.Errorf("[GetMe] Request failed! Error: %s", err)
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[GetMe] Request failed with status: %d", resp.StatusCode)
	}

	var me Me
	err = json.NewDecoder(resp.Body).Decode(&me)
	if err != nil {
		return nil, fmt.Errorf("[GetMe] Error unmarshalling response: %s", err)
	}

	return &me, resp.Body.Close()
}

type TimeEntry struct {
	Id          int
	Pid         int
	Start       time.Time
	Stop        time.Time
	Duration    int
	Description string
	Tags        []string
}

func (toggl *TogglApiHttpClient) GetTimeEntries(start time.Time, end time.Time) ([]TimeEntry, error) {
	params := map[string]string{
		"start_date": start.Format(time.RFC3339),
		"end_date":   end.Format(time.RFC3339),
	}

	resp, err := toggl.getAuthenticatedWithQueryParams("/time_entries", params)
	if err != nil {
		return nil, fmt.Errorf("[GetTimeEntries] Request failed! Error: %s", err)
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[GetTimeEntries] Request failed with status: %d", resp.StatusCode)
	}

	var entries []TimeEntry
	err = json.NewDecoder(resp.Body).Decode(&entries)
	if err != nil {
		return nil, fmt.Errorf("[GetTimeEntries] Error unmarshalling response: %s", err)
	}

	return entries, resp.Body.Close()
}

type Project struct {
	Data ProjectData
}

type ProjectData struct {
	Id   int
	Name string
}

func (toggl *TogglApiHttpClient) GetProjectById(pid int) (*Project, error) {
	resp, err := toggl.getAuthenticated("/projects/" + strconv.Itoa(pid))
	if err != nil {
		return nil, fmt.Errorf("[GetProjectById] Request failed! Error: %s", err)
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[GetProjectById] Request failed with status: %d", resp.StatusCode)
	}

	var data Project
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("[GetProjectById] Error unmarshalling response: %s", err)
	}

	return &data, resp.Body.Close()
}

func (toggl *TogglApiHttpClient) getAuthenticatedWithQueryParams(path string, params map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", toggl.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(config.GetTogglUsername(), config.GetTogglPassword())

	q := req.URL.Query()
	for p := range params {
		q.Add(p, params[p])
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Accept", "application/json")
	return toggl.client.Do(req)
}

func (toggl *TogglApiHttpClient) getAuthenticated(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", toggl.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(config.GetTogglUsername(), config.GetTogglPassword())

	req.Header.Add("Accept", "application/json")
	return toggl.client.Do(req)
}
