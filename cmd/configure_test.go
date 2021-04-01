package cmd

import (
	"errors"
	"github.com/javicg/toggl-sync/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigureCmd_ErrorInitialisingConfig(t *testing.T) {
	configManager := &MockConfigManager{
		InitError: errors.New("stub error initialising config"),
	}
	inputCtrl := &MockInputController{}

	cmd := NewConfigureCmd(configManager, inputCtrl)
	err := cmd.Execute()

	assert.NotNil(t, err)
}

func TestConfigureCmd(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput: "value",
		Password:  "secret",
	})
	err := cmd.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "value", config.GetTogglUsername())
	assert.Equal(t, "secret", config.GetTogglPassword())
	assert.Equal(t, "https://www.toggl.com/api/v8", config.GetTogglServerUrl())
	assert.Equal(t, "value", config.GetJiraUsername())
	assert.Equal(t, "secret", config.GetJiraPassword())
	assert.Equal(t, "value", config.GetJiraServerUrl())
	assert.Equal(t, "value", config.GetJiraProjectKey())
}

func TestConfigureCmd_TrimInputValues(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput: "\n\t value \t\n",
		Password:  "\n\t secret \t\n",
	})
	err := cmd.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "value", config.GetTogglUsername())
	assert.Equal(t, "secret", config.GetTogglPassword())
	assert.Equal(t, "https://www.toggl.com/api/v8", config.GetTogglServerUrl())
	assert.Equal(t, "value", config.GetJiraUsername())
	assert.Equal(t, "secret", config.GetJiraPassword())
	assert.Equal(t, "value", config.GetJiraServerUrl())
	assert.Equal(t, "value", config.GetJiraProjectKey())
}

func TestConfigureCmd_OverrideExistingValues(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput: "value",
		Password:  "secret",
	})
	err := cmd.Execute()
	assert.Nil(t, err)

	cmd = NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput: "updatedValue",
		Password:  "updatedSecret",
	})
	err = cmd.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "updatedValue", config.GetTogglUsername())
	assert.Equal(t, "updatedSecret", config.GetTogglPassword())
	assert.Equal(t, "https://www.toggl.com/api/v8", config.GetTogglServerUrl())
	assert.Equal(t, "updatedValue", config.GetJiraUsername())
	assert.Equal(t, "updatedSecret", config.GetJiraPassword())
	assert.Equal(t, "updatedValue", config.GetJiraServerUrl())
	assert.Equal(t, "updatedValue", config.GetJiraProjectKey())
}

func TestConfigureCmd_PreserveExistingValuesOnEmptyInput(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput: "value",
		Password:  "secret",
	})
	err := cmd.Execute()
	assert.Nil(t, err)

	cmd = NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput: "updatedValue",
		Password:  "",
	})
	err = cmd.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "updatedValue", config.GetTogglUsername())
	assert.Equal(t, "secret", config.GetTogglPassword())
	assert.Equal(t, "https://www.toggl.com/api/v8", config.GetTogglServerUrl())
	assert.Equal(t, "updatedValue", config.GetJiraUsername())
	assert.Equal(t, "secret", config.GetJiraPassword())
	assert.Equal(t, "updatedValue", config.GetJiraServerUrl())
	assert.Equal(t, "updatedValue", config.GetJiraProjectKey())
}

func TestConfigureCmd_OverrideOverheadKeys(t *testing.T) {
	config.SetOverheadKey("meetings", "ENG-1234")
	config.SetOverheadKey("cooking", "ENG-1007")
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput: "value",
		Password:  "secret",
	})
	err := cmd.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "value", config.GetOverheadKey("meetings"))
	assert.Equal(t, "value", config.GetOverheadKey("cooking"))
}

func TestConfigureCmd_PreserveOverheadKeysOnEmptyInput(t *testing.T) {
	config.SetOverheadKey("meetings", "ENG-1234")
	config.SetOverheadKey("cooking", "ENG-1007")
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput: "",
		Password:  "secret",
	})
	err := cmd.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "ENG-1234", config.GetOverheadKey("meetings"))
	assert.Equal(t, "ENG-1007", config.GetOverheadKey("cooking"))
}

func TestConfigureCmd_PropagateErrorWhenReadTextInputFails(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInputError: errors.New("stub error"),
		Password:       "secret",
	})
	err := cmd.Execute()

	assert.NotNilf(t, err, "Input errors should be propagated back to the client")
}

func TestConfigureCmd_PropagateErrorWhenReadPasswordFails(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput:     "value",
		PasswordError: errors.New("stub error"),
	})
	err := cmd.Execute()

	assert.NotNilf(t, err, "Input errors should be propagated back to the client")
}

func TestConfigureCmd_ErrorPersistingConfig(t *testing.T) {
	configManager := &MockConfigManager{
		PersistError: errors.New("stub error persisting config"),
	}
	inputCtrl := &MockInputController{
		TextInput: "value",
		Password:  "secret",
	}
	cmd := NewConfigureCmd(configManager, inputCtrl)
	err := cmd.Execute()

	assert.NotNil(t, err)
}
