.PHONY: all test

test_integ:
	set -e ;\
	docker-compose -f docker-compose.test.yml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f docker-compose.test.yml run integration_tests go test || test_status_code=$$? ;\
	docker-compose -f docker-compose.test.yml down ;\
	exit $$test_status_code ;

test_integ-cleanup:
	docker-compose -f docker-compose.test.yml down \
        --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
  	docker-compose rm -f

# Запуск docker-compose в фоновом режиме
up:
	docker-compose up -d

# Просмотр логов
log: 
	docker-compose logs

# Перезапуск контейнеров
re:
	docker-compose up -d --build

# Останов docker-compose 
dw:
	docker-compose down

cli:
	$(MAKE) --directory=cli_service  cli

del:
	$(MAKE) --directory=cli_service  del

tree:
	tree | xclip -selection clipboard

test:
	$(MAKE) -C bruteforce test
	
install-lint-deps:
	
lint: 
