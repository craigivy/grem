# make file for sky
OS := $(shell uname)

.DEFAULT: default
default: clean protos install fmt lint vet testv

int: clean protos install fmt lint vet int_test

# Must be run with 'sudo -E'
install_env: install_lint install_gotests install_protobuf

.PHONY: install_lint
install_lint:
	@echo "Running: $@"

	@go get -u github.com/golang/lint/golint
	@cd $$GOPATH/src/github.com/golang/lint/golint && go install .

.PHONY: clean
clean:
	@echo "Running: $@"
	@find . -name "*.pb.go" -type f -not -path "./vendor/*" -delete

.PHONY: protos
protos:
	@echo "Running: $@"
	@./compile_protos

.PHONY: lint
lint:
	@echo "Running: $@"
#   _sting.go is generated by go's stringer.
	@test -z "$$(golint ./... | grep -v vendor/ | grep -v "_string.go" | tee /dev/stderr)"


.PHONY: fmt
fmt:
	@echo "Running: $@"
	@test -z "$$(gofmt -l . |  grep -v vendor/ | tee /dev/stderr)"  || (echo "'gofmt' required on the above files" && false)

.PHONY: fmtv
fmtv:
	@echo "Running: $@"
	@test -z "$$(gofmt -d . |  grep -v vendor/ | tee /dev/stderr)"  || (echo "'gofmt' required on the above files" && false)

.PHONY: install
install:
	@echo "Running: $@"
	@go install -a ./...

.PHONY: vet
vet:
	@echo "Running: $@"
	@go list ./... | grep -v /vendor/ | xargs -L1 go vet

.PHONY: test
test:
	@echo "Running: $@"
	@go test -timeout 1m ./...

.PHONY: testv
testv:
	@echo "Running: $@"
	go test -timeout 1m ./... -cover

.PHONY: int_test
int_test:
	@echo "Running: $@"
	go test -tags="int" -timeout 1m ./... -v -cover
