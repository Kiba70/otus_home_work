package hw02unpackstring

import (
	"errors"
	"testing"

	//nolint: depguard
	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "aaa0b3", expected: "aabbb"},
		{input: "\n3aaa0b3", expected: "\n\n\naabbb"},
		{input: "&3aaa0b3", expected: "&&&aabbb"},
		{input: "&1aaa0b0", expected: "&aa"},
		{input: "&1f3#5aa0b0", expected: "&fff#####a"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		// tc := tc // Не понятно, для чего была сделана данная конструкция с затенением
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}

	// Тест функции сына
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := UnpackFromMySon(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `e\fs`, `s3\n8`}

	for _, tc := range invalidStrings {
		// tc := tc // Не понятно, для чего была сделана данная конструкция с затенением
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}

	// Тест функции сына
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := UnpackFromMySon(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
