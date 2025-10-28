### Run

To run this project:

```sh
make
```

to run by handle
```sh 
go run cmd/api/main.go --config-file=config/config_test.yaml
```

### Compile

手动编译二进制包

```sh
go build -o bin/staker-api ./cmd/api
```

自动生成二进制包

```sh
make build
```

运行二进制

```sh
./bin/staker-api --config-file=config/config_test.yaml
```

或
```sh
make start
```