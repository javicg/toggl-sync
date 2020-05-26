package config

import (
	"github.com/spf13/viper"
	"testing"
)

func TestGetTogglUsername(t *testing.T) {
	viper.Set("toggl.username", "togglUser")
	assertSame(t, GetTogglUsername(), "togglUser")
}

func TestSetTogglUsername(t *testing.T) {
	SetTogglUsername("togglUser2")
	assertSame(t, viper.GetString("toggl.username"), "togglUser2")
}

func TestGetTogglPassword(t *testing.T) {
	viper.Set("toggl.password", "togglPassword")
	assertSame(t, GetTogglPassword(), "togglPassword")
}

func TestSetTogglPassword(t *testing.T) {
	SetTogglPassword("togglPassword2")
	assertSame(t, viper.GetString("toggl.password"), "togglPassword2")
}

func TestGetTogglServerUrl(t *testing.T) {
	viper.Set("toggl.server.url", "http://localhost:1234")
	assertSame(t, GetTogglServerUrl(), "http://localhost:1234")
}

func TestSetTogglServerUrl(t *testing.T) {
	SetTogglServerUrl("http://localhost:12345")
	assertSame(t, viper.GetString("toggl.server.url"), "http://localhost:12345")
}

func TestGetJiraServerUrl(t *testing.T) {
	viper.Set("jira.server.url", "http://localhost:4321")
	assertSame(t, GetJiraServerUrl(), "http://localhost:4321")
}

func TestSetJiraServerUrl(t *testing.T) {
	SetJiraServerUrl("http://localhost:43210")
	assertSame(t, viper.GetString("jira.server.url"), "http://localhost:43210")
}

func TestGetJiraUsername(t *testing.T) {
	viper.Set("jira.username", "jiraUser")
	assertSame(t, GetJiraUsername(), "jiraUser")
}

func TestSetJiraUsername(t *testing.T) {
	SetJiraUsername("jiraUser2")
	assertSame(t, viper.GetString("jira.username"), "jiraUser2")
}

func TestGetJiraPassword(t *testing.T) {
	viper.Set("jira.password", "jiraPassword")
	assertSame(t, GetJiraPassword(), "jiraPassword")
}

func TestSetJiraPassword(t *testing.T) {
	SetJiraPassword("jiraPassword2")
	assertSame(t, viper.GetString("jira.password"), "jiraPassword2")
}

func TestGetJiraProjectKey(t *testing.T) {
	viper.Set("jira.project.key", "projectKey")
	assertSame(t, GetJiraProjectKey(), "projectKey")
}

func TestSetJiraProjectKey(t *testing.T) {
	SetJiraProjectKey("projectKey2")
	assertSame(t, viper.GetString("jira.project.key"), "projectKey2")
}

func TestGetOverheadKey(t *testing.T) {
	viper.Set("jira.overhead.meetings", "someValue")
	assertSame(t, GetOverheadKey("meetings"), "someValue")
}

func TestSetOverheadKey(t *testing.T) {
	SetOverheadKey("meetings", "someValue2")
	assertSame(t, viper.GetString("jira.overhead.meetings"), "someValue2")
}

func TestGetAllOverheadKeys(t *testing.T) {
	viper.Set("something", "value")
	viper.Set("jira.overhead.meetings", "overhead1")
	viper.Set("somethingElse", "otherValue")
	viper.Set("jira.overhead.cooking", "overhead2")

	overheadKeys := GetAllOverheadKeys()
	if len(overheadKeys) != 2 {
		t.Errorf("There should be only %d overhead keys, but there were %d", 2, len(overheadKeys))
	}
	for _, k := range overheadKeys {
		if k != "meetings" && k != "cooking" {
			t.Errorf("Expecting %s as overhead keys but got %s", [2]string{"meetings", "cooking"}, overheadKeys)
		}
	}
}

func TestReset(t *testing.T) {
	viper.Set("something", "value")
	Reset()
	if len(viper.AllKeys()) != 0 {
		t.Errorf("Reset should have cleared all keys, but it did not")
	}
}

func assertSame(t *testing.T, first string, second string) {
	if first != second {
		t.Errorf("Expected [%s] to equal [%s] but it did not", first, second)
	}
}
