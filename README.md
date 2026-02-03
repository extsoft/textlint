# prosefmt

![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)

`prosefmt` is the simplest file formatter for when you just want your text to look right. No complex rules, no massive configuration files — just clean text.

## Usage

```bash
prosefmt [--check|--write] [--format=compact|tap|json] [-q|--quiet|--verbose|--debug] <path> [path ...]
```

- **Check** (default): report issues to stdout. Exit code 1 if any issue is found.
- **Write**: write fixes in place. Prints how many files were written. Exit code 0.

You must pass at least one file or directory. Directories are scanned recursively.

### Options

| Option | Description |
|--------|-------------|
| `--check` | Check only, report issues (default if neither `--check` nor `--write` is set). |
| `--write` | Write fixes in place. |
| `--format=compact` | One line per issue plus summary (default). |
| `--format=tap` | TAP 13 output for test runners. |
| `--format=json` | JSON object with `files` and `summary`. |
| `-q`, `--quiet` | Quiet: no stdout; only fatal errors on stderr. Exit code still 1 when issues found. |
| `--verbose` | Verbose: steps, skipped (non-text) paths, and timing on stderr. |
| `--debug` | Debug: verbose plus non-text files skipped with reason, scanner accepted/rejected, rules per file, write steps. |

Exactly one of `--check` or `--write` is allowed. If multiple verbosity flags are set, the noisiest wins (debug > verbose > normal > quiet).

### Rules

| ID | Description |
|----|-------------|
| **TL001** | File must end with exactly one newline (LF or CRLF). |
| **TL010** | No trailing spaces or tabs at the end of a line. |

Both LF and CRLF line endings are supported; the tool preserves the detected style when writing.

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

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
