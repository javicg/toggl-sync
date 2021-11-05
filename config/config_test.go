package config

import (
	"github.com/spf13/viper"
	"reflect"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	viper.Set(TogglUsername, "togglUser")
	assertSame(t, Get(TogglUsername), "togglUser")
}

func TestGetSlice(t *testing.T) {
	input := "ENG,MGMT"
	viper.Set(JiraProjectKey, strings.Split(input, ","))
	assertSameSlice(t, GetSlice(JiraProjectKey), []string{"ENG", "MGMT"})
}

func TestSet(t *testing.T) {
	Set(TogglUsername, "togglUser2")
	assertSame(t, viper.GetString("toggl.username"), "togglUser2")
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

func TestFileUsed(t *testing.T) {
	viper.SetConfigFile("test-config.yml")
	assertSame(t, FileUsed(), "test-config.yml")
}

func assertSame(t *testing.T, first interface{}, second interface{}) {
	if first != second {
		t.Errorf("Expected [%s] to equal [%s] but it did not", first, second)
	}
}

func assertSameSlice(t *testing.T, first []string, second []string) {
	if !reflect.DeepEqual(first, second) {
		t.Errorf("Expected [%s] to equal [%s] but it did not", first, second)
	}
}
