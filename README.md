# rentalgames-server

| Branch | Build Status |
| ------ | ------------ |
| master | [![Build Status](https://travis-ci.com/itslaves/rentalgames-server.svg?branch=master)](https://travis-ci.com/itslaves/rentalgames-server) |
| develop | [![Build Status](https://travis-ci.com/itslaves/rentalgames-server.svg?branch=develop)](https://travis-ci.com/itslaves/rentalgames-server) |

---

## Usages

### create

`curl -X POST -F author='freddie' -F content='content' localhost:8080/articles`

### retrieve

`curl -X GET 'localhost:8080/articles?author=freddie'`

## Reference

- <https://godoc.org/golang.org/x/oauth2#Config.TokenSource>