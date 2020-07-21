# Command runner for Go

## Introduction

recmd is a small tool for running commands stored in a JSON file. You can search tools by hash or part of the command string. Output will also be in JSON format. This tool was inspired by the *docker* command but is more general purpose. Commands are run *in place* and do not involve agents or process monitoring. As such, this tool is best suited for commands that have short execution times.

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

## Sample 

```bash
$ ./recmd-cli list
COMMAND HASH            COMMAND STRING                                  COMMAND COMMENT                                   
4a8a9fc31dc15a4         df                                              Show disk usage                                   
f10cad261de273f         hostname -i | awk -F" " '{print $1}'            Show IP address                                   
0103f67a0bc0b4e         docker images -a | grep "^<none>"               List all images with no tag                       
ccec33fd116787f         docker run hello-world                          Hello world from docker                           
721c0b35b5ab6b5         curl http://worldclockapi.com/api/json/e...     REST Services that will return current date/time i...

$ cat ~/.cmd_history.json
[
	{
		"commandHash": "4a8a9fc31dc15a4b87bb145b05db3a",
		"commandString": "df",
		"comment": "Show disk usage",
		"creationTime": "2020-07-19T14:12:10.588141643-07:00",
		"lastExecutionTime": "2020-07-19T14:12:10.588141793-07:00"
	},
	{
		"commandHash": "f10cad261de273f596a58bc87f683c",
		"commandString": "hostname -i | awk -F\" \" '{print $1}'",
		"comment": "Show IP address",
		"creationTime": "2020-07-20T15:49:24.800631745-07:00",
		"lastExecutionTime": "2020-07-20T15:49:24.800631885-07:00"
	},
	{
		"commandHash": "0103f67a0bc0b4e3d0447cacfe1632",
		"commandString": "docker images -a | grep \"^\u003cnone\u003e\"",
		"comment": "List all images with no tag",
		"creationTime": "2020-07-20T20:32:18.12061626-07:00",
		"lastExecutionTime": "2020-07-20T20:32:18.1206164-07:00"
	},
	{
		"commandHash": "ccec33fd116787f035eb6f92beef28",
		"commandString": "docker run hello-world",
		"comment": "Hello world from docker",
		"creationTime": "2020-07-20T20:58:29.978073568-07:00",
		"lastExecutionTime": "2020-07-20T20:58:29.978073719-07:00"
	},
	{
		"commandHash": "721c0b35b5ab6b5471c7f48a52058e",
		"commandString": "curl http://worldclockapi.com/api/json/est/now",
		"comment": "REST Services that will return current date/time in JSON for any registered time zone.",
		"creationTime": "2020-07-20T21:26:50.522437605-07:00",
		"lastExecutionTime": "2020-07-20T21:26:50.522437736-07:00"
	}
]$

$ ./recmd-cli run 72
{"$id":"1","currentDateTime":"2020-07-21T00:42-04:00","utcOffset":"-04:00:00","isDayLightSavingsTime":true,"dayOfTheWeek":"Tuesday","timeZoneName":"Eastern Standard Time","currentFileTime":132397657727273369,"ordinalDate":"2020-203","serviceResponse":null}

```

Commands that take an additional parameter may either use the command hash or the comamnd string. It is usually best to use the command hash. For example, *./recmd-cli run "docker" would run both commands that begin with *docker*. If we only want to run the hello-world docker container, we could either specify the hash or provide a more complete command string that would distinguish it from any other.

```bash
$ ./recmd-cli run "docker run hello"

Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:
 1. The Docker client contacted the Docker daemon.
 2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
    (amd64)
 3. The Docker daemon created a new container from that image which runs the
    executable that produces the output you are currently reading.
 4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.

To try something more ambitious, you can run an Ubuntu container with:
 $ docker run -it ubuntu bash

Share images, automate workflows, and more with a free Docker ID:
 https://hub.docker.com/

For more examples and ideas, visit:
 https://docs.docker.com/get-started/
 ```

