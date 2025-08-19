# Global Arguments Support

This document explains how to use the new global arguments feature that allows you to apply common arguments, environment variables, and working directory settings to all subcommands.

## Overview

Global arguments are persistent flags that get applied to all tasks when using the `compose run` command. This is useful when you want to:

- Add common flags to all commands (e.g., `--verbose`, `--config`)
- Set environment variables for all processes
- Specify a common working directory for all tasks

## Configuration Methods

Global arguments can be configured in two ways:

1. **YAML Configuration File** (recommended for persistent settings)
2. **Command Line Flags** (useful for temporary overrides)

**Note**: Command line flags override YAML configuration settings.

## Available Global Settings

### `global.args`
Appends additional arguments to all commands. You can specify multiple values.

**YAML Configuration:**
```yaml
global:
  args:
    - "--verbose"
    - "--log-level=info"
    - "--config=/etc/app.conf"
```

**Command Line:**
```bash
compose run --global-args="--verbose" --global-args="--config=/etc/app.conf"
```

### `global.env`
Sets environment variables for all commands. You can specify multiple values.

**YAML Configuration:**
```yaml
global:
  env:
    - "DEBUG=true"
    - "LOG_LEVEL=info"
    - "NODE_ENV=production"
```

**Command Line:**
```bash
compose run --global-env="DEBUG=true" --global-env="LOG_LEVEL=info"
```

### `global.working_dir`
Sets the working directory for all commands. This can be overridden by individual tasks.

**YAML Configuration:**
```yaml
global:
  working_dir: "/opt/app"
```

**Command Line:**
```bash
compose run --global-working-dir="/opt/app"
```

**Per-Task Override:**
```yaml
tasks:
  api-server:
    cmds: "node server.js"
    working_dir: "/opt/app/api"  # Overrides global working directory
  database:
    cmds: "postgres -D /var/lib/postgresql/data"
    working_dir: "/var/lib/postgresql"  # Overrides global working directory
```

## Working Directory Priority

Working directory is resolved in this order (highest to lowest priority):

1. **Task-specific** `working_dir` setting
2. **Global** `working_dir` setting (from YAML or command line)
3. **Current working directory** (default)

## Complete Example Configuration

```yaml
# Global settings that apply to all tasks (can be overridden per task)
global:
  args:
    - "--verbose"
    - "--log-level=info"
  env:
    - "NODE_ENV=production"
    - "LOG_LEVEL=info"
  working_dir: "/opt/app"  # Default working directory for all tasks

tasks:
  web-server:
    cmds: "nginx -g 'daemon off;'"
    # Uses global working_dir: "/opt/app"
    
  api-server:
    cmds: "node server.js --port 3000"
    working_dir: "/opt/app/api"  # Override global working directory
    
  database:
    cmds: "postgres -D /var/lib/postgresql/data"
    working_dir: "/var/lib/postgresql"  # Override global working directory
```

## How It Works

1. **YAML Configuration**: Global settings are read from the `global` section of your compose file
2. **Command Line Override**: If you specify command line flags, they override the YAML settings
3. **Task Execution**: Global arguments are appended to each task's command line, environment variables are merged, and working directory is set
4. **Working Directory Priority**: Task-specific `working_dir` overrides global `working_dir`, which overrides the current working directory

## Execution Examples

With the configuration above:

**Using YAML configuration only:**
```bash
compose run
```

Will execute:
- `nginx -g 'daemon off;' --verbose --log-level=info` (in `/opt/app` with `NODE_ENV=production`, `LOG_LEVEL=info`)
- `node server.js --port 3000 --verbose --log-level=info` (in `/opt/app` with `NODE_ENV=production`, `LOG_LEVEL=info`)
- `postgres -D /var/lib/postgresql/data --verbose --log-level=info` (in `/opt/app` with `NODE_ENV=production`, `LOG_LEVEL=info`)

**Overriding with command line flags:**
```bash
compose run --global-args="--debug" --global-env="DEBUG=true"
```

Will execute with `--debug` instead of `--verbose` and `DEBUG=true` instead of the YAML environment variables.

## Use Cases

- **Development**: Add `--debug` or `--verbose` to all commands
- **Configuration**: Specify a common config file path for all services
- **Environment**: Set development/staging/production environment variables
- **Deployment**: Change working directory for containerized deployments
- **Team Settings**: Share common configuration in version control

## Priority Order

1. **Command line flags** (highest priority)
2. **YAML configuration** (default values)
3. **System defaults** (lowest priority)

## Notes

- Global arguments are applied **after** the task's original arguments
- Environment variables are merged with existing system environment
- Working directory changes are applied before command execution
- These flags work with all existing subcommands that execute tasks
- YAML configuration provides a clean, maintainable way to manage global settings