package version

// -ldflags "-X radicalvpnd/version.version=x.x.x"
var version string

func GetVersion() string {
	if len(version) == 0 {
		return "<unknown-version>"
	}

	return version
}
