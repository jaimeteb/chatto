package version

import (
	"fmt"
)

var (
	version = "v0.0.0"
	commit  = "0000000000000000000000000000000000000000"
	date    = "0001-01-01 00:00:00 +0000 UTC"
	builtBy = "dev"
)

// BuildResponse is a JSON response with the
// current version of the bot or extension
type BuildResponse struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	BuiltAt string `json:"built_at"`
	BuiltBy string `json:"built_by"`
}

// BuildStr returns the build information provided
// by ldflags at build time as a formatted string
func BuildStr() string {
	var result = fmt.Sprintf(
		"version: %s\n"+
			"commit: %s\n"+
			"built at: %s\n"+
			"built by: %s\n",
		version, commit, date, builtBy)

	return result
}

// Build returns the build information provided
// by ldflags at build time as a response struct
func Build() BuildResponse {
	return BuildResponse{
		Version: version,
		Commit:  commit,
		BuiltAt: date,
		BuiltBy: builtBy,
	}
}
