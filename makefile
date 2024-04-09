build:
	go get ./...
	go build . 
run: 
	TZ="America/Bahia" go run .
run-docker: 
	docker-compose up