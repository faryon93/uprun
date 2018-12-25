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
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// ---------------------------------------------------------------------------------------
//  application entry
// ---------------------------------------------------------------------------------------

func main() {
	buf, err := ioutil.ReadFile("uprun.example.hcl")
	if err != nil {
		panic(err)
	}

	conf := Conf{}
	err = hcl.Decode(&conf, string(buf))
	if err != nil {
		panic(err)
	}

	failure := make(chan *Service)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range signals {
			logrus.Infoln("uprun received signal", sig, ": initiating shutdown")
			if sig == syscall.SIGTERM || sig == syscall.SIGINT {
				for _, service := range conf.Services {
					service.IgnoreFailure = true

					// TODO : sigterm -> timeout -> sigkill
					err := service.Signal(syscall.SIGTERM)
					if err != nil {
						logrus.Errorln("failed to signal service", service.Name, err.Error())
					}
				}

			} else {
				// forward the signal to all services
				for _, service := range conf.Services {
					err := service.Signal(sig)
					if err != nil {
						logrus.Errorln("failed to signal service", service.Name, err.Error())
					}
				}
			}
		}
	}()

	go func() {
		for svc := range failure {
			logrus.Warnln("service", svc.Name, "failed: shutting down everything")

			for _, service := range conf.Services {
				err := service.Signal(syscall.SIGTERM)
				if err != nil {
					logrus.Errorln("failed to signal service", service.Name, err.Error())
				}
			}

		}
	}()

	wg := sync.WaitGroup{}
	for _, service := range conf.Services {
		logrus.Infoln("starting service", service.Name)
		err := service.Spawn(&wg, failure)
		if err != nil {
			logrus.Errorln("failed to spawn server", service.Name, err.Error())
			continue
		}
	}

	wg.Wait()

	logrus.Println("all services have ended: exiting uprun")
}
