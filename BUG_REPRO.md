# 修复前故障复现（Docker）

## 项目与标准命令
该 CLI 读取本地账户套餐和用量快照，计算超额用量费用并输出稳定的账单文本报表。在仓库根目录执行：

```sh
go build ./...
go run ./cmd/usage-report -input ./examples/sample.json
go test ./...
```

## 环境构建与编译
以下两个平台均已成功构建镜像，并在容器内成功执行 `go build ./...`：

```sh
./build_benzhi_docker.sh linux/amd64
./build_benzhi_docker.sh linux/arm64
```

## 故障触发步骤
在仓库根目录先构建 amd64 复现镜像：

```sh
docker buildx build --platform linux/amd64 --load --file benzhi.Dockerfile --tag subscription-usage-cli:bug-002-base .
```

再执行：

```sh
docker run --rm --platform linux/amd64 subscription-usage-cli:bug-002-base go test ./...
```

## 实际错误输出
```text
ok  	github.com/1260124186-cc/subscription-usage-cli/cmd/usage-report	0.048s
?   	github.com/1260124186-cc/subscription-usage-cli/internal/domain	[no test files]
ok  	github.com/1260124186-cc/subscription-usage-cli/internal/output	0.043s
--- FAIL: TestGenerateDoesNotReorderCallerUsage (0.00s)
    generator_test.go:53: first caller event = "evt-a", want "evt-z"
FAIL
FAIL	github.com/1260124186-cc/subscription-usage-cli/internal/service	0.051s
--- FAIL: TestSnapshotCloneDoesNotShareUsageEvents (0.00s)
    snapshot_clone_test.go:15: original usage units = 99, want 4
FAIL
FAIL	github.com/1260124186-cc/subscription-usage-cli/internal/store	0.045s
FAIL

exit status: 1
```

## 期望行为
账单生成后，调用方持有的用量事件顺序和事件数据保持不变；对快照副本的修改也不应影响原始快照。相同测试命令应以退出码 0 结束。
