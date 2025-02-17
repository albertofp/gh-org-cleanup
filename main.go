package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/albertofp/gh-org-cleanup/pkg/gh"
	"github.com/albertofp/gh-org-cleanup/pkg/slack"
)

func main() {
	ghUtil := gh.New()
	slackUtil := slack.New()

	ctx := context.Background()

	ghUsers, err := ghUtil.GetOrgMembers()
	if err != nil {
		panic(err)
	}

	slackUsers, err := slackUtil.GetUsers(ctx)
	if err != nil {
		panic(err)
	}

	slackSet := make(map[string]struct{})
	for _, su := range slackUsers {
		ghHandle := su.GithubHandle
		if ghHandle != "" {
			slackSet[ghHandle] = struct{}{}
		}
	}

	usersNotInSlack := make([]string, 0)
	for _, ghUser := range ghUsers {
		if _, found := slackSet[ghUser]; !found {
			usersNotInSlack = append(usersNotInSlack, ghUser)
			fmt.Println(ghUser)
		}
	}

	log.Info(len(usersNotInSlack), "users not in slack")
}
