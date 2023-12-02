dep:
	cd ./backend && go mod tidy

run:
	cd ./backend && go run .

test:
	cd ./backend && go test -short -cover ./...

build:
	npm run build --prefix ./frontend && \
	cd ./backend && go build -o bin/aoc-bingo .
