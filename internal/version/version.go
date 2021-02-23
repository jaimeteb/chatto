package version

import "fmt"

var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

// Build returns the build information provided by ldflags at build time.
func Build() string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}

	return result
}
