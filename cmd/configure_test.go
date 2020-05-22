package cmd

import (
	"errors"
	"github.com/javicg/toggl-sync/config"
	"testing"
)

type MockInputReader struct {
	TextInput      string
	TextInputError error
	Password       string
	PasswordError  error
}

func (mr MockInputReader) RequestTextInput(string) (string, error) {
	if mr.TextInputError != nil {
		return "", mr.TextInputError
	}
	return mr.TextInput, nil
}

func (mr MockInputReader) RequestPassword(string) (string, error) {
	if mr.PasswordError != nil {
		return "", mr.PasswordError
	}
	return mr.Password, nil
}

func TestUpdateConfiguration(t *testing.T) {
	err := updateConfiguration(&MockInputReader{
		TextInput: "value",
		Password:  "secret",
	})

	if err != nil {
		t.Errorf("updateConfiguration failed: %s", err)
	}

	assertConfigValue(t, "Toggl username", config.GetTogglUsername(), "value")
	assertConfigValue(t, "Toggl password", config.GetTogglPassword(), "secret")
	assertConfigValue(t, "Toggl server url", config.GetTogglServerUrl(), "https://www.toggl.com/api/v8")
	assertConfigValue(t, "Jira username", config.GetJiraUsername(), "value")
	assertConfigValue(t, "Jira password", config.GetJiraPassword(), "secret")
	assertConfigValue(t, "Jira server url", config.GetJiraServerUrl(), "value")
	assertConfigValue(t, "Jira server url", config.GetJiraProjectKey(), "value")
}

func TestUpdateConfiguration_TrimInputValues(t *testing.T) {
	err := updateConfiguration(&MockInputReader{
		TextInput: "\n\t value \t\n",
		Password:  "\n\t secret \t\n",
	})

	if err != nil {
		t.Errorf("updateConfiguration failed: %s", err)
	}

	assertConfigValue(t, "Toggl username", config.GetTogglUsername(), "value")
	assertConfigValue(t, "Toggl password", config.GetTogglPassword(), "secret")
	assertConfigValue(t, "Toggl server url", config.GetTogglServerUrl(), "https://www.toggl.com/api/v8")
	assertConfigValue(t, "Jira username", config.GetJiraUsername(), "value")
	assertConfigValue(t, "Jira password", config.GetJiraPassword(), "secret")
	assertConfigValue(t, "Jira server url", config.GetJiraServerUrl(), "value")
	assertConfigValue(t, "Jira server url", config.GetJiraProjectKey(), "value")
}

func TestUpdateConfiguration_OverrideExistingValues(t *testing.T) {
	err := updateConfiguration(&MockInputReader{
		TextInput: "value",
		Password:  "secret",
	})
	if err != nil {
		t.Errorf("updateConfiguration failed: %s", err)
	}

	err = updateConfiguration(&MockInputReader{
		TextInput: "updatedValue",
		Password:  "updatedSecret",
	})
	if err != nil {
		t.Errorf("updateConfiguration failed: %s", err)
	}

	assertConfigValue(t, "Toggl username", config.GetTogglUsername(), "updatedValue")
	assertConfigValue(t, "Toggl password", config.GetTogglPassword(), "updatedSecret")
	assertConfigValue(t, "Toggl server url", config.GetTogglServerUrl(), "https://www.toggl.com/api/v8")
	assertConfigValue(t, "Jira username", config.GetJiraUsername(), "updatedValue")
	assertConfigValue(t, "Jira password", config.GetJiraPassword(), "updatedSecret")
	assertConfigValue(t, "Jira server url", config.GetJiraServerUrl(), "updatedValue")
	assertConfigValue(t, "Jira server url", config.GetJiraProjectKey(), "updatedValue")
}

func TestUpdateConfiguration_PreserveExistingValuesOnEmptyInput(t *testing.T) {
	err := updateConfiguration(&MockInputReader{
		TextInput: "value",
		Password:  "secret",
	})
	if err != nil {
		t.Errorf("updateConfiguration failed: %s", err)
	}

	err = updateConfiguration(&MockInputReader{
		TextInput: "updatedValue",
		Password:  "",
	})
	if err != nil {
		t.Errorf("updateConfiguration failed: %s", err)
	}

	assertConfigValue(t, "Toggl username", config.GetTogglUsername(), "updatedValue")
	assertConfigValue(t, "Toggl password", config.GetTogglPassword(), "secret")
	assertConfigValue(t, "Toggl server url", config.GetTogglServerUrl(), "https://www.toggl.com/api/v8")
	assertConfigValue(t, "Jira username", config.GetJiraUsername(), "updatedValue")
	assertConfigValue(t, "Jira password", config.GetJiraPassword(), "secret")
	assertConfigValue(t, "Jira server url", config.GetJiraServerUrl(), "updatedValue")
	assertConfigValue(t, "Jira server url", config.GetJiraProjectKey(), "updatedValue")
}

func TestUpdateConfiguration_PropagateErrorWhenReadTextInputFails(t *testing.T) {
	err := updateConfiguration(&MockInputReader{
		TextInputError: errors.New("stub error"),
		Password:       "secret",
	})
	if err == nil {
		t.Error("Input errors should be propagated back to the client")
	}
}

func TestUpdateConfiguration_PropagateErrorWhenReadPasswordFails(t *testing.T) {
	err := updateConfiguration(&MockInputReader{
		TextInput:     "value",
		PasswordError: errors.New("stub error"),
	})
	if err == nil {
		t.Error("Input errors should be propagated back to the client")
	}
}

func assertConfigValue(t *testing.T, configName string, value string, expectedValue string) {
	if value != expectedValue {
		t.Errorf("Expecting %s value to equal [%s] but was [%s]", configName, expectedValue, value)
	}
}
