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

build: api-server worker admin s-isync

swagger:
	swagger validate pkg/swagger/swagger.yaml
	go generate github.com/dccn-tg/dr-data-stager/internal/api-server github.com/dccn-tg/dr-data-stager/pkg/swagger

doc: swagger
	swagger serve pkg/swagger/swagger.yaml

api-server:
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go build -a -installsuffix cgo -o build/data-stager-api internal/api-server/main.go

test_api-server:
	@GOPATH=$(GOPATH) GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go test -count=1 -v github.com/dccn-tg/dr-data-stager/internal/api-server/...

test_utility:
	@GOPATH=$(GOPATH) GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go test -count=1 -v github.com/dccn-tg/dr-data-stager/pkg/utility/...

test_worker:
	@GOPATH=$(GOPATH) GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go test -count=1 -v github.com/dccn-tg/dr-data-stager/internal/worker/...

worker:
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go build -a -installsuffix cgo -o build/data-stager-worker internal/worker/main.go

admin:
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go build -a -installsuffix cgo -o build/data-stager-admin internal/admin/main.go

s-isync:
	GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go build -a -installsuffix cgo -o build/s-isync internal/s-isync/*.go

test_crypto:
	@GOPATH=$(GOPATH) GOOS=$(GOOS) GO111MODULE=$(GO111MODULE) go test -v github.com/dccn-tg/dr-data-stager/pkg/utility/... -run TestRsaCryptography
