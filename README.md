```
       _ _        _
  __ _(_) |_     (_) __ _ _ __   ___  _ __ ___
 / _` | | __|____| |/ _` | '_ \ / _ \| '__/ _ \
| (_| | | ||_____| | (_| | | | | (_) | | |  __/
 \__, |_|\__|    |_|\__, |_| |_|\___/|_|  \___|
 |___/              |___/

```

`git-ignore` is a simple tool that creates `.gitignore` file for the
language of your project.

Usage
=====

Download the binary, and put it anywhere on your `$PATH`, then you can
invoke it with `git ignore <options>`

## List available `.gitignore` templates:

    git ignore --list

## Generate language-specific .gitignore file:

    git ignore -g <language>

e.g.,

    git ignore -g python

## Append to existing .gitignore file:

If you already have `.gitignore` file but want to be supplemented with
the popular gitignore template, you can use the `-a` flag to append the
ignore entries to your existing file:

    git ignore -g python -a

or:

    git ignore -g python --append

Development
===========

`git-ignore` is written in [Golang](https://golang.org)

Clone this repository and then run `make build`
