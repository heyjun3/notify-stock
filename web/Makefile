.PHONY: build \
	push

build:
	docker compose build web
	docker image tag web-web us-west1-docker.pkg.dev/kj-project-0/notify-stock/web-server

push:
	docker push us-west1-docker.pkg.dev/kj-project-0/notify-stock/web-server
