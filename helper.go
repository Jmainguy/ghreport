package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

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

func formatTimeElapsed(duration time.Duration) string {
	if duration.Hours() >= 24 {
		return fmt.Sprintf("%.0f days", duration.Hours()/24)
	} else if duration.Hours() >= 1 {
		return fmt.Sprintf("%.0f hours", duration.Hours())
	} else {
		return fmt.Sprintf("%.0f minutes", duration.Minutes())
	}
}

func getReviewDecisionEmoji(decision string) string {
	switch decision {
	case "APPROVED":
		return "✅" // Green checkmark emoji
	case "CHANGES_REQUESTED":
		return "❌" // Red x emoji
	case "REVIEW_REQUIRED":
		return "🔍" // Magnifying glass emoji or any other appropriate emoji
	default:
		return "😅" // Emoji indicating everything is okay, but no review was requested
	}
}

func getMergeableEmoji(mergeable string) string {
	if mergeable == "MERGEABLE" {
		return "✅" // Green checkmark emoji
	} else if mergeable == "CONFLICTING" {
		return "❌" // Red x emoji
	}
	return ""
}
