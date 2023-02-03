REPO=github.com/ocervell/cariddi

fmt:
	@gofmt -s ./*
	@echo "Done."

remod:
	@rm -rf go.*
	@go mod init ${REPO}
	@go get
	@echo "Done."

update:
	@go get -u
	@go mod tidy -v
	@make unlinux
	@git pull
	@make linux
	@echo "Done."

lint:
	@golangci-lint run

linux:
	@go build ./cmd/cariddi
	@sudo mv ./cariddi /usr/bin/
	@echo "Done."

unlinux:
	@sudo rm -rf /usr/bin/cariddi
	@echo "Done."

test:
	@go test -v -race ./...
	@echo "Done."
