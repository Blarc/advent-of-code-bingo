# Advent of Code Bingo
Web application for playing bingo based on Advent of Code.


## Development

### Generating swagger
```bash
cd ./backend
swag init --parseDependency --parseInternal
```

### Running the server
First you need to build the frontend:
```bash
cd ./frontend
npm run build
```
Then you can run the server that serves both backend and frontend:

```bash
cd ./backend
go run .
```
Swagger is available at http://localhost:8080/api/v1/swagger/index.html and the frontend at http://localhost:8080/.


## Acknowledgements
- Inspired by [this reddit post](https://www.reddit.com/r/adventofcode/comments/17icom1/ive_created_bingo_cards_to_have_a_little_fun_at/).
