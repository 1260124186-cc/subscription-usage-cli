# 修复前故障复现（Docker）

## 项目与标准命令
本项目读取账户套餐与用量快照，计算超额用量与费用，并输出供账务系统消费的文本报表。在仓库根目录可执行：

```sh
go build ./...
go run ./cmd/usage-report -input ./examples/sample.json
go test ./...
```

## 环境构建与编译
已实际执行以下 linux/amd64 命令，镜像构建和容器内编译均成功：

```sh
docker buildx build --platform linux/amd64 --load -f benzhi.Dockerfile -t subscription-usage-cli:delivery-003-base-amd64 .
docker run --rm --platform linux/amd64 subscription-usage-cli:delivery-003-base-amd64 go build ./...
```

已实际执行以下 linux/arm64 命令，镜像构建和容器内编译均成功：

```sh
docker buildx build --platform linux/arm64 --load -f benzhi.Dockerfile -t subscription-usage-cli:delivery-003-base-arm64 .
docker run --rm --platform linux/arm64 subscription-usage-cli:delivery-003-base-arm64 go build ./...
```

## 故障触发步骤
在仓库根目录使用 linux/arm64 基线镜像执行：

```sh
docker run --rm --platform linux/arm64 subscription-usage-cli:delivery-003-base-arm64 go test -count=1 ./...
```

## 实际错误输出
```text
--- FAIL: TestRunReturnsOutputWriteError (0.00s)
    main_test.go:27: run() error = nil, want an output write error
FAIL
FAIL	github.com/1260124186-cc/subscription-usage-cli/cmd/usage-report	0.004s
?   	github.com/1260124186-cc/subscription-usage-cli/internal/domain	[no test files]
--- FAIL: TestWriteTextReturnsDetailWriteError (0.00s)
    text_test.go:37: WriteText() error = nil, want the account write error
FAIL
FAIL	github.com/1260124186-cc/subscription-usage-cli/internal/output	0.006s
ok  	github.com/1260124186-cc/subscription-usage-cli/internal/service	0.003s
?   	github.com/1260124186-cc/subscription-usage-cli/internal/store	[no test files]
FAIL
exit status: 1
```

## 期望行为
该测试命令应以退出码 0 完成。对于报表输出目标不可用的情况，调用方应能观察到明确的失败结果，而不是把执行结果视为成功。
