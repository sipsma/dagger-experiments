package actionsrunner

import (
	"dagger.io/dagger"
)

func ActionsRunner(c *dagger.Client, token, repo, runnerName string, labels []string) (*dagger.Container, error) {
	ctr := c.Container().
		From("myoung34/github-runner:latest").
		WithEnvVariable("ACCESS_TOKEN", token).
		WithEnvVariable("REPO_URL", repo).
		WithEnvVariable("LABELS", "dagger-runner").
		WithEnvVariable("RUNNER_NAME", "test-dagger-runner")

	return ctr, nil
}
