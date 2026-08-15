# 修复前故障复现（Docker）

## 项目与标准命令
该项目读取账户套餐和用量快照，计算超额用量与应收费用，并输出供账务系统使用的文本报表。在仓库根目录执行以下标准命令：

```sh
go build ./...
go run ./cmd/usage-report -input ./examples/sample.json
go test ./...
```

前两条命令可以成功结束并输出示例报表。修复前，`go test ./...` 会因超时写入场景的测试失败而以非零退出码结束。

## 环境构建与编译
已在修复前基线实际执行以下两个平台的镜像构建命令：

```sh
./build_benzhi_docker.sh linux/amd64
./build_benzhi_docker.sh linux/arm64
```

两个平台的镜像均成功构建，且均在容器内完成 `go build ./...`。对应的容器内编译命令为：

```sh
docker run --rm --platform linux/amd64 subscription-usage-cli:benzhi-linux-amd64 go build ./...
docker run --rm --platform linux/arm64 subscription-usage-cli:benzhi-linux-arm64 go build ./...
```

## 故障触发步骤
在仓库根目录执行：

```sh
go test ./cmd/usage-report -run TestRunStopsWritingAtTimeout -count=1
```

## 实际错误输出
```text
--- FAIL: TestRunStopsWritingAtTimeout (0.02s)
    main_test.go:26: run() error = nil, want context deadline exceeded
FAIL
FAIL	github.com/1260124186-cc/subscription-usage-cli/cmd/usage-report	1.870s
FAIL
exit code: 1
```

## 期望行为
当调用方配置的超时已到期且输出目标响应缓慢时，命令应停止报表输出并返回超时错误，而不是以成功结果结束或继续等待输出完成。
