package main

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

func getPrFromRepo(client Client, org, repo string) ([]PullRequest, error) {
	var repoQuery struct {
		Repository struct {
			PullRequests struct {
				PageInfo struct {
					HasNextPage githubv4.Boolean
					EndCursor   githubv4.String
				}
				Edges []PullRequestEdge
			} `graphql:"pullRequests(first: 100, states: $states, after: $cursor)"`
		} `graphql:"repository(name: $repo, owner: $org)"`
	}
	states := []githubv4.PullRequestState{
		githubv4.PullRequestStateOpen,
	}

	var PRS []PullRequest

	variables := map[string]interface{}{
		"org":    githubv4.String(org),
		"cursor": (*githubv4.String)(nil), // Null after argument to get first page.
		"states": states,
		"repo":   githubv4.String(repo),
	}

	for {
		err := client.Query(context.Background(), &repoQuery, variables)
		if err != nil {
			return PRS, err
		}

		// The three dots ... is called the ellipsis or "unpacking" operator.
		// It allows a slice to be expanded in place and passed as individual arguments to a variadic function like append().
		// The verbose equivalent would look like
		// prsToAdd := extractPRDataFromEdges(repoQuery.Repository.PullRequests.Edges)
		// PRS = append(PRS, prsToAdd[0], prsToAdd[1], ..., prsToAdd[N])

		PRS = append(PRS, extractPRDataFromEdges(repoQuery.Repository.PullRequests.Edges)...)

		if !repoQuery.Repository.PullRequests.PageInfo.HasNextPage {
			break
		} else {
			variables["cursor"] = githubv4.NewString(repoQuery.Repository.PullRequests.PageInfo.EndCursor)
		}
		// Sleep for at least a second. https://docs.github.com/en/rest/guides/best-practices-for-integrators
		time.Sleep(2 * time.Second)
	}
	return PRS, nil
}

func extractPRDataFromEdges(edges []PullRequestEdge) []PullRequest {
	var PRS []PullRequest

	for _, edge := range edges {
		if !bool(edge.Node.IsDraft) {
			var pr PullRequest
			pr.URL = edge.Node.URL.String()
			pr.CreatedAt = edge.Node.CreatedAt
			pr.Owner = edge.Node.Author.Login
			pr.Mergeable = edge.Node.Mergeable
			pr.ReviewDecision = edge.Node.ReviewDecision
			PRS = append(PRS, pr)
		}
	}

	return PRS
}
