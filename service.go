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
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/kballard/go-shellquote"
	"github.com/sirupsen/logrus"
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

type Service struct {
	Name          string `hcl:",key"`
	Command       string `hcl:"command"`
	CaptureStdOut bool   `hcl:"capture_stdout"`
	CaptureStdErr bool   `hcl:"capture_stderr"`
	IgnoreFailure bool   `hcl:"ignore_failure"`
	StopTimeout   string `hcl:"stop_timeout"`
	SecretPrefix  string `hcl:"secret_prefix"`

	// private members
	cmd *exec.Cmd
}

// ---------------------------------------------------------------------------------------
//  public members
// ---------------------------------------------------------------------------------------

// Spawn starts this service.
// wg is incremented when starting the task and decremented when the task is done
// failure gets a message if the tasks exists and IgnoreFailure is not set
// env are env variables which should be appended to this tasks env
func (s *Service) Spawn(wg *sync.WaitGroup, failure chan *Service, env []string) error {
	cmd, err := shellquote.Split(s.Command)
	if err != nil {
		return err
	}

	s.cmd = newCmd(cmd)
	s.cmd.Env = append(env, os.Environ()...)

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

	// start a watcher for task exit
	wg.Add(1)
	go func() {
		_ = s.cmd.Wait()

		// do the necessary signaling
		wg.Done()
		if !s.IgnoreFailure {
			failure <- s
		}

		logrus.Printf("service \"%s\" exited", s.Name)
	}()

	return nil
}

// Signal sends an os signal to this service.
func (s *Service) Signal(sig os.Signal) error {
	if s.cmd.ProcessState != nil && s.cmd.ProcessState.Exited() {
		return nil
	}

	return s.cmd.Process.Signal(sig)
}

// Shutdown gracefully shutsdown this service. Exit is not reported
// to the failure channel.
func (s *Service) Shutdown() error {
	// TODO : sigterm -> timeout -> sigkill

	// you cannot shutdown a services which is not running anymore
	if s.cmd.ProcessState != nil && s.cmd.ProcessState.Exited() {
		return nil
	}

	// stopping of the services is intented -> do not report
	// failures/exits and do not handle them.
	s.IgnoreFailure = true

	// send SIGTERM in order to ask the application
	// to shutdown gracefully
	err := s.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}

	return nil
}
