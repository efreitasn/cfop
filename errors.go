package cfop

import (
	"errors"
	"fmt"
)

// ErrInvalidOptionNameOrAlias indicates that an invalid option name or invalid option alias was provided.
var ErrInvalidOptionNameOrAlias = errors.New("cfop: invalid option name or option alias")

// ErrInvalidFlagNameOrAlias indicates that an invalid flag name or invalid flag alias was provided.
var ErrInvalidFlagNameOrAlias = errors.New("cfop: invalid flag name or flag alias")

// ErrInvalidTermType indicates that an invalid term type was provided.
var ErrInvalidTermType = errors.New("cfop: invalid term type")

// ErrMissingRootCmdName indicates that a name for the root cmd wasn't provided.
var ErrMissingRootCmdName = errors.New("cfop: missing name for root cmd")

// ErrMissingSubcmdName indicates that a name for a subcmd wasn't provided.
var ErrMissingSubcmdName = errors.New("cfop: missing name for subcmd")

// ErrMissingSubcmdParser indicates that a parser for a subcmd wasn't provided.
var ErrMissingSubcmdParser = errors.New("cfop: missing parser for subcmd")

// ErrMissingCmdFn indicates that a function for a cmd wasn't provided.
var ErrMissingCmdFn = errors.New("cfop: missing function for cmd")

// ErrMissingTermTypeForTerm indicates that a term's type wasn't provided.
type ErrMissingTermTypeForTerm struct {
	Term string
}

func (e ErrMissingTermTypeForTerm) Error() string {
	return fmt.Sprintf("cfop: type not provided for term: %v", e.Term)
}

// ErrInvalidArgumentName indicates that an invalid argument name was provided.
type ErrInvalidArgumentName struct {
	ArgumentPos int
}

func (e ErrInvalidArgumentName) Error() string {
	return fmt.Sprintf("cfop: name of argument at %v is invalid", e.ArgumentPos)
}

// The errors below are those that are shown to the user.

// ErrUnexpectedOption indicates that an unexpected option or flag was provided.
type ErrUnexpectedOption struct {
	OptionName string
	IsAlias    bool
}

func (e ErrUnexpectedOption) Error() string {
	if e.IsAlias {
		return fmt.Sprintf("unexpected -%v option", e.OptionName)
	}

	return fmt.Sprintf("unexpected --%v option", e.OptionName)
}

// ErrUnexpectedOptionOrFlag indicates that an unexpected option or flag was provided.
type ErrUnexpectedOptionOrFlag struct {
	OptionOrFlagName string
	IsAlias          bool
}

func (e ErrUnexpectedOptionOrFlag) Error() string {
	if e.IsAlias {
		return fmt.Sprintf("unexpected -%v option/flag", e.OptionOrFlagName)
	}

	return fmt.Sprintf("unexpected --%v option/flag", e.OptionOrFlagName)
}

// ErrOptionExpectsDifferentValueType indicates that an option expects a value of a type different than the one provided.
type ErrOptionExpectsDifferentValueType struct {
	OptionName   string
	ExpectedType TermType
	IsAlias      bool
}

func (e ErrOptionExpectsDifferentValueType) Error() string {
	if e.IsAlias {
		return fmt.Sprintf("-%v option expects a value of type %v", e.OptionName, e.ExpectedType)
	}

	return fmt.Sprintf("--%v option expects a value of type %v", e.OptionName, e.ExpectedType)
}

// ErrOptionsExpectsAValue indicates that an option expects a value, but one wasn't provided.
type ErrOptionsExpectsAValue struct {
	OptionName string
	IsAlias    bool
}

func (e ErrOptionsExpectsAValue) Error() string {
	if e.IsAlias {
		return fmt.Sprintf("-%v option expects a value", e.OptionName)
	}

	return fmt.Sprintf("--%v option expects a value", e.OptionName)
}

// ErrOptionsExpectsAType indicates that an option expects a type, but one wasn't provided.
type ErrOptionsExpectsAType struct {
	OptionName string
	IsAlias    bool
}

func (e ErrOptionsExpectsAType) Error() string {
	if e.IsAlias {
		return fmt.Sprintf("-%v option expects a value", e.OptionName)
	}

	return fmt.Sprintf("--%v option expects a value", e.OptionName)
}

// ErrUnexpectedArgument indicates that an unexpected argument was provided.
type ErrUnexpectedArgument struct {
	Argument string
}

func (e ErrUnexpectedArgument) Error() string {
	return fmt.Sprintf("unexpected argument: %v", e.Argument)
}

// ErrArgumentExpectsDifferentValueType indicates that an argument expects a value of a type different than the one provided.
type ErrArgumentExpectsDifferentValueType struct {
	ArgumentPos  int
	ArgumentName string
	ExpectedType TermType
	Value        string
}

func (e ErrArgumentExpectsDifferentValueType) Error() string {
	return fmt.Sprintf("the <%v> argument (%v) expects a value of type %v", e.ArgumentName, e.Value, e.ExpectedType)
}

// ErrMissingArguments indicates that not all arguments were provided.
var ErrMissingArguments = errors.New("missing argument(s)")

// ErrRequiredOptionNotProvided indicates that a required option wasn't provided.
type ErrRequiredOptionNotProvided struct {
	OptionName string
}

func (e ErrRequiredOptionNotProvided) Error() string {
	return fmt.Sprintf("--%v option is required", e.OptionName)
}

// ErrMissingSubcmd indicates that a subcmd wasn't provided.
var ErrMissingSubcmd = errors.New("missing subcmd")

// ErrUnknownSubcmd indicates that an unknown subcmd was provided.
type ErrUnknownSubcmd struct {
	SubcmdName string
}

func (e ErrUnknownSubcmd) Error() string {
	return fmt.Sprintf("unknown subcmd: %v", e.SubcmdName)
}
