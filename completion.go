package cfop

import "fmt"

var completionStr = `_%[1]v() 
{
    local opts
    COMPREPLY=()
    opts=$(%[1]v __introspect__ "${COMP_WORDS[@]:1:$COMP_CWORD-1}")

    COMPREPLY=($(compgen -W "${opts}" -- "${COMP_WORDS[1]}"))
}

complete -o default -F _%[1]v %[1]v
`

func newCompletionParser(rootCmdName string) Parser {
	return NewSubcmdsSet(
		Subcmd{
			Name:        "bash",
			Description: "prints the bash completion",
			Parser: NewCmd(CmdConfig{
				Fn: func(terms *CmdTermsSet) {
					fmt.Printf(completionStr, rootCmdName)
				},
			}),
		},
		Subcmd{
			Name:        "zsh",
			Description: "prints the zsh completion",
			Parser: NewCmd(CmdConfig{
				Fn: func(terms *CmdTermsSet) {
					fmt.Printf(completionStr, rootCmdName)
				},
			}),
		},
	)
}
