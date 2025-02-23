package core

import "fmt"

func LogCommits() ([]*Commit, error) {

	var commits []*Commit
	hash, err := getCurrentCommit()
	if err != nil {
		return nil, fmt.Errorf("getting current commit: %v", err)
	}

	for hash != "" {
        commit, err := GetCommit(hash)
        if err != nil {
            return nil, err
        }
        commits = append(commits, commit)
        hash = commit.Parent
    }

	// commit, err := GetCommit("a9fb5a712cc3853a8b60ed509f0e23f02a4dee26")

	// fmt.Printf("%+v\n", *commit)

	return commits, err
}
