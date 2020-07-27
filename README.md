# Command runner for Go

## Introduction

recmd is a small tool for running commands. The tool user interface was inspired by the *docker* command. Commands are run *in place* and do not involve agents or monitoring of external processes. As such, this tool is best suited for commands that have short execution times.

## Quick start

Invoking recmd-cli will display the available commands.

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

The command to add requires two flags.

```bash
$ ./recmd-cli add
Usage: recmd-cli add -c <command> -i <comment>
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
$ ./recmd-cli add -c "sudo ss -tulpn | grep :80" -i "Find what process is listening to port 80"
$ ./recmd-cli list|grep 80
6e2d304e213958e         sudo ss -tulpn | grep :80                       Find what process is listening to port 80               -
$ ./recmd-cli run 6e2d304e213958e
⠋ Scheduling commmand tcp   LISTEN 0      511          0.0.0.0:8000      0.0.0.0:*    users:(("nginx",pid=543724,fd=5),("nginx",pid=543723,fd=5))
tcp   LISTEN 0      511          0.0.0.0:80        0.0.0.0:*    users:(("nginx",pid=543724,fd=4),("nginx",pid=543723,fd=4))

✓ Scheduling commmand 
```

It is possible to add commands that read from environment variables. For example, we could add a command that takes a port variable.

```bash
$ ./recmd-cli add -c 'sudo ss -tulpn | grep :$port' -p "port=$port" -i "Find what process is listenning to port"
```

Then to use this function, we would set the port variable just as we are invoking recmd-cli.

```bash
$ ./recmd-cli list |grep \$port
6d83258f3296b50         sudo ss -tulpn | grep :$port                    Find what process is listenning to port                 0 second(s)

$ export port=8080;./recmd-cli run 6d83258f3296b50;unset port
✓ Scheduling commmand 
tcp   LISTEN 0      4096               *:8080             *:*    users:(("hello",pid=160507,fd=3))  
```

## Gotchas!

If your command uses backticks, use single-quotes around the command instead of double-quotes to prevent the tool from storing the result of the command. For example:

With double quotes:

```bash
$ ./recmd-cli add  -c "(echo `ls -al ~/. | wc -l` - 2) | bc" -i "Count of dot files in home directory"
[taro@zaxman recmd-cli]$ ./recmd-cli list
COMMAND HASH            COMMAND STRING                                  COMMENT                                                 DURATION
e75710f201f513f         (echo 155 - 2) | bc                             Count of dot files in home directory                    -
```

With single quotes:

```bash
$ ./recmd-cli add  -c '(echo `ls -al ~/. | wc -l` - 2) | bc' -i "Count of dot files in home directory"
[taro@zaxman recmd-cli]$ ./recmd-cli list
COMMAND HASH            COMMAND STRING                                  COMMENT                                                 DURATION
2f35231f613da42         (echo `ls -al ~/. | wc -l` - 2) | bc            Count of dot files in home directory                    -
```
