default: build run

build:
	docker-compose build

run:
	docker-compose up -d

stop:
	docker-compose down

.PHONY: test
test:
	cd testing \
		&& docker build --tag traffic-wizard-testing . \
		&& docker run --rm traffic-wizard-testing