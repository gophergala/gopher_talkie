# Gopher Talkie
Secure voice messaging for geeks.

## Features
* Easy-to-use CLI for recording voice messages
* Sending voice message securely to anyone using PGP encryption

## Use Case
<TODO>

## Build

### Setting up the Go development environment.
* Mac OS X: `./bootstrap.darwin.sh`
* Ubuntu: `sudo ./bootstrap.linux.sh`

### Install PortAudio
* Mac: `brew install portaudio`
* Ubuntu: `sudo apt-get install portaudio`

### Build
```
make
```

## Usage
```
NAME:
   talkie - Secure voicing messaging for geeks

USAGE:
   talkie [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
  nklizhe@gmail.com - <unknown@email>

COMMANDS:
   list		list all messages
   send		record and send a voice message
   listen	listen messages
   delete	delete a message
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
 ```
