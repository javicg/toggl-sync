package cmd

import (
	"errors"
	"fmt"
	"github.com/javicg/toggl-sync/api"
	"github.com/javicg/toggl-sync/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRootCmd_MissingDate(t *testing.T) {
	setupBasicConfig()

	cmd := NewRootCmd(&MockConfigManager{}, RejectAllInputController{t: t}, &MockTogglAPI{}, &MockJiraAPI{})
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_ProvidingDateAndSyncingCurrentDate(t *testing.T) {
	setupBasicConfig()

	cmd := NewRootCmd(&MockConfigManager{}, RejectAllInputController{t: t}, &MockTogglAPI{}, &MockJiraAPI{})
	cmd.SetArgs([]string{"2020-05-22", "--current-date"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_ErrorInitialisingConfig(t *testing.T) {
	configManager := &MockConfigManager{
		InitError: errors.New("stub error initialising config"),
	}
	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, &MockTogglAPI{}, &MockJiraAPI{})
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_InitConfigNotOk(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: false,
	}
	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, &MockTogglAPI{}, &MockJiraAPI{})
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_InvalidConfig(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}

	config.Reset()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, &MockTogglAPI{}, &MockJiraAPI{})
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
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
				Pid:         10,
				Duration:    240,
				Description: "ENG-1002",
			},
			{
				Id:          3,
				Pid:         10,
				Duration:    140,
				Description: "ENG-1002",
			},
			{
				Id:          4,
				Duration:    360,
				Description: "ENG-1003",
			},
			{
				Id:          5,
				Pid:         10,
				Duration:    444,
				Description: "ENG-1003",
			},
		},
		Project: api.Project{
			Data: api.ProjectData{
				Id:   1,
				Name: "testing",
			},
		},
	}
	jiraAPI := &MockJiraAPI{}

	setupBasicConfig()
	config.SetOverheadKey("testing", "ENG-1001")

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.Nil(t, err)

	assert.NoError(t, jiraAPI.VerifyWorkLogged("Writing toggl-sync tests", 120))
	assert.NoError(t, jiraAPI.VerifyWorkLogged("ENG-1002", 380))
	assert.NoError(t, jiraAPI.VerifyWorkLogged("ENG-1003", 360))
	assert.NoError(t, jiraAPI.VerifyWorkLogged("ENG-1003", 444))
	assert.NoError(t, jiraAPI.VerifyNoOtherWorkLogged())
}

func TestRootCmd_NoTimeEntries(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
		TimeEntries: []api.TimeEntry{},
	}
	jiraAPI := &MockJiraAPI{}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_CurrentDate_NoTimeEntries(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
		TimeEntries: []api.TimeEntry{},
	}
	jiraAPI := &MockJiraAPI{}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"--current-date"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_DryRun(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
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
	jiraAPI := &RejectAllCallsJiraAPI{t: t}

	setupBasicConfig()
	config.SetOverheadKey("testing", "ENG-1001")

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22", "--dry-run"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_DryRun_NoTimeEntries(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
		TimeEntries: []api.TimeEntry{},
	}
	jiraAPI := &RejectAllCallsJiraAPI{t: t}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22", "--dry-run"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_DryRun_ErrorFetchingUserDetails(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
		MeError: errors.New("stub error"),
	}
	jiraAPI := &MockJiraAPI{}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22", "--dry-run"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_DryRun_ErrorParsingSyncDate(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
	}
	jiraAPI := &MockJiraAPI{}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2nd January 2006", "--dry-run"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_DryRun_ErrorFetchingTimeEntries(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
		TimeEntriesError: errors.New("stub error"),
	}
	jiraAPI := &MockJiraAPI{}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22", "--dry-run"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_DryRun_ValidationFailed_EmptyDescription(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
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
	jiraAPI := &MockJiraAPI{}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22", "--dry-run"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_DryRun_ValidationFailed_OverheadWorkWithoutProjectId(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
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
	jiraAPI := &MockJiraAPI{}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22", "--dry-run"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_DryRun_ValidationFailed_NegativeDuration(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
		TimeEntries: []api.TimeEntry{
			{
				Id:          1,
				Pid:         10,
				Duration:    -1630161422,
				Description: "New project (still working on it!)",
			},
		},
	}
	jiraAPI := &MockJiraAPI{}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22", "--dry-run"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_ErrorLoggingProjectWork_ShouldNotStopSync(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
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
	jiraAPI := &MockJiraAPI{
		APIError: errors.New("stub error"),
	}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_ErrorLoggingOverheadWork_EntryWithoutProjectId_ShouldNotStopSync(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
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
	jiraAPI := &MockJiraAPI{}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_ErrorLoggingOverheadWork_ShouldNotStopSync(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
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
	jiraAPI := &MockJiraAPI{
		APIError: errors.New("stub error"),
	}

	setupBasicConfig()
	config.SetOverheadKey("testing", "ENG-1001")

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_LoggingOverheadWork_RequestOverheadKey(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
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
	jiraAPI := &MockJiraAPI{}
	inputCtrl := &MockInputController{
		TextInput: "ENG-1001",
	}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, inputCtrl, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_LoggingOverheadWork_RequestOverheadKey_ErrorPersistingConfig(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk:       true,
		PersistError: errors.New("stub error persisting configuration"),
	}
	togglAPI := &MockTogglAPI{
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
	jiraAPI := &MockJiraAPI{}
	inputCtrl := &MockInputController{
		TextInput: "ENG-1001",
	}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, inputCtrl, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_LoggingOverheadWork_ErrorRequestingOverheadKey_ShouldNotStopSync(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglAPI := &MockTogglAPI{
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
	jiraAPI := &MockJiraAPI{}
	inputCtrl := &MockInputController{
		TextInputError: errors.New("stub error"),
	}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, inputCtrl, togglAPI, jiraAPI)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

type MockConfigManager struct {
	InitOk       bool
	InitError    error
	PersistError error
}

func (mock MockConfigManager) Init() (ok bool, err error) {
	if mock.InitError != nil {
		return false, mock.InitError
	}
	return mock.InitOk, nil
}

func (mock MockConfigManager) Persist() error {
	return mock.PersistError
}

type MockTogglAPI struct {
	Me               api.Me
	MeError          error
	TimeEntries      []api.TimeEntry
	TimeEntriesError error
	Project          api.Project
	ProjectError     error
}

func (mock MockTogglAPI) GetMe() (*api.Me, error) {
	return &mock.Me, mock.MeError
}

func (mock MockTogglAPI) GetTimeEntries(time.Time, time.Time) ([]api.TimeEntry, error) {
	return mock.TimeEntries, mock.TimeEntriesError
}

func (mock MockTogglAPI) GetProjectById(int) (*api.Project, error) {
	return &mock.Project, mock.ProjectError
}

type LoggedEntry struct {
	Description string
	Duration    time.Duration
}

type MockJiraAPI struct {
	LoggedWork []LoggedEntry
	APIError   error
}

func (mock *MockJiraAPI) LogWork(description string, duration time.Duration) error {
	mock.trackLog(description, duration)
	return mock.APIError
}

func (mock *MockJiraAPI) LogWorkWithUserDescription(_ string, duration time.Duration, description string) error {
	mock.trackLog(description, duration)
	return mock.APIError
}

func (mock *MockJiraAPI) trackLog(description string, duration time.Duration) {
	mock.LoggedWork = append(mock.LoggedWork, LoggedEntry{
		Description: description,
		Duration:    duration,
	})
}

func (mock *MockJiraAPI) VerifyWorkLogged(description string, duration int) error {
	expectedEntry := LoggedEntry{
		Description: description,
		Duration:    time.Duration(duration) * time.Second,
	}
	for i, entry := range mock.LoggedWork {
		if entry == expectedEntry {
			mock.removeFromLog(i)
			return nil
		}
	}
	return fmt.Errorf("work log does not contain [%s - %d]", description, duration)
}

func (mock *MockJiraAPI) removeFromLog(idx int) {
	mock.LoggedWork[idx] = mock.LoggedWork[len(mock.LoggedWork)-1]
	mock.LoggedWork = mock.LoggedWork[:len(mock.LoggedWork)-1]
}

func (mock *MockJiraAPI) VerifyNoOtherWorkLogged() error {
	if len(mock.LoggedWork) == 0 {
		return nil
	}
	return fmt.Errorf("there were unexpected entries logged: %s\n", mock.LoggedWork)
}

type RejectAllInputController struct {
	t *testing.T
}

func (ctrl RejectAllInputController) requestTextInput(string) (input string, err error) {
	ctrl.t.Fatal("no input should be requested")
	return
}

func (ctrl RejectAllInputController) requestPassword(string) (input string, err error) {
	ctrl.t.Fatal("no input should be requested")
	return
}

type RejectAllCallsJiraAPI struct {
	t *testing.T
}

func (mock RejectAllCallsJiraAPI) LogWork(string, time.Duration) (err error) {
	mock.t.Fatal("no API should be called")
	return
}

func (mock RejectAllCallsJiraAPI) LogWorkWithUserDescription(string, time.Duration, string) (err error) {
	mock.t.Fatal("no API should be called")
	return
}

func setupBasicConfig() {
	config.Reset()
	viper.SetConfigFile("test-config.yml")
	config.Set(config.TogglServerURL, "http://localhost/toggl")
	config.Set(config.TogglUsername, "TogglUser")
	config.Set(config.TogglPassword, "TogglPassword")
	config.Set(config.JiraServerURL, "http://localhost/jira")
	config.Set(config.JiraUsername, "JiraUser")
	config.Set(config.JiraPassword, "JiraPassword")
	config.Set(config.JiraProjectKey, "ENG")
}
