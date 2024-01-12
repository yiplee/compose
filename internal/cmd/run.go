/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"os/signal"
	"syscall"

	"github.com/google/shlex"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yiplee/compose"
	"golang.org/x/sync/errgroup"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTasks(cmd)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runTasks(cmd *cobra.Command) error {
	tasks := loadTasks()
	if len(tasks) == 0 {
		return errors.New("no tasks")
	}

	ctx, stop := signal.NotifyContext(
		cmd.Context(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)
	for idx := range tasks {
		t := tasks[idx]
		t.Out = cmd.OutOrStdout()
		t.Err = cmd.ErrOrStderr()
		g.Go(func() error {
			return compose.Run(ctx, t)
		})
	}

	return g.Wait()
}

func loadTasks() []*compose.Task {
	v := viper.Sub("tasks")
	if v == nil {
		return nil
	}

	var tasks []*compose.Task
	for name := range v.AllSettings() {
		log.Debug().Msgf("load task %s", name)
		t := v.Sub(name)
		cmds := t.GetString("cmds")

		if args, err := shlex.Split(cmds); err == nil && len(args) > 0 {
			tasks = append(tasks, &compose.Task{
				Name:  name,
				Cmd:   args[0],
				Args:  args[1:],
				Delay: t.GetDuration("delay"),
			})
		}
	}

	return tasks
}
