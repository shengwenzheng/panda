package version

import "fmt"

// These variables are set via linker flags at build time.
var (
	GitTag    = "v0.0.0-dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// String returns a formatted version string.
func String() string {
	return fmt.Sprintf("Version: %s, Commit: %s, Built at: %s", GitTag, GitCommit, BuildDate)
}

// Short returns the git tag only.
func Short() string {
	return GitTag
}
