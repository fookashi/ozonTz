
.PHONY: run, build

run:
ifeq ($(db),pg)
	@docker-compose -f docker-compose.yml -f docker-compose.postgres.yml --env-file config/.env up --build
else
	@docker-compose -f docker-compose.yml --env-file config/.env up --build
endif

stop:
	@docker-compose -f docker-compose.yml --env-file config/.env down

build:
	@docker-compose -f docker-compose.yml --env-file config/.env build

clean:
	@docker-compose down -v --remove-orphans
	@docker system prune -f

env-template:
	@cp config/.env.example config/.env