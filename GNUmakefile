TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=hashicorp.com
NAMESPACE=sh
NAME=sitehost
BINARY=terraform-provider-${NAME}
VERSION=1.2.0
OS_ARCH=linux_amd64
SRC := go.sum $(shell git ls-files -cmo --exclude-standard -- "*.go")
TESTABLE := ./...

default: install

build:
	go build -o ${BINARY}

release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -parallel 10 -short

bin/golangci-lint: GOARCH =
bin/golangci-lint: GOOS =
bin/golangci-lint: go.sum
	@go build -o $@ github.com/golangci/golangci-lint/cmd/golangci-lint

bin/go-acc: GOARCH =
bin/go-acc: GOOS =
bin/go-acc: go.sum
	@go build -o $@ github.com/ory/go-acc

lint: CGO_ENABLED = 1
lint: GOARCH =
lint: GOOS =
lint: bin/golangci-lint $(SRC)
	$< run --max-same-issues 50

tidy:
	go mod tidy

dirty: tidy
	git status --porcelain
	@[ -z "$$(git status --porcelain)" ]

vet: GOARCH =
vet: GOOS =
vet: CGO_ENABLED =
vet: bin/go-acc $(SRC)
	$< --covermode=atomic $(TESTABLE) -- -race -v