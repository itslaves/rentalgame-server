# rentalgames-server

| Branch  | Build Status                                                                                                                               |
| ------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| master  | [![Build Status](https://travis-ci.com/itslaves/rentalgames-server.svg?branch=master)](https://travis-ci.com/itslaves/rentalgames-server)  |
| develop | [![Build Status](https://travis-ci.com/itslaves/rentalgames-server.svg?branch=develop)](https://travis-ci.com/itslaves/rentalgames-server) |

---

## Project setup

### Compiles

```sh
go build -o ./bin/rg-server
```

### Tests

```sh
go test -v ./... # Run your unit tests
gofmt -l . # Check your code formats
golint ./... # Lints your codes
go vet ./... # Run static analysis for your codes
```

### Run

```sh
./bin/rg-server --help # Prints available commands
./bin/rg-server {command} ({command_args})
# ex) ./bin/rg-server debug -p 8080
```

**(NOTE)**
*RG_ENV* 환경변수를 통해서 서버 환경을 동적으로 선택할 수 있음. (develop|test|staging|production)

## Usages

### create

`curl -X POST -F author='freddie' -F content='content' localhost:8080/articles`

### retrieve

`curl -X GET 'localhost:8080/articles?author=freddie'`

## Reference

- <https://godoc.org/golang.org/x/oauth2#Config.TokenSource>