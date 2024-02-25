.PHONY: docker

docker:
	@rm webook || true
	@docker rmi -f gphughcn/gotest:webook-v0.0.1
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker build -t gphughcn/gotest:webook-v0.0.1 .
