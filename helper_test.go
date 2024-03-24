package main

import (
	"fmt"
	"testing"
)

func TestGetOwnerAndRepo(t *testing.T) {
	testCases := []struct {
		input         string
		expectedOwner string
		expectedRepo  string
		expectedError error
	}{
		{
			input:         "owner/repo",
			expectedOwner: "owner",
			expectedRepo:  "repo",
			expectedError: nil,
		},
		{
			input:         "invalid-repo",
			expectedOwner: "",
			expectedRepo:  "",
			expectedError: fmt.Errorf("invalid-repo is not a valid repo name for this tool. It should be in the form of Owner/Reponame, like Jmainguy/ghreport"),
		},
		{
			input:         "owner/",
			expectedOwner: "",
			expectedRepo:  "",
			expectedError: fmt.Errorf("owner/ is not a valid repo name for this tool. It should be in the form of Owner/Reponame, like Jmainguy/ghreport"),
		},
		{
			input:         "",
			expectedOwner: "",
			expectedRepo:  "",
			expectedError: fmt.Errorf(" is not a valid repo name for this tool. It should be in the form of Owner/Reponame, like Jmainguy/ghreport"),
		},
	}

	for _, testCase := range testCases {
		owner, repo, err := getOwnerAndRepo(testCase.input)

		// Verify the expected owner
		if owner != testCase.expectedOwner {
			t.Errorf("Expected owner: %s, but got: %s", testCase.expectedOwner, owner)
		}

		// Verify the expected repo
		if repo != testCase.expectedRepo {
			t.Errorf("Expected repo: %s, but got: %s", testCase.expectedRepo, repo)
		}

		// Verify the expected error
		if fmt.Sprint(err) != fmt.Sprint(testCase.expectedError) {
			t.Errorf("Expected error: %v, but got: %v", testCase.expectedError, err)
		}
	}
}
