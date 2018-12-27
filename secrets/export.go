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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

// Export reads all secrets in a directory.
func Export(path string) (Environment, error) {
	env := make([]string, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		buf, err := ioutil.ReadFile(path)
		if err != nil {
			logrus.Errorln("failed to read secret:", err.Error())
			return nil
		}

		env = append(env, strings.ToUpper(info.Name())+"="+string(buf))

		return nil
	})
	if err != nil {
		return make([]string, 0), err
	}

	return env, nil
}
