package cmd

import (
	"errors"
	"github.com/javicg/toggl-sync/api"
	"github.com/javicg/toggl-sync/config"
	"testing"
	"time"
)

type MockTogglApi struct {
	Me          api.Me
	TimeEntries []api.TimeEntry
	Project     api.Project
}

func (mock MockTogglApi) GetMe() (*api.Me, error) {
	return &mock.Me, nil
}

func (mock MockTogglApi) GetTimeEntries(time.Time, time.Time) ([]api.TimeEntry, error) {
	return mock.TimeEntries, nil
}

func (mock MockTogglApi) GetProjectById(id int) (*api.Project, error) {
	return &mock.Project, nil
}

type MockJiraApi struct {
	ApiError error
}

func (mock MockJiraApi) LogWork(string, time.Duration) error {
	return mock.ApiError
}

func (mock MockJiraApi) LogWorkWithUserDescription(string, string, time.Duration) error {
	return mock.ApiError
}

func TestSync(t *testing.T) {
	togglApi := &MockTogglApi{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
		TimeEntries: []api.TimeEntry{
			{
				Id:          1,
				Pid:         1,
				Duration:    120,
				Description: "Writing toggl-sync tests",
			},
			{
				Id:          2,
				Duration:    240,
				Description: "ENG-1002",
			},
		},
		Project: api.Project{
			Data: api.ProjectData{
				Id:   1,
				Name: "testing",
			},
		},
	}
	jiraApi := &MockJiraApi{}

	config.SetJiraProjectKey("ENG")
	config.SetOverheadKey("testing", "ENG-1001")

	if err := sync(togglApi, jiraApi, "2020-05-22", false); err != nil {
		t.Errorf("Sync failed with an error: %s", err)
	}
}

func TestSync_DryRun(t *testing.T) {
	togglApi := &MockTogglApi{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
		TimeEntries: []api.TimeEntry{
			{
				Id:          1,
				Pid:         1,
				Duration:    120,
				Description: "Writing toggl-sync tests",
			},
			{
				Id:          2,
				Duration:    240,
				Description: "ENG-1002",
			},
		},
		Project: api.Project{
			Data: api.ProjectData{
				Id:   1,
				Name: "testing",
			},
		},
	}
	jiraApi := &MockJiraApi{
		ApiError: errors.New("no Jira API should be called during a dry run"),
	}

	config.SetJiraProjectKey("ENG")
	config.SetOverheadKey("testing", "ENG-1001")

	if err := sync(togglApi, jiraApi, "2020-05-22", true); err != nil {
		t.Errorf("Sync failed with an error: %s", err)
	}
}
