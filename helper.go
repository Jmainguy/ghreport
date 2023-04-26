package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func getEnvVariables() ([]string, string, error) {
	envSubscribedRepos := os.Getenv("subscribedRepos")
	if envSubscribedRepos == "" {
		return nil, "", fmt.Errorf("env variable subscribedRepos is not defined")
	}

	subscribedRepos := strings.Split(envSubscribedRepos, " ")

	envToken := os.Getenv("ghreportToken")
	if envToken == "" {
		return nil, "", fmt.Errorf("env variable ghreportToken is not defined")
	}

	return subscribedRepos, envToken, nil
}

func createGithubClient(envToken string) *githubv4.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: envToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	return githubv4.NewClient(httpClient)
}

func getOwnerAndRepo(ownerRepo string) (string, string, error) {
	ownerAndRepo := strings.Split(ownerRepo, "/")
	if len(ownerAndRepo) != 2 {
		return "", "", fmt.Errorf("%s is not a valid repo name for this tool. It should be in the form of Owner/Reponame, like Jmainguy/ghreport", ownerRepo)
	}

	for _, key := range ownerAndRepo {
		if key == "" {
			return "", "", fmt.Errorf("%s is not a valid repo name for this tool. It should be in the form of Owner/Reponame, like Jmainguy/ghreport", ownerRepo)
		}
	}

	return ownerAndRepo[0], ownerAndRepo[1], nil
}

func compareSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
