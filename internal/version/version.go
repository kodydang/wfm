package version

import "fmt"

// Build-time variables injected via -ldflags
var (
	Version   = "dev"
	BuildDate = "unknown"
	Commit    = "none"
)

// String returns a formatted version string.
func String() string {
	return fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, BuildDate)
}
