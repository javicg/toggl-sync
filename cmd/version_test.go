package cmd

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCmd(t *testing.T) {
	output := bytes.NewBufferString("")

	cmd := NewVersionCmd()
	cmd.SetOut(output)
	err := cmd.Execute()
	assert.Nil(t, err)

	version, err := ioutil.ReadAll(output)
	assert.Nil(t, err)
	if !isValidVersion(version) {
		t.Error("Version should return a valid version number")
	}
}

func isValidVersion(version []byte) bool {
	trimmedVersion := strings.TrimSpace(string(version))
	validVersionRegex := regexp.MustCompile(`^toggl-sync v([0-9]+).([0-9]+).(.*)$`)
	return validVersionRegex.MatchString(trimmedVersion)
}
