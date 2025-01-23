package hw10programoptimization

import (
	"errors"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go" //nolint: depguard
)

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string `json:"email"`
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

// Using struct tags allows us to parse only relevant field(s)
// saving us time and memory.
// Package Jsoniter offers better performance than encoding/json for the given task.
func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	dStat := make(DomainStat)
	decoder := jsoniter.NewDecoder(r)
	user := &User{}
	var err error
	for {
		err = decoder.Decode(user)
		if errors.Is(err, io.EOF) {
			return dStat, nil
		}
		if err != nil {
			return nil, err
		}
		if strings.Contains(user.Email, "."+domain) {
			// Не читабельно и ужасно. Но позволяет отказаться от промежуточных переменных и сэкономить память.
			dStat[strings.ToLower(user.Email[strings.Index(user.Email, `@`)+1:])]++
		}
	}
}
