package version

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
