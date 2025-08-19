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

	// Get global arguments and environment variables from command line flags
	globalArgs, _ := cmd.Flags().GetStringSlice("global-args")
	globalEnv, _ := cmd.Flags().GetStringSlice("global-env")
	globalWorkingDir, _ := cmd.Flags().GetString("global-working-dir")

	// Get global settings from YAML configuration (only if command line flags are not set)
	if len(globalArgs) == 0 && viper.IsSet("global.args") {
		globalArgs = viper.GetStringSlice("global.args")
		log.Debug().Strs("global_args", globalArgs).Msg("loaded global args from YAML config")
	}
	if len(globalEnv) == 0 && viper.IsSet("global.env") {
		globalEnv = viper.GetStringSlice("global.env")
		log.Debug().Strs("global_env", globalEnv).Msg("loaded global env from YAML config")
	}
	if globalWorkingDir == "" && viper.IsSet("global.working_dir") {
		globalWorkingDir = viper.GetString("global.working_dir")
		log.Debug().Str("global_working_dir", globalWorkingDir).Msg("loaded global working dir from YAML config")
	}

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
				log.Debug().Str("task", name).Strs("env", globalEnv).Msg("applied global env to task")
			}
			
			// Set working directory: task-specific overrides global
			taskWorkingDir := t.GetString("working_dir")
			if taskWorkingDir != "" {
				task.WorkingDir = taskWorkingDir
				log.Debug().Str("task", name).Str("working_dir", taskWorkingDir).Msg("applied task-specific working dir")
			} else if globalWorkingDir != "" {
				task.WorkingDir = globalWorkingDir
				log.Debug().Str("task", name).Str("working_dir", globalWorkingDir).Msg("applied global working dir to task")
			}
			
			// Log when global args are applied
			if len(globalArgs) > 0 {
				log.Debug().Str("task", name).Strs("global_args", globalArgs).Msg("applied global args to task")
			}
			
			tasks = append(tasks, task)
		}
	}

	return tasks
}
