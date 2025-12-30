# scripting run

run:
	echo "Init running...."
	go run cmd/api/main.go

sync:
	echo "Syncing . . ."
	go mod tidy
