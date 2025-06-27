.PHONY: proto clean build run all run-dev run-test run-prod deploy-dev deploy-test deploy-prod

PROTO_DIR=api/study
BIN_DIR=bin
CMD_DIR=cmd/server

PROTO_DIR=api

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
		$(PROTO_DIR)/v1/roland/roland.proto


build:
	go build -o $(BIN_DIR)/study-server ./$(CMD_DIR)

run: build
	./$(BIN_DIR)/study-server

run-dev: build
	cp .env.dev .env
	./$(BIN_DIR)/study-server

run-test: build
	cp .env.test .env
	./$(BIN_DIR)/study-server

run-prod: build
	cp .env.prod .env
	./$(BIN_DIR)/study-server

clean:
	@echo "Cleaning generated files..."
	find $(PROTO_DIR) -type f \( -name '*.pb.go' -o -name '*_grpc.pb.go' \) -delete
	rm -rf $(BIN_DIR)
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

deploy-dev:
	@echo "Deploying to Heroku Dev..."
	heroku container:push web --app studyguides-api-dev --arg ENV=dev
	heroku container:release web --app studyguides-api-dev
	@echo "Dev deployment complete!"

deploy-test:
	@echo "Deploying to Heroku Test..."
	heroku container:push web --app studyguides-api-test --arg ENV=test
	heroku container:release web --app studyguides-api-test
	@echo "Test deployment complete!"

deploy-prod:
	@echo "Deploying to Heroku Prod..."
	heroku container:push web --app studyguides-api-prod --arg ENV=prod
	heroku container:release web --app studyguides-api-prod
	@echo "Prod deployment complete!"
