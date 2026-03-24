# Feature Breakdown: Non-Interactive Task Automation

## 1. Executive Summary

### Feature Overview

Add a fully non-interactive task automation path to the CLI so an agent or script can complete the workflow `search -> read description -> branch guard -> prepend Branch: ... -> push -> implement` without entering a TUI picker or editor. The feature should preserve the current interactive user experience while adding explicit CLI flags for stable task targeting, machine-readable output, non-interactive reads and writes, and CI-safe failure behavior.

### Key Objectives

- Allow `tasks search`, `tasks view`, and `tasks update` to operate on a task by stable Asana GID.
- Expose full task payload fields needed for automation, especially `gid`, `name`, `notes`, `completed`, and `permalink_url`.
- Allow task description updates without opening the interactive editor.
- Provide JSON output for script parsing.
- Add explicit search controls for incomplete-only results, result limits, and pagination.
- Prevent automation runs from hanging on `survey` prompts when stdin is not a TTY or when a non-interactive mode is requested.

### Current-State Observations

- `pkg/cmd/tasks/view/view.go` is fully interactive and selects a task from `QueryTasks(...)` before calling `task.Fetch(...)`.
- `pkg/cmd/tasks/update/update.go` is fully interactive and depends on `survey` selection plus editor/input prompts for all updates.
- `pkg/cmd/tasks/search/search.go` already uses the richer workspace search API, but only prints human-readable numbered output and only requests `name` and `due_on`.
- `internal/api/asana/tasks.go` already contains most of the data model and transport support needed for this feature: `Task.ID`, `Task.Notes`, `Task.Completed`, `Task.PermalinkURL`, `UpdateTaskRequest`, `SearchTasksQuery.Completed`, and `Options{Fields, Limit, Offset}`.
- `pkg/iostreams/iostreams.go` already exposes TTY state, which can be reused for CI-safe behavior.

### Expected Impact

- Agents can use the CLI end-to-end for the branch workflow instead of falling back to REST or MCP for read/update steps.
- Scripts become deterministic because task selection is based on GID rather than a numbered prompt list.
- Human users keep the current interactive flows unless they opt into non-interactive flags.

### Assumptions

- Existing interactive behavior remains the default when no explicit task ID or non-interactive action flag is provided.
- JSON support is only required for commands in the task automation path (`tasks search`, `tasks view`, and optionally `tasks update` responses).
- Description updates only need plain-text notes support; HTML notes editing is out of scope.

### Success Criteria

- A script can run `asana tasks search --query "Agents" --incomplete --output json --limit 10` and receive parseable task objects containing stable IDs.
- A script can run `asana tasks view --task 12001234 --output json` and receive the task notes, completion state, and permalink without any prompt.
- A script can run `asana tasks update --task 12001234 --prepend-notes "Branch: codex/foo\n\n"` and update the task description without opening an editor.
- If a command would otherwise prompt but stdin is not a TTY, the CLI exits with a clear actionable error instead of hanging.

## 2. Component Architecture

### Text Diagram

```text
User / Agent / CI
        |
        v
  Cobra task commands
  - pkg/cmd/tasks/search
  - pkg/cmd/tasks/view
  - pkg/cmd/tasks/update
        |
        v
  Shared task automation layer
  - task ID resolution
  - non-interactive guardrails
  - notes mutation helpers
  - output formatting / JSON encoding
        |
        v
  Asana API client
  - QueryTasks
  - Workspace.SearchTasks
  - Task.Fetch
  - Task.Update
        |
        v
  Asana REST API
```

### Components

#### A. Task Command Surface

Primary files:

- `pkg/cmd/tasks/search/search.go`
- `pkg/cmd/tasks/view/view.go`
- `pkg/cmd/tasks/update/update.go`

Responsibilities:

- Parse new flags such as `--task`, `--output`, `--incomplete`, `--limit`, `--page-offset`, and non-interactive update actions.
- Choose between interactive and non-interactive execution paths.
- Return stable, consistent stdout/stderr behavior.

#### B. Shared Automation Utilities

Likely location:

- Extend `pkg/cmdutils`, or create a focused shared package under `pkg/cmd/tasks/shared`.

Responsibilities:

- Resolve a task either by explicit GID or by existing interactive selection.
- Validate incompatible flag combinations.
- Centralize "interactive required" errors.
- Implement note mutation helpers such as prepend/replace semantics.

#### C. Output Formatting Layer

Likely location:

- Existing `pkg/format` package, or a small new helper adjacent to task commands.

Responsibilities:

- Define JSON response shapes for search and view.
- Keep human-readable output unchanged unless `--output json` is requested.
- Ensure update commands can emit concise machine-readable confirmation if needed.

#### D. Asana API Access

Primary file:

- `internal/api/asana/tasks.go`

Responsibilities:

- Reuse existing search and task fetch/update transport.
- Request the correct `opt_fields`.
- Possibly add a direct fetch helper if task-by-ID resolution becomes repetitive.

### Integration Points

- `tasks search` must request more fields and expose pagination-related inputs already supported by `asana.Options`.
- `tasks view` must bypass `survey` selection when `--task <gid>` is present.
- `tasks update` must bypass both task picker and action picker when task ID plus action flags are present.
- TTY detection from `pkg/iostreams` must gate prompt-based flows.

### External Dependencies

- Asana task search endpoint: `/workspaces/{workspace_gid}/tasks/search`
- Asana task fetch endpoint: `/tasks/{task_gid}`
- Asana task update endpoint: `/tasks/{task_gid}`
- `survey` prompt library for the existing interactive path

## 3. Implementation Tasks

## Task 1: Define the non-interactive CLI contract

**Component**: Task command surface

**Description**: Define and document the exact flag contract for task targeting, output mode, incomplete filtering, pagination, and non-interactive update actions so implementation can stay coherent across commands.

**Dependencies**: None

**Acceptance Criteria**:
- [ ] `tasks view` and `tasks update` accept a stable task identifier flag such as `--task` or `--task-id`.
- [ ] `tasks search` accepts explicit machine-oriented flags for `--output`, `--incomplete`, `--limit`, and pagination.
- [ ] `tasks update` has explicit non-interactive update flags for notes editing, not only the action picker.
- [ ] Incompatible flag combinations are defined up front and documented.

**Implementation Notes**: Favor a single canonical flag name, with aliases only if needed for ergonomics. Reuse patterns from existing `--limit` flags across list commands.

**Estimated Effort**: Small

## Task 2: Add shared task resolution by GID

**Component**: Shared automation utilities, task commands

**Description**: Introduce a reusable task resolver that returns a fetched task when a GID is supplied, while preserving the existing interactive picker for commands invoked without a task ID.

**Dependencies**: Task 1

**Acceptance Criteria**:
- [ ] `tasks view --task <gid>` fetches that task directly without prompting.
- [ ] `tasks update --task <gid>` fetches that task directly without prompting.
- [ ] Existing prompt-based selection continues to work when no task ID is supplied and stdin is interactive.
- [ ] Invalid or missing task IDs produce a clear error message.

**Implementation Notes**: The API layer already supports `task.Fetch(client)` once a `Task{ID: gid}` shell object is available. This task is a good fit for `pkg/cmdutils`.

**Estimated Effort**: Small

## Task 3: Add a shared interactive guard for CI-safe execution

**Component**: Shared automation utilities, iostreams integration

**Description**: Introduce a common guard that fails fast when a command would prompt in a non-TTY context, and optionally allow an explicit no-prompt mode for scripts.

**Dependencies**: Task 1

**Acceptance Criteria**:
- [ ] If a command needs `survey` interaction and `stdin` is not a TTY, it returns an actionable error instead of blocking.
- [ ] The error message tells the user which non-interactive flags to provide.
- [ ] Commands that receive sufficient non-interactive flags run successfully in non-TTY environments.
- [ ] Interactive behavior in terminals is unchanged.

**Implementation Notes**: `pkg/iostreams.IOStreams` already tracks `IsStdinTTY`. Keep this logic centralized rather than re-implementing per command.

**Estimated Effort**: Small

## Task 4: Extend `tasks search` for automation-grade filtering and pagination

**Component**: `pkg/cmd/tasks/search`, Asana search query/options wiring

**Description**: Add explicit incomplete filtering, result limit, and cursor/offset paging to `tasks search`, wiring them to existing Asana query and option types.

**Dependencies**: Task 1

**Acceptance Criteria**:
- [ ] `tasks search` supports an explicit incomplete-only filter such as `--incomplete` or `--completed=false`.
- [ ] `tasks search` supports `--limit` and passes it through to the Asana request.
- [ ] `tasks search` supports pagination input such as `--page-offset` or equivalent cursor flag.
- [ ] Search help output documents these options clearly.

**Implementation Notes**: `asana.SearchTasksQuery.Completed` and `asana.Options.Limit/Offset` already exist; this is primarily command-layer exposure and looping behavior.

**Estimated Effort**: Medium

## Task 5: Add machine-readable output for `tasks search`

**Component**: Output formatting, `pkg/cmd/tasks/search`

**Description**: Add `--output json` to `tasks search` and return enough fields for stable downstream automation.

**Dependencies**: Task 1, Task 4

**Acceptance Criteria**:
- [ ] `tasks search --output json` returns valid JSON to stdout.
- [ ] Each result includes at minimum `gid`, `name`, `completed`, and `permalink_url`, plus any other displayed metadata.
- [ ] Human-readable output remains the default when `--output` is not specified.
- [ ] Error messages remain on stderr and do not corrupt JSON stdout.

**Implementation Notes**: Search currently requests only `name` and `due_on`. This task must expand field selection to include the machine-readable fields.

**Estimated Effort**: Medium

## Task 6: Add non-interactive full task read to `tasks view`

**Component**: `pkg/cmd/tasks/view`, output formatting

**Description**: Allow `tasks view` to read a task directly by GID and emit either the current human-readable details or a JSON representation of the full payload required by the branch guard workflow.

**Dependencies**: Task 2, Task 3

**Acceptance Criteria**:
- [ ] `tasks view --task <gid>` runs without prompting.
- [ ] `tasks view --task <gid> --output json` includes `gid`, `name`, `notes`, `completed`, and `permalink_url`.
- [ ] The existing prompt-driven `tasks view` still works in an interactive terminal.
- [ ] The command exits cleanly with an error when invoked non-interactively without a task ID.

**Implementation Notes**: Reuse the same JSON output conventions as search where possible, but include richer fields for a single task view.

**Estimated Effort**: Medium

## Task 7: Add explicit non-interactive update actions to `tasks update`

**Component**: `pkg/cmd/tasks/update`

**Description**: Add flag-driven update actions so scripts can update a specific task without selecting a task or action through `survey`.

**Dependencies**: Task 2, Task 3

**Acceptance Criteria**:
- [ ] `tasks update --task <gid> --name "..."` updates the task name without prompting.
- [ ] `tasks update --task <gid> --notes "..."` replaces the description without prompting.
- [ ] `tasks update --task <gid> --complete` marks the task complete without prompting.
- [ ] `tasks update --task <gid> --due-on YYYY-MM-DD` updates the due date without prompting.

**Implementation Notes**: Keep the interactive action menu for users who run `asana tasks update` with no non-interactive action flags.

**Estimated Effort**: Medium

## Task 8: Implement safe notes mutation helpers for prepend/merge workflows

**Component**: Shared automation utilities, `pkg/cmd/tasks/update`

**Description**: Add description mutation helpers so agents can prepend `Branch: ...` while preserving the existing body of the task notes.

**Dependencies**: Task 6, Task 7

**Acceptance Criteria**:
- [ ] `tasks update --task <gid> --prepend-notes "Branch: foo\n\n"` prepends text before the existing notes.
- [ ] A file-backed variant such as `--notes-file` or `--prepend-notes-file` is supported if the CLI contract includes it.
- [ ] Prepend logic preserves the existing notes content exactly after the inserted prefix.
- [ ] Empty existing notes are handled without producing malformed extra separators.

**Implementation Notes**: This task is what unlocks the branch guard workflow. Implement the transformation in plain Go string logic so it is easy to unit test.

**Estimated Effort**: Medium

## Task 9: Define machine-readable update responses

**Component**: Output formatting, `pkg/cmd/tasks/update`

**Description**: Define whether update commands should support `--output json` and, if so, return the updated task or a concise confirmation object for scripting.

**Dependencies**: Task 5, Task 7, Task 8

**Acceptance Criteria**:
- [ ] Update command output is deterministic in automation mode.
- [ ] If JSON is supported, it is valid and contains at least the task ID plus the fields changed.
- [ ] Human-readable success messages remain the default for interactive use.
- [ ] Update success output does not require parsing English prose in automation mode.

**Implementation Notes**: This is optional from the user’s gap list, but it materially improves script ergonomics and keeps the output model consistent across commands.

**Estimated Effort**: Small

## Task 10: Add automated test coverage for dual-mode behavior

**Component**: Task command tests, shared helper tests

**Description**: Add unit and command-level tests covering interactive fallback, non-interactive flags, notes mutation behavior, JSON output, pagination inputs, and non-TTY failures.

**Dependencies**: Task 2 through Task 9

**Acceptance Criteria**:
- [ ] Tests cover GID-based task resolution for `view` and `update`.
- [ ] Tests cover JSON output shape and stderr/stdout separation.
- [ ] Tests cover notes replacement and prepend semantics, including empty notes and existing branch headers if deduplication is later added.
- [ ] Tests cover non-TTY failure behavior when required flags are missing.

**Implementation Notes**: Existing command tests and `iostreams.Test()` provide a workable pattern for asserting stdout/stderr and TTY state.

**Estimated Effort**: Medium

## Task 11: Update help text and README examples for automation workflows

**Component**: Command help, README

**Description**: Update command examples and documentation so users understand the new CLI-only automation path.

**Dependencies**: Task 4 through Task 9

**Acceptance Criteria**:
- [ ] `tasks search`, `tasks view`, and `tasks update` help text includes the new automation flags.
- [ ] README includes at least one end-to-end non-interactive agent workflow example.
- [ ] Documentation distinguishes interactive and non-interactive modes clearly.
- [ ] Examples use stable task IDs instead of numbered prompt positions.

**Implementation Notes**: Keep examples grounded in the branch-guard use case because that is the main feature driver.

**Estimated Effort**: Small

### Dependency Graph

```text
Task 1 -> Tasks 2, 3, 4
Task 2 -> Tasks 6, 7
Task 3 -> Tasks 6, 7
Task 4 -> Task 5
Tasks 6 + 7 -> Task 8
Tasks 5 + 7 + 8 -> Task 9
Tasks 2..9 -> Task 10
Tasks 4..9 -> Task 11
```

### Parallelization Opportunities

- Task 2 and Task 3 can proceed in parallel after Task 1.
- Task 4 can proceed in parallel with Task 2 and Task 3.
- Task 6 and Task 7 can proceed in parallel once shared task resolution and non-interactive guardrails exist.
- Task 10 can begin once the first command-level changes land, but should finish after the JSON and notes-mutation work stabilizes.

### Critical Path

`Task 1 -> Task 2 -> Task 6 -> Task 7 -> Task 8 -> Task 10 -> Task 11`

This path covers the minimum required automation workflow of selecting a task, reading notes, prepending branch text, and trusting the behavior in CI.

## 4. Acceptance Criteria Reference

| Task | Primary Outcome | Key Acceptance Criteria |
| --- | --- | --- |
| 1 | Shared CLI contract | Stable task flag, output mode, incomplete filter, pagination contract defined |
| 2 | Direct task fetch by GID | `view`/`update` work with `--task <gid>` and skip prompts |
| 3 | CI-safe interaction | Non-TTY prompt attempts fail fast with actionable error |
| 4 | Search controls | Explicit incomplete filter, `--limit`, pagination flag |
| 5 | Search JSON | Valid JSON with `gid`, `name`, `completed`, `permalink_url` |
| 6 | View JSON/read path | `tasks view --task <gid> --output json` returns notes and permalink |
| 7 | Update flags | Name, notes, complete, and due date can be changed without prompts |
| 8 | Notes prepend | Branch prefix can be prepended while preserving existing notes |
| 9 | Update automation output | Deterministic update response in automation mode |
| 10 | Tests | Coverage for GID targeting, JSON, prepend logic, and non-TTY failure |
| 11 | Docs | Help and README explain the CLI-only automation workflow |

## 5. Validation Plan

### Requirement-to-Task Traceability

| Requirement ID | Requirement | Tasks Covering It | Validation Method |
| --- | --- | --- | --- |
| R1 | Address tasks by stable GID | 1, 2, 6, 7, 10, 11 | Command tests for `--task`; manual CLI smoke test |
| R2 | Read full task payload non-interactively | 2, 5, 6, 10 | JSON output tests; manual fetch of a known task |
| R3 | Update notes non-interactively | 1, 7, 8, 9, 10 | Unit tests for mutation helpers; command tests; manual update of a fixture task |
| R4 | Machine-readable search and view output | 1, 5, 6, 9, 10 | JSON schema/assertion tests; stdout/stderr separation tests |
| R5 | Explicit incomplete-only search | 4, 10, 11 | Query construction tests; manual search against mixed complete/incomplete tasks |
| R6 | Result limits and pagination | 4, 5, 10 | Pagination tests using mocked next-page offsets |
| R7 | Stream-friendly / CI-friendly mode | 3, 6, 7, 10 | Non-TTY command tests; manual run with piped stdin/stdout |

### Testing Strategy

#### Unit Tests

- Notes mutation helpers:
  - prepend into empty notes
  - prepend into multi-line notes
  - replace notes
  - optional file input handling
- Flag validation:
  - incompatible combinations
  - missing task ID for non-interactive actions
  - unsupported output mode values
- Query builder behavior:
  - incomplete filter mapping
  - limit propagation
  - pagination offset propagation

#### Command-Level Tests

- `tasks search --output json`
- `tasks search --incomplete --limit 5`
- `tasks view --task <gid> --output json`
- `tasks update --task <gid> --notes "..."`
- `tasks update --task <gid> --prepend-notes "..."`
- error behavior when stdin is not a TTY and a prompt would be required

#### Integration / Manual Validation

- Search for an open task queue, extract a GID, and fetch that task with `tasks view --output json`.
- Parse the first line of notes to simulate the branch guard.
- Prepend a branch header via `tasks update`.
- Re-run `tasks view` and confirm the new first line and preserved remaining notes.
- Confirm the CLI exits non-zero instead of hanging when the same commands are invoked without required flags in a non-TTY shell.

### Validation Checklist

- [ ] Search results can be consumed without parsing numbered English prose.
- [ ] A single task can be targeted solely by GID.
- [ ] Full task notes are readable from stdout without entering a prompt flow.
- [ ] Notes can be updated or prepended without opening the editor.
- [ ] Incomplete-only search is explicit and test-covered.
- [ ] Result limits and pagination are documented and test-covered.
- [ ] Interactive-only paths fail fast in CI/non-TTY usage.
- [ ] Legacy interactive behavior still works in a terminal.

## 6. Completion Criteria

### Implementation Complete

- [ ] All implementation tasks in this plan are complete.
- [ ] `tasks search`, `tasks view`, and `tasks update` support the documented non-interactive workflow.
- [ ] Stable GID targeting works across all relevant commands.
- [ ] Search supports explicit incomplete filtering, result limits, and pagination inputs.

### Quality Assurance

- [ ] Unit tests pass.
- [ ] Command-level tests pass.
- [ ] No JSON output path writes human-readable noise to stdout.
- [ ] No command enters a prompt flow in non-TTY mode without first failing fast.

### Requirement Validation

- [ ] R1 through R7 are covered by tests or manual validation evidence.
- [ ] The branch guard workflow can be completed using only the CLI.
- [ ] Existing interactive users are not forced into automation-oriented flags.

### Documentation

- [ ] Command help for `tasks search`, `tasks view`, and `tasks update` is updated.
- [ ] README examples include the non-interactive agent workflow.
- [ ] Any new flag aliases and incompatibilities are documented.

### Release Readiness

- [ ] Behavior is verified against a real Asana workspace.
- [ ] Error messages are actionable for both humans and agents.
- [ ] No regression is observed in the current interactive task flows.

## 7. Risk & Mitigation

| Risk | Impact | Probability | Mitigation |
| --- | --- | --- | --- |
| Flag sprawl makes the CLI inconsistent | Medium | Medium | Define the shared CLI contract first and reuse common flag names across commands |
| JSON output becomes unstable or mixes with human text | High | Medium | Centralize output mode handling and add stdout/stderr separation tests |
| Non-interactive path accidentally breaks interactive UX | High | Medium | Keep prompt-based logic as a fallback path and add explicit regression tests |
| Search pagination semantics are unclear to users | Medium | Medium | Use a clear cursor flag name and document examples in help/README |
| Notes prepend corrupts existing descriptions | High | Low | Isolate prepend logic in a small tested helper and validate against multiline notes |
| TTY detection misses some CI edge cases | Medium | Low | Base behavior on `IOStreams.IsStdinTTY` and add manual pipe-based validation |
| Asana field selection is incomplete for JSON output | Medium | Medium | Expand `opt_fields` deliberately and add tests that assert required fields are present |

### Fallback Plans

- If a shared JSON output abstraction is too large for the first pass, implement task-command-local JSON emitters with a consistent schema and refactor later.
- If update JSON output is too much scope, prioritize search/view JSON plus deterministic plain-text update success, but keep Task 9 explicitly tracked.
- If cursor pagination is awkward to expose initially, ship `--limit` plus a raw offset/cursor flag before designing a more ergonomic paging UX.

## 8. Next Steps

1. Finalize the command contract for new flags and incompatible combinations.
2. Implement shared task resolution and interactive guardrails before touching individual commands deeply.
3. Upgrade `tasks search` first, because it establishes JSON output, incomplete filtering, and limit/pagination behavior needed by the full workflow.
4. Implement `tasks view --task` next to unlock branch guard reads.
5. Implement non-interactive `tasks update` and notes prepend support last, then validate the end-to-end workflow.
6. Keep this document updated as implementation decisions narrow the exact flag names or output schema.

### Progress Tracking Guidance

- Mark tasks complete only after their acceptance criteria and validation items are satisfied.
- Record any deviations from the proposed CLI contract directly in this file so docs, tests, and implementation stay aligned.
- Treat Task 8 as the workflow gate: the feature is not complete until prepend behavior is proven against real task notes.
