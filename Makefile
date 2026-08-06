up: docker-compose.yml
	docker compose up -d --build


down: docker-compose.yml
	docker compose down



# 
SUBDIRS := user-base friend-request-base api-rest-gateway


# clean and build all services
# Loop over subdirectories as targets
all: $(SUBDIRS)
.PHONY: all clean $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C services/$@