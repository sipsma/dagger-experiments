package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"dagger.io/dagger"
	"github.com/google/go-github/v48/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/sipsma/dagger-experiments/actionsrunner"
)

func main() {
	ctx := context.Background()

	token, ok := os.LookupEnv("GHA_RUNNER_TOKEN")
	if !ok {
		panic(fmt.Errorf("GHA_RUNNER_TOKEN is not set"))
	}
	repo, ok := os.LookupEnv("GHA_RUNNER_REPO")
	if !ok {
		panic(fmt.Errorf("GHA_RUNNER_REPO is not set"))
	}

	labels := []string{"dagger-runner"}
	if envLabels, ok := os.LookupEnv("GHA_RUNNER_LABELS"); ok {
		labels = strings.Split(envLabels, ",")
	}

	runnerPrefix := "test-dagger-runner"
	if envRunnerPrefix, ok := os.LookupEnv("GHA_RUNNER_RUNNER_PREFIX"); ok {
		runnerPrefix = envRunnerPrefix
	}

	appID, err := strconv.Atoi(os.Getenv("GH_APP_ID"))
	if err != nil {
		panic(err)
	}

	webhookSecret, ok := os.LookupEnv("GH_WEBHOOK_SECRET")
	if !ok {
		panic(fmt.Errorf("GH_WEBHOOK_SECRET is not set"))
	}

	privateKey, ok := os.LookupEnv("GH_PRIVATE_KEY")
	if !ok {
		panic(fmt.Errorf("GH_PRIVATE_KEY is not set"))
	}

	cfg := githubapp.Config{
		V3APIURL: "https://api.github.com/",
		App: struct {
			IntegrationID int64  `yaml:"integration_id" json:"integrationId"`
			WebhookSecret string `yaml:"webhook_secret" json:"webhookSecret"`
			PrivateKey    string `yaml:"private_key" json:"privateKey"`
		}{
			IntegrationID: int64(appID),
			WebhookSecret: webhookSecret,
			PrivateKey:    privateKey,
		},
	}

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	h := &handler{
		daggerClient: c,
		token:        token,
		repo:         repo,
		labels:       labels,
		runnerPrefix: runnerPrefix,
	}

	http.Handle("/", githubapp.NewDefaultEventDispatcher(cfg, h))

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 45363)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}

type handler struct {
	daggerClient *dagger.Client
	token        string
	repo         string
	labels       []string
	runnerPrefix string
}

func (h *handler) Handles() []string {
	return []string{"workflow_job"}
}

func (h *handler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	fmt.Println("Received event", eventType)
	if eventType != "workflow_job" {
		return nil
	}
	var event github.WorkflowJobEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}
	runnerLabels := event.WorkflowJob.Labels
	if len(runnerLabels) == 0 {
		return nil
	}
	var isDaggerRunner bool
	for _, label := range runnerLabels {
		// n^2 but who cares
		for _, otherLabel := range h.labels {
			if label == otherLabel {
				isDaggerRunner = true
				break
			}
		}
	}
	if !isDaggerRunner {
		fmt.Printf("received job for different runner %v, ignoring\n", runnerLabels)
		return nil
	}

	go func() {
		fmt.Println("starting actions runner", eventType)
		err := actionsrunner.Run(context.TODO(), h.daggerClient, actionsrunner.Config{
			Token:            h.token,
			Repo:             h.repo,
			Labels:           h.labels,
			RunnerNamePrefix: h.runnerPrefix,
			Count:            1,
		})
		if err != nil {
			fmt.Printf("runner error: %v\n", err)
		}
	}()

	return nil
}
