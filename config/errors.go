package config

import (
	"fmt"
)

// ConfigError is an error type returned upon configuration errors.
type ConfigError struct {
	err error
}

// ConfigErrorFromString returns a ConfigError from the given error string.
func ConfigErrorFromString(format string, args ...interface{}) *ConfigError {
	return &ConfigError{
		err: fmt.Errorf(format, args...),
	}
}

// ConfigErrorFromError returns a ConfigError from the given error object.
func ConfigErrorFromError(err error) *ConfigError {
	return &ConfigError{
		err: err,
	}
}

func (ce ConfigError) Error() string {
	return fmt.Sprintf("error parsing config: %v", ce.err)
}
