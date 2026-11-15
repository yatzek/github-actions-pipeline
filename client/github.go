package client

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v93/github"
)

type GistClient interface {
	ListUserGists(ctx context.Context, username string) ([]*github.Gist, *github.Response, error)
}

type GitHubClient struct {
	client *github.Client
}

const timeout = 10 * time.Second

func NewGitHubClient() (*GitHubClient, error) {
	client, err := github.NewClient(
		github.WithTimeout(timeout),
	)

	if err != nil {
		return nil, fmt.Errorf("error when initializing github client: %w", err)
	}

	return &GitHubClient{client: client}, nil
}

func (c *GitHubClient) ListUserGists(ctx context.Context, username string) ([]*github.Gist, *github.Response, error) {
	return c.client.Gists.List(ctx, username, nil)
}
