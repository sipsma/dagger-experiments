package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
	"github.com/sipsma/dagger-experiments/actionsrunner"
)

func main() {
	ctx := context.Background()

	token, ok := os.LookupEnv("GHA_ACTIONS_TOKEN")
	if !ok {
		panic(fmt.Errorf("GHA_ACTIONS_TOKEN is not set"))
	}
	repo, ok := os.LookupEnv("GHA_ACTIONS_REPO")
	if !ok {
		panic(fmt.Errorf("GHA_ACTIONS_REPO is not set"))
	}

	labels := []string{"dagger-runner"}
	if envLabels, ok := os.LookupEnv("GHA_ACTIONS_LABELS"); ok {
		labels = strings.Split(envLabels, ",")
	}

	runnerName := "test-dagger-runner"
	if envRunnerName, ok := os.LookupEnv("GHA_ACTIONS_RUNNER_NAME"); ok {
		runnerName = envRunnerName
	}

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	ctr, err := actionsrunner.ActionsRunner(c, token, repo, runnerName, labels)
	if err != nil {
		panic(err)
	}

	_, err = ctr.WithExec(nil, dagger.ContainerWithExecOpts{
		ExperimentalPrivilegedNesting: true,
	}).ExitCode(ctx)
	if err != nil {
		panic(err)
	}
}
