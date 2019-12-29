# cfp
cfp is a package to help build CLIs in go.

## Naming
There's a lot of confusion regarding the terms in a CLI. This package aims to use a nomenclature that is well-known and doesn't present any ambiguity. The nomenclature is as follow:

### Command (CMD)
A command is always mapped to a function. It can have options, flags, arguments and other commands as well, called subcommands.

### Option
An option starts with `-` or `--`. Generally, the `--` is the full version (e.g. `--name`), while the `-` version is the alias version (e.g. `-n`). An option always takes an argument, which can only be of int, float or string type. This argument can be added to the option in two ways:

```
--opt=20
--opt 20
```

In the example, `opt` is the name of the option and `20` is the argument.

### Flag
A flag is an option that can only be of boolean type. That way, it doesn't take an argument.

### Subcommand
A subcommand doesn't start with any special character and it's basically an argument with a function mapped to it. To be a subcommand, it must come right after the command, otherwise it'll be considered an argument to the command. Any options, flags or arguments that come after a subcommand are handled by the subcommand, and not by the command.

### Argument
If the term is not an option or flag, nor a subcommand, it is an argument. It can be an argument to an option, to a command or to a subcommand.

### Example
Let's take the grep command as an example to show how this nomenclature plays out:

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