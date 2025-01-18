# Asana CLI

A command-line interface to manage your Asana tasks and projects directly from your terminal.

<div>
    <a href="https://pkg.go.dev/github.com/timwehrle/asana">
        <img src="https://pkg.go.dev/badge/github.com/timwehrle/asana.svg" alt="Go Reference">
    </a>
    <a href="https://github.com/timwehrle/asana/blob/main/LICENSE">
        <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License">
    </a>
   <a href="https://github.com/timwehrle/asana/actions/workflows/go.yml">
      <img src="https://github.com/timwehrle/asana/actions/workflows/go.yml/badge.svg" alt="Go Pipeline">
   </a>
   <a href="https://goreportcard.com/report/github.com/timwehrle/asana">
      <img src="https://goreportcard.com/badge/github.com/timwehrle/asana" alt="Go Report Card">
   </a>
</div>

# Installation

## Pre-built binaries

Download the latest binary for your platform from the [releases page](https://github.com/timwehrle/asana/releases).

## From Source

```shell
go install github.com/timwehrle/asana@latest
```

## Bash Installation

```shell
curl -sSL https://raw.githubusercontent.com/timwehrle/asana/main/scripts/install.sh | bash
```

## Homebrew Installation

```shell
brew tap timwehrle/asana
brew install --formula asana
```

# Getting started

## Authentication

1. Get your Personal Access Token from Asana (Settings > Apps > Developer Apps)
2. Run the login command:
   ```shell
   asana auth login
   ```
3. Follow the prompts to paste your token and select your default workspace.

To check the current status of your authentication and the Asana API:
```shell
asana auth status
```

## Configuration

Set or get your default workspace:

```shell
asana config set default-workspace
asana config set dw

asana config get default-workspace
asana config get dw
```

## Basic Commands

View your tasks:

```shell
asana tasks list # List all your tasks
asana tasks list --sort due-desc # Sort tasks by descending due date
asana tasks view # Interactive task viewer with details
asana tasks update # Interactive task updater
```

View the projects in your workspace:
```shell
asana projects list # List all the projects
asana projects list -l 25 --sort desc # List with options
```

View the users in your workspace:
```shell
asana users list # List all the users
asana users list -l 25 --sort desc # List with options
```

For more usage:
```shell
asana help # Show all available commands
```

# Contributing

If something feels off, you see an opportunity to improve performance, or think some
functionality is missing, we’d love to hear from you! Please review our [contributing docs][contributing] for
detailed instructions on how to provide feedback or submit a pull request. Thank you!

# License

This project is licensed under the MIT License. See the [LICENSE][license] file for details.

[contributing]: ./.github/CONTRIBUTING.md
[license]: ./LICENSE
