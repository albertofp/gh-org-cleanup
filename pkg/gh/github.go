package gh

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v60/github"
	log "github.com/sirupsen/logrus"
)

type GithubUtil struct {
	client *github.Client
	org    string
}

func New() *GithubUtil {
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	org := os.Getenv("GITHUB_ORG")
	if org == "" {
		log.Error("GITHUB_ORG is not set")
	}

	return &GithubUtil{
		client: client,
		org:    org,
	}
}

func (g *GithubUtil) GetOrgMembers() ([]string, error) {
	users := []string{}
	ctx := context.Background()

	opts := &github.ListMembersOptions{
		// This only lists 30 members maximum by default
		ListOptions: github.ListOptions{PerPage: 100},
	}

	ghUsers, _, err := g.client.Organizations.ListMembers(ctx, g.org, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list members: %w", err)
	}

	for _, user := range ghUsers {
		users = append(users, *user.Login)
	}

	return users, nil
}
