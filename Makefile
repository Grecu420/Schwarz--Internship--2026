up: docker-compose.yml
	docker compose up -d --build

down: docker-compose.yml
	docker compose down


rebuild:
	docker compose up --build --no-cache -d


fresh:
	docker compose down -v
	docker compose build --no-cache
	docker compose up -d



# 
SUBDIRS := user-base property-base friend-request-base conversation-base message-base auth-base api-rest-gateway 


# clean and build all services
# Loop over subdirectories as targets
all: $(SUBDIRS)
.PHONY: all clean $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C services/$@