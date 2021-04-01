package cmd

import (
	"errors"
	"github.com/javicg/toggl-sync/api"
	"github.com/javicg/toggl-sync/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRootCmd_MissingDate(t *testing.T) {
	setupBasicConfig()

	cmd := NewRootCmd(&MockConfigManager{}, RejectAllInputController{t: t}, &MockTogglApi{}, &MockJiraApi{})
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_ProvidingDateAndSyncingCurrentDate(t *testing.T) {
	setupBasicConfig()

	cmd := NewRootCmd(&MockConfigManager{}, RejectAllInputController{t: t}, &MockTogglApi{}, &MockJiraApi{})
	cmd.SetArgs([]string{"2020-05-22", "--current-date"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_ErrorInitialisingConfig(t *testing.T) {
	configManager := &MockConfigManager{
		InitError: errors.New("stub error initialising config"),
	}
	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, &MockTogglApi{}, &MockJiraApi{})
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_InitConfigNotOk(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: false,
	}
	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, &MockTogglApi{}, &MockJiraApi{})
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_InvalidConfig(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}

	config.Reset()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, &MockTogglApi{}, &MockJiraApi{})
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
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

	setupBasicConfig()
	config.SetOverheadKey("testing", "ENG-1001")

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglApi, jiraApi)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_NoTimeEntries(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
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

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglApi, jiraApi)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_CurrentDate_NoTimeEntries(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
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

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglApi, jiraApi)
	cmd.SetArgs([]string{"--current-date"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_DryRun(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
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
	jiraApi := &RejectAllCallsJiraApi{t: t}

	setupBasicConfig()
	config.SetOverheadKey("testing", "ENG-1001")

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglApi, jiraApi)
	cmd.SetArgs([]string{"2020-05-22", "--dry-run"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_DryRun_NoTimeEntries(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglApi := &MockTogglApi{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
		TimeEntries: []api.TimeEntry{},
	}
	jiraApi := &RejectAllCallsJiraApi{t: t}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglApi, jiraApi)
	cmd.SetArgs([]string{"2020-05-22", "--dry-run"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_DryRun_ErrorFetchingUserDetails(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglApi := &MockTogglApi{
		MeError: errors.New("stub error"),
	}
	jiraApi := &MockJiraApi{}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglApi, jiraApi)
	cmd.SetArgs([]string{"2020-05-22", "--dry-run"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_DryRun_ErrorParsingSyncDate(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
	togglApi := &MockTogglApi{
		Me: api.Me{
			Data: api.PersonalInfo{
				Email:    "tester@toggl-sync.com",
				Fullname: "TogglSync Tester",
			},
		},
	}
	jiraApi := &MockJiraApi{}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglApi, jiraApi)
	cmd.SetArgs([]string{"2nd January 2006", "--dry-run"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_DryRun_ErrorFetchingTimeEntries(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
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

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglApi, jiraApi)
	cmd.SetArgs([]string{"2020-05-22", "--dry-run"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_DryRun_ValidationFailed_EmptyDescription(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
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

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglApi, jiraApi)
	cmd.SetArgs([]string{"2020-05-22", "--dry-run"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_DryRun_ValidationFailed_OverheadWorkWithoutProjectId(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
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

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglApi, jiraApi)
	cmd.SetArgs([]string{"2020-05-22", "--dry-run"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_ErrorLoggingProjectWork_ShouldNotStopSync(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
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

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglApi, jiraApi)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_ErrorLoggingOverheadWork_EntryWithoutProjectId_ShouldNotStopSync(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
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

	setupBasicConfig()

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglApi, jiraApi)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_ErrorLoggingOverheadWork_ShouldNotStopSync(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
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

	setupBasicConfig()
	config.SetOverheadKey("testing", "ENG-1001")

	cmd := NewRootCmd(configManager, RejectAllInputController{t: t}, togglApi, jiraApi)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_LoggingOverheadWork_RequestOverheadKey(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
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
	inputCtrl := &MockInputController{
		TextInput: "ENG-1001",
	}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, inputCtrl, togglApi, jiraApi)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.Nil(t, err)
}

func TestRootCmd_LoggingOverheadWork_RequestOverheadKey_ErrorPersistingConfig(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk:       true,
		PersistError: errors.New("stub error persisting configuration"),
	}
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
	inputCtrl := &MockInputController{
		TextInput: "ENG-1001",
	}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, inputCtrl, togglApi, jiraApi)
	cmd.SetArgs([]string{"2020-05-22"})
	err := cmd.Execute()
	assert.NotNil(t, err)
}

func TestRootCmd_LoggingOverheadWork_ErrorRequestingOverheadKey_ShouldNotStopSync(t *testing.T) {
	configManager := &MockConfigManager{
		InitOk: true,
	}
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
	inputCtrl := &MockInputController{
		TextInputError: errors.New("stub error"),
	}

	setupBasicConfig()

	cmd := NewRootCmd(configManager, inputCtrl, togglApi, jiraApi)
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

func (mock MockJiraApi) LogWorkWithUserDescription(string, time.Duration, string) error {
	return mock.ApiError
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

type RejectAllCallsJiraApi struct {
	t *testing.T
}

func (mock RejectAllCallsJiraApi) LogWork(string, time.Duration) (err error) {
	mock.t.Fatal("no API should be called")
	return
}

func (mock RejectAllCallsJiraApi) LogWorkWithUserDescription(string, time.Duration, string) (err error) {
	mock.t.Fatal("no API should be called")
	return
}

func setupBasicConfig() {
	config.Reset()
	viper.SetConfigFile("test-config.yml")
	config.SetTogglServerUrl("http://localhost/toggl")
	config.SetTogglUsername("TogglUser")
	config.SetTogglPassword("TogglPassword")
	config.SetJiraServerUrl("http://localhost/jira")
	config.SetJiraUsername("JiraUser")
	config.SetJiraPassword("JiraPassword")
	config.SetJiraProjectKey("ENG")
}
