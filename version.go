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
//  global variables
// ---------------------------------------------------------------------------------------

var (
	// release information
	AppName    = "uprun"
	AppVersion = "0.1.0"

	// filled by build tool
	GitCommit   string
	GitBranch   string
	BuildTime   string
	BuildNumber string
)

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

// Returns the application version string.
func GetAppVersion() string {
	str := AppName + " " + AppVersion
	if len(BuildNumber) > 0 {
		str += "-" + BuildNumber
	}

	if len(GitCommit) > 0 {
		str += " (#" + GitCommit
	}

	if len(GitBranch) > 0 {
		str += "-" + GitBranch
	}

	if len(BuildTime) > 0 {
		str += " / " + BuildTime + ")"
	} else {
		if len(GitCommit) > 0 {
			str += ")"
		}
	}

	return str
}
