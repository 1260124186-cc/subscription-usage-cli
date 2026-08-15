# 修复前故障复现（Docker）

## 项目与标准命令

本项目是一个根据订阅账户和用量快照生成使用情况报告的 Go CLI。在仓库根目录可执行以下标准命令：

```sh
go build ./...
go run ./cmd/usage-report -input ./examples/sample.json
go test ./...
```

## 环境构建与编译

已实际执行以下命令构建 linux/amd64 镜像并在容器内编译：

```sh
docker buildx build --platform linux/amd64 --load --file benzhi.Dockerfile --tag subscription-usage-cli:delivery-base-005-amd64 .
docker run --rm --platform linux/amd64 subscription-usage-cli:delivery-base-005-amd64 go build ./...
```

已实际执行以下命令构建 linux/arm64 镜像并在容器内编译：

```sh
docker buildx build --platform linux/arm64 --load --file benzhi.Dockerfile --tag subscription-usage-cli:delivery-base-005-arm64 .
docker run --rm --platform linux/arm64 subscription-usage-cli:delivery-base-005-arm64 go build ./...
```

两个平台的镜像构建和容器内 `go build ./...` 均成功。目标故障在下节命令中触发。

## 故障触发步骤

在仓库根目录先按上节命令构建 linux/arm64 镜像，再执行：

```sh
docker run --rm --platform linux/arm64 subscription-usage-cli:delivery-base-005-arm64 go test ./...
```

## 实际错误输出

```text
--- FAIL: TestRunWithInputClosesInputOnce (0.00s)
    main_test.go:30: input Close() calls = 2, want 1
FAIL
FAIL	github.com/1260124186-cc/subscription-usage-cli/cmd/usage-report	0.001s
?   	github.com/1260124186-cc/subscription-usage-cli/internal/domain	[no test files]
ok  	github.com/1260124186-cc/subscription-usage-cli/internal/output	0.001s
ok  	github.com/1260124186-cc/subscription-usage-cli/internal/service	0.002s
--- FAIL: TestLoadSnapshotFileReturnsCloseError (0.00s)
    snapshot_test.go:14: LoadSnapshotFile() error = nil, want the input close error
FAIL
FAIL	github.com/1260124186-cc/subscription-usage-cli/internal/store	0.001s
FAIL
```

## 期望行为

同一仓库与输入在容器内执行 `go test ./...` 时应以退出状态 0 完成，并且不报告任何失败测试。
