package appleopensource

import "go-darwin.dev/appleopensource/semver"

func TrimZeros(version semver.Version) semver.Version {
	return trimZeros(version)
}
