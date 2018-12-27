BUILD_TIME=$(shell date -Iminutes)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

all:
	go build -ldflags '-X "main.BuildTime='$(BUILD_TIME)'" -X "main.GitCommit='$(GIT_COMMIT)'" -X "main.GitBranch='$(GIT_BRANCH)'" -s -w' -v .
