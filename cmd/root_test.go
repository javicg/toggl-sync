package cmd

import (
	"errors"
	"github.com/javicg/toggl-sync/api"
	"github.com/javicg/toggl-sync/config"
	"testing"
	"time"
)

type MockTogglApi struct {
	Me               api.Me
	MeError          error
	TimeEntries      []api.TimeEntry
	TimeEntriesError error
	Project          api.Project
	ProjectError     error
}

func (mock MockTogglApi) GetMe() (*api.Me, error) {
	return &mock.Me, mock.MeError
}

func (mock MockTogglApi) GetTimeEntries(time.Time, time.Time) ([]api.TimeEntry, error) {
	return mock.TimeEntries, mock.TimeEntriesError
}

func (mock MockTogglApi) GetProjectById(int) (*api.Project, error) {
	return &mock.Project, mock.ProjectError
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

var expectNoInput = MockInputController{
	TextInputError: errors.New("no input should be requested"),
	PasswordError:  errors.New("no input should be requested"),
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

	config.Reset()
	config.SetJiraProjectKey("ENG")
	config.SetOverheadKey("testing", "ENG-1001")

	if err := sync(expectNoInput, togglApi, jiraApi, "2020-05-22", false); err != nil {
		t.Errorf("Sync failed with an error: %s", err)
	}
}

func TestSync_NoTimeEntries(t *testing.T) {
	togglApi := &MockTogglApi{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
		TimeEntries: []api.TimeEntry{},
	}
	jiraApi := &MockJiraApi{}

	config.Reset()
	config.SetJiraProjectKey("ENG")
	config.SetOverheadKey("testing", "ENG-1001")

	if err := sync(expectNoInput, togglApi, jiraApi, "2020-05-22", false); err != nil {
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

	config.Reset()
	config.SetJiraProjectKey("ENG")
	config.SetOverheadKey("testing", "ENG-1001")

	if err := sync(expectNoInput, togglApi, jiraApi, "2020-05-22", true); err != nil {
		t.Errorf("Sync failed with an error: %s", err)
	}
}

func TestSync_DryRun_NoTimeEntries(t *testing.T) {
	togglApi := &MockTogglApi{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
		TimeEntries: []api.TimeEntry{},
	}
	jiraApi := &MockJiraApi{
		ApiError: errors.New("no Jira API should be called during a dry run"),
	}

	config.Reset()
	config.SetJiraProjectKey("ENG")
	config.SetOverheadKey("testing", "ENG-1001")

	if err := sync(expectNoInput, togglApi, jiraApi, "2020-05-22", true); err != nil {
		t.Errorf("Sync failed with an error: %s", err)
	}
}

func TestSync_DryRun_ErrorFetchingUserDetails(t *testing.T) {
	togglApi := &MockTogglApi{
		MeError: errors.New("stub error"),
	}
	jiraApi := &MockJiraApi{}

	config.Reset()

	if err := sync(expectNoInput, togglApi, jiraApi, "2020-05-22", true); err == nil {
		t.Error("Sync should have failed with an error")
	}
}

func TestSync_DryRun_ErrorParsingSyncDate(t *testing.T) {
	togglApi := &MockTogglApi{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
	}
	jiraApi := &MockJiraApi{}

	config.Reset()

	if err := sync(expectNoInput, togglApi, jiraApi, "2nd January 2006", true); err == nil {
		t.Error("Sync should have failed with an error")
	}
}

func TestSync_DryRun_ErrorFetchingTimeEntries(t *testing.T) {
	togglApi := &MockTogglApi{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
		TimeEntriesError: errors.New("stub error"),
	}
	jiraApi := &MockJiraApi{}

	config.Reset()

	if err := sync(expectNoInput, togglApi, jiraApi, "2020-05-22", true); err == nil {
		t.Error("Sync should have failed with an error")
	}
}

func TestSync_DryRun_ValidationFailed_EmptyDescription(t *testing.T) {
	togglApi := &MockTogglApi{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
		TimeEntries: []api.TimeEntry{
			{
				Id:       1,
				Duration: 120,
			},
		},
	}
	jiraApi := &MockJiraApi{}

	config.Reset()
	config.SetJiraProjectKey("ENG")

	if err := sync(expectNoInput, togglApi, jiraApi, "2020-05-22", true); err == nil {
		t.Error("Sync should have failed with an error")
	}
}

func TestSync_DryRun_ValidationFailed_OverheadWorkWithoutProjectId(t *testing.T) {
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
				Duration:    300,
				Description: "Coffee break",
			},
		},
	}
	jiraApi := &MockJiraApi{}

	config.Reset()
	config.SetJiraProjectKey("ENG")

	if err := sync(expectNoInput, togglApi, jiraApi, "2020-05-22", true); err == nil {
		t.Error("Sync should have failed with an error")
	}
}

func TestSync_ErrorLoggingProjectWork_ShouldNotStopSync(t *testing.T) {
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
				Duration:    240,
				Description: "ENG-1001",
			},
		},
	}
	jiraApi := &MockJiraApi{
		ApiError: errors.New("stub error"),
	}

	config.Reset()
	config.SetJiraProjectKey("ENG")

	if err := sync(expectNoInput, togglApi, jiraApi, "2020-05-22", false); err != nil {
		t.Errorf("Sync failed with an error: %s", err)
	}
}

func TestSync_ErrorLoggingOverheadWork_EntryWithoutProjectId_ShouldNotStopSync(t *testing.T) {
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
		},
		ProjectError: errors.New("stub error"),
	}
	jiraApi := &MockJiraApi{}

	config.Reset()
	config.SetJiraProjectKey("ENG")

	if err := sync(expectNoInput, togglApi, jiraApi, "2020-05-22", false); err != nil {
		t.Errorf("Sync failed with an error: %s", err)
	}
}

func TestSync_ErrorLoggingOverheadWork_ShouldNotStopSync(t *testing.T) {
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
		},
		Project: api.Project{
			Data: api.ProjectData{
				Id:   1,
				Name: "testing",
			},
		},
	}
	jiraApi := &MockJiraApi{
		ApiError: errors.New("stub error"),
	}

	config.Reset()
	config.SetJiraProjectKey("ENG")
	config.SetOverheadKey("testing", "ENG-1001")

	if err := sync(expectNoInput, togglApi, jiraApi, "2020-05-22", false); err != nil {
		t.Errorf("Sync failed with an error: %s", err)
	}
}

func TestSync_LoggingOverheadWork_RequestOverheadKey(t *testing.T) {
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
		},
		Project: api.Project{
			Data: api.ProjectData{
				Id:   1,
				Name: "testing",
			},
		},
	}
	jiraApi := &MockJiraApi{}
	inputCtrl := MockInputController{
		TextInput: "ENG-1001",
	}

	config.Reset()
	config.SetJiraProjectKey("ENG")

	if err := sync(inputCtrl, togglApi, jiraApi, "2020-05-22", false); err != nil {
		t.Errorf("Sync failed with an error: %s", err)
	}
}

func TestSync_LoggingOverheadWork_ErrorRequestingOverheadKey_ShouldNotStopSync(t *testing.T) {
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
		},
		Project: api.Project{
			Data: api.ProjectData{
				Id:   1,
				Name: "testing",
			},
		},
	}
	jiraApi := &MockJiraApi{}
	inputCtrl := MockInputController{
		TextInputError: errors.New("stub error"),
	}

	config.Reset()
	config.SetJiraProjectKey("ENG")

	if err := sync(inputCtrl, togglApi, jiraApi, "2020-05-22", false); err != nil {
		t.Errorf("Sync failed with an error: %s", err)
	}
}
