# STC
Standard template construct (https://wh40k.lexicanum.com/wiki/Standard_Template_Construct)

## Commands
Build & Run
```
go build -o dist/stc cmd/stc/main.go && ./dist/stc
``` 
Intigration tests
```
go test -t 10m test/int/*
```

## Structure
```
cmd - entry mains
dist - build output
docs - all kind of docs and pictures
internal - internal code
test - frontend and integration tests
web - statics & templates

```