将configs中的configs_example.yaml修改为configs_example，修改配置文件

- 编译到linux

```bash
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o sensitive 
```

- 开发环境运行

```bash
go run .
```