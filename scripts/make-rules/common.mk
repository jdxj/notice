FILENAME := notice.out
DOCKER := docker
OUTPUT := output
MAIN := cmd/main
DEPLOY := deploy

GIT_TAG := $(shell git describe --tags --abbrev=0)
DOCKER_TAG := jdxj/notice:$(GIT_TAG)

.PHONY: clean
clean:
	@rm -rf $(OUTPUT)
	@rm -rf $(DOCKER)/$(FILENAME)