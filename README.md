# Compose

Compose is a tool for defining and running multi commands. With Compose, you use a YAML file to configure your application's services. Then, with a single command, you create and start all the services from your configuration.

### Configure

Place compose.yaml in the working dir, then run ```compose``` to start the services.

```yaml
tasks:
    service_1: # service name
        cmds: sleep 5
    service_2:
        cmds: sleep 5
        delay: 1s # delay before start
```

### Systemd

Run ```compose systemd``` in the working dir to generate a systemd service file.
