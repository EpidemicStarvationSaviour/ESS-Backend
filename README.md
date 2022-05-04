# ESS-Backend

## Deploy

### Configuration

``` bash
cp ./conf/app.sample.ini ./conf/app.ini
```

Then modify the configuration file

### Build

``` bash
go build .
```

### Run

``` bash
./ess
```

## Develop

### Quick Run

``` bash
go run ./main.go
```

### Api Doc

#### Read
Visit `/swagger/index.html` to see the api doc.

#### Update
Use the following command to update docs:
``` bash
swag init
```

If you don't have such a command. Please install swag:
``` bash
go install github.com/swaggo/swag/cmd/swag@latest
```

And add `$GOPATH/bin` to your `$PATH`.

#### Write
For declarative comments format, please refer to [Swaggo Doc](https://github.com/swaggo/swag/blob/master/README.md#declarative-comments-format).

### Git Commit Message

Please use the following format:

[Semantic Commit Messages](https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716)
