// Contains functions to simplify checking the required pre-conditions before executing an operation
package diagnostics

import (
	"reflect"
	"strings"

	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

// Makes sure provided value is nil, otherwise it panics.
// value: the value to be checked
// valueName: Optional. Contains the name of the value the Assert was called on
// message: Optional. Contains the message to be displated if the Assert failed.
// If checks passed, same provided value will be returned
func IsNil(value interface{}, valueName, message string) interface{} {
	if value != nil {
		panic(getMessage(valueName, message))
	}

	return value
}

// Makes sure provided value is nil or empty, otherwise it panics.
// value: the value to be checked
// valueName: Optional. Contains the name of the value the Assert was called on
// message: Optional. Contains the message to be displated if the Assert failed.
// If checks passed, same provided value will be returned
func IsNilOrEmpty(value interface{}, valueName, message string) interface{} {
	if value != nil {
		switch val := value.(type) {
		case system.UUID:
			if val != system.EmptyUUID {
				panic(getMessage(valueName, message))
			}
		case string:
			if len(val) != 0 {
				panic(getMessage(valueName, message))
			}
		default:
			panic(getMessage(valueName, message))
		}
	}

	return value
}

// Makes sure provided value is nil, empty or contains whitespace only, otherwise it panics.
// value: the value to be checked
// valueName: Optional. Contains the name of the value the Assert was called on
// message: Optional. Contains the message to be displated if the Assert failed.
// If checks passed, same provided value will be returned
func IsNilOrEmptyOrWhitespace(value interface{}, valueName, message string) interface{} {
	if value != nil {
		switch val := value.(type) {
		case system.UUID:
			if val != system.EmptyUUID {
				panic(getMessage(valueName, message))
			}
		case string:
			if len(strings.TrimSpace(val)) != 0 {
				panic(getMessage(valueName, message))
			}
		default:
			panic(getMessage(valueName, message))
		}
	}

	return value
}

// Makes sure provided value is not nil, otherwise it panics.
// value: the value to be checked
// valueName: Optional. Contains the name of the value the Assert was called on
// message: Optional. Contains the message to be displated if the Assert failed.
// If checks passed, same provided value will be returned
func IsNotNil(value interface{}, valueName, message string) interface{} {
	if value == nil {
		panic(getMessage(valueName, message))
	} else {
		valueType := reflect.TypeOf(value).Kind()

		if valueType == reflect.UnsafePointer || valueType == reflect.Ptr {
			if reflect.ValueOf(value).IsNil() {
				panic(getMessage(valueName, message))
			}
		}
	}

	return value
}

// Makes sure provided value is not nil or empty, otherwise it panics.
// value: the value to be checked
// valueName: Optional. Contains the name of the value the Assert was called on
// message: Optional. Contains the message to be displated if the Assert failed.
// If checks passed, same provided value will be returned
func IsNotNilOrEmpty(value interface{}, valueName, message string) interface{} {
	IsNotNil(value, valueName, message)

	switch val := value.(type) {
	case system.UUID:
		if val == system.EmptyUUID {
			panic(getMessage(valueName, message))
		}
	case string:
		if len(val) == 0 {
			panic(getMessage(valueName, message))
		}
	default:
		panic(getMessage(valueName, message))
	}

	return value
}

// Makes sure provided value is not nil, empty or contains whitesapce only, otherwise it panics.
// value: the value to be checked
// valueName: Optional. Contains the name of the value the Assert was called on
// message: Optional. Contains the message to be displated if the Assert failed.
// If checks passed, same provided value will be returned
func IsNotNilOrEmptyOrWhitespace(value interface{}, valueName, message string) interface{} {
	IsNotNilOrEmpty(value, valueName, message)

	if val, ok := value.(string); ok {
		if len(strings.TrimSpace(val)) == 0 {
			panic(getMessage(valueName, message))
		}
	}

	return value
}

// Returns the message to be displayed when calling different Assert functions.
// valueName: Optional. Contains the name of value the Assert was called on.
// message: Optional. Contains the message to be displayed if the Assert failed.
// Remarks: Both valueName and message will be trimmed before running the checks
func getMessage(valueName, message string) string {
	valueName = strings.TrimSpace(valueName)
	message = strings.TrimSpace(message)

	if len(valueName) != 0 && len(message) != 0 {
		return message
	} else if len(valueName) == 0 && len(message) != 0 {
		return message
	} else if len(valueName) != 0 && len(message) == 0 {
		return valueName
	}

	return ""
}
