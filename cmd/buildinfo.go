package cmd

// Field injected by ldflag. See goreleaser.yml
var (
	version    = "<version>"
	commitDate = "<commitDate>"
	commit     = "<commit>"
	target     = "<target>"
)

func Version() string {
	return version
}

func CommitDate() string {
	return commitDate
}

func Commit() string {
	return commit
}

func Target() string {
	return target
}
