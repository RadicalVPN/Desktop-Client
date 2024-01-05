package version

import (
	"encoding/json"
	"net/http"
)

type TagCommitInfo struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

type Tag struct {
	Name   string        `json:"name"`
	Commit TagCommitInfo `json:"commit"`
}

func getAllTags() ([]Tag, error) {
	resp, err := http.Get("https://api.github.com/repos/RadicalVPN/desktop-client/tags")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	tags := []Tag{}

	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, err
	}

	return tags, nil
}
