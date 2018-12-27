package secrets

// uprun
// Copyright (C) 2018 Maximilian Pachl

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// ---------------------------------------------------------------------------------------
//  imports
// ---------------------------------------------------------------------------------------

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------------------
//  variables
// ---------------------------------------------------------------------------------------

var (
	env = Environment{"FOO=BAR", "TEST=BLUB", "MY_VAR=VARIABLE"}
)

// ---------------------------------------------------------------------------------------
//  tests
// ---------------------------------------------------------------------------------------

func TestEnvironment_WithPrefix(t *testing.T) {
	prefix := env[0][0:1]
	for _, val := range env.WithPrefix(prefix) {
		if !strings.HasPrefix(val, prefix) {
			t.Error("WithPrefix() returned non prefix value", val)
		}
	}
}

func TestEnvironment_WithPrefixEmpty(t *testing.T) {
	allEnv := env.WithPrefix("")
	for i := range env {
		if env[i] != allEnv[i] {
			t.Error("WithPrefix(\"\") does not match the original env")
		}
	}
}

func TestEnvironment_WithPrefixNone(t *testing.T) {
	noEnv := env.WithPrefix("-")
	if len(noEnv) > 0 {
		t.Error("WithPrefix(\"-\") should return empty env")
	}
}

func TestEnvironment_WithPrefixCase(t *testing.T) {
	prefix := env[0][0:1]

	for _, val := range env.WithPrefix(strings.ToLower(prefix)) {
		if !strings.HasPrefix(val, strings.ToUpper(prefix)) {
			t.Error("WithPrefix() should not be case sensitive", val)
		}
	}
}
