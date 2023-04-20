package main

import "github.com/shurcooL/githubv4"

// PR : A pullRequest
type PR struct {
	CreatedAt githubv4.DateTime `json:"createdAt"`
	URL       string            `json:"url"`
}
