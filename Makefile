DOCKERFILE := deploy/Dockerfile
TARGETS := daily-bacon daily-bacon-gateway openmeteo
GIT_REF := $(shell \
	if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then \
		tag=$$(git describe --tags --exact-match 2>/dev/null); \
		if [ -n "$${tag}" ]; then \
			printf '%s' "$${tag}"; \
		else \
			branch=$$(git rev-parse --abbrev-ref HEAD 2>/dev/null); \
			if [ -n "$${branch}" ] && [ "$${branch}" != "HEAD" ]; then \
				date=$$(git show -s --format=%cd --date=format:%Y%m%d 2>/dev/null); \
				hash=$$(git rev-parse --short HEAD 2>/dev/null); \
				printf '%s-%s-%s' "$${branch}" "$${date}" "$${hash}"; \
			else \
				git rev-parse --short HEAD 2>/dev/null; \
			fi; \
		fi; \
	else \
		printf 'unknown'; \
	fi)

DOCKER_TARGETS := $(addprefix docker-,$(TARGETS))
IMAGES := $(foreach target,$(TARGETS),$(target):$(GIT_REF))

.PHONY: docker-all $(DOCKER_TARGETS)

docker-all: $(DOCKER_TARGETS)
	@printf 'built images:\n'
	@printf '%s\n' $(IMAGES)

$(foreach target,$(TARGETS),$(eval docker-$(target):; \
	@echo "building $(target):$(GIT_REF)"; \
	docker build \
		--build-arg TARGET=$(target) \
		-f $(DOCKERFILE) \
		-t $(target):$(GIT_REF) \
		. \
	))
