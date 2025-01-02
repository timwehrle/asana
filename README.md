# Asana CLI

A command-line interface to manage your Asana tasks and projects directly from your terminal.

<div>
    <a href="https://pkg.go.dev/github.com/timwehrle/asana">
        <img src="https://pkg.go.dev/badge/github.com/timwehrle/asana.svg" alt="Go Reference">
    </a>
    <a href="https://github.com/timwehrle/asana/blob/main/LICENSE">
        <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License">
    </a>
</div>

# Installation

## Pre-built binaries

Download the latest binary for your platform from the [releases page](https://github.com/timwehrle/asana/releases).

## From Source

```shell
go install github.com/timwehrle/asana
```

# Getting started

## Authentication

1. Get your Personal Access Token from Asana (Settings > Apps > Developer Apps)
2. Run the login command:
   ```shell
   asana auth login
   ```
3. Follow the prompts to paste your token and select your default workspace.

## Configuration

Set or get your default workspace:

```shell
asana config set default-workspace # or the shorter alias
asana config set dw

asana config get default-workspace # or the shorter alias
asana config get dw
```

## Basic Commands

View your tasks:

```shell
asana tasks list # List all your tasks
asana tasks view # Interactive task viewer with details
```

## Advanced Usage

```shell
asana help # Show all available commands
asana workspaces list # List all workspaces
```

# Contributing

Contributions are welcome! Here's how you can help:

- Report bugs and feature requests through issues
- Submit pull requests for bug fixes or new features
- Improve documentation
- Share feedback on usability

Please ensure your pull requests adhere to the following:

- Follow the existing code style
- Include tests for new functionality
- Update documentation as needed

# License

This project is licensed under the MIT License. See the [LICENSE][license] file for details.

[license]: ./LICENSE
