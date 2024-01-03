package version

import (
	"encoding/json"
	"net/http"
	"strings"
)

type TagCommitInfo struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

type Tag struct {
	Name   string        `json:"name"`
	Commit TagCommitInfo `json:"commit"`
}

// -ldflags "-X radicalvpnd/version.version=x.x.x"
var version string
var commitHash string

func GetVersion() string {
	if len(version) == 0 {
		return "<unknown-version>"
	}

	return version
}

func GetCommitHash() string {
	if len(commitHash) == 0 {
		return "<unknown-commit-hash>"
	}

	return commitHash
}

func IsNightlyOutdated() bool {
	if !IsNightly() {
		return false
	}

	latestHash, err := GetLatestNightlyCommitHash()
	if err != nil {
		return false
	}

	return latestHash != GetCommitHash()
}

func IsNightly() bool {
	return strings.HasPrefix(version, "nightly")
}

func GetLatestNightlyCommitHash() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/RadicalVPN/desktop-client/tags")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	tags := []Tag{}

	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "", err
	}

	for _, tag := range tags {
		if tag.Name == "nightly" {
			return tag.Commit.Sha, nil
		}
	}

	//should never happen
	return "", nil
}
