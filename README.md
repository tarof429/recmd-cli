# Command runner for Go

## Introduction

recmd is a small tool for running commands. The tool user interface was inspired by the *docker* command. Commands are run *in place* and do not involve agents or monitoring of external processes. As such, this tool is best suited for commands that have short execution times.

## Quick start

```bash
$ ./recmd-cli 

        recmd-cli is a command runner which manages commands. You can store
        commands in-line or execute scripts. It supports simple CRUD operations. 
        Commands can be modified by editing the JSON configuration file.

Usage:
  recmd-cli [command]

Available Commands:
  add         Add a command
  delete      Delete a command
  help        Help about any command
  init        Initialize the application
  inspect     Inspect the command
  list        List commands
  run         Run a command
  search      Search for a command

Flags:
      --config string   config file (default is $HOME/.recmd-cli.yaml)
  -h, --help            help for recmd-cli
  -t, --toggle          Help message for toggle

Use "recmd-cli [command] --help" for more information about a command.
```

## Configuration

recmd-cli stores commands in $HOME/.cmd_history.json. 

## Usage 

```bash
$ ./recmd-cli list
COMMAND HASH            COMMAND STRING                                  COMMENT                                                 DURATION
4a8a9fc31dc15a4         df                                              Show disk usage                                         0 second(s)
f10cad261de273f         hostname -i | awk -F" " '{print $1}'            Show IP address                                         0 second(s)
0103f67a0bc0b4e         docker images -a | grep "^<none>"               List all images with no tag                             0 second(s)
ccec33fd116787f         docker run hello-world                          Hello world from docker                                 1 second(s)
721c0b35b5ab6b5         curl http://worldclockapi.com/api/json/e...     REST Services that will return current date/time i...   0 second(s)
d8e7c90b269a1ca         sleep 3; echo hello!                            Sleep...                                                3 second(s)
37fa265330ad83e         pwd                                             List current directory                                  0 second(s)
ebfdec641529d4b         ls                                              List files                                              0 second(s)
920f8f5815b381e         env                                             env                                                     0 second(s) 

$ ./recmd-cli run d8e7c90b269a1ca
⠋ Scheduling commmand hello!

✓ Scheduling commmand 
$ ./recmd-cli add "sudo ss -tulpn | grep :80" "Find what process is listening to port 80"
$ ./recmd-cli list|grep 80
6e2d304e213958e         sudo ss -tulpn | grep :80                       Find what process is listening to port 80               -
$ ./recmd-cli run 6e2d304e213958e
⠋ Scheduling commmand tcp   LISTEN 0      511          0.0.0.0:8000      0.0.0.0:*    users:(("nginx",pid=543724,fd=5),("nginx",pid=543723,fd=5))
tcp   LISTEN 0      511          0.0.0.0:80        0.0.0.0:*    users:(("nginx",pid=543724,fd=4),("nginx",pid=543723,fd=4))

✓ Scheduling commmand 

