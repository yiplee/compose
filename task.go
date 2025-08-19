package compose

import (
	"io"
	"time"
)

type Task struct {
	Name string
	Cmd  string
	Args []string

	Delay time.Duration

	Out io.Writer
	Err io.Writer
	
	// Global settings
	Env        []string
	WorkingDir string
}
