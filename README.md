# Genie
Store and execute commands based on yaml configuration files.

## Installation

### Prebuilt binaries
- Download the binary for your platform from the Github releases page
- Copy the executable to a location within your PATH (ex: /usr/local/bin)

### From Source
The tool is written in Golang. To compile and run the program first
install your Golang development environment and run the following commands

```shell
go build
mv genie /usr/local/bin
```

## Usage

Genie looks for command.yaml files in your current directory and its
children
### Create commands.yaml file
```shell
genie init
```

### List detected commands
```shell
genie
```

#### Execute command
```shell
genie {command name}
```

### Global commands.yaml file
Genie will load global commands from a `.genie-commands.yaml` file
in your home directory
```yaml
---
shell: /bin/zsh
commands:
  hello:
    - command: echo "From the home directory!"
    - command: echo done
```

## Example commands.yaml
```yaml
---
shell: /bin/zsh
commands:
  hello:
    - command: echo "Hello world, how're we today!"
    - command: echo done
  something:
    - command: echo
  test-environment:
    - command: echo It is $ENV before $ENV2
      environment:
        - name: ENV
          value: ITWORKS
        - name: ENV2
          value: ITWORKS2
```
