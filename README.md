# Command runner for Go

## Introduction

`recmd` is a small tool for running commands. The tool's user interface is inspired by the *docker* command.

## Quick start

Invoking recmd-cli will display the available commands.

```bash
$ ./recmd-cli

	recmd-cli is a command runner which manages commands. You can store commands in-line or execute scripts. It supports simple CRUD operations.

Usage:
  recmd-cli [command]

Available Commands:
  add         Add a command
  delete      Delete a command
  help        Help about any command
  list        List commands
  run         Run a command
  search      Search for a command by its comment
  select      Select a command by is hash

Flags:
      --config string   config file (default is $HOME/.recmd-cli.yaml)
  -h, --help            help for recmd-cli
  -t, --toggle          Help message for toggle

Use "recmd-cli [command] --help" for more information about a command.
```

The command to add requires two flags.

```bash
$ ./recmd-cli add
Usage: recmd-cli add -c <command> -d <description>
```

## Configuration

recmd-cli stores commands in $HOME/.cmd_history.json. Use recmd-cli init to create it.

## Usage 

First start `recmd-dmn`. 

```bash
$ ./recmd-dmn
2020/10/23 17:47:08 Starting server on :8999
```

```bash
$ ./recmd-cli list
HASH                    COMMAND                                         DESCRIPTION                                             DURATION
4a8a9fc31dc15a4         df                                              Show disk usage                                         0 second(s)
f10cad261de273f         hostname -i | awk -F" " '{print $1}'            Show IP address                                         0 second(s)
0103f67a0bc0b4e         docker images -a | grep "^<none>"               List all images with no tag                             0 second(s)
ccec33fd116787f         docker run hello-world                          Hello world from docker                                 1 second(s)
721c0b35b5ab6b5         curl http://worldclockapi.com/api/json/e...     REST Services that will return current date/time i...   0 second(s)
d8e7c90b269a1ca         sleep 3; echo hello!                            Sleep...                                                3 second(s)
37fa265330ad83e         pwd                                             List current directory                                  0 second(s)
ebfdec641529d4b         ls                                              List files                                              0 second(s)
920f8f5815b381e         env                                             env                                                     0 second(s) 

$  ./recmd-cli run 18b63bce19510d0
✓ Scheduling commmand
Using default tag: latest
latest: Pulling from ubuntucore/jenkins-ubuntu
Digest: sha256:662cdd29e6b2b8d50f2c4e2ffa5121740c40bc518848ff912065d6a163846e65
Status: Image is up to date for ubuntucore/jenkins-ubuntu:latest
docker.io/ubuntucore/jenkins-ubuntu:latest
```

If you know the command will run for a long period of time, and you do not want `recmd` to block, there is a -b option which returns control back to the user after a 1 second delay. The command below likewise pulls the ubuntu image.

```bash
$   ./recmd-cli run 18b63bce19510d0 -b
✓ Scheduling commmand
```

## Gotchas!

If your command uses backticks, use single-quotes around the command instead of double-quotes to prevent the tool from storing the result of the command. For example:

With double quotes:

```bash
$ ./recmd-cli add  -c '(echo `ls -al ~/. | wc -l` - 2) | bc' -d "Count of dot files in home directory"
$ ./recmd-cli search count
[
	{
		"commandHash": "2f35231f613da4276feb1c4274375c",
		"commandString": "(echo `ls -al ~/. | wc -l` - 2) | bc",
		"description": "Count of dot files in home directory",
		"duration": -1
	}
]
```
