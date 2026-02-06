# prosefmt

![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)

`prosefmt` is the simplest text formatter for making your files look correct. No complex rules, no massive configuration files — just clean text.

## CLI reference

**Synopsis**

```bash
prosefmt [command] [flags] [path...]
```

Pass at least one file or directory; directories are scanned recursively. By default the tool runs in check mode (report only). Use `--write` to apply fixes in place.

**Commands**

- [version](#version)
- [completion](#completion)

**Options**

- [--check](#--check)
- [--write](#--write)

**Output**

- [--format](#--format)
- [--quiet](#--quiet)
- [--verbose](#--verbose)
- [--debug](#--debug)

### `version`

Print the version number. Run: `prosefmt version`.

### `completion`

Generate a shell completion script. Usage: `prosefmt completion <shell>` with one of `bash`, `zsh`, `fish`, or `powershell`. See [Shell completion](#shell-completion) below for install steps.

### Options

#### `--check`

Check only: scan paths and report issues to stdout. Exit code is 1 if any issue is found, 0 otherwise. This is the default when neither `--check` nor `--write` is set. Exactly one of `--check` or `--write` is allowed.

#### `--write`

Write fixes in place. Files with issues are modified on disk. Prints how many files were written; exit code is 0. Exactly one of `--check` or `--write` is allowed.

### Output

#### `--format`

Output format for check mode. One of:

- **compact** (default): one line per issue as `file:line:col: rule: message`, then a summary line `N file(s) scanned, M issue(s).`
- **tap**: TAP 13 for test runners (e.g. `1..M`, `not ok N - file:line:col rule message`).
- **json**: JSON with `files` (path → list of `{line, column, rule, message}`) and `summary` (`files`, `issues`).

#### `--quiet`

Quiet (`-q`): no normal stdout; only fatal errors on stderr. Exit code is still 1 when issues are found in check mode.

#### `--verbose`

Verbose: emit steps, skipped (non-text) paths, and timing on stderr.

#### `--debug`

Debug: same as verbose, plus non-text files skipped with reason, scanner accepted/rejected list, rules per file, and write steps. If multiple verbosity flags are set, the noisiest wins (debug > verbose > normal > quiet).

## Implementatio Notes

### Rules

| ID | Description |
|----|-------------|
| **TL001** | File must end with exactly one newline (LF or CRLF). |
| **TL010** | No trailing spaces or tabs at the end of a line. |

Both LF and CRLF line endings are supported; the tool preserves the detected style when writing.

### Text vs binary

Files are included only if they are valid UTF-8 and contain no null bytes. Binary and invalid-encoding files are skipped. When no text files are found, the summary includes "No text files found." (and "0 file(s) scanned, 0 issue(s)." in compact format).

## Development

Install [mise](https://mise.jdx.dev/) (dev tool version manager), then run `mise init` in the repo so the project’s tools and env are activated. See [mise docs](https://mise.jdx.dev/) for install and usage.

This project uses [hk](https://hk.jdx.dev) for code checks and git hooks. Use `mise check` or `mise fix` to check or autofix.

`mise build` builds the CLI binary.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
