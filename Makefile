HOSTNAME := https://sw.bakerstreet.io
DRY_RUN := true
RUNNER_TAG := ''

.PHONY: test-ubuntu-arm64
test-ubuntu-arm64:
	# SHA can be calculated like this:
    #echo -n "payload=$(cat github-runner-provisioner/test/ubuntu-arm64_payload.json)" | openssl dgst -sha1 -hmac FAKE_TOKEN
	make test-github-provisioner SHA1=8e39e0658c5eacf3a3e006a46ef46092cbccb5ec RUNNER_TAG=ubuntu-arm64

.PHONY: test-macOS-arm64
test-macOS-arm64:
	# SHA can be calculated like this:
	#echo -n "payload=$(cat github-runner-provisioner/test/macos-arm64_payload.json)" | openssl dgst -sha1 -hmac FAKE_TOKEN

	make test-github-provisioner SHA1=e504cfa93721fbea2a394d4de9c9be7d5270fc19 RUNNER_TAG=macOS-arm64

.PHONY: test-github-provisioner
test-github-provisioner:
	curl -v $(HOSTNAME)/github-runner-provisioner/?dry-run=$(DRY_RUN) -d "payload=$$(cat github-runner-provisioner/test/$(RUNNER_TAG)_payload.json)" -H 'X-Hub-Signature-256: sha1=$(SHA1)'

.PHONY: download-go-modules
download-go-modules:
	cd github-runner-provisioner; \
    go mod download

.PHONY: build
build: download-go-modules
	cd github-runner-provisioner; \
    go mod download; \
	go build .

.PHONY: go-unit-tests
go-unit-tests: download-go-modules
	cd github-runner-provisioner; \
	go test ./...
