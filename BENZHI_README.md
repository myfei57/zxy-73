基于 Go 实现的 RubberCure 项目，一款后端服务，管理橡胶硫化机群与模具调度控制。

## 运行

使用固定 Go 工具链构建并启动控制台：

```sh
go build -mod=vendor ./cmd/rubbercure
./rubbercure
```

控制台默认监听 127.0.0.1:8080，提供机群状态、硫化启动、停机与审计查询接口。
