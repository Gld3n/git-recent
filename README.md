# git-recent
CLI tool to quickly (git) switch back to recent branches. Leverages [Huh](https://github.com/charmbracelet/huh) to provide a nice and simple selection

## Features
- Shows your last **5** visited branches (deduplicated) by default.
- Navigate with **Arrow keys** or Vim keys (`j`/`k`).
- Works automatically as a native git command (`git recent`).

## Installation

### Option A: Using Go (Recommended)
```bash
go install [github.com/YOUR_USERNAME/git-recent@latest](https://github.com/YOUR_USERNAME/git-recent@latest)
```

### Option B: Manual
Download the binary for your OS from the Releases page.

Move it to a folder in your $PATH (e.g., /usr/local/bin or C:\Windows\System32).

### Usage
Since the tool is named git-recent, Git automatically detects it as a subcommand.

```bash
# Show last 5 branches (default)
git recent

# Show last 10 branches
git recent -n 10
```

### How it works
It parses your git reflog to find where you've actually been, filtering out duplicates and deleted branches, then uses a TUI to let you select one.
