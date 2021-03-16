package common

import "fmt"

var (
	version = "dev"
	commit  = "local"
	date    = "unknown"
)

func GetVersion() string {
	return fmt.Sprintf("%s-%s", version, commit)
}

func GetBuildDate() string {
	return date
}
