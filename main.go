// Command gh-share is a gh extension that prints a short markdown (or
// plain) summary of a pull request and copies it to the clipboard, so it can
// be pasted directly into Slack or other chat clients.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	gh "github.com/cli/go-gh/v2"
	"github.com/spf13/pflag"
)

// defaultTemplate renders a Slack-friendly summary using colon-emoji
// shortcodes (`:pull_request:`, `:octocat:`), which Slack expands client
// side. Override it with --template or GH_PR_SHARE_TEMPLATE for clients that
// don't support shortcodes, or to match an org's preferred format.
//
// Links end in a trailing "/#" fragment, which keeps the link clickable
// while stopping Slack (and similar clients) from unfurling a preview card
// for it.
const defaultTemplate = `{{printf ":pull_request: *[%s](%s/#)* (+%.0f/-%.0f)\n:octocat: [%s/%s](https://github.com/%s/%s/#)" .title .url .additions .deletions .headRepositoryOwner.login .headRepository.name .headRepositoryOwner.login .headRepository.name}}`

const jsonFields = "title,url,additions,deletions,headRepositoryOwner,headRepository"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	// pflag (unlike the stdlib flag package) supports flags interspersed
	// with positional arguments, e.g. `gh share 14616 --no-clipboard`.
	fs := pflag.NewFlagSet("gh-share", pflag.ContinueOnError)
	fs.Usage = func() { fmt.Print(usage()) }

	template := fs.String("template", "", "Go template used to render the summary (overrides $GH_PR_SHARE_TEMPLATE)")
	repo := fs.StringP("repo", "R", "", "Select another repository using the [HOST/]OWNER/REPO format")
	noClipboard := fs.Bool("no-clipboard", false, "Print the summary without copying it to the clipboard")

	if err := fs.Parse(args); err != nil {
		if err == pflag.ErrHelp {
			return nil
		}
		return err
	}

	tmpl := resolveTemplate(*template)

	ghArgs := []string{"pr", "view", "--json", jsonFields, "--template", tmpl}
	if *repo != "" {
		ghArgs = append(ghArgs, "-R", *repo)
	}
	ghArgs = append(ghArgs, fs.Args()...)

	stdout, stderr, err := gh.Exec(ghArgs...)
	if err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("%s", strings.TrimRight(stderr.String(), "\n"))
		}
		return err
	}

	summary := stdout.String()
	fmt.Print(summary)

	if !*noClipboard {
		if err := clipboard.WriteAll(summary); err != nil {
			return fmt.Errorf("could not copy to clipboard: %w", err)
		}
		fmt.Fprintln(os.Stderr, "✓ Copied to clipboard")
	}

	return nil
}

// resolveTemplate picks, in order of precedence: the --template flag, the
// GH_PR_SHARE_TEMPLATE environment variable (handy for setting an org-wide
// default via MDM/shell profile), then the embedded default.
func resolveTemplate(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if envValue := os.Getenv("GH_PR_SHARE_TEMPLATE"); envValue != "" {
		return envValue
	}
	return defaultTemplate
}

func usage() string {
	return `Print a short markdown summary of a pull request and copy it to the clipboard.

USAGE
  gh share [<number> | <url> | <branch>] [flags]

FLAGS
  -R, --repo string      Select another repository using the [HOST/]OWNER/REPO format
      --template string   Go template used to render the summary (overrides $GH_PR_SHARE_TEMPLATE)
      --no-clipboard      Print the summary without copying it to the clipboard

ENVIRONMENT
  GH_PR_SHARE_TEMPLATE   Default template to use, e.g. set org-wide via shell profile

EXAMPLES
  $ gh share
  $ gh share 14616 -R cli/cli
  $ gh share https://github.com/cli/cli/pull/14616
`
}
