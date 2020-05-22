package api

import (
	"github.com/javicg/toggl-sync/config"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestTogglApi_GetMe(t *testing.T) {
	expectedMe := Me{
		Data: PersonalInfo{
			Email:    "tester@toggl-sync.com",
			Fullname: "TogglSync Tester",
		},
	}

	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint:     "/me",
			ResponseCode: http.StatusOK,
			ResponseBody: AsJsonString(expectedMe),
		}).
		Create()
	defer server.Close()

	config.SetTogglServerUrl(server.URL)

	togglApi := NewTogglApi()
	me, err := togglApi.GetMe()
	if err != nil {
		t.Errorf("API call failed with an error: %s", err)
	} else if *me != expectedMe {
		t.Errorf("Unexpected response: was [%#v] instead of [%#v]", *me, expectedMe)
	}
}

func TestTogglApi_GetMe_ErrorWhenRequestFails(t *testing.T) {
	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint:     "/me",
			ResponseCode: http.StatusBadGateway,
		}).
		Create()
	defer server.Close()

	config.SetTogglServerUrl(server.URL)

	togglApi := NewTogglApi()
	_, err := togglApi.GetMe()
	if err == nil {
		t.Error("API errors should be returned to the client")
	}
}

func TestTogglApi_GetMe_ErrorWhenResponseHasUnexpectedFormat(t *testing.T) {
	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint:     "/me",
			ResponseCode: http.StatusOK,
			ResponseBody: AsJsonString("Bogus!"),
		}).
		Create()
	defer server.Close()

	config.SetTogglServerUrl(server.URL)

	togglApi := NewTogglApi()
	_, err := togglApi.GetMe()
	if err == nil {
		t.Error("JSON marshalling errors should be returned to the client")
	}
}

func TestTogglApi_GetTimeEntries(t *testing.T) {
	expectedTimeEntries := []TimeEntry{
		{
			Id:          1,
			Duration:    120,
			Description: "Writing toggl-sync tests",
		},
		{
			Id:          2,
			Duration:    240,
			Description: "Increasing toggl-sync test coverage",
		},
	}

	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint:     "/time_entries",
			ResponseCode: http.StatusOK,
			ResponseBody: AsJsonString(expectedTimeEntries),
		}).
		Create()
	defer server.Close()

	config.SetTogglServerUrl(server.URL)

	togglApi := NewTogglApi()
	startDate, _ := time.Parse("2006-01-02", "2020-05-08")
	endDate, _ := time.Parse("2006-01-02", "2020-05-09")
	entries, err := togglApi.GetTimeEntries(startDate, endDate)
	if err != nil {
		t.Errorf("API call failed with an error: %s", err)
	} else if !reflect.DeepEqual(entries, expectedTimeEntries) {
		t.Errorf("Unexpected response: was [%#v] instead of [%#v]", entries, expectedTimeEntries)
	}
}

func TestTogglApi_GetTimeEntries_ErrorWhenRequestFails(t *testing.T) {
	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint:     "/time_entries",
			ResponseCode: http.StatusBadGateway,
		}).
		Create()
	defer server.Close()

	config.SetTogglServerUrl(server.URL)

	togglApi := NewTogglApi()
	startDate, _ := time.Parse("2006-01-02", "2020-05-08")
	endDate, _ := time.Parse("2006-01-02", "2020-05-09")
	_, err := togglApi.GetTimeEntries(startDate, endDate)
	if err == nil {
		t.Error("API errors should be returned to the client")
	}
}

func TestTogglApi_GetTimeEntries_ErrorWhenResponseHasUnexpectedFormat(t *testing.T) {
	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint:     "/time_entries",
			ResponseCode: http.StatusOK,
			ResponseBody: AsJsonString("Bogus!"),
		}).
		Create()
	defer server.Close()

	config.SetTogglServerUrl(server.URL)

	togglApi := NewTogglApi()
	startDate, _ := time.Parse("2006-01-02", "2020-05-08")
	endDate, _ := time.Parse("2006-01-02", "2020-05-09")
	_, err := togglApi.GetTimeEntries(startDate, endDate)
	if err == nil {
		t.Error("JSON marshalling errors should be returned to the client")
	}
}

func TestTogglApi_GetProjectById(t *testing.T) {
	projectId := 10
	expectedProject := Project{
		Data: ProjectData{
			Id:   projectId,
			Name: "Top Secret",
		},
	}

	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint:     "/projects/" + strconv.Itoa(projectId),
			ResponseCode: http.StatusOK,
			ResponseBody: AsJsonString(expectedProject),
		}).
		Create()
	defer server.Close()

	config.SetTogglServerUrl(server.URL)

	togglApi := NewTogglApi()
	project, err := togglApi.GetProjectById(projectId)
	if err != nil {
		t.Errorf("API call failed with an error: %s", err)
	} else if *project != expectedProject {
		t.Errorf("Unexpected response: was [%#v] instead of [%#v]", *project, expectedProject)
	}
}

func TestTogglApi_GetProjectById_ErrorWhenRequestFails(t *testing.T) {
	projectId := 10

	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint:     "/projects/" + strconv.Itoa(projectId),
			ResponseCode: http.StatusBadGateway,
		}).
		Create()
	defer server.Close()

	config.SetTogglServerUrl(server.URL)

	togglApi := NewTogglApi()
	_, err := togglApi.GetProjectById(projectId)
	if err == nil {
		t.Error("API errors should be returned to the client")
	}
}

func TestTogglApi_GetProjectById_ErrorWhenResponseHasUnexpectedFormat(t *testing.T) {
	projectId := 10

	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint:     "/projects/" + strconv.Itoa(projectId),
			ResponseCode: http.StatusOK,
			ResponseBody: AsJsonString("Bogus!"),
		}).
		Create()
	defer server.Close()

	config.SetTogglServerUrl(server.URL)

	togglApi := NewTogglApi()
	_, err := togglApi.GetProjectById(projectId)
	if err == nil {
		t.Error("JSON marshalling errors should be returned to the client")
	}
}
