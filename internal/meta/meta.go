package meta

var (
	Version    string
	CommitHash string
	BuildTime  string
)

func GetVersion() string {
	if Version == "" {
		Version = "development"
	}
	return Version
}

func GetCommitHash() string {
	if CommitHash == "" {
		CommitHash = "unknown"
	}
	return CommitHash
}

func GetBuildTime() string {
	if BuildTime == "" {
		BuildTime = "unknown"
	}
	return BuildTime
}
