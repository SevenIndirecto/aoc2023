# Advent of Code 2023 - golang

## Initial setup

First create your .env file by copying .env.example

```bash
cp .env.example .env
```

1. Login on https://adventofcode.com/ and grab your session cookie ID.
2. Input into .env file

## Start new day

```bash
./next.sh
```

This prepares a new dir, copies templates and downloads your input.

To run tests in a module.

```bash
go test -run ''
```

To run a particular day, cd to dir and:

```bash
got run aoc.go
```

## Summary

Ran out of steam for the last 2 days, but pretty ok otherwise :)