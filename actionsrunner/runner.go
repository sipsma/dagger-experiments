package actionsrunner

import (
	"context"
	"strings"

	"dagger.io/dagger"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Token            string
	Repo             string
	RunnerNamePrefix string
	Labels           []string
	Count            int
}

func Run(ctx context.Context, c *dagger.Client, cfg Config) error {
	ctr := c.Container().
		From("myoung34/github-runner:latest").
		WithEnvVariable("ACCESS_TOKEN", cfg.Token).
		WithEnvVariable("REPO_URL", cfg.Repo).
		WithEnvVariable("LABELS", strings.Join(cfg.Labels, ",")).
		WithEnvVariable("EPHEMERAL", "true")

	eg, ctx := errgroup.WithContext(ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for i := 0; i < cfg.Count; i++ {
		eg.Go(func() error {
			defer cancel()
			for { // loop exits when context is canceled which results in non-nil error below
				// get a random uuid
				id, err := uuid.NewRandom()
				if err != nil {
					return err
				}
				_, err = ctr.
					WithEnvVariable("RUNNER_NAME", cfg.RunnerNamePrefix+"-"+id.String()).
					WithExec(nil, dagger.ContainerWithExecOpts{
						ExperimentalPrivilegedNesting: true,
					}).ExitCode(ctx)
				if err != nil {
					return err
				}
			}
		})
	}

	return eg.Wait()
}
