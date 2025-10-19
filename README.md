# Flexdot

[![Test](https://github.com/hidakatsuya/flexdot-go/actions/workflows/test.yml/badge.svg)](https://github.com/hidakatsuya/flexdot-go/actions/workflows/test.yml)

A flexible dotfile manager: Go implementation of [flexdot](https://github.com/hidakatsuya/flexdot).

## Features

- Declarative dotfile management using YAML index files
- Safe symlink creation with backup and restore
- Backup retention and cleanup
- Colorized output

## Getting Started

### Installation

Download the latest release from the [GitHub Releases page](https://github.com/hidakatsuya/flexdot-go/releases) and extract the archive for your platform. The binary name is `flexdot`.

#### Option: Build from Source

If you want to build from source (requires Go 1.24 or later):

```sh
git clone https://github.com/hidakatsuya/flexdot-go.git
cd flexdot-go
go build -o flexdot ./cmd
```

### Directory Structure Example

```
$HOME/dotfiles/
├── common/
│   ├── bin/
│   │   └── myscript
│   └── vim/
│       └── .vimrc
├── macOS/
│   ├── bash/
│   │   └── .bash_profile
│   └── codex/
│       └── prompts/
│           ├── code.md
│           ├── debug.md
│           └── test.md
├── ubuntu/
│   └── bash/
│       └── .bashrc
├── macOS.yml
├── ubuntu.yml
```

### Example Index File (`macOS.yml`)

```yaml
common:
  bin:
    myscript: bin
  vim:
    .vimrc: .
macOS:
  bash:
    .bash_profile: .
  codex:
    prompts:
      "*.md": .codex/prompts
```

This will link:
- `$HOME/dotfiles/common/bin/myscript` to `$HOME/bin/myscript`
- `$HOME/dotfiles/common/vim/.vimrc` to `$HOME/.vimrc`
- `$HOME/dotfiles/macOS/bash/.bash_profile` to `$HOME/.bash_profile`
- All `.md` files in `$HOME/dotfiles/macOS/codex/prompts/` to `$HOME/.codex/prompts/`

**Wildcard patterns**: You can use the `*` wildcard to match multiple files of the same type. For example, `"*.md"` matches all Markdown files in the directory.

### Usage

#### Install dotfiles

```sh
flexdot install [-H|--home_dir path] <index.yml>
```

- Use `--home_dir` or `-H` to specify the home directory.
- If `<index.yml>` or `--home_dir` is omitted, the value from `config.yml` will be used.

#### Clear all backups

```sh
flexdot clear-backups
```

### Command Reference

- `install [-H|--home_dir path] <index.yml>`
  Install dotfiles as specified in the index file.
  - `--home_dir`/`-H`: Set the home directory (overrides config.yml)
  - `<index.yml>`: Path to the index YAML file (overrides config.yml)
  - If omitted, values are taken from `config.yml`.
  - Both must be set either via CLI or config.yml.
- `clear-backups`
  Remove all backup directories under `./backup/`.

### config.yml

You can place a `config.yml` in your dotfiles directory to set default options:

You can also generate a template `config.yml` with default values by running:

```
flexdot init
```

```yaml
keep_max_count: 10         # (optional) Number of backup directories to keep (default: 10)
home_dir: /home/yourname   # (optional) Default home directory for install command
index_yml: ubuntu.yml      # (optional) Default index YAML file for install command
```

- CLI options take precedence over config.yml.
- If `keep_max_count` is omitted, the default value 10 is used.

### Backup

When a file is replaced, it is moved to a timestamped backup directory under `./backup/YYYYMMDDHHMMSS/`.

## Development

### Testing

Run all tests:

```sh
go test ./...
```

## Contributing

Contributions are welcome! Please open issues or pull requests.

## License

This project is licensed under the MIT License.

## Acknowledgements

- [flexdot](https://github.com/hidakatsuya/flexdot) - Original Ruby implementation and inspiration.
