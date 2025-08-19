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
	cmd.Stderr = NewPrefixWriter(t.Err, prefix)
	cmd.Stdout = NewPrefixWriter(t.Out, prefix)

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

	// Pre-allocate strings to avoid repeated allocations
	startMsg := fmt.Sprintf("starting %s", t.Name)
	finishedMsg := fmt.Sprintf("finished %s", t.Name)
	exitedMsg := fmt.Sprintf("exited %s with status code", t.Name)

	for {
		log.Info().Str("cmd", t.Cmd).Strs("args", t.Args).Msg(startMsg)

		err := run(ctx, t)

		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			log.Info().Err(err).Msg(finishedMsg)
			return err
		}

		log.Info().Err(err).Msgf("%s %d", exitedMsg, exitErr.ExitCode())

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
}
