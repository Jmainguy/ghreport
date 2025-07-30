package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// Flags : Define a struct to hold command-line flags
type Flags struct {
	OutputFormat string
}

func main() {
	// Parse command-line flags
	flags := parseFlags()

	config, err := getConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Use config.DefaultOutput if flag not set
	if flags.OutputFormat == "" && config.DefaultOutput != "" {
		flags.OutputFormat = config.DefaultOutput
	}

	client := createGithubClient(config.Token)

	// Determine if we should use config or fallback to user's owned repos
	hasSubscribed := len(config.SubscribedRepos) > 0
	hasAutoDiscover := len(config.AutoDiscover.Organizations) > 0 || len(config.AutoDiscover.Users) > 0

	var repoList []string

	if hasSubscribed || hasAutoDiscover {
		// Use only subscribedRepos and autoDiscover
		repoMap := make(map[string]bool)
		// Add subscribedRepos
		for _, repoString := range config.SubscribedRepos {
			repoMap[strings.ToLower(repoString)] = true
		}
		// Add autoDiscover organizations
		for _, org := range config.AutoDiscover.Organizations {
			repos, err := getReposFromOrganization(client, org.Name, org.Topic)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, repo := range repos {
				repoMap[strings.ToLower(string(repo.NameWithOwner))] = true
			}
		}
		// Add autoDiscover users
		for _, user := range config.AutoDiscover.Users {
			repos, err := getReposFromUser(client, user.Name, user.Topic)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, repo := range repos {
				repoMap[strings.ToLower(string(repo.NameWithOwner))] = true
			}
		}
		// Build final repo list
		for repo := range repoMap {
			repoList = append(repoList, repo)
		}
	} else {
		// Fallback: get all repos owned by the authenticated user
		username, err := getAuthenticatedUsername(client)
		if err != nil {
			fmt.Println("Could not determine authenticated username:", err)
			os.Exit(1)
		}
		repos, err := getReposFromUser(client, username, "")
		if err != nil {
			fmt.Println("Could not get repos for authenticated user:", err)
			os.Exit(1)
		}
		for _, repo := range repos {
			repoList = append(repoList, strings.ToLower(string(repo.NameWithOwner)))
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// List of repos to watch

	for _, ownerRepo := range repoList {
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

			// Choose output format based on the flag
			switch flags.OutputFormat {
			case "json":
				outputJSON(pr, timeLabel, reviewDecisionEmoji, mergeableEmoji)
			case "singleline":
				outputSingleLine(pr, timeLabel, reviewDecisionEmoji, mergeableEmoji)
			default:
				outputDefault(w, pr, timeLabel, reviewDecisionEmoji, mergeableEmoji)
			}
		}
	}

	if flags.OutputFormat == "" {
		if err := w.Flush(); err != nil {
			fmt.Println("error flushing output:", err)
		}
	}
}

// Function to parse command-line flags
func parseFlags() Flags {
	var outputFormat string
	flag.StringVar(&outputFormat, "output", "", "Output format (default, singleline, json)")
	flag.Parse()
	return Flags{OutputFormat: outputFormat}
}

// Function to output data in default format
func outputDefault(w *tabwriter.Writer, pr PullRequest, timeLabel, reviewDecisionEmoji, mergeableEmoji string) {
	if _, err := fmt.Fprintf(w, "%s\n\tauthor: %s\n\tAge: %s \n\treviewDecision: %s\n\tmergeable %s\n", pr.URL, pr.Owner, timeLabel, reviewDecisionEmoji, mergeableEmoji); err != nil {
		fmt.Println("error writing output:", err)
	}
}

// Function to output data in single-line format
func outputSingleLine(pr PullRequest, timeLabel, reviewDecisionEmoji, mergeableEmoji string) {
	fmt.Printf("%s author: %s Age: %s reviewDecision: %s mergeable: %s\n", pr.URL, pr.Owner, timeLabel, reviewDecisionEmoji, mergeableEmoji)
}

// Function to output data in JSON format
func outputJSON(pr PullRequest, timeLabel, reviewDecisionEmoji, mergeableEmoji string) {
	data := struct {
		URL            string `json:"url"`
		Author         string `json:"author"`
		Age            string `json:"age"`
		ReviewDecision string `json:"review_decision"`
		Mergeable      string `json:"mergeable"`
	}{
		URL:            pr.URL,
		Author:         string(pr.Owner),
		Age:            timeLabel,
		ReviewDecision: reviewDecisionEmoji,
		Mergeable:      mergeableEmoji,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))
}

// Helper to get the authenticated user's login name
func getAuthenticatedUsername(client Client) (string, error) {
	var query struct {
		Viewer struct {
			Login string
		}
	}
	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		return "", err
	}
	return query.Viewer.Login, nil
}
