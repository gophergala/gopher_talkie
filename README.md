# Gopher Talkie
`talkie` is a command-line tool for sending secure voice messages.

## Features
* Easy-to-use CLI for recording voice messages
* Sending voice message securely to anyone using PGP encryption


## Use Case
<TODO>

## Install

* Install PortAudio
  * Mac: `brew install portaudio`
  * Ubuntu: `apt-get install portaudio19-dev`

* Install GunPG and configure it
  * Mac: `brew install gpg`
  * Ubuntu: `apt-get install gnu-pg`

* Run `go get github.com/gophergala/gopher_talkie/src/talkie`

## Usage
`talkie` assumes [GnuPG](https://www.gnupg.org/) is already installed on your Linux/Mac machine.

```
NAME:
   talkie - Secure voicing messaging for geeks

USAGE:
   talkie [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
  Tom Li - <nklizhe@gmail.com>

COMMANDS:
   list   list all messages
   send   record and send a voice message
   help, h  Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --server "130.211.156.226:3333"  
   --help, -h       show help
   --version, -v      print the version
   
 ```

* Connect to your
## Build
We are using [GPM](https://github.com/pote/gpm) and [GVP](https://github.com/pote/gvp) to manage Go packages.

* Setting up the Go development environment.
  * Mac OS X: `./bootstrap.darwin.sh`
  * Ubuntu: `sudo ./bootstrap.linux.sh`

* Run

```
cd gopher_talkie/
make deps
make
```

## Start a 'talkie' server
* Build the server: `make server`
* Generate a new PGP key for the server: `gpg --gen-key` (Note: use an empty passphrase)
* Start the server: `talkie-server --server-key <PGP Key>`


