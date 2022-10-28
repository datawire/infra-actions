HOSTNAME := https://sw.bakerstreet.io

.PHONY: test-github-provisioner
test-github-provisioner:
	curl -v $(HOSTNAME)/github-runner-provisioner/ -d "payload=$$(cat github-runner-provisioner/test/payload.json)"
