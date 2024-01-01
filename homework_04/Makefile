compose-up:
	docker-compose -f docker-compose.yml up --remove-orphans -d

compose-down:
	docker-compose -f docker-compose.yml down

docker-build:
	docker build \
		--platform linux/amd64 \
		-t "astrviktor/producer:1.0.0" \
		-f ./golang/docker/producer/Dockerfile .

	docker build \
		--platform linux/amd64 \
		-t "astrviktor/consumer:1.0.0" \
		-f ./golang/docker/consumer/Dockerfile .