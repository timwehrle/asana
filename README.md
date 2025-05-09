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
go install github.com/timwehrle/asana/cmd/asana@latest
```

## Bash Installation

This installation option is experimental, because the OS and architecture get detected automatically.

```shell
curl -sSL https://raw.githubusercontent.com/timwehrle/asana/main/scripts/install.sh | bash
```

## Homebrew Installation

```shell
brew tap timwehrle/asana
brew install --formula asana
```

## Having troubles with keyrings on WSL2?

If you're running into issues with keyring access on WSL2, there's a simple workaround!
You can find a detailed explanation here: [https://github.com/XeroAPI/xoauth/issues/25#issuecomment-2364599936](https://github.com/XeroAPI/xoauth/issues/25#issuecomment-2364599936)

To make development smoother, we've also provided a setup script.
It installs the necessary packages and configures the GNOME keyring automatically. You probably have to do this every time you start your WSL2 environment. 
```shell
chmod +x scripts/setup-wsl-keyring.sh
./scripts/setup-wsl-keyring.sh
```
After running the script, keyring functionality should be available in your WSL2 environment.

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

View tasks with filters:

```shell
asana tasks search --assignee me,12345678 # Search tasks by assignee and more filters
```

View the projects in your workspace:
```shell
asana projects list # List all the projects
asana projects list -l 25 --sort desc # List with options
```

View the teams in your workspace:

```shell
asana teams list # List all teams
```

View the users in your workspace:
```shell
asana users list # List all the users
asana users list -l 25 --sort desc # List with options
```

View tags of your workspace:

```shell
asana tags list # List all tags
asana tags list --favorite # List tags that you marked as favorite
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
