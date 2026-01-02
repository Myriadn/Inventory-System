# scripting run
COVER_OUT := coverage.out
COVER_HTML := coverage.html

run:
	@echo Running ...
	go run cmd/api/main.go

sync:
	@echo Syncing . . .
	go mod tidy

coverage:
	@echo Calculating coverage...
	go test -coverprofile=$(COVER_OUT) ./internal/repository/... ./internal/service/...
	go tool cover -html=$(COVER_OUT) -o $(COVER_HTML)

clear:
	@echo Cleaning up...
	go clean
	rm -f $(COVER_OUT) $(COVER_HTML)
