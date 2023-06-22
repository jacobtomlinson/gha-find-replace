# Find and Replace Action

[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-Find%20and%20Replace-blue.svg?colorA=24292e&colorB=0366d6&style=flat&longCache=true&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAA4AAAAOCAYAAAAfSC3RAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAM6wAADOsB5dZE0gAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAERSURBVCiRhZG/SsMxFEZPfsVJ61jbxaF0cRQRcRJ9hlYn30IHN/+9iquDCOIsblIrOjqKgy5aKoJQj4O3EEtbPwhJbr6Te28CmdSKeqzeqr0YbfVIrTBKakvtOl5dtTkK+v4HfA9PEyBFCY9AGVgCBLaBp1jPAyfAJ/AAdIEG0dNAiyP7+K1qIfMdonZic6+WJoBJvQlvuwDqcXadUuqPA1NKAlexbRTAIMvMOCjTbMwl1LtI/6KWJ5Q6rT6Ht1MA58AX8Apcqqt5r2qhrgAXQC3CZ6i1+KMd9TRu3MvA3aH/fFPnBodb6oe6HM8+lYHrGdRXW8M9bMZtPXUji69lmf5Cmamq7quNLFZXD9Rq7v0Bpc1o/tp0fisAAAAASUVORK5CYII=)](https://github.com/jacobtomlinson/gha-find-replace)
[![Actions Status](https://github.com/jacobtomlinson/gha-find-replace/workflows/Build/badge.svg)](https://github.com/jacobtomlinson/gha-find-replace/actions)
[![Actions Status](https://github.com/jacobtomlinson/gha-find-replace/workflows/Integration%20Test/badge.svg)](https://github.com/jacobtomlinson/gha-find-replace/actions)

This action will find and replace strings in your project files.

## Usage

### Example workflow

This example replaces `hello` with `world` in all of your project files.

```yaml
name: My Workflow
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Find and Replace
        uses: jacobtomlinson/gha-find-replace@v3
        with:
          find: "hello"
          replace: "world"
          regex: false
```

### Inputs

| Input                  | Description                                                                                                                            |
| ---------------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| `find`                 | A string to find and replace in your project files. _(Defaults to a [regular expression](https://github.com/google/re2/wiki/Syntax).)_ |
| `replace`              | The string to replace it with.                                                                                                         |
| `include` _(optional)_ | A glob of files to include. _Defaults to `**`._                                                                                        |
| `exclude` _(optional)_ | A glob of files to exclude. _Defaults to `.git/**` ._                                                                                  |
| `regex` _(optional)_   | Whether to match with.find as a regular expression instead of a fixed string. _Defaults to `true`._                                    |

### Outputs

| Output          | Description                                 |
| --------------- | ------------------------------------------- |
| `modifiedFiles` | The number of files that have been modified |

## Examples

### Including a subdirectory

You can limit your find and replace to a directory.

```yaml
name: My Workflow
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Find and Replace
        uses: jacobtomlinson/gha-find-replace@v3
        with:
          find: "hello"
          replace: "world"
          include: "justthisdirectory/**"
          regex: true
```

### Filter by file name

You can limit your find and replace to just files with a specific name.

```yaml
name: My Workflow
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Find and Replace
        uses: jacobtomlinson/gha-find-replace@v3
        with:
          find: "hello"
          replace: "world"
          include: "**README.md" # Will match all README.md files in any nested directory
```

### Exclude by file type

You can set your find and replace to ignore certain file types.

```yaml
name: My Workflow
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Find and Replace
        uses: jacobtomlinson/gha-find-replace@v3
        with:
          find: "hello"
          replace: "world"
          exclude: "**/*.py" # Do not modify Python files
```

If you are getting an error message that reads, "refusing to allow a GitHub App to create or update workflow ..." it means that your GitHub action may be trying to edit your new workflow file. Omit this file from the find-and-replace search with `exclude: .`.

```yaml
name: My Workflow
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Find and Replace
        uses: jacobtomlinson/gha-find-replace@v3
        with:
          find: "hello"
          replace: "world"
          exclude: . # Do not modify this file
```

### Pushing changes back

Any modifications during a GitHub Actions workflow are only made to the working copy checked out by the `actions/checkout` step. If you want those changes to be pushed back to the repository you'll need to add a final step that does this. You will have to give your workflow write permissions.

```yaml
name: My Workflow
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v2
      - name: Find and Replace
        uses: jacobtomlinson/gha-find-replace@v3
        with:
          find: "hello"
          replace: "world"
          regex: false
      - name: Push changes
        uses: stefanzweifel/git-auto-commit-action@v4
```

_If you need the push event to trigger other workflows, use a `repo` scoped [Personal Access Token](https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line)._

```yaml
      - name: Push changes
        uses: ad-m/github-push-action@v0.6.0
        with:
          github_token: ${{ secrets.MY_PAT }}
          branch: ${{ github.ref }}
```

## Publishing

To publish a new version of this Action we need to update the Docker image tag in `action.yml` and also create a new release on GitHub.

- Work out the next tag version number.
- Update the Docker image in `action.yml`.
- Create a new release on GitHub with the same tag.
