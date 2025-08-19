# Global Arguments Support

This document explains how to use the new global arguments feature that allows you to apply common arguments, environment variables, and working directory settings to all subcommands.

## Overview

Global arguments are persistent flags that get applied to all tasks when using the `compose run` command. This is useful when you want to:

- Add common flags to all commands (e.g., `--verbose`, `--config`)
- Set environment variables for all processes
- Specify a common working directory for all tasks

## Available Global Flags

### `--global-args`
Appends additional arguments to all commands. You can specify multiple values.

**Example:**
```bash
compose run --global-args="--verbose" --global-args="--config=/etc/app.conf"
```

This will append `--verbose --config=/etc/app.conf` to every task's command line.

### `--global-env`
Sets environment variables for all commands. You can specify multiple values.

**Example:**
```bash
compose run --global-env="DEBUG=true" --global-env="LOG_LEVEL=info"
```

This will set `DEBUG=true` and `LOG_LEVEL=info` for every task.

### `--global-working-dir`
Sets the working directory for all commands.

**Example:**
```bash
compose run --global-working-dir="/opt/app"
```

This will change to `/opt/app` before executing each task.

## How It Works

1. **Global arguments** are appended to the end of each task's argument list
2. **Global environment variables** are merged with the system environment
3. **Global working directory** overrides the current working directory for each task

## Example Configuration

Given this `compose.yaml`:

```yaml
tasks:
  web-server:
    cmds: "nginx -g 'daemon off;'"
  api-server:
    cmds: "node server.js --port 3000"
```

Running with global arguments:

```bash
compose run --global-args="--verbose" --global-args="--log-level=debug"
```

Will execute:
- `nginx -g 'daemon off;' --verbose --log-level=debug`
- `node server.js --port 3000 --verbose --log-level=debug`

## Use Cases

- **Development**: Add `--debug` or `--verbose` to all commands
- **Configuration**: Specify a common config file path for all services
- **Environment**: Set development/staging/production environment variables
- **Deployment**: Change working directory for containerized deployments

## Notes

- Global arguments are applied **after** the task's original arguments
- Environment variables are merged with existing system environment
- Working directory changes are applied before command execution
- These flags work with all existing subcommands that execute tasks