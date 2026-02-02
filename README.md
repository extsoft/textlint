# textlint

CLI tool to check and fix text files. Scans files or directories (recursively), reports or fixes style issues. Only text files are processed (content-based: valid UTF-8, no null bytes).

## Usage

```bash
textlint [--check|--fix] [--format=compact|tap|json] [-q|--quiet|--verbose|--debug] <path> [path ...]
```

- **Check** (default): report issues to stdout. Exit code 1 if any issue is found.
- **Fix**: apply fixes in place. Prints how many files were fixed. Exit code 0.

You must pass at least one file or directory. Directories are scanned recursively.

### Options

| Option | Description |
|--------|-------------|
| `--check` | Check only, report issues (default if neither `--check` nor `--fix` is set). |
| `--fix` | Fix issues in place. |
| `--format=compact` | One line per issue plus summary (default). |
| `--format=tap` | TAP 13 output for test runners. |
| `--format=json` | JSON object with `files` and `summary`. |
| `-q`, `--quiet` | Quiet: no stdout; only fatal errors on stderr. Exit code still 1 when issues found. |
| `--verbose` | Verbose: steps, skipped (non-text) paths, and timing on stderr. |
| `--debug` | Debug: verbose plus non-text files skipped with reason, scanner accepted/rejected, rules per file, fix steps. |

Exactly one of `--check` or `--fix` is allowed. If multiple verbosity flags are set, the noisiest wins (debug > verbose > normal > quiet).

### Rules

| ID | Description |
|----|-------------|
| **TL001** | File must end with exactly one newline (LF or CRLF). |
| **TL010** | No trailing spaces or tabs at the end of a line. |

Both LF and CRLF line endings are supported; the tool preserves the detected style when fixing.

### Output formats

- **compact**: `file:line:col: rule: message` per issue, then `N file(s) scanned, M issue(s).`
- **tap**: TAP 13 (e.g. `1..M`, `not ok N - file:line:col rule message`).
- **json**: `{"files": {"path": [{"line", "column", "rule", "message"}]}, "summary": {"files", "issues"}}`.

### Text vs binary

Files are included only if they are valid UTF-8 and contain no null bytes. Binary and invalid-encoding files are skipped. When no text files are found, the summary includes "No text files found." (and "0 file(s) scanned, 0 issue(s)." in compact format).

## Development

Install [mise](https://mise.jdx.dev/) (dev tool version manager), then run `mise init` in the repo so the project’s tools and env are activated. See [mise docs](https://mise.jdx.dev/) for install and usage.

This project uses [hk](https://hk.jdx.dev) for code checks and git hooks. Use `mise check` or `mise fix` to check or autofix.

`mise build` builds the CLI binary.
