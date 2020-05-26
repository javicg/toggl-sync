package cmd

import (
	"errors"
	"github.com/javicg/toggl-sync/config"
	"testing"
)

func TestUpdateConfiguration(t *testing.T) {
	err := updateConfiguration(&MockInputController{
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
	err := updateConfiguration(&MockInputController{
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
	err := updateConfiguration(&MockInputController{
		TextInput: "value",
		Password:  "secret",
	})
	if err != nil {
		t.Errorf("updateConfiguration failed: %s", err)
	}

	err = updateConfiguration(&MockInputController{
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
	err := updateConfiguration(&MockInputController{
		TextInput: "value",
		Password:  "secret",
	})
	if err != nil {
		t.Errorf("updateConfiguration failed: %s", err)
	}

	err = updateConfiguration(&MockInputController{
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

func TestUpdateConfiguration_OverrideOverheadKeys(t *testing.T) {
	config.SetOverheadKey("meetings", "ENG-1234")
	config.SetOverheadKey("cooking", "ENG-1007")
	err := updateConfiguration(&MockInputController{
		TextInput: "value",
		Password:  "secret",
	})
	if err != nil {
		t.Errorf("updateConfiguration failed: %s", err)
	}

	assertConfigValue(t, "jira.overhead.meetings", config.GetOverheadKey("meetings"), "value")
	assertConfigValue(t, "jira.overhead.cooking", config.GetOverheadKey("cooking"), "value")
}

func TestUpdateConfiguration_PreserveOverheadKeysOnEmptyInput(t *testing.T) {
	config.SetOverheadKey("meetings", "ENG-1234")
	config.SetOverheadKey("cooking", "ENG-1007")
	err := updateConfiguration(&MockInputController{
		TextInput: "",
		Password:  "secret",
	})
	if err != nil {
		t.Errorf("updateConfiguration failed: %s", err)
	}

	assertConfigValue(t, "jira.overhead.meetings", config.GetOverheadKey("meetings"), "ENG-1234")
	assertConfigValue(t, "jira.overhead.cooking", config.GetOverheadKey("cooking"), "ENG-1007")
}

func TestUpdateConfiguration_PropagateErrorWhenReadTextInputFails(t *testing.T) {
	err := updateConfiguration(&MockInputController{
		TextInputError: errors.New("stub error"),
		Password:       "secret",
	})
	if err == nil {
		t.Error("Input errors should be propagated back to the client")
	}
}

func TestUpdateConfiguration_PropagateErrorWhenReadPasswordFails(t *testing.T) {
	err := updateConfiguration(&MockInputController{
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
