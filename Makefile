build-apps:
	@docker-compose -f docker-compose.yml build

run-debug:
	@docker-compose -f docker-compose.yml up

run:
	@docker-compose -f docker-compose.yml up -d

down:
	@docker-compose -f docker-compose.yml down

clean-images:
	@echo "---------------- Cleaning dangling Docker images ----------------"
	@docker images -f "dangling=true" -q | xargs --no-run-if-empty docker rmi -f

build: build-apps clean-images

clean: down clean-images