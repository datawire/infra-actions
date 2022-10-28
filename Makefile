HOSTNAME := https://sw.bakerstreet.io
DRY_RUN := true

.PHONY: test-github-provisioner
test-github-provisioner:
	curl -v $(HOSTNAME)/github-runner-provisioner/?dry-run=$(DRY_RUN) -d "payload=$$(cat github-runner-provisioner/test/payload.json)"
