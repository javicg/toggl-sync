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

func TestConfigureCmd_UseCorrectTogglAPIVersion(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput: "value",
		Password:  "secret",
	})
	err := cmd.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "https://api.track.toggl.com/api/v8", config.Get(config.TogglServerURL))
}

func TestConfigureCmd(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput: "value",
		Password:  "secret",
	})
	err := cmd.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "value", config.Get(config.TogglUsername))
	assert.Equal(t, "secret", config.Get(config.TogglPassword))
	assert.Equal(t, "value", config.Get(config.JiraUsername))
	assert.Equal(t, "secret", config.Get(config.JiraPassword))
	assert.Equal(t, "value", config.Get(config.JiraServerURL))
	assert.Equal(t, "value", config.GetSlice(config.JiraProjectKey)[0])
}

func TestConfigureCmd_TrimInputValues(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput: "\n\t value \t\n",
		Password:  "\n\t secret \t\n",
	})
	err := cmd.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "value", config.Get(config.TogglUsername))
	assert.Equal(t, "secret", config.Get(config.TogglPassword))
	assert.Equal(t, "value", config.Get(config.JiraUsername))
	assert.Equal(t, "secret", config.Get(config.JiraPassword))
	assert.Equal(t, "value", config.Get(config.JiraServerURL))
	assert.Equal(t, "value", config.GetSlice(config.JiraProjectKey)[0])
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
	assert.Equal(t, "updatedValue", config.Get(config.TogglUsername))
	assert.Equal(t, "updatedSecret", config.Get(config.TogglPassword))
	assert.Equal(t, "updatedValue", config.Get(config.JiraUsername))
	assert.Equal(t, "updatedSecret", config.Get(config.JiraPassword))
	assert.Equal(t, "updatedValue", config.Get(config.JiraServerURL))
	assert.Equal(t, "updatedValue", config.GetSlice(config.JiraProjectKey)[0])
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
	assert.Equal(t, "updatedValue", config.Get(config.TogglUsername))
	assert.Equal(t, "secret", config.Get(config.TogglPassword))
	assert.Equal(t, "updatedValue", config.Get(config.JiraUsername))
	assert.Equal(t, "secret", config.Get(config.JiraPassword))
	assert.Equal(t, "updatedValue", config.Get(config.JiraServerURL))
	assert.Equal(t, "updatedValue", config.GetSlice(config.JiraProjectKey)[0])
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

func TestConfigureCmd_PropagateErrorWhenReadingTogglUsernameFails(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInputError: errors.New("stub error"),
		Password:       "secret",
	})
	err := cmd.Execute()

	assert.NotNilf(t, err, "Input errors should be propagated back to the client")
}

func TestConfigureCmd_PropagateErrorWhenReadingJiraServerUrlFails(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput:          "value",
		FailTextInputAfter: 1,
		TextInputError:     errors.New("stub error"),
		Password:           "secret",
	})
	err := cmd.Execute()

	assert.NotNilf(t, err, "Input errors should be propagated back to the client")
}

func TestConfigureCmd_PropagateErrorWhenReadingJiraUsernameFails(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput:          "value",
		FailTextInputAfter: 2,
		TextInputError:     errors.New("stub error"),
		Password:           "secret",
	})
	err := cmd.Execute()

	assert.NotNilf(t, err, "Input errors should be propagated back to the client")
}

func TestConfigureCmd_PropagateErrorWhenReadingJiraProjectKeyFails(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput:          "value",
		FailTextInputAfter: 3,
		TextInputError:     errors.New("stub error"),
		Password:           "secret",
	})
	err := cmd.Execute()

	assert.NotNilf(t, err, "Input errors should be propagated back to the client")
}

func TestConfigureCmd_PropagateErrorWhenReadingOverheadKey(t *testing.T) {
	config.SetOverheadKey("meetings", "ENG-1234")
	config.SetOverheadKey("cooking", "ENG-1007")
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput:          "value",
		FailTextInputAfter: 4,
		TextInputError:     errors.New("stub error"),
		Password:           "secret",
	})
	err := cmd.Execute()

	assert.NotNilf(t, err, "Input errors should be propagated back to the client")
}

func TestConfigureCmd_PropagateErrorWhenReadingTogglPasswordFails(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput:     "value",
		PasswordError: errors.New("stub error"),
	})
	err := cmd.Execute()

	assert.NotNilf(t, err, "Input errors should be propagated back to the client")
}

func TestConfigureCmd_PropagateErrorWhenReadingJiraPasswordFails(t *testing.T) {
	cmd := NewConfigureCmd(&MockConfigManager{}, &MockInputController{
		TextInput:         "value",
		Password:          "secret",
		FailPasswordAfter: 1,
		PasswordError:     errors.New("stub error"),
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
