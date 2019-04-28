## Test
### create
`curl -X POST -F author='freddie' -F content='content' localhost:8080/articles`

### retrieve
`curl -X GET 'localhost:8080/articles?author=freddie'`


## Reference
https://godoc.org/golang.org/x/oauth2#Config.TokenSource