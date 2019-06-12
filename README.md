# rentalgames-server

| Branch  | Build Status                                                                                                                               |
| ------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| master  | [![Build Status](https://travis-ci.com/itslaves/rentalgames-server.svg?branch=master)](https://travis-ci.com/itslaves/rentalgames-server)  |
| develop | [![Build Status](https://travis-ci.com/itslaves/rentalgames-server.svg?branch=develop)](https://travis-ci.com/itslaves/rentalgames-server) |

---

## Command lines

### Build

```sh
$ go build -o ./bin/rg-server
```

### Run

```sh
$ ./bin/rg-server --help
$ RG_ENV=develop ./bin/rg-server debug -p 8080 # RG_ENV 로 구동 환경을 선택할 수 있음 (develop|test|staging|production)
```

## Usages

### create

`curl -X POST -F author='freddie' -F content='content' localhost:8080/articles`

### retrieve

`curl -X GET 'localhost:8080/articles?author=freddie'`

## Reference

- <https://godoc.org/golang.org/x/oauth2#Config.TokenSource>