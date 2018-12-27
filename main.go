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
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/faryon93/uprun/secrets"
)

// ---------------------------------------------------------------------------------------
//  variables
// ---------------------------------------------------------------------------------------

var (
	ForceColors bool
	ConfigPath  string

	Config *Conf
)

// ---------------------------------------------------------------------------------------
//  application entry
// ---------------------------------------------------------------------------------------

func main() {
	var err error

	flag.BoolVar(&ForceColors, "colors", false, "force logging with colors")
	flag.StringVar(&ConfigPath, "conf", "uprun.chl", "path to config file")
	flag.Parse()

	// setup logger
	formater := logrus.TextFormatter{ForceColors: ForceColors}
	logrus.SetFormatter(&formater)
	logrus.SetOutput(os.Stdout)

	// load the configuration file
	Config, err = LoadConf(ConfigPath)
	if err != nil {
		logrus.Errorln("failed to read config file:", err.Error())
		os.Exit(-1)
	}

	failure := make(chan *Service)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range signals {
			logrus.Infoln("uprun received signal", sig, ": initiating shutdown")
			if sig == syscall.SIGTERM || sig == syscall.SIGINT {
				for _, service := range Config.Services {
					service.IgnoreFailure = true

					// TODO : sigterm -> timeout -> sigkill
					err := service.Signal(syscall.SIGTERM)
					if err != nil {
						logrus.Errorln("failed to signal service", service.Name, err.Error())
					}
				}

			} else {
				// forward the signal to all services
				for _, service := range Config.Services {
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

			for _, service := range Config.Services {
				err := service.Signal(syscall.SIGTERM)
				if err != nil {
					logrus.Errorln("failed to signal service", service.Name, err.Error())
				}
			}

		}
	}()

	// gather all the secrets from filesystem
	exportedSecrets, err := secrets.Export(Config.SecretDir)
	if err != nil {
		logrus.Errorln("failed to export secrets:", err.Error())
	}

	// all configured services should be started one at a time
	wg := sync.WaitGroup{}
	for _, service := range Config.Services {
		logrus.Infoln("starting service", service.Name)

		// filter all secerts for this service
		secretsEnv := exportedSecrets.WithPrefix(service.SecretPrefix)
		logrus.Debugln("exported secrets:", secretsEnv)

		// spawn the service
		err := service.Spawn(&wg, failure, secretsEnv)
		if err != nil {
			logrus.Errorln("failed to spawn server", service.Name, err.Error())
			continue
		}
	}

	// wait until all services have been terminated
	wg.Wait()
	logrus.Println("all services have ended: exiting uprun")
}
