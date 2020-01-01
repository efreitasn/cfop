# cfop
cfop is a package for helping build CLIs in go.

<a href="https://godoc.org/github.com/efreitasn/cfop"><img src="https://godoc.org/github.com/efreitasn/cfop?status.svg" alt="GoDoc"></a>

## Features
* Shell completion.
* Built-in error and help messages.
* Validation of argument types (int, float, string).

## Naming
There's a lot of confusion regarding the terms in a CLI. This package aims to use a nomenclature that is well-known and doesn't present any ambiguity. The nomenclature is as follow:

### Term
Each word received when executing a command is called a term, including the command's name. It's also called arguments or args (in go, it's `os.Args`).

### Command
A command is always mapped to a parser. It can have options, flags, arguments and other commands as well. There's only one command per CLI and it's called the root command. Any other command after that is called a subcommand.

#### Subcommand
To be a subcommand, it must come right after the command, otherwise it'll be considered an argument to the command or to one of the command's options. Any options, flags or arguments that come after a subcommand are handled by the subcommand, and not by any previous command or subcommand.

### Option
An option starts with `-` or `--`. Generally, the `--` is the full version (e.g. `--name`), while the `-` version is the alias version (e.g. `-n`). An option always takes an argument, which can be added to the option in two ways:

```
--opt=20
--opt 20
```

In the example, `opt` is the name of the option and `20` is the argument.

### Flag
A flag is an option that can only be of boolean type. That way, it doesn't take an argument.

### Argument
If the term is not an option or flag, nor a subcommand, it is an argument. It can be an argument to an option, to a command or to a subcommand. An argument has a type, which is `TermInt`, `TermFloat` or `TermString`.

### Example
Let's take the `grep` command as an example to show how this nomenclature is applied:

```bash
grep foo -n -B=20
```

* `foo` is an argument.
* `-n` is a flag.
* `-B` is an option.
* `20` is the argument of the `-B` option.

Now with the `openssl` command:

```bash
openssl x509 -text -in root-cert.pem
```

* `x509` is a subcommand.
* `-text` is a flag.
* `-in` is an option.
* `root-cert.pem` is the argument of the `-in` option.

## The package
For the full documentation, see the [godoc page](https://godoc.org/github.com/efreitasn/cfop).

### Parsers
Each command is mapped to a parser. A parser takes a list of terms and performs parsing and validation on them. There are at least two parsers, the `rootCmd`, which is created implicitly, and the parser provided to the `Init()` function.

It all starts with the `rootCmd` parser. This parser doesn't apply any logic on the list of terms, it just calls the next parser, which is the one provided to the `Init()` function. After this, the next one depends on the current one and the next term in the list, if there's one. These steps are repeated until a `Cmd` parser is reached. Once a `Cmd` parser is reached, it's just a matter of parsing flags, options and arguments, if there's any.

There are three types of parsers:

* **rootCmd**: represents the root command. It's always used as the first parser. As the first letter of its name implies, this parser is not supposed to be used explicitly by the users of this package. Instead, users should use the `Init()` function, which takes the name and the description of the root command, a list of terms (e.g. os.Args) and a parser to parse the terms coming after the root command.

* **Cmd**: represents a command and contains a list of options, flags and arguments, as well as a function to be called when parsing the command.

* **SubcmdsSet**: represents a map of subcommands to parsers. Besides a parser, each subcommand also has a name and a description.

### Completion
This packages provides support for shell completion. To get a space-separated list of available subcommands, options or flags, run `<rootCmd> __introspect__ <strs>` where `<rootCmd>` is the root command's name and `<strs>` is a space-separated list of terms already typed by the user up to, but not including, the current one (in bash, `$COMP_WORDS[$COMP_CWORD]`). With the return of this command, the only step left is to match the current term with each tern in the list. This can be done with the `compgen` bash function. A completion bash script that makes use of the completion features provided by cfop, defaulting to filename completion in case there's no match, is as follows:

```bash
_<rootCmd>() 
{
    local cur opts
    COMPREPLY=()
    cur=$COMP_WORDS[$COMP_CWORD]
    opts=$(<rootCmd> __introspect__ "${COMP_WORDS[@]:1:${COMP_CWORD}-1}"|awk -v OFS="\n" '$1=$1')

    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
}

complete -o default -F _<rootCmd> <rootCmd>
```

`<rootCmd>` should be replaced by the root command's name, the same name that is passed as the first argument of the `Init()` function. Note that this script still needs to be placed in the `/etc/bash_completion.d` directory.

A better introduction to bash completion can be found [here](https://www.gnu.org/software/bash/manual/html_node/Programmable-Completion-Builtins.html).

### Example

```go
package main

import (
  "fmt"
  "os"

  "github.com/efreitasn/cfop"
)

func main() {
  set := cfop.NewSubcmdsSet()

  set.Add(
    "add",
    "Adds something",
    cfop.NewSubcmdsSet(
      cfop.Subcmd{
        Name:        "user",
        Description: "Adds user",
        Parser: cfop.NewCmd(
          cfop.CmdConfig{
            Fn: func(cts *cfop.CmdTermsSet) {},
            Flags: []cfop.CmdFlag{
              cfop.CmdFlag{
                Name:        "admin",
                Description: "whether the new user is an admin",
              },
            },
            Args: []cfop.CmdArg{
              cfop.CmdArg{
                Name:        "name",
                Description: "the name of the user",
                T:           cfop.TermString,
              },
              cfop.CmdArg{
                Name:        "email",
                Description: "the email of the user",
                T:           cfop.TermString,
              },
            },
          },
        ),
      },
    ),
  )

  set.Add(
    "update",
    "Updates something",
    cfop.NewSubcmdsSet(
      cfop.Subcmd{
        Name:        "user",
        Description: "Updates user",
        Parser: cfop.NewCmd(
          cfop.CmdConfig{
            Fn: func(cts *cfop.CmdTermsSet) {},
            Options: []cfop.CmdOption{
              cfop.CmdOption{
                Name:        "name",
                Description: "the name of the user",
                T:           cfop.TermString,
              },
              cfop.CmdOption{
                Name:        "email",
                Description: "the email of the user",
                T:           cfop.TermString,
              },
            },
          },
        ),
      },
    ),
  )

  set.Add(
    "remove",
    "Removes something",
    cfop.NewSubcmdsSet(
      cfop.Subcmd{
        Name:        "user",
        Description: "Removes user",
        Parser: cfop.NewCmd(
          cfop.CmdConfig{
            Fn: func(cts *cfop.CmdTermsSet) {},
            Args: []cfop.CmdArg{
              cfop.CmdArg{
                Name:        "name",
                Description: "the name of the user",
                T:           cfop.TermString,
              },
            },
          },
        ),
      },
    ),
  )

  err := cfop.Init(
    "app",
    "An app",
    os.Args,
    set,
  )
  if err != nil {
    fmt.Println(err)
  }
}
```

A graph that represents this CLI's structure is
<div align="center">
  <img alt="CLI graph" src="https://raw.githubusercontent.com/efreitasn/cfop/master/example_cli_graph.jpg">
</div>