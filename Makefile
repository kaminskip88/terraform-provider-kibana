# TEST?=$$(go list ./... | grep -v 'vendor')
TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
HOSTNAME=hashicorp.com
NAMESPACE=kaminskip88
NAME=kibana
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=darwin_amd64

KIBANA_URL ?= http://127.0.0.1:5601
KIBANA_USERNAME ?= elastic
KIBANA_PASSWORD ?= changeme

default: build

build: fmtcheck
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	KIBANA_URL=${KIBANA_URL} KIBANA_USERNAME=${KIBANA_USERNAME} KIBANA_PASSWORD=${KIBANA_PASSWORD} TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 2m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

docker-elk:
	git clone https://github.com/deviantony/docker-elk.git

compose.up: docker-elk
	docker-compose -f docker-elk/docker-compose.yml up -d elasticsearch kibana setup

compose.wait: compose.up
	until curl -s -X POST -u elastic:changeme -H "Content-Type: application/json" \
	  	http://localhost:9200/_security/user/kibana_system/_password \
		-d "{\"password\":\"changeme\"}" | grep -q "^{}"; do sleep 10; \
	done

compose.down: docker-elk
	docker-compose -f docker-elk/docker-compose.yml down -v

compose.logs: docker-elk
	docker-compose -f docker-elk/docker-compose.yml logs -f

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile test-serv compose.up compose.down compose.logs
