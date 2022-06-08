.PHONY: tag.release
tag.release:
	@git tag -f -a v$(shell gsemver bump) -m release

.PHONY: tag.push
tag.push: tag.release
	@git push origin $(GIT_TAG)

.PHONY: tag.test
tag.test:
	@git tag -f -a test-$(shell git rev-parse --short HEAD) -m test-tag