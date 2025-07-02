.PHONY: proto clean build run all run-dev run-test run-prod

# Directory configuration
PROTO_DIR=api
BIN_DIR=bin
CMD_DIR=cmd/server

# Create a new version tag
.PHONY: tag
tag:
	@echo "Current version: $$(make get-version)"
	@read -p "Enter new version (e.g., 0.2.4): " new_version; \
	git tag -a "v$$new_version" -m "Release v$$new_version"; \
	echo "Created tag v$$new_version"

# Get current version or default to 0.0.0
get-version:
	@latest_tag=$$(git tag --sort=-v:refname | head -n1 | sed 's/^v//'); \
	if [ -z "$$latest_tag" ]; then \
		echo "0.0.0"; \
	else \
		echo "$$latest_tag"; \
	fi

# Increment patch version
.PHONY: bump-patch
bump-patch:
	@echo "Bumping patch version..."
	@current=$$(make get-version); \
	major=$$(echo $$current | cut -d. -f1); \
	minor=$$(echo $$current | cut -d. -f2); \
	patch=$$(echo $$current | cut -d. -f3); \
	new_patch=$$((patch + 1)); \
	new_version="v$$major.$$minor.$$new_patch"; \
	git tag -a "$$new_version" -m "Release $$new_version"; \
	echo "Created tag $$new_version"

# Increment minor version
.PHONY: bump-minor
bump-minor:
	@echo "Bumping minor version..."
	@current=$$(make get-version); \
	major=$$(echo $$current | cut -d. -f1); \
	minor=$$(echo $$current | cut -d. -f2); \
	new_minor=$$((minor + 1)); \
	new_version="v$$major.$$new_minor.0"; \
	git tag -a "$$new_version" -m "Release $$new_version"; \
	echo "Created tag $$new_version"

# Increment major version
.PHONY: bump-major
bump-major:
	@echo "Bumping major version..."
	@current=$$(make get-version); \
	major=$$(echo $$current | cut -d. -f1); \
	new_major=$$((major + 1)); \
	new_version="v$$new_major.0.0"; \
	git tag -a "$$new_version" -m "Release $$new_version"; \
	echo "Created tag $$new_version"

# Git workflow: merge dev to test, then test to main
.PHONY: git-workflow
git-workflow:
	@echo "Starting git workflow: dev → test → main"
	git checkout test && git merge dev && git push && git checkout main && git merge test && git push && git checkout dev
	@echo "Git workflow completed successfully"

proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/v1/search/search.proto \
		$(PROTO_DIR)/v1/health/health.proto \
		$(PROTO_DIR)/v1/user/user.proto \
		$(PROTO_DIR)/v1/shared/tag.proto \
		$(PROTO_DIR)/v1/shared/user.proto \
		$(PROTO_DIR)/v1/shared/tagsearchresult.proto \
		$(PROTO_DIR)/v1/shared/contexttype.proto \
		$(PROTO_DIR)/v1/shared/usersearchresult.proto \
		$(PROTO_DIR)/v1/shared/tagtype.proto \
		$(PROTO_DIR)/v1/shared/contentrating.proto \
		$(PROTO_DIR)/v1/shared/question.proto \
		$(PROTO_DIR)/v1/shared/studymethod.proto \
		$(PROTO_DIR)/v1/shared/interactiontype.proto \
		$(PROTO_DIR)/v1/shared/deckassignment.proto \
		$(PROTO_DIR)/v1/interaction/interaction.proto \
		$(PROTO_DIR)/v1/tag/tag.proto \
		$(PROTO_DIR)/v1/question/question.proto \
		$(PROTO_DIR)/v1/shared/reporttype.proto \
		$(PROTO_DIR)/v1/chat/chat.proto \
		$(PROTO_DIR)/v1/shared/bundle.proto \
		$(PROTO_DIR)/v1/shared/parsertype.proto \
		$(PROTO_DIR)/v1/shared/exporttype.proto \
		$(PROTO_DIR)/v1/shared/contentdescriptortype.proto \
		$(PROTO_DIR)/v1/roland/roland.proto \
		$(PROTO_DIR)/v1/admin/admin.proto \
		$(PROTO_DIR)/v1/shared/metadata.proto \
		$(PROTO_DIR)/v1/shared/passage.proto \
		$(PROTO_DIR)/v1/shared/taginfo.proto \
		$(PROTO_DIR)/v1/shared/tagindexresult.proto \
		$(PROTO_DIR)/v1/shared/node.proto \
		$(PROTO_DIR)/v1/shared/prompt.proto \
		$(PROTO_DIR)/v1/shared/questiontag.proto \
		$(PROTO_DIR)/v1/shared/tagrow.proto \
		$(PROTO_DIR)/v1/shared/tagnode.proto \
		$(PROTO_DIR)/v1/shared/ancestor.proto \
		$(PROTO_DIR)/v1/shared/section.proto \
		$(PROTO_DIR)/v1/shared/guide.proto \

build:
	go build -o ./bin/server ./cmd/server

run: build
	./bin/server

run-dev: build
	cp .env.dev .env
	./bin/server

run-test: build
	cp .env.test .env
	./bin/server

run-prod: build
	cp .env.prod .env
	./bin/server

clean:
	@echo "Cleaning generated files..."
	find api -type f \( -name '*.pb.go' -o -name '*_grpc.pb.go' \) -delete
	rm -f server
	rm -f .env
	@echo "Done."

fmt:
	find . -type f -name '*.go' ! -name '*.pb.go' -exec gofmt -s -w {} +

lint:
	golangci-lint run --timeout=2m --skip-dirs-use-default --skip-files='.*\.pb\.go'

auth-evans-dev:
	cp .env.dev .env
	evans --host localhost --port 1973 --header "Authorization=Bearer $$(cat .jwt.dev)" -r

auth-evans-test:
	cp .env.test .env
	evans --host localhost --port 1973 --header "Authorization=Bearer $$(cat .jwt.test)" -r

auth-evans-prod:
	cp .env.prod .env
	evans --host localhost --port 1973 --header "Authorization=Bearer $$(cat .jwt.prod)" -r

evans-dev:
	cp .env.dev .env
	evans --host localhost --port 1973 -r

evans-test:
	cp .env.test .env
	evans --host localhost --port 1973 -r

evans-prod:
	cp .env.prod .env
	evans --host localhost --port 1973 -r

generate-tokens:
	node scripts/generate_jwt.js > .jwt.dev
	node scripts/generate_jwt.js > .jwt.test
	node scripts/generate_jwt.js > .jwt.prod

all: proto build
