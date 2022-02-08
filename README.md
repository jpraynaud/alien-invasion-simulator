# Alien Invasion Simulator
[![Go](https://github.com/jpraynaud/alien-invasion-simulator/actions/workflows/go.yml/badge.svg)](https://github.com/jpraynaud/alien-invasion-simulator/actions/workflows/go.yml)
[![CodeQL](https://github.com/jpraynaud/alien-invasion-simulator/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/jpraynaud/alien-invasion-simulator/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jpraynaud/alien-invasion-simulator)](https://goreportcard.com/report/github.com/jpraynaud/alien-invasion-simulator)
[![GoDoc](https://godoc.org/github.com/jpraynaud/alien-invasion-simulator?status.svg)](https://godoc.org/github.com/jpraynaud/alien-invasion-simulator)

This project implements a simple **Alien Invasion** simulator CLI. 

---

* [Principle](#principle)
* [Assumptions](#assumptions)
* [Parameters](#parameters)
* [Install](#install)
* [Run](#run)
* [Examples](#examples)
* [Tests](#tests)
* [Help](#help)
* [Documentation](#documentation)

---

## Principle

The principle of the **simulation** is:
* a **world** describes a list of **cities** and their possible **links** to other **cities** 
* a **link** can be defined in any **direction** of this set: **{North, East, South, West}**
* some **aliens** are spawned randomly in the **world**
* the **aliens** move randomly from one **city** to another **city** using an existing **link**
* when two **aliens** meet in a **city** they fight so that:
    * the **city** gets destroyed (so do the links to this **city**)
    * the **aliens** are trapped (so that they are not able to move anymore)
* the **simulation** ends when any of the conditions below is met:
    * all the **cities** are destroyed
    * all the **aliens** are trapped
    * a maximum number of **steps** is reached
    
---

## Assumptions

The following assumptions have been made :
* the **city** names don't include any space (which should be replaced by any other character). For example, use ***New-York*** instead of ***New York***.
* **aliens** are spawned once at the beginning of the simulation
* the validity of the **links** is not checked (meaning that a **city** may be linked to the same city through several directions)

---

## Parameters

The following parameters are available :
* **aliens** (shorthanded to **n**) the number of aliens spawned at startup (defaults to **5**)
* **steps** (shorthanded to **s**) the number of maximum steps allowed (defaults to **10,000**)
* **file** (shorthanded to **m**) the path of the world map file (defaults to **map.txt**)

---

## Install

### From Source

**Step 1: Install Golang**

- Install a [correctly configured](https://golang.org/doc/install) Go toolchain (version 1.17+). 
- Make sure that your `GOPATH` and `GOBIN` environment variables are properly set up.

**Step 2: Get source code**

```bash
#Download sources from github
git clone https://github.com/jpraynaud/alien-invasion-simulator

# Go to sources directory
cd alien-invasion-simulator

# Checkout master branch
git checkout master
```

**Step 3 : Build binary**

```bash
# Build
go build -o bin/alien-invasion cmd/cli/main.go
```

**Step 4 : Verify**

```bash
# Verify
./bin/alien-invasion --help

# or Build and Run
go run cmd/cli/main.go --help
```

That should output something like:

```bash
An alien invasion simulator.
More informations available at: https://github.com/jpraynaud/alien-invasion-simulator

Usage:
  alien-invasion [flags]

Flags:
  -n, --aliens uint   total number of aliens (default 5)
  -m, --file string   world map file path (default "map.txt")
  -h, --help          help for alien-invasion
  -s, --steps uint    maximum number of steps (default 10000)
```

---

## Run
```bash
# Run
./bin/alien-invasion

# or
go run cmd/cli/main.go
```

That should output something like:

```bash
London has been destroyed by Alien #3 and Alien #1
Warsaw has been destroyed by Alien #4 and Alien #2

Roma north=Geneva west=Barcelona
Athens
Stockholm
Geneva
Paris north=Brussels south=Barcelona east=Berlin
Brussels
Berlin north=Stockholm
Barcelona north=Paris east=Roma
```

---

## Examples

- Set number of spawned **aliens**:
```bash
# Run
./bin/alien-invasion -n 100

# or
go run cmd/cli/main.go --aliens 100
```


- Set maximum number of steps:
```bash
# Run
./bin/alien-invasion -s 5

# or
go run cmd/cli/main.go --steps 5
```

- Set map file path:
```bash
# Run
./bin/alien-invasion -m ../maps-directory/other-map.txt

# or
go run cmd/cli/main.go --file ../maps-directory/other-map.txt
```

- Combine options:
```bash
# Run
./bin/alien-invasion -n 4 -s 10

# or
go run cmd/cli/main.go --aliens 4 --steps 10
```

---

## Tests

Run unit tests:
```sh
# Test with code coverage
go test -cover ./...
```

```sh
# Test with verbose output
go test -cover -v ./...
```

---

## Help

Get help:

```sh
# Help with executable
./bin/alien-invasion -h

# or
go run cmd/cli/main.go --help
```

## Documentation

[![GoDoc](https://godoc.org/github.com/jpraynaud/alien-invasion-simulator?status.svg)](https://godoc.org/github.com/jpraynaud/alien-invasion-simulator)
