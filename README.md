# gh-share

A [`gh`](https://cli.github.com) extension that prints a short markdown summary
of a pull request and copies it to the clipboard, so it can be pasted
straight into Slack or other chat clients.

## Install

```
gh extension install edjeffreys/gh-share
```

## Usage

```
gh share [<number> | <url> | <branch>] [flags]
```

With no argument, it uses the pull request associated with the current
branch (same resolution `gh pr view` uses).

```
$ gh share
:pull_request: *[Clarify issue view comment output modes](https://github.com/cli/cli/pull/14616/#)* (+1/-0)
:octocat: [cli/cli](https://github.com/cli/cli/#)
✓ Copied to clipboard
```

### Flags

| Flag              | Description                                                        |
| ------------------ | ------------------------------------------------------------------ |
| `-R, --repo`       | Select another repository using the `[HOST/]OWNER/REPO` format     |
| `--template`       | Go template used to render the summary (see below)                 |
| `--no-clipboard`   | Print the summary without copying it to the clipboard              |

### Customizing the template

The default template uses Slack colon-emoji shortcodes (`:pull_request:`,
`:octocat:`), which Slack expands client side but other tools may render as
literal text. Override it per invocation with `--template`, or set a
consistent default across an organization's machines with the
`GH_PR_SHARE_TEMPLATE` environment variable (e.g. in a shared shell profile
or MDM-managed dotfiles):

```sh
export GH_PR_SHARE_TEMPLATE='{{printf "[%s](%s) (+%.0f/-%.0f)" .title .url .additions .deletions}}'
```

Precedence: `--template` flag > `GH_PR_SHARE_TEMPLATE` env var > built-in
default.

The template is executed against the same JSON fields as
`gh pr view --json title,url,additions,deletions,headRepositoryOwner,headRepository`.

## Why an extension rather than a `gh` built-in?

This started as a proposal to add `gh pr share` to the `gh` CLI itself. The
maintainers suggested this kind of user/org-specific workflow is a better fit
for an extension, which can be installed consistently across an
organization's machines without needing to land in core `gh`.

## License

[MIT](LICENSE)
