package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"os"
	"strings"
)

func main() {

	subscribedRepos := strings.Split(os.Getenv("subscribedRepos"), " ")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("ghreportToken")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// List of repos to watch
	for _, ownerRepo := range subscribedRepos {
		owner := strings.Split(ownerRepo, "/")[0]
		repo := strings.Split(ownerRepo, "/")[1]
		// List of pull requests for specific repo
		pullRequests, _, err := client.PullRequests.List(ctx, owner, repo, nil)
		if err != nil {
			panic(err)
		}
		for _, pr := range pullRequests {
			// We only care about open PR's
			if *pr.State == "open" {
				// If not a draft PR, so ready to be looked at
				if !*pr.Draft {
					fmt.Println(*pr.HTMLURL)
				}
			}
		}
	}

}
