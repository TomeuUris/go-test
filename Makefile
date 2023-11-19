.PHONY: build-dev build-prod run-dev run-prod
build-dev:
	docker build --target development -t my-go-app-dev .

build-prod:
	docker build --target production -t my-go-app-prod .

run-dev:
	docker run -p 8080:8080 my-go-app-dev

run-prod:
	docker run -p 8080:8080 my-go-app-prod