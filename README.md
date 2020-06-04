# Scaff
Store the files you use to bootstrap projects(docs, design notes, READMEs, common starter dependencies etc) in one place and access them quickly when they are needed.

# Installation
### Using go get

`go get github.com/dmralev/scaff`

# Usage

```
Usage:
  scaff [command]

Available Commands:
  add         Store a file or directory under a single namespace
  get         Copy files from a given namespace to the current directory.
  help        Help about any command
  list        List stored namespaces.
  remove      Remove a filepath or namespace.
  show        See the insides of a given namespace in a tree format.

Flags:
  -h, --help   help for scaff
  ```
  
# How it works
One of the main reasons created this tool was to solve my need to easily reuse files when starting projects. Many times I would like to reuse documents, , and other commonly used things when starting a project. It's also forcing me to standardise how I approach starting a project projects, from the specs, to the architecture plans.

## Adding files
Add expects to receive a path to directory or file, and namespace under which to store the added files.

`scaff add [directory|filepath] [namespace] [flags]`

## List and Show

List is used to get an overview over the namespaces you have.

`scaff list`

With `show` you can see the files from a specific namespace in a `tree` like format.

`scaff show [namespace]`



## Removing files and namespaces
With remove you are deleting given directory or filepath from a namespace, or the namespace itself.

If the command receives one argument, it is assumed that it has received namespace to delete.

`scaff remove test_namespace`

Two provided arguments means that a [file or directory] should be looked up in the given [namespace] and removed from there.

`scaff remove README.md test_namespace`

# Caveats
## It's local
- One thing I really hate is that for now, the namespaces are stored in the home directory. This has many issues, but mostly I don't want to pollute other people's root folder, or mine for that matter, but right now I'm my only target audience with this, so it's good enough.

Although I'll definitely create a DropBox adapter in the future.