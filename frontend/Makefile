# Variables
PROJECT_NAME = cloud-consulting-business
DOCKER_IMAGE = $(PROJECT_NAME):latest
DOCKER_REGISTRY = docker.io/tetianamost

# Development commands
.PHONY: start
start:
	npm start

.PHONY: build
build:
	npm run build

.PHONY: test
test:
	npm test

# Docker commands
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-push
docker-push:
	docker tag $(DOCKER_IMAGE) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)

.PHONY: helm-lint
helm-lint:
	helm lint helm/cloud-consulting

.PHONY: helm-template
helm-template:
	helm template cloud-consulting helm/cloud-consulting

.PHONY: helm-install
helm-install:
	helm install cloud-consulting helm/cloud-consulting

.PHONY: helm-upgrade
helm-upgrade:
	helm upgrade cloud-consulting helm/cloud-consulting

.PHONY: helm-uninstall
helm-uninstall:
	helm uninstall cloud-consulting

# Update the deploy target
.PHONY: deploy
deploy: build docker-build docker-push helm-upgrade