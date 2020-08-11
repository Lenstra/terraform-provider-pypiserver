TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
GOPATH = $(shell go env GOPATH)

.PHONY=build
build:
	GOBIN=$(GOPATH)/bin go install


.PHONY=testacc
testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

.PHONY=fmt
fmt:
	gofmt -w $(GOFMT_FILES)

.PHONY=fmtcheck
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"
