HOSTNAME := https://sw.bakerstreet.io
DRY_RUN := true
RUNNER_TAG := ''

.PHONY: test-runners
test-runners: /tmp/test-ubuntu-arm64 /tmp/test-macOS-arm64
	@if [ "$$(cat /tmp/test-ubuntu-arm64)" != "200" || "$$(cat /tmp/test-macOS-arm64)" != "200" ]; then \
		@echo "Test failed"; \
		exit 1; \
	fi
	@echo "Test passed"

.PHONY: /tmp/test-ubuntu-arm64
/tmp/test-ubuntu-arm64:
	make -s run-ubuntu-arm64 HOSTNAME=http://localhost:8080 DRY_RUN=true | tail -n 1 > /tmp/test-ubuntu-arm64
	echo "Response code: $$(cat /tmp/test-ubuntu-arm64)"

.PHONY: /tmp/test-macOS-arm64
/tmp/test-macOS-arm64:
	make -s run-macOS-arm64 HOSTNAME=http://localhost:8080 DRY_RUN=true | tail -n 1 > /tmp/test-macOS-arm64
	echo "Response code: $$(cat /tmp/test-macOS-arm64)"

.PHONY: run-ubuntu-arm64
run-ubuntu-arm64:
	# SHA can be calculated like this:
    #echo -n "payload=$(cat github-runner-provisioner/test/ubuntu-arm64_payload.json)" | openssl dgst -sha1 -hmac FAKE_TOKEN
	make test-github-provisioner SHA1=8e39e0658c5eacf3a3e006a46ef46092cbccb5ec RUNNER_TAG=ubuntu-arm64

.PHONY: run-macOS-arm64
run-macOS-arm64:
	# SHA can be calculated like this:
	#echo -n "payload=$(cat github-runner-provisioner/test/macOS-arm64_payload.json)" | openssl dgst -sha1 -hmac FAKE_TOKEN

	make test-github-provisioner SHA1=e504cfa93721fbea2a394d4de9c9be7d5270fc19 RUNNER_TAG=macOS-arm64



.PHONY: test-github-provisioner
test-github-provisioner:
	curl -v -o /dev/null -s $(HOSTNAME)/github-runner-provisioner/?dry-run=$(DRY_RUN) -d "payload=$$(cat test/$(RUNNER_TAG)_payload.json)" -H 'X-Hub-Signature-256: sha1=$(SHA1)' -w "%{http_code}\n"

.PHONY: download-go-modules
download-go-modules:
	go mod download

.PHONY: build
build: download-go-modules
	go mod download; \
	go build .

build-image:
	docker build -t github-runner-provisioner .

.PHONY: go-unit-tests
go-unit-tests: download-go-modules
	go test ./...

.PHONY: update-go-mocks
update-go-mocks:
	go install github.com/golang/mock/mockgen@v1.6.0; \
	mockgen --source internal/aws/iface.go -destination  internal/aws/mocks/mocks.go
