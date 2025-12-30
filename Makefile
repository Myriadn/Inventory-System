# scripting run

run:
	echo "Running ..."
	go run cmd/api/main.go

sync:
	echo "Syncing . . ."
	go mod tidy
