package main

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
	"github.com/hashicorp/hcl"
	"io/ioutil"
)

// ---------------------------------------------------------------------------------------
//  constants
// ---------------------------------------------------------------------------------------

const (
	DefaultSecretPath = "/run/secrets"
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

type Conf struct {
	SecretDir    string `hcl:"secret_dir"`
	SecretPrefix string `hcl:"secret_prefix"`

	Services []*Service `hcl:"service"`
}

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

// LoadConf loads a configuration file which is stored in
// the given filesystem path.
func LoadConf(path string) (*Conf, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := Conf{}
	err = hcl.Decode(&conf, string(buf))
	if err != nil {
		return nil, err
	}

	if conf.SecretDir == "" {
		conf.SecretDir = DefaultSecretPath
	}

	// if no secret prefix is configured per service
	// we should use the global one
	for _, service := range conf.Services {
		if service.SecretPrefix == "" {
			service.SecretPrefix = Config.SecretPrefix
		}
	}

	return &conf, nil
}
