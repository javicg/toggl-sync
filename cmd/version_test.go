package cmd

import (
	"regexp"
	"testing"
)

func TestVersion(t *testing.T) {
	validVersion := regexp.MustCompile(`^toggl-sync v([0-9]+).([0-9]+).(.*)$`)
	if !validVersion.MatchString(getVersion()) {
		t.Error("Version should return a valid version number")
	}
}
