package cfp

import (
	"errors"
	"fmt"
)

// ErrInvalidOptionNameOrAlias indicates that an invalid option name or invalid option alias was provided.
var ErrInvalidOptionNameOrAlias = errors.New("cfp: invalid option name or option alias")

// ErrInvalidFlagNameOrAlias indicates that an invalid flag name or invalid flag alias was provided.
var ErrInvalidFlagNameOrAlias = errors.New("cfp: invalid flag name or flag alias")

// ErrInvalidTermType indicates that an invalid term type was provided.
var ErrInvalidTermType = errors.New("cfp: invalid term type")

// ErrUnexpectedOptionOrFlag indicates that an invalid option or flag was provided.
type ErrUnexpectedOptionOrFlag struct {
	OptionOrFlagName string
}

func (e ErrUnexpectedOptionOrFlag) Error() string {
	return fmt.Sprintf("cfp: %v is not a valid option or flag", e.OptionOrFlagName)
}

// ErrInvalidValueForOption indicates that an invalid value was provided to an option.
type ErrInvalidValueForOption struct {
	OptionName string
	Value      string
}

func (e ErrInvalidValueForOption) Error() string {
	return fmt.Sprintf("cfp: %v is an invalid value for %v option", e.Value, e.OptionName)
}

// ErrValueNotProvidedForOption indicates that an option is passed but its value isn't.
type ErrValueNotProvidedForOption struct {
	OptionName string
}

func (e ErrValueNotProvidedForOption) Error() string {
	return fmt.Sprintf("cfp: value not provided for %v option", e.OptionName)
}

// ErrUnexpectedArgument indicates that an unexpected argument was provided.
type ErrUnexpectedArgument struct {
	Argument string
}

func (e ErrUnexpectedArgument) Error() string {
	return fmt.Sprintf("cfp: unexpected argument: %v", e.Argument)
}

// ErrInvalidValueForArgument indicates that an invalid value was provided to an argument.
type ErrInvalidValueForArgument struct {
	ArgumentPos int
	Value       string
}

func (e ErrInvalidValueForArgument) Error() string {
	return fmt.Sprintf("cfp: %v is an invalid value for argument at %v", e.Value, e.ArgumentPos)
}

// ErrMissingArguments indicates that not all arguments were provided.
var ErrMissingArguments = errors.New("cfp: there are missing arguments")

// ErrRequiredOptionNotProvided indicates that a required option wasn't provided.
type ErrRequiredOptionNotProvided struct {
	OptionName string
}

func (e ErrRequiredOptionNotProvided) Error() string {
	return fmt.Sprintf("cfp: required option %v not provided", e.OptionName)
}
