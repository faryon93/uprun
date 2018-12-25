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
	"github.com/kballard/go-shellquote"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

type Conf struct {
	Services []*Service `hcl:"service"`
}

type Service struct {
	Name string `hcl:",key"`

	Command       string `hcl:"command"`
	CaptureStdOut bool   `hcl:"capture_stdout"`
	CaptureStdErr bool   `hcl:"capture_stderr"`
	IgnoreFailure bool   `hcl:"ignore_failure"`
	StopTimeout   string `hcl:"stop_timeout"`

	cmd *exec.Cmd
}

// ---------------------------------------------------------------------------------------
//  public members
// ---------------------------------------------------------------------------------------

func (s *Service) Spawn(wg *sync.WaitGroup, failure chan *Service) error {
	cmd, err := shellquote.Split(s.Command)
	if err != nil {
		return err
	}

	s.cmd = exec.Command(cmd[0], cmd[1:]...)
	s.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if s.CaptureStdOut {
		s.cmd.Stdout = os.Stdout
	}

	if s.CaptureStdErr {
		s.cmd.Stderr = os.Stdout
	}

	err = s.cmd.Start()
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		s.cmd.Wait()
		logrus.Printf("service %s exited", s.Name)

		wg.Done()

		if !s.IgnoreFailure {
			failure <- s
		}
	}()

	return nil
}

func (s *Service) Signal(sig os.Signal) error {
	if s.cmd.ProcessState != nil && s.cmd.ProcessState.Exited() {
		return nil
	}

	return s.cmd.Process.Signal(sig)
}
