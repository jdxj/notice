.PHONY: image.build
image.build: go.build
	@cp $(OUTPUT)/$(FILENAME) $(DOCKER)
	@docker build -t $(DOCKER_TAG) $(DOCKER)

.PHONY: image.push
image.push: image.build
	@docker push $(DOCKER_TAG)