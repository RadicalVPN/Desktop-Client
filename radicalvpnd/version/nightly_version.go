package version

import "strings"

func IsNightlyOutdated() bool {
	if !IsNightly() {
		return false
	}

	latestHash, err := getLatestNightlyCommitHash()
	if err != nil {
		return false
	}

	return latestHash != GetCommitHash()
}

func IsNightly() bool {
	return strings.HasPrefix(version, "nightly")
}

func getLatestNightlyCommitHash() (string, error) {
	tags, err := getAllTags()
	if err != nil {
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
