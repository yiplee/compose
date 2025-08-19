/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"

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
	tasks := loadTasks(cmd)
	if len(tasks) == 0 {
		return errors.New("no tasks")
	}

	var g errgroup.Group
	for idx := range tasks {
		t := tasks[idx]
		t.Out = cmd.OutOrStdout()
		t.Err = cmd.ErrOrStderr()
		g.Go(func() error {
			return compose.Run(cmd.Context(), t)
		})
	}

	return g.Wait()
}

func loadTasks(cmd *cobra.Command) []*compose.Task {
	v := viper.Sub("tasks")
	if v == nil {
		return nil
	}

	// Get global arguments and environment variables
	globalArgs, _ := cmd.Flags().GetStringSlice("global-args")
	globalEnv, _ := cmd.Flags().GetStringSlice("global-env")
	globalWorkingDir, _ := cmd.Flags().GetString("global-working-dir")

	// Pre-allocate slice with known capacity
	settings := v.AllSettings()
	tasks := make([]*compose.Task, 0, len(settings))
	
	for name := range settings {
		log.Debug().Msgf("load task %s", name)
		t := v.Sub(name)
		cmds := t.GetString("cmds")

		if args, err := shlex.Split(cmds); err == nil && len(args) > 0 {
			// Apply global arguments to each task
			finalArgs := make([]string, 0, len(args)+len(globalArgs))
			finalArgs = append(finalArgs, args...)
			finalArgs = append(finalArgs, globalArgs...)
			
			task := &compose.Task{
				Name:  name,
				Cmd:   args[0],
				Args:  finalArgs[1:], // Skip the command name, include original args + global args
				Delay: t.GetDuration("delay"),
			}
			
			// Set global environment variables if specified
			if len(globalEnv) > 0 {
				task.Env = globalEnv
			}
			
			// Set global working directory if specified
			if globalWorkingDir != "" {
				task.WorkingDir = globalWorkingDir
			}
			
			tasks = append(tasks, task)
		}
	}

	return tasks
}
