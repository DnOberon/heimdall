#### Installation

If you have the Go toolchain installed you can simply run
`go get github.com/dnoberon/heimdall` to install heimdall to $GOPATH/bin. If you want to run heimdall outside `$GOPATH/bin` make sure that `$GOPATt/bin` is [included in your \$PATH](https://golang.org/doc/code.html#GOPATH).
If not, you can download any of heimdall’s releases for your platform [here](https://github.com/DnOberon/heimdall/releases). Future plans for heimdall include creating homebrew recipes as well as a few other package managers, but for now you’ll have to either build from source, use Go’s `get` command, or download the application manually.

The easiest way to get started with heimdall is to ask for its help menu.

```
> heimdall -h

Heimdall gives you a quick way to monitor, repeat, and selectively
log a CLI application. Quick configuration options allow
you to effectively test and monitor a CLI application in
development

Usage:
  heimdall [flags]

Flags:
  -h, --help               help for heimdall
  -l, --log                Toggle logging of provided program's stdout and stderr output to file
  -f, --logFilter string   Allows for log filtering via regex string. Use only valid with log flag
  -r, --repeat int         Designate how many times to repeat your program with supplied arguments
  -t, --timeout duration   Designate when to kill your provided program
  -v, --verbose            Toggle display of provided program's stdout and stderr output while heimdall runs

```

Let’s run through a quick example based on the problem that started this whole thing - a console application managing a third-party, hidden application. I want heimdall to filter the logs that both my application and the hidden one outputs as well as killing my application if it hangs.

Telling heimdall to do that is easy -
`heimdall --timeout=30m --log --logFilter=<[^<>]+> exportApplication`

</br>
