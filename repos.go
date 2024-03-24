package main

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

func getReposFromOrganization(client Client, org string, topic string) ([]Repo, error) {
	var repoQuery struct {
		Organization struct {
			Repositories struct {
				PageInfo struct {
					HasNextPage githubv4.Boolean
					EndCursor   githubv4.String
				}
				Edges []RepositoryEdge
			} `graphql:"repositories(first: 100, after: $cursor, orderBy: {field: CREATED_AT, direction: DESC}, isArchived: false)"`
		} `graphql:"organization(login: $org)"`
	}

	var repos []Repo

	variables := map[string]interface{}{
		"org":    githubv4.String(org),
		"cursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	for {
		err := client.Query(context.Background(), &repoQuery, variables)
		if err != nil {
			return repos, err
		}

		repos = append(repos, extractRepoDataFromEdges(repoQuery.Organization.Repositories.Edges, topic)...)

		if !repoQuery.Organization.Repositories.PageInfo.HasNextPage {
			break
		} else {
			variables["cursor"] = githubv4.NewString(repoQuery.Organization.Repositories.PageInfo.EndCursor)
		}
		// Sleep for at least a second. https://docs.github.com/en/rest/guides/best-practices-for-integrators
		time.Sleep(2 * time.Second)
	}

	return repos, nil
}

func getReposFromUser(client Client, user string, topic string) ([]Repo, error) {
	var repoQuery struct {
		User struct {
			Repositories struct {
				PageInfo struct {
					HasNextPage githubv4.Boolean
					EndCursor   githubv4.String
				}
				Edges []RepositoryEdge
			} `graphql:"repositories(first: 100, after: $cursor, orderBy: {field: CREATED_AT, direction: DESC}, isArchived: false)"`
		} `graphql:"user(login: $user)"`
	}

	var repos []Repo

	variables := map[string]interface{}{
		"user":   githubv4.String(user),
		"cursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	for {
		err := client.Query(context.Background(), &repoQuery, variables)
		if err != nil {
			return repos, err
		}

		repos = append(repos, extractRepoDataFromEdges(repoQuery.User.Repositories.Edges, topic)...)

		if !repoQuery.User.Repositories.PageInfo.HasNextPage {
			break
		} else {
			variables["cursor"] = githubv4.NewString(repoQuery.User.Repositories.PageInfo.EndCursor)
		}
		// Sleep for at least a second. https://docs.github.com/en/rest/guides/best-practices-for-integrators
		time.Sleep(2 * time.Second)
	}

	return repos, nil
}

func extractRepoDataFromEdges(edges []RepositoryEdge, topic string) []Repo {
	var repos []Repo

	for _, edge := range edges {
		repo := edge.Node
		var r Repo
		r.NameWithOwner = repo.NameWithOwner
		if topic != "" {
			if containsTopic(repo.RepositoryTopics.Nodes, topic) {
				repos = append(repos, r)
			}
		} else {
			repos = append(repos, r)
		}
	}

	return repos
}

func containsTopic(nodes []RepositoryTopicNode, topic string) bool {
	for _, node := range nodes {
		if node.Topic.Name == topic {
			return true
		}
	}
	return false
}
