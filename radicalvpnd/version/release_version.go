package version

import (
	"strings"

	"github.com/samber/lo"
	"golang.org/x/mod/semver"
)

func IsReleaseOutdated() bool {
	if !IsRelease() {
		return false
	}

	latestRelease, err := getLatestRelease()
	if err != nil {
		return false
	}

	return semver.Compare(version, latestRelease) == -1
}

func getLatestRelease() (string, error) {
	releases, err := getAllReleases()
	if err != nil {
		return "", err
	}

	semver.Sort(releases)

	return releases[len(releases)-1], nil
}

func getAllReleases() ([]string, error) {
	tags, err := getAllTags()
	if err != nil {
		return nil, err
	}

	releases := lo.FilterMap(tags, func(tag Tag, _ int) (string, bool) {
		tagName := tag.Name

		if strings.HasPrefix(tagName, "v") {
			return tagName, true
		}

		return "", false
	})

	return releases, nil
}

func IsRelease() bool {
	return strings.HasPrefix(version, "v")
}
