.PHONY: deploy

STAGE = develop
TS = 15684867164

build: test
	env CGO_ENABLED=0 env GOARCH=amd64 env GO111MODULE=on env GOOS=linux go build -o ./bin/gsuite-password-change ./cmd/gsuite_change_password/main.go

test:
	go fmt ./...
	go vet ./...
	go test -v ./...

clean:
	rm -rf ./bin

deploy: build
	sls deploy --verbose --force

deploy-list:
	sls deploy list

rollback:
	sls rollback
