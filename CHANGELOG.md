# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

#### Commands

- **`jira field list`** — list all Jira fields including custom field IDs (handy for figuring out `--custom` keys on `issue create`/`edit`). Supports `--custom-only`.
- **`jira issue comment list ISSUE-KEY`** — list comments on an issue, newest first. Supports `--limit N`.
- **`jira issue link list`** (alias `types`) — list the global issue link types (the names accepted as the third arg to `jira issue link`).
- **`jira issue worklog list ISSUE-KEY`** — list worklogs on an issue, newest first. Supports `--limit N`.
- **View subcommands** for the resource types that only had `list`:
  - `jira epic view EPIC-KEY` (wraps `issue view` since epics are issues).
  - `jira sprint view SPRINT-ID`
  - `jira project view PROJECT-KEY`
  - `jira board view BOARD-ID`
  - `jira release view VERSION-ID`

#### Flags

- **`-w, --web`** on `jira issue view` (and all new `view` subcommands) — opens the resource in the default browser. URL is printed to stdout. Honors `JIRA_BROWSER` / `BROWSER` env vars.
- **`--plain`**, **`--no-headers`**, **`--delimiter`**, **`--raw`**, **`--csv`** added to **`project list`**, **`board list`**, **`release list`** (previously only available on `issue`/`epic`/`sprint list`).
- **`--paginate`** added to every list command. When set, walks all pages internally and returns the full result set:
  - `jira issue list --paginate` / `jira epic list --paginate` — token-based on Cloud (v3), startAt on on-prem (v2).
  - `jira sprint list --paginate` / `jira sprint list SPRINT_ID --paginate` — startAt paging.
  - `jira board list --paginate` — startAt paging.
  - `jira epic list EPIC-KEY --paginate` — startAt paging.
  - `jira project list --paginate` / `jira release list --paginate` — documented no-op (endpoints already return all results).
- **`--raw`** on every new `view` subcommand.

#### Packaging

- **Nix flake package output**: `nix build .` produces `./result/bin/jira`; `nix run . -- issue list` works.

### Changed

- **`--paginate` semantics**: now a **boolean** flag everywhere. When set, walks all pages internally and returns the full set. Previously took `<from>:<limit>` string for client-side per-call sizing.
- **Spinner output**: progress spinners (e.g. "Fetching issues...") are auto-suppressed when stdout is not a TTY. Output piped to another command or redirected to a file no longer contains spinner garbage on stderr. No flag needed.
- **`jira sprint list --paginate`** (no-args) now fetches all sprints across all pages, ordered newest-first.

### Breaking changes

- **`--paginate` no longer accepts `<from>:<limit>`**. Scripts passing `--paginate 0:50` or `--paginate 20` will fail at flag parsing. Drop the value (`--paginate`) or omit the flag for default single-page behavior.
