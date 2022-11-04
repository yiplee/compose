package compose

import "io"

type Task struct {
	Name string
	Cmd  string
	Args []string

	Out io.Writer
	Err io.Writer
}
