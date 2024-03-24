package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

func main() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	config, err := getConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := createGithubClient(config.Token)

	// If autoDiscover is configured and there are organizations specified, get repos from those organizations
	if len(config.AutoDiscover.Organizations) > 0 {
		for _, org := range config.AutoDiscover.Organizations {
			repos, err := getReposFromOrganization(client, org.Name, org.Topic)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Append discovered repos to subscribedRepos list
			var repoNames []string
			for _, repo := range repos {
				repoNames = append(repoNames, string(repo.NameWithOwner))
			}
			config.SubscribedRepos = append(config.SubscribedRepos, repoNames...)
		}
	}

	if len(config.AutoDiscover.Users) > 0 {
		for _, user := range config.AutoDiscover.Users {
			repos, err := getReposFromUser(client, user.Name, user.Topic)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Append discovered repos to subscribedRepos list
			var repoNames []string
			for _, repo := range repos {
				repoNames = append(repoNames, string(repo.NameWithOwner))
			}
			config.SubscribedRepos = append(config.SubscribedRepos, repoNames...)
		}
	}

	// Ensure list of repos is unique
	repoMap := make(map[string]bool)
	for _, repoString := range config.SubscribedRepos {
		repoMap[strings.ToLower(repoString)] = true
	}

	// List of repos to watch

	for ownerRepo := range repoMap {
		owner, repo, err := getOwnerAndRepo(ownerRepo)
		if err != nil {
			fmt.Println(err)
			continue
		}

		pullRequests, err := getPrFromRepo(client, owner, repo)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, pr := range pullRequests {
			timeElapsed := time.Since(pr.CreatedAt.Time)
			timeLabel := formatTimeElapsed(timeElapsed)

			reviewDecisionEmoji := getReviewDecisionEmoji(string(pr.ReviewDecision))

			mergeableEmoji := getMergeableEmoji(string(pr.Mergeable))

			fmt.Fprintf(w, "%s\n\tauthor: %s\n\tAge: %s \n\treviewDecision: %s\n\tmergeable %s\n", pr.URL, pr.Owner, timeLabel, reviewDecisionEmoji, mergeableEmoji)

			//fmt.Fprintf(w, "%s:\tauthor: %s\tcreatedAt %s\treviewDecison %s\tmergeable %s\n", pr.URL, pr.Owner, pr.CreatedAt, pr.ReviewDecision, pr.Mergeable)
		}
	}

	w.Flush()
}
