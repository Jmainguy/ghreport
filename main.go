package main

import (
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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

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
		w.Flush()
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
	fmt.Fprintf(w, "%s\n\tauthor: %s\n\tAge: %s \n\treviewDecision: %s\n\tmergeable %s\n", pr.URL, pr.Owner, timeLabel, reviewDecisionEmoji, mergeableEmoji)
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
