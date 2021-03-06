# Stufy [![Build Status](https://travis-ci.com/ArthurHlt/stufy.svg?branch=master)](https://travis-ci.com/ArthurHlt/stufy)

Stufy is a standalone cli for managing [statusfy](https://statusfy.co) deployment.

You will able to create incident events and scheduleds tasks and managed them without having nodejs stack installed.

[statusfy](https://statusfy.co) has the good idea to let you deploy status page as a static webpage by simply pushing 
commits on a git repository and be available as gitlab/github pages. With this in mind, with this cli you could directly 
create/update/delete incidents and scheduleds tasks by simply set a git repo as a target. You don't even need git installed 
on your machine.

**Important**: For now stufy only support config file in yaml.

## install

### On *nix system

You can install this via the command-line with either `curl` or `wget`.

#### via curl

```bash
$ bash -c "$(curl -fsSL https://raw.github.com/ArthurHlt/stufy/master/bin/install.sh)"
```

#### via wget

```bash
$ bash -c "$(wget https://raw.github.com/ArthurHlt/stufy/master/bin/install.sh -O -)"
```

### On windows

You can install it by downloading the `.exe` corresponding to your cpu from releases page: https://github.com/ArthurHlt/stufy/releases .
Alternatively, if you have terminal interpreting shell you can also use command line script above, it will download file in your current working dir.

### From go command line

Simply run in terminal:

```bash
$ go get github.com/ArthurHlt/stufy
```

## Usage 

```
Usage:
  cli [OPTIONS] <command>

Application Options:
  -t, --target=  Set a target, this can be a directory path or a git repo (e.g.: git@github.com:ArthurHlt/stufy-test.git or
                 https://user:password@github.com/ArthurHlt/stufy-test.git)
  -v, --version  Show version

Help Options:
  -h, --help     Show this help message

Available commands:
  add-alias         Add an alias to your current target to use instead of plain target (aliases: a)
  delete-incident   Delete an existing incident (aliases: d)
  delete-scheduled  Delete an existing scheduled task (aliases: ds)
  finish-scheduled  Finish a scheduled task (aliases: fs)
  list-incidents    List incidents (aliases: li)
  list-scheduleds   List scheduleds (aliases: ls)
  new-incident      Create a new incident (aliases: n)
  new-scheduled     Create a new scheduled task (aliases: ns)
  remove-alias      Remove an alias (aliases: ra)
  resync            Resynchronize your target (useful when merging issue on git repo) (aliases: r)
  update-incident   Update an existing incident (aliases: u)
  update-scheduled  Update an existing scheduled task (aliases: us)
```

**Note**: When you choose to open editor you can also set default editor with env var `EDITOR`.

## Auto completion

Here is an example how to activate bash auto completion for stufy:

```bash
_completion_stufy() {
    # All arguments except the first one
    args=("${COMP_WORDS[@]:1:$COMP_CWORD}")

    # Only split on newlines
    local IFS=$'\n'

    # Call completion (note that the first element of COMP_WORDS is
    # the executable itself)
    COMPREPLY=($(GO_FLAGS_COMPLETION=1 ${COMP_WORDS[0]} "${args[@]}"))
    return 0
}

complete -F _completion_stufy stufy
```