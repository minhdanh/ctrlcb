[![codecov](https://codecov.io/gh/minhdanh/ctrlcb/branch/master/graph/badge.svg?token=F41311IRII)](https://codecov.io/gh/minhdanh/ctrlcb) [![wercker status](https://app.wercker.com/status/5d32f8d4752bdfd67baaa44a3b70c108/s/master "wercker status")](https://app.wercker.com/project/byKey/5d32f8d4752bdfd67baaa44a3b70c108)

# ctrlcb

A tool to help copy files or directories without having to type the full paths of the source and destination, by utilizing the clipboard to store the source paths. Think of `Ctrl + c` and `Ctrl + v` to copy files, but using command line. Yes we have `cp` command to do this but in some cases typing the long source and destination paths is tedious.

## Install

There're two commands to install: `ctrlcb-copy` to copy paths to clipboard, and `ctrlcb-paste` to actually copy the files (or directories) to your current working directory:

```
go get github.com/minhdanh/ctrlcb/cmd/ctrlcb-paste
go get github.com/minhdanh/ctrlcb/cmd/ctrlcb-copy
```

## Usage

Let's say you have opened two terminals with two different working directories and you want to copy a directory (for example: `tests`) from one to another.

On the terminal that have the directory you want to copy:
```
ctrlcb-copy tests
```

Then on the other terminal:
```
ctrlcb-paste
```

## Options

When pasting with `ctrlcb-paste`, there're some options available:
- `-k`: keep the paths of the source. For example if you run `ctrlcb-copy a/b/c/file.txt`:
    - without this flag, `file.txt` will be copied to current working directory
    - with this flag, `file.txt` will be copied to `a/b/c/file.txt` in current working directory
- `-f`:  if same file or directory already exist on the destination, use this flag to overwrite them