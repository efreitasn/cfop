/*
Package cfop provides a library for building CLIs.

To start parsing a CLI, just call the Init function passing the name and description of the root command, the
os.Args slice and a parser.

There are three parsers: rootCmd, Cmd and SubcmdsSet. The first one is just a reference to the root command's
name and description (the data passed to the Init function) and is created implicitly by the same function.
The second is a struct holding information about flags, options and arguments. The last one is a map of name
and description to another Parser, and doesn't accept any options or flags, only subcommands.

When parsing, the first parser called is the one provided to the Init function and any subsequent parsers are
called according to the terms (os.Args) until a Cmd is reached. After that, any options, arguments and/or flags
are parsed.

This means that a subcmd can have subcommands and any one of them can also have subcommands and so on until a
Cmd is reached. An example is:

	set := cfop.NewSubcmdsSet()

	set.Add(
		"foo",
		"foo's description",
		cfop.NewSubcmdsSet(
			cfop.Subcmd{
				Name:        "bar",
				Description: "foo",
				Parser: cfop.NewCmd(cfop.CmdConfig{
					Fn: func(cts *cfop.CmdTermsSet) {
						fmt.Println("hello world")
					},
				}),
			},
		),
	)

	set.Add(
		"foobar",
		"",
		cfop.NewCmd(cfop.CmdConfig{
			Fn: func(cts *cfop.CmdTermsSet) {
				fmt.Println("foobar")
			},
			Options: []cfop.CmdOption{
				cfop.CmdOption{
					Name:  "name",
					Alias: "n",
					T:     cfop.TermString,
				},
			},
		}),
	)

	err := cfop.Init(
		"testing",
		"A test",
		os.Args,
		set,
	)
	if err != nil {
		logs.Err.Println(err)
	}

In this example, testing is the root command's name and it has two subcommands: foo and foobar. The foobar
subcommand uses a Cmd parser, which has only one option, called name. The foo subcommand is a SubcmdsSet parser,
which means that it also has one or more subcommands. In this case, it has one subcommand named bar, which is a
Cmd. Now, when run

	testing foo bar

the root command's parser (SubcmdsSet) is called, which, by checking the terms (os.Args), will conclude that the
next parser is foo's (SubcmdsSet). This parser, by the same means, will conclude that the next one is bar's (Cmd).
Since bar's doesn't have any options, flags or aguments, it will only call the function provided to the Fn field,
which will print hello world to the user.

This library also provides completion features. For more info, see the GitHub page of this package.
*/
package cfop
