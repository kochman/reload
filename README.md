# reload

reload is a tool written in Go that watches a Go source code file for changes and automatically runs it.

## Installation

`go get github.com/kochman/reload`

reload should run anywhere Go runs, though it has only been tested on OS X.

## Usage

`reload path-to-source.go`

reload will pass through any arguments.

`reload path-to-source.go -exampleArg="example"`
