# Quote
- [gin-react-boilerplate](https://github.com/wadahiro/gin-react-boilerplate)

# setup

```shell
go get -u github.com/jteeuwen/go-bindata/...
go get -u github.com/elazarl/go-bindata-assetfs/...
```

# package assets

    $GOPATH/bin/go-bindata -o ./bindata.go ../assets/...

# test

    go test -failfast -v -cover -coverprofile=cover.out ./...
    go tool cover -func=cover.out
    go tool cover -html=cover.out -o cover.html
	go tool cover -func=cover.out -o cover.all

# config

```
title = "xxx"

# Trace, Debug, Info, Warning, Error, Fatal, Panic
log = "info"

[csv]
path = "." # load & save path

    [csv.types] # rename csv types
    O = "xx"
    I = "xx"
    L = "xx"
    B = "xx"
    R = "xx"
    X = "xx"
    U = "xx"

    [csv.tags] # rename csv tags
    member = "xx"
    proj   = "xx"
    unit   = "xx"
    unkown = "xx"
```