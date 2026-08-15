# 订阅用量结算 CLI

该项目读取账户套餐与用量快照，计算每个账户的超额用量和应收费用，并输出可供账务系统消费的稳定文本报表。运行时只依赖本地 JSON 输入，不访问外部网络或数据库。

在仓库根目录执行：

```sh
go build ./...
go run ./cmd/usage-report -input ./examples/sample.json
go test ./...
```

## Docker 环境

Docker 构建使用实际文件 `benzhi.Dockerfile`。`go.mod` 固定 Go 语言版本为 `1.26.0`，镜像固定使用 `golang:1.26.2-bookworm`，并通过 `GOTOOLCHAIN=local` 禁止在容器内自动切换工具链。镜像始终在容器内从源码执行 `go mod download` 和 `go build ./...`，不使用宿主机二进制。

分别验证 `linux/amd64` 与 `linux/arm64`：

```sh
./build_benzhi_docker.sh linux/amd64
./build_benzhi_docker.sh linux/arm64
```

脚本会依次为指定平台构建镜像、在容器内执行 `go build ./...`，然后启动默认命令 `go run ./cmd/usage-report -input ./examples/sample.json`。

也可以手工执行：

```sh
docker buildx build --platform linux/amd64 --load -f benzhi.Dockerfile -t subscription-usage-cli:manual-amd64 .
docker run --rm subscription-usage-cli:manual-amd64 go build ./...
docker run --rm subscription-usage-cli:manual-amd64
docker run --rm subscription-usage-cli:manual-amd64 go test ./...
```

将上述命令中的平台和镜像名替换为 `linux/arm64` 与 `subscription-usage-cli:manual-arm64`，即可完成 arm64 的同等验证。

通过标准是每条构建和测试命令以退出码 0 结束，启动命令输出账户明细和 `total_charge_cents` 总计行；示例输入应包含 `acme` 的 25 单位超额和总费用 175 分。
