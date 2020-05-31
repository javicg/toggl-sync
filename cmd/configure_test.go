package cmd

import (
	"errors"
	"github.com/javicg/toggl-sync/config"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, "value", config.GetTogglUsername())
	assert.Equal(t, "secret", config.GetTogglPassword())
	assert.Equal(t, "https://www.toggl.com/api/v8", config.GetTogglServerUrl())
	assert.Equal(t, "value", config.GetJiraUsername())
	assert.Equal(t, "secret", config.GetJiraPassword())
	assert.Equal(t, "value", config.GetJiraServerUrl())
	assert.Equal(t, "value", config.GetJiraProjectKey())
}

func TestUpdateConfiguration_TrimInputValues(t *testing.T) {
	err := updateConfiguration(&MockInputController{
		TextInput: "\n\t value \t\n",
		Password:  "\n\t secret \t\n",
	})

	if err != nil {
		t.Errorf("updateConfiguration failed: %s", err)
	}

	assert.Equal(t, "value", config.GetTogglUsername())
	assert.Equal(t, "secret", config.GetTogglPassword())
	assert.Equal(t, "https://www.toggl.com/api/v8", config.GetTogglServerUrl())
	assert.Equal(t, "value", config.GetJiraUsername())
	assert.Equal(t, "secret", config.GetJiraPassword())
	assert.Equal(t, "value", config.GetJiraServerUrl())
	assert.Equal(t, "value", config.GetJiraProjectKey())
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

	assert.Equal(t, "updatedValue", config.GetTogglUsername())
	assert.Equal(t, "updatedSecret", config.GetTogglPassword())
	assert.Equal(t, "https://www.toggl.com/api/v8", config.GetTogglServerUrl())
	assert.Equal(t, "updatedValue", config.GetJiraUsername())
	assert.Equal(t, "updatedSecret", config.GetJiraPassword())
	assert.Equal(t, "updatedValue", config.GetJiraServerUrl())
	assert.Equal(t, "updatedValue", config.GetJiraProjectKey())
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

	assert.Equal(t, "updatedValue", config.GetTogglUsername())
	assert.Equal(t, "secret", config.GetTogglPassword())
	assert.Equal(t, "https://www.toggl.com/api/v8", config.GetTogglServerUrl())
	assert.Equal(t, "updatedValue", config.GetJiraUsername())
	assert.Equal(t, "secret", config.GetJiraPassword())
	assert.Equal(t, "updatedValue", config.GetJiraServerUrl())
	assert.Equal(t, "updatedValue", config.GetJiraProjectKey())
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

	assert.Equal(t, "value", config.GetOverheadKey("meetings"))
	assert.Equal(t, "value", config.GetOverheadKey("cooking"))
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

	assert.Equal(t, "ENG-1234", config.GetOverheadKey("meetings"))
	assert.Equal(t, "ENG-1007", config.GetOverheadKey("cooking"))
}

func TestUpdateConfiguration_PropagateErrorWhenReadTextInputFails(t *testing.T) {
	err := updateConfiguration(&MockInputController{
		TextInputError: errors.New("stub error"),
		Password:       "secret",
	})
	assert.NotNilf(t, err, "Input errors should be propagated back to the client")
}

func TestUpdateConfiguration_PropagateErrorWhenReadPasswordFails(t *testing.T) {
	err := updateConfiguration(&MockInputController{
		TextInput:     "value",
		PasswordError: errors.New("stub error"),
	})
	assert.NotNilf(t, err, "Input errors should be propagated back to the client")
}
