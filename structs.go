package main

import (
	"context"

	"github.com/shurcooL/githubv4"
)

// PR : A pullRequest
type PR struct {
	CreatedAt githubv4.DateTime `json:"createdAt"`
	URL       string            `json:"url"`
	Owner     githubv4.String   `json:"owner"`
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
	}
}

// You can write a function that accepts an interface as an argument,
// and then pass either MockClient or githubv4.Client to it.

// Client : A arbitrary interface to support either MockClient or githubv4.Client
type Client interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}
