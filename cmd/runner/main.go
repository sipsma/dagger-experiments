package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"dagger.io/dagger"
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

	count := 2
	if envCount, ok := os.LookupEnv("GHA_RUNNER_COUNT"); ok {
		var err error
		count, err = strconv.Atoi(envCount)
		if err != nil {
			panic(fmt.Errorf("GHA_RUNNER_COUNT is not a valid integer: %w", err))
		}
	}

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	err = actionsrunner.Run(ctx, c, actionsrunner.Config{
		Token:            token,
		Repo:             repo,
		Labels:           labels,
		RunnerNamePrefix: runnerPrefix,
		Count:            count,
	})
	if err != nil {
		panic(err)
	}
}
