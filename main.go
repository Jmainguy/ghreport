package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {

	envSubscribedRepos := os.Getenv("subscribedRepos")
	if envSubscribedRepos == "" {
		fmt.Println("Env variable subscribedRepos is not defined")
		os.Exit(1)
	}

	subscribedRepos := strings.Split(envSubscribedRepos, " ")

	envToken := os.Getenv("ghreportToken")
	if envToken == "" {
		fmt.Println("Env variable ghreportToken is not defined")
		os.Exit(1)
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: envToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	// List of repos to watch
	for _, ownerRepo := range subscribedRepos {
		ownerAndRepo := strings.Split(ownerRepo, "/")
		if len(ownerAndRepo) == 2 {
			owner := ownerAndRepo[0]
			repo := ownerAndRepo[1]
			// List of pull requests for specific repo
			pullRequests, err := getPrFromRepo(client, owner, repo)
			if err != nil {
				panic(err)
			}
			for _, pr := range pullRequests {
				fmt.Printf("%s: createdAt %s\n", pr.URL, pr.CreatedAt)
			}
		} else {
			fmt.Printf("%s is not a valid repo name for this tool. It should be in the form of Owner/Reponame, like Jmainguy/ghreport\n", ownerRepo)
		}
	}

}
