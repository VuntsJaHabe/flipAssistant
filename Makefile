go: 
	go run . &

build: 
	cd flip-assistant && npm run dev &

run: build go
	wait
