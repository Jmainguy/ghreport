package main

import (
	"context"

	"github.com/shurcooL/githubv4"
)

// PullRequest : A pullRequest
type PullRequest struct {
	CreatedAt      githubv4.DateTime `json:"createdAt"`
	URL            string            `json:"url"`
	Owner          githubv4.String   `json:"owner"`
	ReviewDecision githubv4.String   `json:"reviewDecision"`
	Mergeable      githubv4.String   `json:"mergeable"`
}

// Repo : Struct for repo providing NameWithOwner
type Repo struct {
	NameWithOwner githubv4.String `json:"nameWithOwner"`
}

// PullRequestEdge : A PullRequestEdge
type PullRequestEdge struct {
	Node struct {
		URL       githubv4.URI
		CreatedAt githubv4.DateTime
		IsDraft   githubv4.Boolean
		Author    struct {
			Login githubv4.String
		}
		ReviewDecision githubv4.String
		Mergeable      githubv4.String
	}
}

// RepositoryEdge : Graphql edge for repository
type RepositoryEdge struct {
	Node RepositoryNode
}

// RepositoryNode : Graphql node for repository
type RepositoryNode struct {
	NameWithOwner    githubv4.String
	RepositoryTopics struct {
		Nodes []RepositoryTopicNode
	} `graphql:"repositoryTopics(first: 100)"`
}

// RepositoryTopicNode : Struct for the repository topic
type RepositoryTopicNode struct {
	Topic struct {
		Name string
	}
}

// You can write a function that accepts an interface as an argument,
// and then pass either MockClient or githubv4.Client to it.

// Client : A arbitrary interface to support either MockClient or githubv4.Client
type Client interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}
