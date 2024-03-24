package main

import (
	"context"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClient that implements the Client interface
type MockClient struct {
	mock.Mock
}

func (m *MockClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	args := m.Called(ctx, q, variables)
	return args.Error(0)
}

func TestExtractPRDataFromEdges(t *testing.T) {
	uri1 := githubv4.URI{URL: &url.URL{Scheme: "https", Host: "github.com", Path: "/my-org/my-repo/pull/1"}}
	uri2 := githubv4.URI{URL: &url.URL{Scheme: "https", Host: "github.com", Path: "/my-org/my-repo/pull/2"}}

	// Define input data
	edges := []PullRequestEdge{
		{
			Node: struct {
				URL       githubv4.URI
				CreatedAt githubv4.DateTime
				IsDraft   githubv4.Boolean
				Author    struct {
					Login githubv4.String
				}
				ReviewDecision githubv4.String
				Mergeable      githubv4.String
			}{
				URL:       uri1,
				CreatedAt: githubv4.DateTime{Time: time.Now().UTC().Truncate(time.Hour)},
				IsDraft:   githubv4.Boolean(false),
				Author: struct {
					Login githubv4.String
				}{
					Login: githubv4.String("john_doe"),
				},
			},
		},
		{
			Node: struct {
				URL       githubv4.URI
				CreatedAt githubv4.DateTime
				IsDraft   githubv4.Boolean
				Author    struct {
					Login githubv4.String
				}
				ReviewDecision githubv4.String
				Mergeable      githubv4.String
			}{
				URL:       uri2,
				CreatedAt: githubv4.DateTime{Time: time.Now().UTC().Truncate(time.Hour)},
				IsDraft:   githubv4.Boolean(true),
				Author: struct {
					Login githubv4.String
				}{
					Login: githubv4.String("john_doe"),
				},
			},
		},
	}

	// Define expected output
	expected := []PR{
		{
			URL:       "https://github.com/my-org/my-repo/pull/1",
			CreatedAt: githubv4.DateTime{Time: time.Now().UTC().Truncate(time.Hour)},
			Owner:     "john_doe",
		},
	}

	// Call function and check output
	actual := extractPRDataFromEdges(edges)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("extractPRDataFromEdges() returned incorrect result.\nexpected: %v\nactual: %v", expected, actual)
	}
}

func TestGetPrFromRepo(t *testing.T) {
	// Mock data for testing
	org := "testOrg"
	repo := "testRepo"
	prs := []PR{
		{CreatedAt: githubv4.DateTime{}, URL: "https://github.com/testOrg/testRepo/pull/1", Owner: "john_doe"},
	}

	// Mock client
	mockClient := new(MockClient)

	uri := githubv4.URI{URL: &url.URL{Scheme: "https", Host: "github.com", Path: "/testOrg/testRepo/pull/1"}}

	// First call to Query returns no error and hasNextPage is true
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Modify the result of the query to have a next page
		repoQuery := args.Get(1).(*struct {
			Repository struct {
				PullRequests struct {
					PageInfo struct {
						HasNextPage githubv4.Boolean
						EndCursor   githubv4.String
					}
					Edges []PullRequestEdge
				} `graphql:"pullRequests(first: 100, states: $states, after: $cursor)"`
			} `graphql:"repository(name: $repo, owner: $org)"`
		})
		repoQuery.Repository.PullRequests.PageInfo.HasNextPage = false
		repoQuery.Repository.PullRequests.PageInfo.EndCursor = "endCursor"
		repoQuery.Repository.PullRequests.Edges = []PullRequestEdge{
			{
				Node: struct {
					URL       githubv4.URI
					CreatedAt githubv4.DateTime
					IsDraft   githubv4.Boolean
					Author    struct {
						Login githubv4.String
					}
					ReviewDecision githubv4.String
					Mergeable      githubv4.String
				}{
					URL:       uri,
					CreatedAt: githubv4.DateTime{},
					IsDraft:   githubv4.Boolean(false),
					Author: struct {
						Login githubv4.String
					}{
						Login: githubv4.String("john_doe"),
					},
				},
			},
		}
	})

	// Second call to Query returns no error and hasNextPage is false
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Modify the result of the query to have no more pages
		repoQuery := args.Get(1).(*struct {
			Repository struct {
				PullRequests struct {
					PageInfo struct {
						HasNextPage githubv4.Boolean
						EndCursor   githubv4.String
					}
					Edges []PullRequestEdge
				} `graphql:"pullRequests(first: 100, states: $states, after: $cursor)"`
			} `graphql:"repository(name: $repo, owner: $org)"`
		})
		repoQuery.Repository.PullRequests.PageInfo.HasNextPage = false
		repoQuery.Repository.PullRequests.PageInfo.EndCursor = "endCursor"
		repoQuery.Repository.PullRequests.Edges = []PullRequestEdge{}
	})

	// Test with mocked client
	prsResult, err := getPrFromRepo(mockClient, org, repo)
	if assert.NoError(t, err) {
		assert.Equal(t, prs, prsResult)
	}
	mockClient.AssertExpectations(t)
}
