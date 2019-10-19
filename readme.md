 Heimdall
----
[![GoDoc](https://godoc.org/github.com/DnOberon/heimdall/bifrost?status.svg)](https://godoc.org/github.com/DnOberon/heimdall/bifrost)

Heimdall allows you to manage and monitor command line applications. Provided with an executable or command to run `heimdall` can perform the following operations:

* Repeat the command n times
* Run command in parallel in n instances
* Filter and log command's output
* Kill hung applications through user specified timeout



## Install

```console
go get github.com/dnoberon/heimdall 
```

If you want to run heimdall outside `$GOPATH/bin` make sure that `$GOPATH/bin` is [included in your \$PATH](https://golang.org/doc/code.html#GOPATH).
If not, you can download  heimdall’s releases for your platform [here](https://github.com/DnOberon/heimdall/releases).

## How to use `heimdall`

The easiest way to get started with heimdall after installation is to ask for its help menu.

```console
> heimdall -h

Heimdall gives you a quick way to monitor, repeat, and selectively
log a CLI application. Quick configuration options allow
you to effectively test and monitor a CLI application in
development

Usage:
  heimdall [flags]
  
Available Commands:
  help        Help about any command
  init        Create a configuration for heimdall to replace command flag arguments
  run         Run heimdall using the "heimdall_config.json" file in the current directory 

Flags:
  -h, --help                help for heimdall
  -l, --log                 Toggle logging of provided program's stdout and stderr output to file, appends if file exists
      --logFilter string    Allows for log filtering via regex string. Use only valid with log flag
      --logName string      Specify the log file name, defaults to heimdall.log (default "heimdall.log")
      --logOverwrite        Toggle logging of provided program's stdout and stderr output to file
  -p, --parallelCount int   Designate how many instances of your should run in parallel at one time 
  -r, --repeat int          Designate how many times to repeat your program with supplied arguments (default 1)
  -t, --timeout duration    Designate when to kill your provided program
  -v, --verbose             Toggle display of provided program's stdout and stderr output while heimdall runs

```

Let’s run through a quick example based on the problem that started sthis whole thing - a console application managing a third-party, hidden application. 

I want heimdall to filter the logs that both my application and the hidden one outputs as well (here we filter for < and > characters as long as there is at least 1 preceding character) as killing my application if it hangs.

Telling heimdall to do that is easy -

`heimdall --timeout=30m --log --logFilter=<[^<>]+> exportApplication`

</br>

## Running `heimdall` with a configuration file

This tool provides the option of generating a json configuration file for ease of use. All command line flag arguments are available and represented inside the configuration file.

First, generate your configuration file using the `heimdall init` command and following the interactive prompts.

```console
> heimdall init
Executable path: tester
✔ Program arguments separated by comma: █

```

Once your configuration file is generated you can run heimdall with a single command -

`heimdall run`

Your configuration file, `heimdall_config.json` is available in the directory you first ran `heimdall init` in. To make changes you can edit this file directly or rerun `heimdall init`

</br>

## Can't you do this with a bash or powershell script?
You sure could - but you'd have to spend the time to build it, have separate scripts for at least windows and linux systems, or if you're using someone else's script, modify it to the point that it's going to work for your particular situation. You'll also have to handle multi-threading and logging yourself.

`heimdall` is meant to be extremely easy to use across platforms and applications by design. You shouldn't have to write some complicated scripts to achieve what `heimdall` can do with a single command.

</br>

## More information

I've written [an article](https://notyourlanguage.com/post/heimdall/) about the "why" of heimdall as well as stepping through source code from the earliest version.
