ENVIRONMENT?=local

CICD_FOLDER=./cicd/$(ENVIRONMENT)
ENV_FILE := $(CICD_FOLDER)/.env
COMPOSE_FILE := $(CICD_FOLDER)/compose.yaml
PWD := $(shell pwd)
include $(ENV_FILE)

MIGRATION_PATH=$(shell realpath $(MIGRATION_RELATIVE_PATH))

export MIGRATION_PATH
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
	$(call validate_env)
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