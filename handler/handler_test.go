package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gist-service/client"

	"github.com/google/go-github/v93/github"
)

type clientStub struct {
	responseFunc func(ctx context.Context, username string) ([]*github.Gist, *github.Response, error)
}

func (m *clientStub) ListUserGists(ctx context.Context, username string) ([]*github.Gist, *github.Response, error) {
	return m.responseFunc(ctx, username)
}

func TestRootEndpoint(t *testing.T) {
	server := NewServer(&clientStub{})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("expected valid JSON response: %v", err)
	}
	if res["message"] == "" {
		t.Fatalf("expected message in root response")
	}
}

// https://docs.github.com/en/rest/gists/gists?apiVersion=2022-11-28#list-gists-for-a-user--status-codes
func TestGetUserGists_200(t *testing.T) {
	gistID := "6cad326836d38bd3a7ae"
	desc := "Hello world!"
	isPublic := true
	login := "octocat"

	sampleGist := &github.Gist{
		ID:          &gistID,
		Description: &desc,
		Public:      &isPublic,
		Owner: &github.User{
			Login: &login,
		},
	}

	stubClient := &clientStub{
		responseFunc: func(ctx context.Context, username string) ([]*github.Gist, *github.Response, error) {
			httpResp := &http.Response{
				StatusCode: http.StatusOK,
			}
			resp := &github.Response{Response: httpResp}
			return []*github.Gist{sampleGist}, resp, nil
		},
	}

	server := NewServer(stubClient)
	req := httptest.NewRequest(http.MethodGet, "/octocat", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var gists []*github.Gist
	if err := json.Unmarshal(rec.Body.Bytes(), &gists); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(gists) != 1 || *gists[0].ID != gistID || *gists[0].Owner.Login != "octocat" {
		t.Fatalf("unexpected gists response: %+v", gists)
	}
}

// https://docs.github.com/en/rest/gists/gists?apiVersion=2022-11-28#list-gists-for-a-user--status-codes
func TestGetUserGists_422(t *testing.T) {
	stubClient := &clientStub{
		responseFunc: func(ctx context.Context, username string) ([]*github.Gist, *github.Response, error) {
			httpResp := &http.Response{StatusCode: http.StatusUnprocessableEntity}
			resp := &github.Response{Response: httpResp}
			return nil, resp, fmt.Errorf("validation failed, or the endpoint has been spammed")
		},
	}

	server := NewServer(stubClient)
	req := httptest.NewRequest(http.MethodGet, "/octocat", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", rec.Code)
	}

	var errResp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse JSON error: %v", err)
	}

	if errResp.Error != "validation failed, or the endpoint has been spammed" {
		t.Errorf("expected error message 'user not found', got %q", errResp.Error)
	}
}

// runs against live github api
func TestGetUserGists_LiveGitHubOctocat(t *testing.T) {
	ghClient, err := client.NewGitHubClient()
	if err != nil {
		t.Fatalf("failed to create github client: %v", err)
	}
	server := NewServer(ghClient)

	req := httptest.NewRequest(http.MethodGet, "/octocat", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 from live GitHub API for octocat, got %d: %s", rec.Code, rec.Body.String())
	}

	var gists []*github.Gist
	if err := json.Unmarshal(rec.Body.Bytes(), &gists); err != nil {
		t.Fatalf("expected JSON array response: %v", err)
	}

	if len(gists) == 0 {
		t.Fatalf("expected octocat to have at least one gist")
	}

	firstGist := gists[0]
	if firstGist.Owner == nil || firstGist.Owner.Login == nil || *firstGist.Owner.Login != "octocat" {
		t.Fatalf("expected owner login to be octocat, got %v", firstGist.Owner)
	}
}
