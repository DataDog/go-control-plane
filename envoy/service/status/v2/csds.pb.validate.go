//go:build !disable_pgv
// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/service/status/v2/csds.proto

package statusv2

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on ClientStatusRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ClientStatusRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ClientStatusRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ClientStatusRequestMultiError, or nil if none found.
func (m *ClientStatusRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ClientStatusRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetNodeMatchers() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ClientStatusRequestValidationError{
						field:  fmt.Sprintf("NodeMatchers[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ClientStatusRequestValidationError{
						field:  fmt.Sprintf("NodeMatchers[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ClientStatusRequestValidationError{
					field:  fmt.Sprintf("NodeMatchers[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ClientStatusRequestMultiError(errors)
	}

	return nil
}

// ClientStatusRequestMultiError is an error wrapping multiple validation
// errors returned by ClientStatusRequest.ValidateAll() if the designated
// constraints aren't met.
type ClientStatusRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ClientStatusRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ClientStatusRequestMultiError) AllErrors() []error { return m }

// ClientStatusRequestValidationError is the validation error returned by
// ClientStatusRequest.Validate if the designated constraints aren't met.
type ClientStatusRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClientStatusRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClientStatusRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ClientStatusRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClientStatusRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClientStatusRequestValidationError) ErrorName() string {
	return "ClientStatusRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ClientStatusRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClientStatusRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClientStatusRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClientStatusRequestValidationError{}

// Validate checks the field values on PerXdsConfig with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PerXdsConfig) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PerXdsConfig with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PerXdsConfigMultiError, or
// nil if none found.
func (m *PerXdsConfig) ValidateAll() error {
	return m.validate(true)
}

func (m *PerXdsConfig) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Status

	switch v := m.PerXdsConfig.(type) {
	case *PerXdsConfig_ListenerConfig:
		if v == nil {
			err := PerXdsConfigValidationError{
				field:  "PerXdsConfig",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetListenerConfig()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PerXdsConfigValidationError{
						field:  "ListenerConfig",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PerXdsConfigValidationError{
						field:  "ListenerConfig",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetListenerConfig()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PerXdsConfigValidationError{
					field:  "ListenerConfig",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *PerXdsConfig_ClusterConfig:
		if v == nil {
			err := PerXdsConfigValidationError{
				field:  "PerXdsConfig",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetClusterConfig()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PerXdsConfigValidationError{
						field:  "ClusterConfig",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PerXdsConfigValidationError{
						field:  "ClusterConfig",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetClusterConfig()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PerXdsConfigValidationError{
					field:  "ClusterConfig",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *PerXdsConfig_RouteConfig:
		if v == nil {
			err := PerXdsConfigValidationError{
				field:  "PerXdsConfig",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetRouteConfig()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PerXdsConfigValidationError{
						field:  "RouteConfig",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PerXdsConfigValidationError{
						field:  "RouteConfig",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetRouteConfig()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PerXdsConfigValidationError{
					field:  "RouteConfig",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *PerXdsConfig_ScopedRouteConfig:
		if v == nil {
			err := PerXdsConfigValidationError{
				field:  "PerXdsConfig",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetScopedRouteConfig()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PerXdsConfigValidationError{
						field:  "ScopedRouteConfig",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PerXdsConfigValidationError{
						field:  "ScopedRouteConfig",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetScopedRouteConfig()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PerXdsConfigValidationError{
					field:  "ScopedRouteConfig",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		_ = v // ensures v is used
	}

	if len(errors) > 0 {
		return PerXdsConfigMultiError(errors)
	}

	return nil
}

// PerXdsConfigMultiError is an error wrapping multiple validation errors
// returned by PerXdsConfig.ValidateAll() if the designated constraints aren't met.
type PerXdsConfigMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PerXdsConfigMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PerXdsConfigMultiError) AllErrors() []error { return m }

// PerXdsConfigValidationError is the validation error returned by
// PerXdsConfig.Validate if the designated constraints aren't met.
type PerXdsConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PerXdsConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PerXdsConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PerXdsConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PerXdsConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PerXdsConfigValidationError) ErrorName() string { return "PerXdsConfigValidationError" }

// Error satisfies the builtin error interface
func (e PerXdsConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPerXdsConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PerXdsConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PerXdsConfigValidationError{}

// Validate checks the field values on ClientConfig with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ClientConfig) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ClientConfig with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ClientConfigMultiError, or
// nil if none found.
func (m *ClientConfig) ValidateAll() error {
	return m.validate(true)
}

func (m *ClientConfig) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetNode()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ClientConfigValidationError{
					field:  "Node",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ClientConfigValidationError{
					field:  "Node",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetNode()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ClientConfigValidationError{
				field:  "Node",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetXdsConfig() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ClientConfigValidationError{
						field:  fmt.Sprintf("XdsConfig[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ClientConfigValidationError{
						field:  fmt.Sprintf("XdsConfig[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ClientConfigValidationError{
					field:  fmt.Sprintf("XdsConfig[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ClientConfigMultiError(errors)
	}

	return nil
}

// ClientConfigMultiError is an error wrapping multiple validation errors
// returned by ClientConfig.ValidateAll() if the designated constraints aren't met.
type ClientConfigMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ClientConfigMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ClientConfigMultiError) AllErrors() []error { return m }

// ClientConfigValidationError is the validation error returned by
// ClientConfig.Validate if the designated constraints aren't met.
type ClientConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClientConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClientConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ClientConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClientConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClientConfigValidationError) ErrorName() string { return "ClientConfigValidationError" }

// Error satisfies the builtin error interface
func (e ClientConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClientConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClientConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClientConfigValidationError{}

// Validate checks the field values on ClientStatusResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ClientStatusResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ClientStatusResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ClientStatusResponseMultiError, or nil if none found.
func (m *ClientStatusResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ClientStatusResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetConfig() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ClientStatusResponseValidationError{
						field:  fmt.Sprintf("Config[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ClientStatusResponseValidationError{
						field:  fmt.Sprintf("Config[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ClientStatusResponseValidationError{
					field:  fmt.Sprintf("Config[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ClientStatusResponseMultiError(errors)
	}

	return nil
}

// ClientStatusResponseMultiError is an error wrapping multiple validation
// errors returned by ClientStatusResponse.ValidateAll() if the designated
// constraints aren't met.
type ClientStatusResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ClientStatusResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ClientStatusResponseMultiError) AllErrors() []error { return m }

// ClientStatusResponseValidationError is the validation error returned by
// ClientStatusResponse.Validate if the designated constraints aren't met.
type ClientStatusResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClientStatusResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClientStatusResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ClientStatusResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClientStatusResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClientStatusResponseValidationError) ErrorName() string {
	return "ClientStatusResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ClientStatusResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClientStatusResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClientStatusResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClientStatusResponseValidationError{}
