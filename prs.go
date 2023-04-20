package main

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

func getPrFromRepo(client *githubv4.Client, org, repo string) ([]PR, error) {
	var repoQuery struct {
		Repository struct {
			PullRequests struct {
				PageInfo struct {
					HasNextPage githubv4.Boolean
					EndCursor   githubv4.String
				}
				Edges []struct {
					Node struct {
						URL       githubv4.URI
						CreatedAt githubv4.DateTime
						IsDraft   githubv4.Boolean
					}
				}
			} `graphql:"pullRequests(first: 100, states: $states, after: $cursor)"`
		} `graphql:"repository(name: $repo, owner: $org)"`
	}
	states := []githubv4.PullRequestState{
		githubv4.PullRequestStateOpen,
	}

	var PRS []PR

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

		// Loop thru each repo, and add it to []repos
		for _, edge := range repoQuery.Repository.PullRequests.Edges {
			if !bool(edge.Node.IsDraft) {
				var pr PR
				pr.URL = edge.Node.URL.String()
				pr.CreatedAt = edge.Node.CreatedAt
				PRS = append(PRS, pr)
			}
		}

		if !repoQuery.Repository.PullRequests.PageInfo.HasNextPage {
			break
		} else {
			variables["cursor"] = githubv4.NewString(repoQuery.Repository.PullRequests.PageInfo.EndCursor)
		}
		// Sleep for at least a second. https://docs.github.com/en/rest/guides/best-practices-for-integrators
		time.Sleep(time.Second * 2)

	}
	return PRS, nil
}
