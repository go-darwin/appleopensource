// Copyright 2020 The appleopensource Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"go.uber.org/multierr"
)

var (
	completionShells = map[string]func(cmd *cobra.Command, out io.Writer, boilerPlate string) error{
		"bash": runCompletionBash,
		"zsh":  runCompletionZsh,
	}
)

type completion struct {
	*aos

	ioStreams *IOStreams
	shell     string
}

// newCmdVersions creates the versions command.
func (a *aos) newCompletion(ctx context.Context, ioStreams *IOStreams) *cobra.Command {
	completion := &completion{
		aos:       a,
		ioStreams: ioStreams,
	}

	cmd := &cobra.Command{
		Use:   "completion",
		Short: "generate completion scrpit",
		RunE: func(cmd *cobra.Command, args []string) error {
			completion.shell = args[0]
			return completion.runCompletion(cmd, completion.ioStreams.Out)
		},
	}

	return cmd
}

func UsageErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}

// runCompletion checks given arguments and executes command.
func (c *completion) runCompletion(cmd *cobra.Command, out io.Writer) error {
	runFunc, found := completionShells[c.shell]
	if !found {
		return UsageErrorf(cmd, "Unsupported shell type %q.", c.shell)
	}

	return runFunc(cmd, out, "")
}

func runCompletionBash(cmd *cobra.Command, out io.Writer, boilerPlate string) error {
	if _, err := out.Write([]byte(boilerPlate)); err != nil {
		return err
	}

	return cmd.GenBashCompletion(out)
}

const (
	zshHead           = "#compdef aos\n"
	zshInitialization = `
__aos_bash_source() {
	alias shopt=':'
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__aos_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__aos_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__aos_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__aos_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__aos_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__aos_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__aos_filedir() {
	# Don't need to do anything here.
	# Otherwise we will get trailing space without "compopt -o nospace"
	true
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --version 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi
__aos_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__aos_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__aos_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__aos_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__aos_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__aos_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__aos_type/g" \
	<<'BASH_COMPLETION_EOF'
`

	zshTail = `
BASH_COMPLETION_EOF
}

__aos_bash_source <(__aos_convert_bash_to_zsh)
_complete aos 2>/dev/null
`
)

func runCompletionZsh(cmd *cobra.Command, out io.Writer, boilerPlate string) (errs error) {
	buf := new(bytes.Buffer)

	_, err := buf.WriteString(zshHead)
	errs = multierr.Append(errs, err)

	_, err = buf.WriteString(zshHead)
	errs = multierr.Append(errs, err)

	if boilerPlate != "" {
		_, err = buf.WriteString(boilerPlate)
		errs = multierr.Append(errs, err)
	}

	_, err = buf.WriteString(zshInitialization)
	errs = multierr.Append(errs, err)

	err = cmd.GenBashCompletion(buf)
	errs = multierr.Append(errs, err)

	_, err = buf.WriteString(zshTail)
	errs = multierr.Append(errs, err)

	_, err = out.Write(buf.Bytes())
	errs = multierr.Append(errs, err)

	return multierr.Combine(errs)
}
