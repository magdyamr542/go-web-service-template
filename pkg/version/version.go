package version

// Injected at build time
var (
	Version      string
	Description  = "Manage learning resources based on tags"
	BuildTime    string
	CommitTime   string
	CommitSHA    string
	CommitBranch string
)
