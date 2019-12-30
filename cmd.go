package cfp

import (
	"fmt"
	"strings"

	"github.com/efreitasn/customo"
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
	args          []interface{}
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
func (ct *CmdTermsSet) GetArgString(n int) string {
	arg := ct.cmd.getArg(n)
	if arg == nil || arg.T != TermString {
		return ""
	}

	return (ct.args[n]).(string)
}

// GetArgInt returns the value of the argument at n of type intger.
// If the argument doesn't exist, its zero value is returned.
func (ct *CmdTermsSet) GetArgInt(n int) int {
	arg := ct.cmd.getArg(n)
	if arg == nil || arg.T != TermInt {
		return 0
	}

	return (ct.args[n]).(int)
}

// GetArgFloat returns the value of the argument at n of type float.
// If the argument doesn't exist, its zero value is returned.
func (ct *CmdTermsSet) GetArgFloat(n int) float64 {
	arg := ct.cmd.getArg(n)
	if arg == nil || arg.T != TermFloat {
		return 0
	}

	return (ct.args[n]).(float64)
}

// CmdOption is a cmd option.
type CmdOption struct {
	// Name is used with -- and is case-insenstive.
	Name string
	// Alias is used with -.
	Alias       string
	Description string
	T           TermType
	Required    bool
}

// CmdFlag is a cmd flag.
type CmdFlag struct {
	// Name is used with -- and is case-insenstive.
	Name string
	// Alias is used with -.
	Alias       string
	Description string
}

// CmdArg is a cmd argument.
type CmdArg struct {
	// Name is used when printing the help message and is case-insenstive.
	Name        string
	Description string
	T           TermType
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
	args            []*CmdArg
}

// NewCmd creates a cmd.
// If an invalid flag name or option name is passed, it panics.
func NewCmd(cc CmdConfig) *Cmd {
	options := make(map[string]*CmdOption)
	requiredOptions := make([]*CmdOption, 0)
	optionsByAlias := make(map[string]*CmdOption)
	flags := make(map[string]*CmdFlag)
	flagsByAlias := make(map[string]*CmdFlag)
	args := make([]*CmdArg, 0, len(cc.Args))

	if cc.Options != nil {
		for i := range cc.Options {
			opt := cc.Options[i]

			if opt.Name == "" ||
				!isOptionWithoutValue("--"+opt.Name) ||
				(opt.Alias != "" && (strings.HasPrefix(opt.Alias, "-") || !isOptionWithoutValue("-"+opt.Alias))) {
				panic(ErrInvalidOptionNameOrAlias)
			}

			opt.Name = strings.ToLower(opt.Name)

			if opt.T == "" {
				panic(ErrMissingTermTypeForTerm{Term: opt.Name})
			}

			options[opt.Name] = &opt

			if opt.Required {
				requiredOptions = append(requiredOptions, &opt)
			}

			if opt.Alias != "" {
				opt.Alias = strings.ToLower(opt.Alias)
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

			flag.Name = strings.ToLower(flag.Name)
			flags[flag.Name] = &flag

			if flag.Alias != "" {
				flag.Alias = strings.ToLower(flag.Alias)
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

			arg.Name = strings.ToLower(arg.Name)

			if arg.T == "" {
				panic(ErrMissingTermTypeForTerm{Term: arg.Name})
			}

			args = append(args, &arg)
		}
	}

	return &Cmd{
		fn:              cc.Fn,
		options:         options,
		requiredOptions: requiredOptions,
		optionsByAlias:  optionsByAlias,
		flags:           flags,
		flagsByAlias:    flagsByAlias,
		args:            args,
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

func (c *Cmd) getArg(n int) *CmdArg {
	if len(c.args) >= (n + 1) {
		return c.args[n]
	}

	return nil
}

// Parse parses a slice of strings.
func (c *Cmd) Parse(pp parentParser, strs []string) error {
	tSet := &CmdTermsSet{
		cmd:           c,
		optionsValues: make(map[string]interface{}),
		args:          make([]interface{}, 0),
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
		nextArgPos := len(tSet.args)

		if len(c.args) >= (nextArgPos + 1) {
			arg := c.args[nextArgPos]
			argVal, valid := isValueValidForTermType(arg.T, str)

			if valid {
				tSet.args = append(tSet.args, argVal)

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

	if len(tSet.args) != len(c.args) {
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

	// Parent cmd's description.
	switch p := pp.parser.(type) {
	case *SubcmdsSet:
		if item := p.items[pp.cmds[len(pp.cmds)-1]]; item.Description != "" {
			sb.WriteString(item.Description + "\n\n")
		}
	case *rootCmd:
		if p.description != "" {
			sb.WriteString(p.description + "\n\n")
		}
	}

	sb.WriteString(fmt.Sprintf("Usage: %v", strings.Join(pp.cmds, " ")))

	hasArgs := c.args != nil && len(c.args) > 0
	hasRequiredOptions := c.requiredOptions != nil && len(c.requiredOptions) > 0
	hasOptionalOptions := c.options != nil && len(c.options) > 0 && len(c.requiredOptions) != len(c.options)
	hasFlags := c.flags != nil && len(c.flags) > 0

	largestArgHelpNameLen := 0

	if hasArgs {
		for _, arg := range c.args {
			helpName := "<" + arg.Name + ">"

			if len(helpName) > largestArgHelpNameLen {
				largestArgHelpNameLen = len(helpName)
			}

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

	sb.WriteString("\n")

	// Arguments
	sb.WriteString("\n")

	for _, arg := range c.args {
		argNameUnstyled := "<" + arg.Name + ">"
		argNameStyled := customo.Format(argNameUnstyled, customo.AttrBold)
		sb.WriteString(argNameStyled)

		if arg.Description != "" {
			descripFormatted := breakStringWithPadding(
				numberOfSpacesNameAndDescription+largestArgHelpNameLen,
				numCols,
				' ',
				arg.Description,
			)

			sb.Write([]byte(descripFormatted[len(argNameUnstyled):]))
		}

		sb.WriteString("\n")
	}

	// Largest option or flag.
	largestOptionOrFlagHelpNameUnstyledLen := 0

	for _, option := range c.options {
		_, helpNameUnstyled := buildOptionOrFlagHelpName(option.Name, option.Alias)
		if len(helpNameUnstyled) > largestOptionOrFlagHelpNameUnstyledLen {
			largestOptionOrFlagHelpNameUnstyledLen = len(helpNameUnstyled)
		}
	}

	for _, flag := range c.flags {
		_, helpNameUnstyled := buildOptionOrFlagHelpName(flag.Name, flag.Alias)
		if len(helpNameUnstyled) > largestOptionOrFlagHelpNameUnstyledLen {
			largestOptionOrFlagHelpNameUnstyledLen = len(helpNameUnstyled)
		}
	}

	// Required options
	if hasRequiredOptions {
		sb.WriteString("\n")
		sb.WriteString("OPTIONS is one or more of:\n")

		for _, option := range c.requiredOptions {
			helpNameStyled, helpNameUnstyled := buildOptionOrFlagHelpName(option.Name, option.Alias)
			sb.WriteString(helpIndentation + helpNameStyled)

			if option.Description != "" {
				descripFormatted := breakStringWithPadding(
					len(helpIndentation)+numberOfSpacesNameAndDescription+largestOptionOrFlagHelpNameUnstyledLen,
					numCols,
					' ',
					option.Description,
				)

				sb.Write([]byte(descripFormatted[len(helpNameUnstyled)+len(helpIndentation):]))
			}

			sb.WriteString("\n")
		}
	}

	// Optional options
	if hasOptionalOptions {
		sb.WriteString("\n")
		sb.WriteString("[OPTIONS] is one or more of:\n")

		for _, option := range c.options {
			if option.Required {
				continue
			}

			helpNameStyled, helpNameUnstyled := buildOptionOrFlagHelpName(option.Name, option.Alias)
			sb.WriteString(helpIndentation + helpNameStyled)

			if option.Description != "" {
				descripFormatted := breakStringWithPadding(
					len(helpIndentation)+numberOfSpacesNameAndDescription+largestOptionOrFlagHelpNameUnstyledLen,
					numCols,
					' ',
					option.Description,
				)

				sb.Write([]byte(descripFormatted[len(helpNameUnstyled)+len(helpIndentation):]))
			}

			sb.WriteString("\n")
		}
	}

	// Flags
	if hasFlags {
		sb.WriteString("\n")
		sb.WriteString("[FLAGS] is one or more of:\n")

		for _, flag := range c.flags {
			helpNameStyled, helpNameUnstyled := buildOptionOrFlagHelpName(flag.Name, flag.Alias)
			sb.WriteString(helpIndentation + helpNameStyled)

			if flag.Description != "" {
				descripFormatted := breakStringWithPadding(
					len(helpIndentation)+numberOfSpacesNameAndDescription+largestOptionOrFlagHelpNameUnstyledLen,
					numCols,
					' ',
					flag.Description,
				)

				sb.Write([]byte(descripFormatted[len(helpNameUnstyled)+len(helpIndentation):]))
			}

			sb.WriteString("\n")
		}
	}

	return sb.String()
}
