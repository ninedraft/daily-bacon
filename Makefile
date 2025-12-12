DOCKERFILE := deploy/Dockerfile
TARGETS := daily-bacon daily-bacon-gateway openmeteo
DOCKER_BUILD := DOCKER_BUILDKIT=1 docker build
DOCKER_PLUGIN_BUILD := DOCKER_BUILDKIT=1 docker buildx build --load
PLATFORMS := linux/amd64 linux/arm64

define sanitize-platform
$(subst /,-,$(1))
endef

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
DOCKER_MULTI_TARGETS := $(addprefix docker-,$(TARGETS:%=%-multi))
MULTI_IMAGES := $(foreach target,$(TARGETS),$(foreach platform,$(PLATFORMS),$(target):$(GIT_REF)-$(call sanitize-platform,$(platform))))

docker-all: $(DOCKER_TARGETS)
	@printf 'built images:\n'
	@printf '%s\n' $(IMAGES)

.PHONY: docker-all $(DOCKER_TARGETS) docker-all-multi $(DOCKER_MULTI_TARGETS)

docker-all-multi: $(DOCKER_MULTI_TARGETS)
	@printf 'built multi-platform images:\n'
	@printf '%s\n' $(MULTI_IMAGES)

$(foreach target,$(TARGETS),$(eval docker-$(target):; \
	@echo "building $(target):$(GIT_REF)"; \
	$(DOCKER_BUILD) \
		--build-arg TARGET=$(target) \
		-f $(DOCKERFILE) \
		-t $(target):$(GIT_REF) \
		. \
))

$(foreach target,$(TARGETS),$(eval docker-$(target)-multi:; \
	@echo "building $(target) for $(PLATFORMS)"; \
	$(foreach platform,$(PLATFORMS),\
		$(DOCKER_PLUGIN_BUILD) \
			--platform=$(platform) \
			--build-arg TARGET=$(target) \
			-f $(DOCKERFILE) \
			-t $(target):$(GIT_REF)-$(call sanitize-platform,$(platform)) \
			. ;\
	) \
))
