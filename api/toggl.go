package api

import (
	"encoding/json"
	"github.com/javicg/toggl-sync/config"
	"log"
	"net/http"
	"time"
)

type TogglApi struct {
	baseUrl string
	client  *http.Client
}

func NewTogglApi() (api *TogglApi) {
	api = &TogglApi{}
	api.baseUrl = "https://www.toggl.com/api/v8"
	api.client = &http.Client{}
	return api
}

type Me struct {
	Data struct {
		Email    string
		Fullname string
	}
}

func (toggl *TogglApi) GetMe() (me Me, err error) {
	resp, err := toggl.getAuthenticated("/me")
	if err != nil {
		log.Fatalln("[GetMe] Request failed! Error:", err)
	} else if resp.StatusCode != 200 {
		log.Fatalf("[GetMe] Request failed with status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&me)
	if err != nil {
		log.Fatalln("[GetMe] Error unmarshalling response:", err)
	}

	return
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

func (toggl *TogglApi) GetTimeEntries(start time.Time, end time.Time) (entries []TimeEntry, err error) {
	params := map[string]string{
		"start_date": start.Format(time.RFC3339),
		"end_date":   end.Format(time.RFC3339),
	}

	resp, err := toggl.getAuthenticatedWithQueryParams("/time_entries", params)
	if err != nil {
		log.Fatalln("[GetTimeEntries] Request failed! Error:", err)
	} else if resp.StatusCode != 200 {
		log.Fatalf("[GetTimeEntries] Request failed with status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&entries)
	if err != nil {
		log.Fatalln("[GetTimeEntries] Error unmarshalling response:", err)
	}

	return
}

func (toggl *TogglApi) getAuthenticatedWithQueryParams(path string, params map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", toggl.baseUrl+path, nil)
	if err != nil {
		return
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

func (toggl *TogglApi) getAuthenticated(path string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", toggl.baseUrl+path, nil)
	if err != nil {
		return
	}

	req.SetBasicAuth(config.GetTogglUsername(), config.GetTogglPassword())

	req.Header.Add("Accept", "application/json")
	return toggl.client.Do(req)
}
