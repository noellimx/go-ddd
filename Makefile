ENVIRONMENT?=local

CICD_FOLDER=./cicd/$(ENVIRONMENT)
DOCKER_FOLDER=./cicd
ENV_FILE := $(CICD_FOLDER)/.env
PWD := $(shell pwd)
include $(ENV_FILE)

COMPOSE_FILE := $(shell realpath $(CICD_FOLDER)/compose.yaml)

MIGRATION_PATH := $(shell realpath $(MIGRATION_RELATIVE_PATH))
export MIGRATION_PATH

APP_MAIN_PATH := $(APP_MAIN_RELATIVE_PATH)
export APP_MAIN_PATH

APP_ROOT := $(shell realpath $(APP_RELATIVE_ROOT))
export APP_ROOT

APP_DOCKERFILE := $(DOCKER_FOLDER)/${APP_DOCKERFILE_NAME}.Dockerfile
export APP_DOCKERFILE

# List of allowed environments
VALID_ENVIRONMENTS := local

COMPOSE_BIN := $(shell which docker-compose || command -v docker-compose)

# Function to check if ENVIRONMENT is valid
define validate_env
  $(if $(filter $(ENVIRONMENT),$(VALID_ENVIRONMENTS)),,\
    $(error Invalid ENVIRONMENT '$(ENVIRONMENT)'. Allowed values: $(VALID_ENVIRONMENTS)))
endef

# Function to validate Docker Compose existence
define validate_compose
  $(if  $(COMPOSE_BIN),,$(error Docker Compose not found. Install 'docker-compose' or Docker with 'docker compose'.))
endef

.PHONY: validate
validate:
	@echo validating env...
	$(call validate_env)
	@echo validating compose...
	$(call validate_compose)
	@echo validation completed. ENVIRONMENT=$(ENVIRONMENT) COMPOSE_BIN=$(COMPOSE_BIN)

.PHONY: hello
hello:
	$(MAKE) validate
	@echo $(ENV_FILE)

.PHONY: colima-start
colima-start:
	@echo ${PWD}
	colima start --network-address --mount ${PWD}:w

.PHONY: setup
setup:
	$(MAKE) validate
	@echo "setup: APP_DOCKERFILE=${APP_DOCKERFILE}"
	$(COMPOSE_BIN) --env-file $(ENV_FILE) -f $(COMPOSE_FILE) up -d  --remove-orphans
	$(MAKE) migrate

# tear down all resources including volumes
.PHONY: teardown
teardown:
	$(MAKE) validate
	$(COMPOSE_BIN) -f $(COMPOSE_FILE) down -v

COMPOSE_MIGRATE:=$(COMPOSE_BIN) -f $(COMPOSE_FILE) run --rm migrate -path=/migrations -database=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable
.PHONY: migrate
migrate:
	$(MAKE) validate
	${COMPOSE_MIGRATE} up

.PHONY: migrate-down
migrate-down:
	$(MAKE) validate
	${COMPOSE_MIGRATE} down -all