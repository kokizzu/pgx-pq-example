GO ?= go
GOVULNCHECK ?= govulncheck
CMD ?=

.PHONY: test integration-test vulncheck verify-dependency-security run

test:
	$(GO) test ./...

integration-test:
	RUN_POSTGRES_INTEGRATION=1 $(GO) test -v ./...

vulncheck:
	$(GOVULNCHECK) ./...

verify-dependency-security: vulncheck

run:
	@test -n "$(CMD)" || (echo "usage: make run CMD='go test ./...'" >&2; exit 2)
	$(CMD)
