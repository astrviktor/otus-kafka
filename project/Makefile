compose-up:
	docker-compose -f build/docker-compose.yml up --remove-orphans -d

compose-down:
	docker-compose -f build/docker-compose.yml down

docker-build-receiver:
	docker build --no-cache \
		--platform linux/amd64 \
		-t "astrviktor/receiver:1.0.0" \
		-f ./build/receiver/Dockerfile .

docker-build-processor:
	docker build --no-cache \
		--platform linux/amd64 \
		-t "astrviktor/processor:1.0.0" \
		-f ./build/processor/Dockerfile .

docker-build-converter:
	docker build --no-cache \
		--platform linux/amd64 \
		-t "astrviktor/converter:1.0.0" \
		-f ./build/converter/Dockerfile .

docker-build-informer:
	docker build --no-cache \
		--platform linux/amd64 \
		-t "astrviktor/informer:1.0.0" \
		-f ./build/informer/Dockerfile .

docker-push:
	docker push astrviktor/receiver:1.0.0
	docker push astrviktor/processor:1.0.0
	docker push astrviktor/converter:1.0.0
	docker push astrviktor/informer:1.0.0

yandex-tank:
	docker exec -it yandex-tank /bin/bash

postgres:
	docker exec -ti postgres psql -U postgres

clickhouse:
	docker exec -ti clickhouse clickhouse-client