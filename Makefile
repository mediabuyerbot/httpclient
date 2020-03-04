.PHONY: deps test cover sync-coveralls mocks

deps:
	go mod download

test: deps
	go test ./...

covertest: deps
	go test  -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

sync-coveralls: deps
	go test  -coverprofile=coverage.out ./...
	goveralls -coverprofile=coverage.out -reponame=httpclient -repotoken=${COVERALLS_HTTPCLIENT_TOKEN} -service=local

mocks: deps
	mockgen -package=httpclient  -destination=client_mock.go . Doer

