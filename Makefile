build:
	@go build -o bin/api

seed:
	@go run scripts/seed.go
docker:
	echo "building docker file"
	@docker build -t api .
	echo "running api inside Docker container"
	@docker run -p 3000:3000 api

	
run: build
	@./bin/api
test: 
	@go test -v ./...