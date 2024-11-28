.PHONY: docker
docker: 
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f cp-qwq/webook:v0.0.1
	@docker build --network=host -t cp-qwq/webook:v0.0.1 .