package cfop

import (
	"fmt"
	"strings"
)

// TermType is the type of a cmd term.
type TermType string

// Term types.
const (
	TermInt    TermType = "int"
	TermFloat  TermType = "float"
	TermString TermType = "string"
)

// CmdTermsSet is a set of terms passed alongside a cmd.
type CmdTermsSet struct {
	cmd           *Cmd
	optionsValues map[string]interface{}
	flagsValues   map[string]bool
	argsValues    map[string]interface{}
}

// GetOptString returns the value of an option of type string.
// name can be either the option's name or the option's alias.
// If the option doesn't exist, its zero value is returned.
func (ct *CmdTermsSet) GetOptString(name string) string {
	opt := ct.cmd.getOption(name)
	if opt == nil || opt.T != TermString {
		return ""
	}

	return (ct.optionsValues[opt.Name]).(string)
}

// GetOptInt returns the value of an option of type integer.
// name can be either the option's name or the option's alias.
// If the option doesn't exist, its zero value is returned.
func (ct *CmdTermsSet) GetOptInt(name string) int {
	opt := ct.cmd.getOption(name)
	if opt == nil || opt.T != TermInt {
		return 0
	}

	return (ct.optionsValues[opt.Name]).(int)
}

// GetOptFloat returns the value of an option of type float.
// name can be either the option's name or the option's alias.
// If the option doesn't exist, its zero value is returned.
func (ct *CmdTermsSet) GetOptFloat(name string) float64 {
	opt := ct.cmd.getOption(name)
	if opt == nil || opt.T != TermFloat {
		return 0
	}

	return (ct.optionsValues[opt.Name]).(float64)
}

// GetFlag returns the value of a flag.
// name can be either the flag's name or the flag's alias.
// If the flag doesn't exist, its zero value is returned.
func (ct *CmdTermsSet) GetFlag(name string) bool {
	f := ct.cmd.getFlag(name)
	if f == nil {
		return false
	}

	return ct.flagsValues[f.Name]
}

// GetArgString returns the value of the argument at n of type string.
// If the argument doesn't exist, its zero value is returned.
func (ct *CmdTermsSet) GetArgString(name string) string {
	arg := ct.cmd.getArgByName(name)
	if arg == nil || arg.T != TermString {
		return ""
	}

	return (ct.argsValues[name]).(string)
}

// GetArgInt returns the value of the argument at n of type intger.
// If the argument doesn't exist, its zero value is returned.
func (ct *CmdTermsSet) GetArgInt(name string) int {
	arg := ct.cmd.getArgByName(name)
	if arg == nil || arg.T != TermInt {
		return 0
	}

	return (ct.argsValues[name]).(int)
}

// GetArgFloat returns the value of the argument at n of type float.
// If the argument doesn't exist, its zero value is returned.
func (ct *CmdTermsSet) GetArgFloat(name string) float64 {
	arg := ct.cmd.getArgByName(name)
	if arg == nil || arg.T != TermFloat {
		return 0
	}

	return (ct.argsValues[name]).(float64)
}

// CmdOption is a cmd option.
type CmdOption struct {
	// Name is used with --, is case-sensitive and cannot start with -.
	Name string
	// Alias is used with -, is case-senstitive and cannot start with -.
	Alias       string
	Description string
	// T is the type of the option.
	T        TermType
	Required bool
}

// CmdFlag is a cmd flag.
type CmdFlag struct {
	// Name is used with --, is case-senstive and cannot start with.
	Name string
	// Alias is used with -, is case-senstive and cannot start with -.
	Alias       string
	Description string
}

// CmdArg is a cmd argument.
type CmdArg struct {
	// Name is used when printing the help message and is case-senstive.
	Name        string
	Description string
	// T is the type of the argument.
	T TermType
}

// CmdConfig is a config used to create a cmd.
type CmdConfig struct {
	Fn      func(*CmdTermsSet)
	Options []CmdOption
	Flags   []CmdFlag
	Args    []CmdArg
}

// Cmd is a command.
type Cmd struct {
	fn              func(*CmdTermsSet)
	options         map[string]*CmdOption
	optionsByAlias  map[string]*CmdOption
	requiredOptions []*CmdOption
	flags           map[string]*CmdFlag
	flagsByAlias    map[string]*CmdFlag
	argsByPos       []*CmdArg
	argsByName      map[string]*CmdArg
}

// NewCmd creates a cmd.
// If an invalid flag name or option name is passed or a function isn't passed, it panics.
func NewCmd(cc CmdConfig) *Cmd {
	options := make(map[string]*CmdOption)
	requiredOptions := make([]*CmdOption, 0)
	optionsByAlias := make(map[string]*CmdOption)
	flags := make(map[string]*CmdFlag)
	flagsByAlias := make(map[string]*CmdFlag)
	argsByName := make(map[string]*CmdArg, len(cc.Args))
	argsByPos := make([]*CmdArg, 0, len(cc.Args))

	if cc.Fn == nil {
		panic(ErrMissingCmdFn)
	}

	if cc.Options != nil {
		for i := range cc.Options {
			opt := cc.Options[i]

			if opt.Name == "" ||
				!isOptionWithoutValue("--"+opt.Name) ||
				(opt.Alias != "" && strings.HasPrefix(opt.Alias, "-")) {
				panic(ErrInvalidOptionNameOrAlias)
			}

			if opt.T == "" {
				panic(ErrMissingTermTypeForTerm{Term: opt.Name})
			}

			options[opt.Name] = &opt

			if opt.Required {
				requiredOptions = append(requiredOptions, &opt)
			}

			if opt.Alias != "" {
				optionsByAlias[opt.Alias] = &opt
			}
		}
	}

	if cc.Flags != nil {
		for i := range cc.Flags {
			flag := cc.Flags[i]

			if flag.Name == "" || !isOptionWithoutValue("--"+flag.Name) || (flag.Alias != "" && !isOptionWithoutValue("-"+flag.Alias)) {
				panic(ErrInvalidFlagNameOrAlias)
			}

			flags[flag.Name] = &flag

			if flag.Alias != "" {
				flagsByAlias[flag.Alias] = &flag
			}
		}
	}

	if cc.Args != nil {
		for i := range cc.Args {
			arg := cc.Args[i]

			if arg.Name == "" {
				panic(ErrInvalidArgumentName{ArgumentPos: i})
			}

			if arg.T == "" {
				panic(ErrMissingTermTypeForTerm{Term: arg.Name})
			}

			argsByName[arg.Name] = &arg
			argsByPos = append(argsByPos, &arg)
		}
	}

	return &Cmd{
		fn:              cc.Fn,
		options:         options,
		requiredOptions: requiredOptions,
		optionsByAlias:  optionsByAlias,
		flags:           flags,
		flagsByAlias:    flagsByAlias,
		argsByPos:       argsByPos,
		argsByName:      argsByName,
	}
}

func (c *Cmd) getFlag(nameOrAlias string) *CmdFlag {
	f, ok := c.flags[nameOrAlias]
	if !ok {
		f, ok := c.flagsByAlias[nameOrAlias]
		if !ok {
			return nil
		}

		return f
	}

	return f
}

func (c *Cmd) getOption(nameOrAlias string) *CmdOption {
	opt, ok := c.options[nameOrAlias]
	if !ok {
		opt, ok := c.optionsByAlias[nameOrAlias]
		if !ok {
			return nil
		}

		return opt
	}

	return opt
}

func (c *Cmd) getArgByPos(n int) *CmdArg {
	if len(c.argsByPos) >= (n + 1) {
		return c.argsByPos[n]
	}

	return nil
}

func (c *Cmd) getArgByName(name string) *CmdArg {
	arg, ok := c.argsByName[name]
	if !ok {
		return nil
	}

	return arg
}

// Parse parses a slice of strings.
func (c *Cmd) Parse(pp parentParser, strs []string) error {
	tSet := &CmdTermsSet{
		cmd:           c,
		optionsValues: make(map[string]interface{}),
		argsValues:    make(map[string]interface{}),
		flagsValues:   make(map[string]bool),
	}
	i := 0

	for i < len(strs) {
		str := strs[i]

		if isHelpFlag(str) {
			printHelp(c, pp)

			return nil
		}

		if isOptionWithValue(str) {
			optName, isAlias := extractOptionName(str)

			var opt *CmdOption

			if isAlias {
				opt = c.optionsByAlias[optName]
			} else {
				opt = c.options[optName]
			}

			if opt == nil {
				return ErrUnexpectedOption{
					OptionName: optName,
					IsAlias:    isAlias,
				}
			}

			optValueStr := extractOptionValue(str)
			if optValueStr == "" {
				return ErrOptionsExpectsAValue{
					OptionName: optName,
					IsAlias:    isAlias,
				}
			}

			optValue, validValue := isValueValidForTermType(opt.T, optValueStr)
			if !validValue {
				return ErrOptionExpectsDifferentValueType{
					OptionName:   optName,
					IsAlias:      isAlias,
					ExpectedType: opt.T,
				}
			}

			tSet.optionsValues[opt.Name] = optValue

			i++
			continue
		}

		if isOptionWithoutValue(str) {
			optName, isAlias := extractOptionName(str)

			var opt *CmdOption

			if isAlias {
				opt = c.optionsByAlias[optName]
			} else {
				opt = c.options[optName]
			}

			if opt == nil {
				// An option without value could be a flag
				var f *CmdFlag

				if isAlias {
					f = c.flagsByAlias[optName]
				} else {
					f = c.flags[optName]
				}

				if f == nil {
					return ErrUnexpectedOptionOrFlag{
						OptionOrFlagName: optName,
						IsAlias:          isAlias,
					}
				}

				tSet.flagsValues[f.Name] = true

				i++
				continue
			}

			if len(strs) > (i+1) && !isOptionWithValue(strs[i+1]) && !isOptionWithoutValue(strs[i+1]) {
				optValueStr := strs[i+1]
				optValue, validValue := isValueValidForTermType(opt.T, optValueStr)
				if !validValue {
					return ErrOptionExpectsDifferentValueType{
						OptionName:   optName,
						IsAlias:      isAlias,
						ExpectedType: opt.T,
					}
				}

				tSet.optionsValues[opt.Name] = optValue

				i += 2

				continue
			}

			return ErrOptionsExpectsAValue{
				OptionName: optName,
				IsAlias:    isAlias,
			}
		}

		// If it reaches this part, it means it's not an option with value
		// (--opt=value) nor an option without value or flag (--opt). This
		// way, we consider it as an argument.
		nextArgPos := len(tSet.argsValues)

		if len(c.argsByName) >= (nextArgPos + 1) {
			arg := c.argsByPos[nextArgPos]
			argVal, valid := isValueValidForTermType(arg.T, str)

			if valid {
				tSet.argsValues[arg.Name] = argVal

				i++
				continue
			} else {
				return ErrArgumentExpectsDifferentValueType{
					ArgumentPos:  nextArgPos,
					ArgumentName: arg.Name,
					ExpectedType: arg.T,
					Value:        str,
				}
			}
		} else {
			return ErrUnexpectedArgument{Argument: str}
		}
	}

	if len(tSet.argsValues) != len(c.argsByName) {
		return ErrMissingArguments
	}

	for _, opt := range c.requiredOptions {
		if _, ok := tSet.optionsValues[opt.Name]; !ok {
			return ErrRequiredOptionNotProvided{OptionName: opt.Name}
		}
	}

	c.fn(tSet)

	return nil
}

func (c *Cmd) help(pp parentParser) string {
	numCols, err := getTermNumCols()
	if err != nil {
		numCols = 67
	}

	sb := strings.Builder{}

	ppDescription := getParentParserDescription(pp)
	if ppDescription != "" {
		sb.WriteString(ppDescription + "\n\n")
	}

	sb.WriteString(fmt.Sprintf("Usage: %v", strings.Join(pp.cmds, " ")))

	hasArgs := c.argsByPos != nil && len(c.argsByPos) > 0
	hasRequiredOptions := c.requiredOptions != nil && len(c.requiredOptions) > 0
	hasOptionalOptions := c.options != nil && len(c.options) > 0 && len(c.requiredOptions) != len(c.options)
	hasFlags := c.flags != nil && len(c.flags) > 0

	if hasArgs {
		for _, arg := range c.argsByPos {
			helpName := "<" + arg.Name + ">"

			sb.WriteString(" " + helpName)
		}
	}

	if hasRequiredOptions {
		sb.WriteString(" OPTIONS")
	}

	if hasOptionalOptions {
		sb.WriteString(" [OPTIONS]")
	}

	if hasFlags {
		sb.WriteString(" [FLAGS]")
	}

	sb.WriteRune('\n')

	// Arguments
	if hasArgs {
		biggestArgHelpNameLen := findBiggestArgHelpNameLen(c.argsByName)
		sb.WriteRune('\n')

		for _, arg := range c.argsByPos {
			argNameStyled, argNameUnstyled := buildArgumentHelpName(arg.Name)
			sb.WriteString(argNameStyled)

			if arg.Description != "" {
				descripFormatted := breakStringIntoPaddedLines(
					numSpacesHelpNameAndDescription+biggestArgHelpNameLen,
					' ',
					numCols,
					arg.Description,
				)

				// the new slice was created so that the help name could
				// align with the description.
				sb.Write([]byte(descripFormatted[len(argNameUnstyled):]))
			}

			sb.WriteRune('\n')
		}
	}

	// Biggest option or flag's name length.
	biggestOptionOrFlagHelpNameLen := findBiggestOptionOrFlagHelpNameLen(c.options, c.flags)

	// Required options
	if hasRequiredOptions {
		sb.WriteRune('\n')
		sb.WriteString("OPTIONS is one or more of:\n")

		for _, option := range c.requiredOptions {
			helpNameStyled, helpNameUnstyled := buildOptionOrFlagHelpName(option.Name, option.Alias)
			sb.WriteString(helpIndentationSpaces + helpNameStyled)

			if option.Description != "" {
				descripFormatted := breakStringIntoPaddedLines(
					helpIndentationNumSpaces+
						numSpacesHelpNameAndDescription+
						biggestOptionOrFlagHelpNameLen,
					' ',
					numCols,
					option.Description,
				)

				sb.Write([]byte(descripFormatted[len(helpNameUnstyled)+helpIndentationNumSpaces:]))
			}

			sb.WriteRune('\n')
		}
	}

	// Optional options
	if hasOptionalOptions {
		sb.WriteRune('\n')
		sb.WriteString("[OPTIONS] is one or more of:\n")

		for _, option := range c.options {
			if option.Required {
				continue
			}

			helpNameStyled, helpNameUnstyled := buildOptionOrFlagHelpName(option.Name, option.Alias)
			sb.WriteString(helpIndentationSpaces + helpNameStyled)

			if option.Description != "" {
				descripFormatted := breakStringIntoPaddedLines(
					helpIndentationNumSpaces+numSpacesHelpNameAndDescription+biggestOptionOrFlagHelpNameLen,
					' ',
					numCols,
					option.Description,
				)

				sb.Write([]byte(descripFormatted[len(helpNameUnstyled)+helpIndentationNumSpaces:]))
			}

			sb.WriteRune('\n')
		}
	}

	// Flags
	if hasFlags {
		sb.WriteRune('\n')
		sb.WriteString("[FLAGS] is one or more of:\n")

		for _, flag := range c.flags {
			helpNameStyled, helpNameUnstyled := buildOptionOrFlagHelpName(flag.Name, flag.Alias)
			sb.WriteString(helpIndentationSpaces + helpNameStyled)

			if flag.Description != "" {
				descripFormatted := breakStringIntoPaddedLines(
					helpIndentationNumSpaces+numSpacesHelpNameAndDescription+biggestOptionOrFlagHelpNameLen,
					' ',
					numCols,
					flag.Description,
				)

				sb.Write([]byte(descripFormatted[len(helpNameUnstyled)+helpIndentationNumSpaces:]))
			}

			sb.WriteRune('\n')
		}
	}

	return sb.String()
}
