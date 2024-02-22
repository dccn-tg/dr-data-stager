ifndef GOPATH
	GOPATH := $(HOME)/go
endif

ifndef GOOS
	GOOS := linux
endif

ifndef GO111MODULE
	GO111MODULE := on
endif

all: build

build: api-server worker s-isync

swagger:
	swagger validate pkg/swagger/swagger.yaml
	go generate github.com/Donders-Institute/dr-data-stager/internal/api-server github.com/Donders-Institute/dr-data-stager/pkg/swagger

doc: swagger
	swagger serve pkg/swagger/swagger.yaml

api-server:
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go build -a -installsuffix cgo -o bin/data-stager-api internal/api-server/main.go

worker:
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go build -a -installsuffix cgo -o bin/data-stager-worker internal/worker/main.go

s-isync:
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go build -a -installsuffix cgo -o bin/s-isync internal/s-isync/main.go
