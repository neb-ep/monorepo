// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: protos/auth/v1/auth.proto

package authv1

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

// Validate checks the field values on SignUpRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SignUpRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SignUpRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SignUpRequestMultiError, or
// nil if none found.
func (m *SignUpRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *SignUpRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetUsername()); l < 6 || l > 255 {
		err := SignUpRequestValidationError{
			field:  "Username",
			reason: "value length must be between 6 and 255 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetFirstName()); l < 6 || l > 255 {
		err := SignUpRequestValidationError{
			field:  "FirstName",
			reason: "value length must be between 6 and 255 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetLastName()); l < 6 || l > 255 {
		err := SignUpRequestValidationError{
			field:  "LastName",
			reason: "value length must be between 6 and 255 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateEmail(m.GetEmail()); err != nil {
		err = SignUpRequestValidationError{
			field:  "Email",
			reason: "value must be a valid email address",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetPassword()); l < 8 || l > 64 {
		err := SignUpRequestValidationError{
			field:  "Password",
			reason: "value length must be between 8 and 64 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return SignUpRequestMultiError(errors)
	}

	return nil
}

func (m *SignUpRequest) _validateHostname(host string) error {
	s := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(host) > 253 {
		return errors.New("hostname cannot exceed 253 characters")
	}

	for _, part := range strings.Split(s, ".") {
		if l := len(part); l == 0 || l > 63 {
			return errors.New("hostname part must be non-empty and cannot exceed 63 characters")
		}

		if part[0] == '-' {
			return errors.New("hostname parts cannot begin with hyphens")
		}

		if part[len(part)-1] == '-' {
			return errors.New("hostname parts cannot end with hyphens")
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return fmt.Errorf("hostname parts can only contain alphanumeric characters or hyphens, got %q", string(r))
			}
		}
	}

	return nil
}

func (m *SignUpRequest) _validateEmail(addr string) error {
	a, err := mail.ParseAddress(addr)
	if err != nil {
		return err
	}
	addr = a.Address

	if len(addr) > 254 {
		return errors.New("email addresses cannot exceed 254 characters")
	}

	parts := strings.SplitN(addr, "@", 2)

	if len(parts[0]) > 64 {
		return errors.New("email address local phrase cannot exceed 64 characters")
	}

	return m._validateHostname(parts[1])
}

// SignUpRequestMultiError is an error wrapping multiple validation errors
// returned by SignUpRequest.ValidateAll() if the designated constraints
// aren't met.
type SignUpRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SignUpRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SignUpRequestMultiError) AllErrors() []error { return m }

// SignUpRequestValidationError is the validation error returned by
// SignUpRequest.Validate if the designated constraints aren't met.
type SignUpRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SignUpRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SignUpRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SignUpRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SignUpRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SignUpRequestValidationError) ErrorName() string { return "SignUpRequestValidationError" }

// Error satisfies the builtin error interface
func (e SignUpRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSignUpRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SignUpRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SignUpRequestValidationError{}

// Validate checks the field values on SignUpResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SignUpResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SignUpResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SignUpResponseMultiError,
// or nil if none found.
func (m *SignUpResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *SignUpResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Username

	// no validation rules for FirstName

	// no validation rules for LastName

	// no validation rules for Email

	if len(errors) > 0 {
		return SignUpResponseMultiError(errors)
	}

	return nil
}

// SignUpResponseMultiError is an error wrapping multiple validation errors
// returned by SignUpResponse.ValidateAll() if the designated constraints
// aren't met.
type SignUpResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SignUpResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SignUpResponseMultiError) AllErrors() []error { return m }

// SignUpResponseValidationError is the validation error returned by
// SignUpResponse.Validate if the designated constraints aren't met.
type SignUpResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SignUpResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SignUpResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SignUpResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SignUpResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SignUpResponseValidationError) ErrorName() string { return "SignUpResponseValidationError" }

// Error satisfies the builtin error interface
func (e SignUpResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSignUpResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SignUpResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SignUpResponseValidationError{}

// Validate checks the field values on SignInRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SignInRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SignInRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SignInRequestMultiError, or
// nil if none found.
func (m *SignInRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *SignInRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Email

	// no validation rules for Password

	if len(errors) > 0 {
		return SignInRequestMultiError(errors)
	}

	return nil
}

// SignInRequestMultiError is an error wrapping multiple validation errors
// returned by SignInRequest.ValidateAll() if the designated constraints
// aren't met.
type SignInRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SignInRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SignInRequestMultiError) AllErrors() []error { return m }

// SignInRequestValidationError is the validation error returned by
// SignInRequest.Validate if the designated constraints aren't met.
type SignInRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SignInRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SignInRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SignInRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SignInRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SignInRequestValidationError) ErrorName() string { return "SignInRequestValidationError" }

// Error satisfies the builtin error interface
func (e SignInRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSignInRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SignInRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SignInRequestValidationError{}

// Validate checks the field values on SignInResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SignInResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SignInResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SignInResponseMultiError,
// or nil if none found.
func (m *SignInResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *SignInResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccessToken

	// no validation rules for RefreshToken

	if len(errors) > 0 {
		return SignInResponseMultiError(errors)
	}

	return nil
}

// SignInResponseMultiError is an error wrapping multiple validation errors
// returned by SignInResponse.ValidateAll() if the designated constraints
// aren't met.
type SignInResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SignInResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SignInResponseMultiError) AllErrors() []error { return m }

// SignInResponseValidationError is the validation error returned by
// SignInResponse.Validate if the designated constraints aren't met.
type SignInResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SignInResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SignInResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SignInResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SignInResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SignInResponseValidationError) ErrorName() string { return "SignInResponseValidationError" }

// Error satisfies the builtin error interface
func (e SignInResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSignInResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SignInResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SignInResponseValidationError{}
