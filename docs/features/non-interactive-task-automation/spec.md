# Non-Interactive Task Automation Specification

## Summary

This feature adds a fully non-interactive task automation path to the Asana CLI for agentic and script-driven workflows. Today, `asana tasks search` is human-readable only, while `asana tasks view` and `asana tasks update` require interactive task selection and prompt-based input. That blocks a CLI-only workflow where an agent searches for a task, reads its description, checks for a branch header, prepends `Branch: <name>`, and continues execution without any human at the terminal.

The feature introduces stable task targeting by GID, machine-readable output, explicit non-interactive read and write paths, explicit incomplete filtering and pagination for task search, and CI-safe prompt guardrails. Interactive behavior remains the default for human users running commands in a TTY unless explicit automation flags are provided.

## Problem Statement

The CLI cannot currently support a fully non-interactive development workflow for Asana tasks because the following capabilities are missing or incomplete:

- task selection by stable ID instead of prompt selection
- non-interactive access to full task payloads, especially notes/description
- non-interactive description updates, including prepend semantics
- machine-readable output from search and view operations
- explicit automation-safe search filters and pagination controls
- fail-fast behavior when a prompt-based command is invoked in CI or another non-TTY context

As a result, automation must fall back to direct API access, MCP, or manual intervention even though the CLI already contains most of the underlying Asana transport capabilities.

## Target Users

- AI agents operating in Codex, Cursor, CI runners, shell scripts, or other non-interactive environments
- developers who want deterministic, parseable CLI behavior for automation
- interactive CLI users who should retain the current prompt-based experience by default

## Success Criteria

- An automation can search for open tasks and identify a task by GID without parsing numbered prompt-oriented output.
- An automation can fetch a specific task’s `name`, `notes`, `completed` state, and `permalink_url` using only the CLI.
- An automation can prepend `Branch: <branch-name>` to a task’s notes while preserving the existing description body.
- Commands that would otherwise prompt fail fast with a clear error when invoked in a non-TTY context without sufficient non-interactive flags.
- Existing interactive task flows continue to work for users running the CLI in a terminal.

## Requirements

| ID | Requirement | Priority | Rationale |
| --- | --- | --- | --- |
| R1 | The CLI MUST allow `tasks view` and `tasks update` to target a task by stable Asana GID. | MUST | Automation cannot rely on prompt ordering or numbered lists. |
| R2 | The CLI MUST provide a non-interactive task read path that returns at least `gid`, `name`, `notes`, `completed`, and `permalink_url`. | MUST | The branch guard workflow depends on description parsing and task identity. |
| R3 | The CLI MUST provide a non-interactive task description write path. | MUST | Automation must be able to set or modify notes without opening an editor. |
| R4 | The CLI MUST support prepend semantics for notes so branch metadata can be inserted above existing descriptions. | MUST | The target workflow requires preserving description content while adding a header. |
| R5 | `tasks search` MUST provide machine-readable output. | MUST | Search results must be parseable without scraping human-oriented text. |
| R6 | `tasks view` SHOULD provide machine-readable output. | SHOULD | Symmetry with search improves automation and reduces ad hoc parsing. |
| R7 | `tasks update` SHOULD provide deterministic machine-readable success output in automation mode. | SHOULD | Scripts benefit from stable update confirmations, though this is secondary to the write itself. |
| R8 | `tasks search` MUST expose explicit incomplete-only filtering. | MUST | Agent queues generally need open tasks only; implicit API defaults are not sufficient. |
| R9 | `tasks search` MUST expose result limiting and pagination controls. | MUST | Large workspaces need bounded, scriptable result retrieval. |
| R10 | Commands MUST fail fast instead of entering prompt flows when stdin is not a TTY and non-interactive inputs are missing. | MUST | CI and agent environments must not hang. |
| R11 | Existing interactive behavior MUST remain the default when commands are run in a TTY without automation flags. | MUST | This feature should extend the CLI, not replace the current UX. |
| R12 | Help text and README examples SHOULD document both interactive and non-interactive usage. | SHOULD | New behavior is only useful if discoverable. |
| R13 | The implementation COULD standardize JSON output schemas across all task commands. | COULD | Consistency improves long-term maintainability but is not required for the first release. |
| R14 | The implementation WON'T add HTML notes editing or generalized bulk task update operations in this feature. | WON'T | The feature is scoped to plain-text notes and single-task automation. |

## Assumptions and Constraints

### Assumptions

- The existing Asana client implementation remains the transport layer; no new HTTP stack is required.
- The CLI continues to use Cobra for command parsing and `survey` for interactive prompts.
- The current `task.Fetch(client)` pattern is sufficient for GID-based resolution once the task ID is known.
- The immediate automation need is single-task operation, not bulk mutation.

### Constraints

- The feature must fit the current command layout under `pkg/cmd/tasks/*`.
- Changes must preserve the existing default user-facing behavior when no automation flags are passed.
- stdout must remain clean in JSON mode; human-readable log lines cannot be mixed into machine-readable output.
- The CLI must operate correctly both in TTY and non-TTY environments.

## Scope

### In Scope

- GID-based task selection for `tasks view` and `tasks update`
- machine-readable output for `tasks search`
- machine-readable output for `tasks view`
- explicit non-interactive update flags for task name, notes, completion state, and due date
- notes prepend helpers for the branch-header workflow
- explicit search filters for incomplete-only behavior
- search limit and pagination controls
- fail-fast prompt guardrails in non-TTY environments
- documentation and tests for the above

### Out of Scope

- HTML notes editing
- bulk task updates
- new Asana API endpoints beyond existing task fetch, search, and update capabilities
- changes to unrelated CLI command groups
- deduplication or normalization of multiple existing `Branch:` headers unless explicitly added later
- cross-command global output-format infrastructure if a smaller task-command-local implementation is sufficient

## User Stories

### Story 1: Search open tasks non-interactively

AS AN AI agent
I WANT TO search for incomplete tasks and receive parseable results
SO THAT I can select the correct task by stable ID without a prompt

ACCEPTANCE CRITERIA:
- GIVEN a workspace containing incomplete and completed tasks
  WHEN the user runs `asana tasks search --query "Agents" --incomplete --output json`
  THEN stdout contains valid JSON describing only incomplete matching tasks
- GIVEN `--limit 10`
  WHEN the command runs successfully
  THEN no more than 10 results are returned
- GIVEN a pagination cursor flag such as `--page-offset <cursor>`
  WHEN the command runs successfully
  THEN the request uses that cursor and returns the corresponding page of results
- GIVEN no matching tasks
  WHEN the command runs in JSON mode
  THEN stdout contains a valid empty JSON result rather than human-readable prose

EDGE CASES:
- Empty query with filters only: search still returns valid filtered results.
- Invalid cursor/offset: the command exits non-zero with a clear error.
- Completed filter omitted: default behavior remains documented and stable.

TECHNICAL NOTES:
- Reuse `asana.SearchTasksQuery.Completed` and `asana.Options{Limit, Offset}`.
- Expand requested fields beyond `name` and `due_on` to include stable IDs and automation fields.

### Story 2: View a task by GID without prompts

AS AN automation script
I WANT TO fetch a task directly by GID
SO THAT I can read its description and metadata in a deterministic way

ACCEPTANCE CRITERIA:
- GIVEN a valid task GID
  WHEN the user runs `asana tasks view --task <gid> --output json`
  THEN stdout contains valid JSON with `gid`, `name`, `notes`, `completed`, and `permalink_url`
- GIVEN a valid task GID and no `--output json`
  WHEN the command runs in a TTY
  THEN it prints the existing human-readable task details without prompting for selection
- GIVEN an invalid or inaccessible task GID
  WHEN the command runs
  THEN it exits non-zero with an actionable error
- GIVEN no task GID in a non-TTY environment
  WHEN the command would otherwise need a prompt
  THEN it fails fast instead of hanging

EDGE CASES:
- Task with empty notes: JSON returns an empty string or omitted field consistently per schema.
- Task not found: surface the Asana error clearly.
- Task found but fields missing due to narrow `opt_fields`: command treats this as an implementation defect and tests must cover it.

TECHNICAL NOTES:
- A shared resolver should choose between GID fetch and the current interactive selection path.

### Story 3: Update a task non-interactively

AS AN AI agent
I WANT TO update a task using flags instead of prompts
SO THAT the CLI can run unattended in CI and editor-integrated agents

ACCEPTANCE CRITERIA:
- GIVEN `asana tasks update --task <gid> --name "New name"`
  WHEN the command runs successfully
  THEN the task name is updated without any prompt
- GIVEN `asana tasks update --task <gid> --notes "New notes"`
  WHEN the command runs successfully
  THEN the task notes are replaced without opening an editor
- GIVEN `asana tasks update --task <gid> --complete`
  WHEN the command runs successfully
  THEN the task becomes completed without task or action selection
- GIVEN `asana tasks update --task <gid> --due-on 2026-03-20`
  WHEN the command runs successfully
  THEN the task due date is updated using the provided date

EDGE CASES:
- Multiple conflicting update flags: return a validation error if combinations are unsupported.
- Invalid due date format: exit non-zero with a date-format error.
- No-op update where the new value matches the existing value: return a deterministic success or explicit no-op response, but do not prompt.

TECHNICAL NOTES:
- Interactive action menus remain as the fallback mode for terminal users who invoke `asana tasks update` without automation flags.

### Story 4: Prepend branch metadata to notes

AS AN automation script
I WANT TO prepend a branch header to existing task notes
SO THAT I can mark the active development branch without destroying the task description

ACCEPTANCE CRITERIA:
- GIVEN existing notes `"Implement parser"`
  WHEN the user runs `asana tasks update --task <gid> --prepend-notes "Branch: codex/foo\n\n"`
  THEN the resulting notes are `"Branch: codex/foo\n\nImplement parser"`
- GIVEN empty existing notes
  WHEN the same prepend command runs
  THEN the resulting notes contain only the prepended content without malformed separators
- GIVEN multiline existing notes
  WHEN prepend runs
  THEN all original content is preserved byte-for-byte after the inserted prefix
- GIVEN a file-backed input option such as `--prepend-notes-file`
  WHEN the command runs successfully
  THEN the file contents are prepended using the same semantics as inline input

EDGE CASES:
- Prefix text without trailing newline: resulting notes remain predictable and documented.
- Large note bodies: prepend logic still preserves full content.
- Existing `Branch:` line already present: current release does not deduplicate unless explicitly specified later.

TECHNICAL NOTES:
- Implement notes mutation in a dedicated helper that can be unit tested independently from Asana network calls.

### Story 5: Fail fast in CI and non-TTY environments

AS A developer running automation in CI
I WANT prompt-based commands to fail fast when required non-interactive inputs are missing
SO THAT jobs do not hang indefinitely

ACCEPTANCE CRITERIA:
- GIVEN stdin is not a TTY
  WHEN the user runs `asana tasks view` with no `--task`
  THEN the command exits non-zero with a message explaining that `--task` is required in non-interactive mode
- GIVEN stdin is not a TTY
  WHEN the user runs `asana tasks update` with neither `--task` nor a non-interactive action flag
  THEN the command exits non-zero with guidance on the required flags
- GIVEN sufficient non-interactive inputs
  WHEN the same commands run in CI
  THEN they succeed without trying to invoke `survey`

EDGE CASES:
- stdout redirected but stdin is still a TTY: behavior follows stdin TTY status, not stdout.
- shell pipes and subshells: commands remain deterministic and do not attempt hidden prompts.

TECHNICAL NOTES:
- Use `IOStreams.IsStdinTTY` as the source of truth for whether prompting is allowed.

## Functional Specification

## Command Contract

### `asana tasks search`

#### New/Changed Flags

| Flag | Type | Required | Purpose |
| --- | --- | --- | --- |
| `--output` | string | No | Select output mode. Initial supported values: `text`, `json`. Default: `text`. |
| `--incomplete` | bool | No | Limit results to incomplete tasks only. |
| `--completed` | bool | No, optional alternative to `--incomplete` | Optional explicit completed-state filter if implemented instead of or in addition to `--incomplete`. |
| `--limit`, `-l` | int | No | Maximum number of results to return. |
| `--page-offset` | string | No | Cursor/offset for fetching a specific page of results. |

#### Behavior

- In text mode, preserve the current human-readable list behavior.
- In JSON mode, emit only valid JSON on stdout.
- When `--incomplete` is provided, set the search query to return only incomplete tasks.
- When pagination is provided, pass the cursor/offset to the Asana request.
- When `--limit` is provided, limit output deterministically.

#### JSON Response Shape

```json
{
  "tasks": [
    {
      "gid": "12001234",
      "name": "Example task",
      "completed": false,
      "due_on": "2026-03-20",
      "permalink_url": "https://app.asana.com/0/..."
    }
  ],
  "next_page_offset": "opaque-cursor-or-empty"
}
```

### `asana tasks view`

#### New/Changed Flags

| Flag | Type | Required | Purpose |
| --- | --- | --- | --- |
| `--task`, `--task-id` | string | Required in non-interactive mode | Stable task GID to fetch directly. |
| `--output` | string | No | Select output mode: `text` or `json`. Default: `text`. |

#### Behavior

- If `--task` is provided, fetch that task directly and skip prompt selection.
- If `--task` is omitted and stdin is a TTY, preserve current interactive selection.
- If `--task` is omitted and stdin is not a TTY, exit with an actionable error.
- In JSON mode, emit the full automation-relevant task payload.

#### JSON Response Shape

```json
{
  "task": {
    "gid": "12001234",
    "name": "Example task",
    "notes": "Branch: codex/foo\n\nImplement parser",
    "completed": false,
    "due_on": "2026-03-20",
    "permalink_url": "https://app.asana.com/0/..."
  }
}
```

### `asana tasks update`

#### New/Changed Flags

| Flag | Type | Required | Purpose |
| --- | --- | --- | --- |
| `--task`, `--task-id` | string | Required in non-interactive mode | Stable task GID to update directly. |
| `--name` | string | No | Update the task name. |
| `--notes` | string | No | Replace task notes with the provided value. |
| `--notes-file` | string | No | Replace task notes with file contents. |
| `--prepend-notes` | string | No | Prepend text to the current notes. |
| `--prepend-notes-file` | string | No | Prepend file contents to the current notes. |
| `--complete` | bool | No | Mark the task complete. |
| `--due-on` | string | No | Set the task due date in `YYYY-MM-DD`. |
| `--output` | string | No | Optional deterministic output mode. Recommended values: `text`, `json`. |

#### Behavior

- If `--task` and one or more non-interactive update flags are provided, skip all prompts and apply the update directly.
- If no non-interactive update flags are provided and stdin is a TTY, preserve the existing interactive task and action selection flow.
- If no non-interactive update flags are provided and stdin is not a TTY, exit with an actionable error.
- If both replace and prepend notes flags are provided together, validation must either reject the combination or define a deterministic precedence. The recommended behavior is to reject the combination.
- If file-backed and inline values are provided for the same update operation, validation must reject the combination.

#### JSON Response Shape

Recommended shape if JSON update output is implemented:

```json
{
  "task": {
    "gid": "12001234",
    "name": "Example task",
    "notes": "Branch: codex/foo\n\nImplement parser",
    "completed": false,
    "due_on": "2026-03-20",
    "permalink_url": "https://app.asana.com/0/..."
  },
  "updated_fields": [
    "notes"
  ]
}
```

If JSON update output is deferred, text-mode success output must still be deterministic and concise.

## Data Models

### Internal Command Options

The task command option structs should be extended to represent:

- task identifier
- output mode
- non-interactive update inputs
- pagination cursor
- explicit incomplete/completed filter state

The representation should distinguish between:

- flag omitted
- flag explicitly provided
- conflicting flags

This matters because search filtering and update validation depend on whether a value was intentionally set.

### Output Models

Recommended internal response types:

- `TaskListOutput`
- `TaskOutput`
- `TaskUpdateOutput`

Common task fields:

- `gid`
- `name`
- `notes`
- `completed`
- `due_on`
- `permalink_url`

Optional metadata:

- `next_page_offset`
- `updated_fields`

### Notes Mutation Model

Notes update operations should be represented as mutually exclusive intents:

- replace notes from inline text
- replace notes from file
- prepend notes from inline text
- prepend notes from file

This avoids hidden precedence rules and makes validation straightforward.

## State Changes

### Search

- No server-side state changes.
- Only query and output behavior changes.

### View

- No server-side state changes.
- Command may shift from prompt-driven task selection to direct fetch by ID.

### Update

- Server-side task fields may change:
  - `name`
  - `notes`
  - `completed`
  - `due_on`
- In prepend mode, the command performs a read-modify-write operation on `notes`.

## Dependencies

### Existing Code Dependencies

- `pkg/cmd/tasks/search/search.go`
- `pkg/cmd/tasks/view/view.go`
- `pkg/cmd/tasks/update/update.go`
- `pkg/cmdutils/select_task.go` or a new shared utility package
- `pkg/iostreams/iostreams.go`
- `internal/api/asana/tasks.go`
- `internal/prompter/prompter.go`

### External Dependencies

- Cobra for command parsing
- Survey for interactive prompts
- Asana REST API via the existing client wrapper

## Security and Reliability Considerations

### Input Validation

- Validate GID input is non-empty before task fetch.
- Validate due dates use `YYYY-MM-DD`.
- Validate output mode against an allowed enum.
- Validate mutually exclusive update flags.
- Validate file paths and surface read errors clearly.

### Output Safety

- In JSON mode, stdout must contain only JSON.
- Human-readable warnings and errors must go to stderr.
- Success and error semantics must be stable enough for automation to rely on exit codes.

### Failure Modes

- network or Asana API failures must propagate as non-zero exits
- invalid task ID must not fall back to interactive selection
- non-TTY prompt attempts must fail fast
- prepend operations must not silently drop or reorder existing notes content

## Edge Cases

### Search Edge Cases

- no results in text mode
- no results in JSON mode
- invalid limit values
- invalid page offsets
- mixed completed and incomplete result sets
- large result sets requiring pagination

### View Edge Cases

- task GID not found
- task belongs to inaccessible workspace
- empty notes
- null due date
- non-TTY invocation without a task ID

### Update Edge Cases

- invalid due date format
- unsupported combination of `--notes` and `--prepend-notes`
- unsupported combination of inline and file-backed variants
- empty inline notes input
- empty prepend content
- large existing note bodies
- no-op updates
- Asana update failure after local mutation succeeds

## Technical Specification

## Overview

The implementation should add a dual-mode execution model to task commands:

- interactive mode for existing terminal users
- non-interactive mode for scripts and agents

Mode selection should be explicit where possible:

- presence of `--task` and automation flags forces non-interactive execution
- absence of those flags falls back to interactive behavior only when stdin is a TTY

## Proposed Architecture

### Shared Task Resolver

Introduce a reusable helper that:

1. checks whether `--task` was provided
2. if yes, constructs a task with that GID and fetches details directly
3. if no, either:
   - uses the existing interactive selection path when prompting is allowed
   - returns an interactive-required error when prompting is not allowed

### Shared Non-Interactive Guard

Introduce a helper such as:

- `RequireInteractive(io, contextMessage)`
- or `ValidateInteractiveAllowed(io, fallbackHint)`

This helper should be used before any `survey` prompt path is entered.

### Output Mode Handling

Add a small output mode abstraction for task commands:

- parse output mode flags
- encode JSON payloads
- keep text rendering local to existing command implementations

Global CLI output standardization is not required for the first iteration, but task commands should share schema conventions.

### Notes Mutation Helper

Create a dedicated helper for notes transformations:

- `ReplaceNotes(existing, replacement) string`
- `PrependNotes(existing, prefix) string`

The prepend helper must preserve the existing notes exactly after the inserted prefix. Separator behavior should be deterministic and documented.

## API Usage

### Existing Endpoints to Reuse

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/workspaces/{workspace_gid}/tasks/search` | Search tasks with filters |
| `GET` | `/tasks/{task_gid}` | Fetch a task by GID |
| `PUT` | `/tasks/{task_gid}` | Update task fields |

### Required Fields

Commands must request sufficient `opt_fields` to fulfill their contracts.

Search minimum:

- `gid`
- `name`
- `completed`
- `due_on`
- `permalink_url`

View/update minimum:

- `gid`
- `name`
- `notes`
- `completed`
- `due_on`
- `permalink_url`

## Validation Rules

### Search Validation

- `--limit` must be zero or positive.
- `--output` must be one of the supported values.
- If both `--incomplete` and an explicit conflicting completed-state flag are supported, reject conflicting combinations.

### View Validation

- `--output` must be one of the supported values.
- In non-TTY mode, missing `--task` is an error.

### Update Validation

- At least one update operation must be specified in non-interactive mode.
- `--notes` and `--notes-file` are mutually exclusive.
- `--prepend-notes` and `--prepend-notes-file` are mutually exclusive.
- replace-notes and prepend-notes operations are mutually exclusive unless a deterministic composition rule is explicitly documented.
- `--due-on` must parse as `YYYY-MM-DD`.
- In non-TTY mode, missing `--task` is an error when a prompt would otherwise be required.

## Testing Strategy

### Unit Tests

- output mode parsing
- validation of conflicting flags
- due date parsing
- notes replace behavior
- notes prepend behavior with:
  - empty existing notes
  - single-line notes
  - multiline notes
  - empty prepend content
- search query construction for incomplete filters, limits, and offsets

### Integration Tests

- `tasks search --output json`
- `tasks search --incomplete --limit 5`
- `tasks view --task <gid> --output json`
- `tasks update --task <gid> --notes "..."`
- `tasks update --task <gid> --prepend-notes "..."`
- file-backed notes update and prepend
- non-TTY failure behavior for interactive-only invocation paths

### Manual / End-to-End Validation

1. Search for a task queue using JSON output.
2. Extract a GID from the results.
3. Fetch the target task with `tasks view --task <gid> --output json`.
4. Parse the first line of `notes` to simulate branch guard logic.
5. Prepend `Branch: <branch-name>` with `tasks update`.
6. Re-fetch the task and verify the inserted prefix and preserved body.
7. Run prompt-based commands in a non-TTY shell and confirm they fail fast.

## Acceptance Criteria

### Feature-Level Acceptance

- `tasks search` can be used as the first step of an unattended workflow.
- `tasks view` can return description data without interactive selection.
- `tasks update` can update notes without interactive prompts.
- prepend semantics support the branch-header use case.
- the feature behaves correctly in both TTY and non-TTY environments.

### Release Gates

- automated tests cover the new non-interactive paths
- interactive regression behavior is verified
- help text is updated
- README documents at least one end-to-end CLI-only workflow

## Implementation Notes

- Prefer a single canonical `--task` flag with `--task-id` only as an alias if needed.
- Prefer `--output json` over bespoke flags like `--json` to leave room for future formats.
- Prefer explicit validation errors over hidden precedence when multiple note-update flags are supplied.
- Prefer command-local output adapters if a global output system would delay delivery.

## Open Questions

- Should `tasks update --output json` return the full updated task or only a minimal confirmation object?
- Should `tasks search` support both `--incomplete` and `--completed`, or only one canonical completion filter?
- What should the exact pagination flag name be: `--page-offset`, `--offset`, or a cursor-specific label?
- Should prepend behavior automatically normalize missing trailing newlines, or should it preserve input exactly as provided?
- Should the CLI eventually deduplicate existing `Branch:` headers, or is simple prepend-only behavior sufficient for v1?

## Recommended Next Step

Use this spec as the implementation contract and keep [breakdown.md](/Users/gosev/REPOS/ET/asana-cli/docs/features/non-interactive-task-automation/breakdown.md) as the execution plan. The spec defines what the feature must do; the breakdown defines how to sequence the work.
