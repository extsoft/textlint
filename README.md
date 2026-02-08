# prosefmt

![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)

`prosefmt` is the simplest text formatter for making your files look correct. No complex rules, no massive configuration files — just clean text.

## CLI reference

**Synopsis**

```bash
prosefmt [command] [path...]
```

Pass at least one file or directory; directories are scanned recursively. With no command, runs **check** by default (report only). Use the **write** command to apply fixes in place. Output options (`--silent`, `--compact`, `--verbose`) apply only to the **check** and **write** commands.

**Commands**

- [check](#check) (default)
- [write](#write)
- [version](#version)

### `check`

Check files and report issues. Scan paths and report issues to stdout. Exit code is 1 if any issue is found, 0 otherwise. This is the default when no command is specified (e.g. `prosefmt path...`).

**Output** (only for this command):

- [--silent](#--silent): No standard output printed. Exit code is still 1 when issues are found.
- [--compact](#--compact): Show report / formatted or errored files (default when no output flag is set).
- [--verbose](#--verbose): Print debug output on stderr (steps, scanning summary, rules per file, timing).

### `write`

Write fixes in place. Files with issues are modified on disk. Prints how many files were written and lists each path; exit code is 0.

**Output** (only for this command): same as [check](#check) — `--silent`, `--compact`, `--verbose`.

### `version`

Print the version number. Run: `prosefmt version`.

### Output (check and write only)

Check prints a compact report: one line per issue as `file:line:col: rule: message`, grouped by file then rule; then a summary line `N file(s) scanned, M issue(s).`

By default output is **compact**: report (or "No text files found.", or "Wrote N file(s):" plus paths in write mode). If multiple output flags are set, the noisiest wins (verbose > compact > silent).

#### `--silent`

No standard output printed. Exit code is still 1 when issues are found in check mode.

#### `--compact`

Show formatted or errored files: report in check mode, "No text files found." when applicable, "Wrote N file(s):" plus one path per line in write mode. Default when no output level flag is set.

#### `--verbose`

Print debug output: steps, scanning summary, scanner accepted/rejected with reasons, rules per file, write steps, and timing on stderr.

## Implementation Notes

### Rules

| ID | Description |
|----|-------------|
| **TL001** | File must end with exactly one newline (LF or CRLF). |
| **TL010** | No trailing spaces or tabs at the end of a line. |

Both LF and CRLF line endings are supported; the tool preserves the detected style when writing.

### Text vs binary

Files are included only if they are valid UTF-8 and contain no null bytes. Binary and invalid-encoding files are skipped. When no text files are found, the summary includes "No text files found." (and "0 file(s) scanned, 0 issue(s).").

## Development

Install [mise](https://mise.jdx.dev/) (dev tool version manager), then run `mise run init` in the repo so the project’s tools and env are activated. See [mise docs](https://mise.jdx.dev/) for install and usage.

This project uses [hk](https://hk.jdx.dev) for code checks and git hooks. Use `mise run check` or `mise run fix` to check or autofix.

`mise run build` builds the CLI binary.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
