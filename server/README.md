# Quote
- [gin-react-boilerplate](https://github.com/wadahiro/gin-react-boilerplate)

# Setup

```shell
go get -u github.com/go-bindata/go-bindata/...
go get -u github.com/elazarl/go-bindata-assetfs/...
```

# Package Assets

    $GOPATH/bin/go-bindata -o ./bindata.go ../assets/...

# Test

    go test -failfast -v -cover -coverprofile=cover.out ./...
    go tool cover -func=cover.out
    go tool cover -html=cover.out -o cover.html
	go tool cover -func=cover.out -o cover.all

# Config

```
title = "xxx"

# Trace, Debug, Info, Warning, Error, Fatal, Panic
log = "info"

# server
ip = 192.168.1.1
port = 3000 # http server port

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
    nuv    = "xx"
    unkown = "xx"
```

# SDS
## DB
### Account

| TIME         | ID          | TYPE       | UNIT          | NUV          | CLASS      | INPUT         | DEADLINE         |
| :----------- | :---------- | :--------- | :------------ | :----------- | :--------- | :------------ | :--------------- |
| 最后更新时间 | 账户名      | 币种       | 份额          | 单位净值     | 类型       | 总流入        | 理财产品到期时间 |
| TIMESTAMP    | VARCHAR(20) | VARCHAR(3) | DECIMAL(32,3) | DECIMAL(8,3) | VARCHAR(5) | DECIMAL(32,3) | TIMESTAMP        |
| timestamp    | XXXX1234    | CNY etc.   | 150.0         | 1.0          |            | 160.0         | timestamp        |

> total=nuv*unit | input
> diff=input-total

### Borrow

| TIME         | ID          | NAME        | AMOUNT        | ACCOUNT     | NOTE        | DEADLINE  |
| :----------- | :---------- | :---------- | :------------ | :---------- | :---------- | :-------- |
| 最后更新时间 | 借款 ID     | 借款人      | 金额          | 账户        | 备注        | 还款日期  |
| TIMESTAMP    | VARCHAR(32) | VARCHAR(20) | DECIMAL(32,3) | VARCHAR(20) | VARCHAR(32) | TIMESTAMP |

> borrow <--
> lend   -->

## Json
### Record

``` json
{
    "uuid":  "xxx",                   // UUID 唯一标识
    "type":  "xxx",                   // 流水类型，如：收入/支出
    "time":  "yyyy-mm-dd hh:mm:ss",   // 生成时间
    "amount": 0.0,                    // 金额
    "account":["xxx",...],            // 账户
    "nuv":    0.0,                    // 净值, omitempty
    "unit":   0.0,                    // 份额, omitempty
    "class":  ["xxx",...],            // 交易分类, omitempty
    "member": "xxx",                  // 成员, omitempty
    "proj":   "xxx",                  // 项目, omitempty
    "note":   "hello",                // 备注, omitempty
    "deadline":"yyyy-mm-dd hh:mm:ss", // 截止时间, omitempty
}
```

### Account

``` json
{
    "time":    "yyyy-mm-dd hh:mm:ss", // 最后更新时间
    "id":      "xxx",                 // 账户唯一标识
    "type":    "xxx",                 // 货币类型|借贷类型：[B]借入, [L]借出
    "amount":  0.0,                   // 金额
    "unit":    0.0,                   // 份额, omitempty
    "nuv":     0.0,                   // 单位净值, omitempty
    "class":   "xxx",                 // 账户类型, omitempty
    "deadline":"yyyy-mm-dd hh:mm:ss", // 截止日期, omitempty
    "member":  "xxx",                 // 借贷人, omitempty
    "account": "xxx",                 // 借贷转出/转入账户, omitempty
    "note":    "xxx",                 // 借贷备注, omitempty
}
```

## URL

### GET
- `/api/finance/record?start=yy-mm-dd hh:mm:ss&end=yy-mm-dd hh:mm:ss` 获取流水
- `/api/finance/account` 获取账户
  + `/api/finance/account?list` 获取账户名列表？

### POST
- `/api/finance/record` 修改记录
- `/api/finance/account` 修改账户