# 修复前故障复现（Docker）

## 项目与标准命令

订阅用量结算 CLI 读取账户套餐和用量快照，计算每个账户的超额用量与应收费用，并输出稳定的文本报表。在仓库根目录执行：

```sh
go build ./...
go run ./cmd/usage-report -input ./examples/sample.json
go test ./...
```

## 环境构建与编译

已在两个平台完成镜像构建和容器内编译：

```sh
docker buildx build --platform linux/amd64 --load --file benzhi.Dockerfile --tag subscription-usage-cli:delivery-base-001 .
docker run --rm --platform linux/amd64 subscription-usage-cli:delivery-base-001 go build ./...
docker buildx build --platform linux/arm64 --load --file benzhi.Dockerfile --tag subscription-usage-cli:delivery-base-arm64-001 .
docker run --rm --platform linux/arm64 subscription-usage-cli:delivery-base-arm64-001 go build ./...
```

两个平台的镜像构建和容器内 `go build ./...` 均成功。目标故障在下节的输入触发步骤中出现。

## 故障触发步骤

在仓库根目录执行以下命令：

```sh
printf '%s\n' '{"accounts":[null],"usage":[]}' > /tmp/subscription-usage-null-snapshot.json
docker run --rm --platform linux/amd64 --mount type=bind,src=/tmp/subscription-usage-null-snapshot.json,dst=/tmp/null-snapshot.json,readonly subscription-usage-cli:delivery-base-001 go run ./cmd/usage-report -input /tmp/null-snapshot.json
```

## 实际错误输出

```text
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x4ed93c]

goroutine 1 [running]:
github.com/1260124186-cc/subscription-usage-cli/internal/store.Snapshot.Validate({{0x1ca417602048, 0x1, 0x1}, {0x63f0c0, 0x0, 0x0}})
	/workspace/internal/store/snapshot.go:32 +0x13c
github.com/1260124186-cc/subscription-usage-cli/internal/store.LoadSnapshot({0x52e9c8, 0x1ca417602030})
	/workspace/internal/store/snapshot.go:23 +0x16b
main.run({0x7ffffffffef0?, 0x1ca41762c0d0?}, 0x12a05f200, {0x52e9a8, 0x1ca417602020})
	/workspace/cmd/usage-report/main.go:41 +0x151
main.main()
	/workspace/cmd/usage-report/main.go:25 +0x110
exit status 2
```

## 期望行为

同一输入应以非零退出状态结束，并向调用方输出可读的输入错误，明确说明账户记录无效；进程不应输出运行时 panic 栈，也不应异常中断后续账单任务。
