package main

import (
	"fmt"
	"os"
	"testing"
)

func TestGetEnvVariablesSuccess(t *testing.T) {
	os.Setenv("subscribedRepos", "owner1/repo1 owner2/repo2")
	os.Setenv("ghreportToken", "testToken")
	defer os.Unsetenv("subscribedRepos")
	defer os.Unsetenv("ghreportToken")

	expectedSubscribedRepos := []string{"owner1/repo1", "owner2/repo2"}
	expectedToken := "testToken"

	subscribedRepos, token, err := getEnvVariables()

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if !compareSlices(subscribedRepos, expectedSubscribedRepos) {
		t.Errorf("Expected subscribedRepos %v, but got %v", expectedSubscribedRepos, subscribedRepos)
	}

	if token != expectedToken {
		t.Errorf("Expected token %s, but got %s", expectedToken, token)
	}
}

func TestGetEnvVariablesMissingSubscribedRepos(t *testing.T) {
	os.Setenv("ghreportToken", "testToken")
	defer os.Unsetenv("ghreportToken")

	expectedErrorMessage := "env variable subscribedRepos is not defined"

	subscribedRepos, token, err := getEnvVariables()

	if subscribedRepos != nil {
		t.Errorf("Expected subscribedRepos to be nil, but got %v", subscribedRepos)
	}

	if token != "" {
		t.Errorf("Expected token to be empty string, but got %s", token)
	}

	if err == nil {
		t.Error("Expected error, but got nil")
	} else if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', but got '%v'", expectedErrorMessage, err)
	}
}

func TestGetEnvVariablesMissingToken(t *testing.T) {
	os.Setenv("subscribedRepos", "owner1/repo1 owner2/repo2")
	defer os.Unsetenv("subscribedRepos")

	expectedErrorMessage := "env variable ghreportToken is not defined"

	subscribedRepos, token, err := getEnvVariables()

	if subscribedRepos != nil {
		t.Errorf("Expected subscribedRepos to be nil, but got %v", subscribedRepos)
	}

	if token != "" {
		t.Errorf("Expected token to be empty string, but got %s", token)
	}

	if err == nil {
		t.Error("Expected error, but got nil")
	} else if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', but got '%v'", expectedErrorMessage, err)
	}
}

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
