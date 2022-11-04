package compose

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func run(ctx context.Context, t *Task) error {
	cmd := exec.Command(t.Cmd, t.Args...)

	prefix := fmt.Sprintf("[%s] ", t.Name)
	cmd.Stderr = &PrefixWriter{Out: t.Err, Prefix: prefix}
	cmd.Stdout = &PrefixWriter{Out: t.Out, Prefix: prefix}

	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan struct{})
	defer func() {
		close(done)
	}()

	go func() {
		select {
		case <-ctx.Done():
			cmd.Process.Signal(syscall.SIGTERM)
			select {
			case <-time.After(3 * time.Second):
				cmd.Process.Kill()
			case <-done:
			}
		case <-done:
		}
	}()

	return cmd.Wait()
}

func Run(ctx context.Context, t *Task) error {
	if t.Delay > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(t.Delay):
		}
	}

	for {
		log.Info().Str("cmd", t.Cmd).Strs("args", t.Args).Msgf("starting %s", t.Name)

		err := run(ctx, t)

		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			log.Info().Err(err).Msgf("finished %s", t.Name)
			return err
		}

		log.Info().Err(err).Msgf("exited %s with status code %d", t.Name, exitErr.ExitCode())

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
}
