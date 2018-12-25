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
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

type Environment []string

// ---------------------------------------------------------------------------------------
//  public members
// ---------------------------------------------------------------------------------------

func (e Environment) WithPrefix(prefix string) Environment {
	prefix = strings.ToUpper(prefix)
	env := make(Environment, 0)
	if prefix == "-" {
		return env
	}

	for _, key := range e {
		if strings.HasPrefix(key, prefix) {
			env = append(env, key)
		}
	}

	return env
}
