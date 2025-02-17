package slack

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	sl "github.com/slack-go/slack"
)

const (
	GithubHandleLabelId = "Xf06P38LGMN3"
)

type SlackUser struct {
	SlackID      string `yaml:"slack_id"`
	Name         string `yaml:"name"`
	Email        string `yaml:"email"`
	GithubHandle string `yaml:"github_handle"`
}

type SlackUtil struct {
	client *sl.Client
	config *SlackUtilConfig
}

type SlackUtilConfig struct {
	GithubHandleLabelId string
}

func New() *SlackUtil {
	client := sl.New(os.Getenv("SLACK_TOKEN"))

	return &SlackUtil{
		client: client,
		config: &SlackUtilConfig{
			GithubHandleLabelId: GithubHandleLabelId,
		},
	}
}

func (s *SlackUtil) getActiveUsers(ctx context.Context) ([]sl.User, error) {
	users, err := s.client.GetUsersContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active slack users: %w", err)
	}
	var active []sl.User
	for _, u := range users {
		if !u.Deleted {
			active = append(active, u)
		}
	}
	return active, nil
}

func (s *SlackUtil) GetUsers(ctx context.Context) ([]SlackUser, error) {
	start := time.Now()
	slackUsers := make([]SlackUser, 0)
	users, err := s.getActiveUsers(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.IsBot {
			continue
		}

		profile, err := s.client.GetUserProfile(&sl.GetUserProfileParameters{
			UserID:        user.ID,
			IncludeLabels: false,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get user profile: %w", err)
		}

		ghHandle := profile.FieldsMap()[s.config.GithubHandleLabelId].Value
		slackUser := SlackUser{
			SlackID:      user.ID,
			Name:         user.Name,
			Email:        user.Profile.Email,
			GithubHandle: ghHandle,
		}
		slackUsers = append(slackUsers, slackUser)
	}

	log.Info("Time: ", time.Since(start).Seconds())
	return slackUsers, nil
}
