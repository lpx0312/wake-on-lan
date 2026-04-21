# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build -o wake-on-lan .

# Run tests
go test ./...

# Run a single test
go test -run TestParseMAC
go test -run TestCreateMagicPacket

# Release (triggers via git tag push)
git tag v1.x.x && git push origin v1.x.x
```

## Architecture

单文件 CLI 工具，无外部依赖（纯 stdlib）。

### 核心文件

- **`main.go`** — CLI 入口 + 全平台核心逻辑：
  - `parseMAC()`：支持 `:`、`-`、无分隔符三种 MAC 格式
  - `createMagicPacket()`：生成 102 字节 WOL 魔术包（6×0xFF + 16×MAC）
  - `sendWOL()`：先尝试 `targetIP`，失败则 fallback 到 `255.255.255.255`

- **`broadcast_windows.go`**（build tag: `windows`）— 使用 `net.DialUDP` 发送 UDP 包（Windows 不支持 `syscall.Sendto`）

- **`broadcast_unix.go`**（build tag: `!windows`）— 使用 `syscall.Socket` + `SO_BROADCAST` + `syscall.Sendto` 发送

### 平台差异关键点

Windows 下 `syscall.Sendto` 不可用，必须用 `net.DialUDP`。两个文件通过 build tag 互斥，共同实现 `sendWolBroadcast(packet []byte, targetIP string, targetPort int) error` 接口。修改发包逻辑时两个文件都要同步考虑。

### CLI 参数风格（2026-04-21 重构）

- 使用 `flag.String("m", "", desc)` + `flag.String("mac", "", desc)` 实现短/长 flag
- `flag.Parse()` 前必须手动检查 `-h`/`--help`（Go 标准库不支持 combinable help）
- `flag.NFlag() == 0 && flag.NArg() == 0` 判断是否无任何参数
- `sendWOL` 返回实际发送的 IP（targetIP 或 fallback 后的广播地址），main 层统一打印

### CLI 用法

```
wake-on-lan -m <MAC_ADDRESS> [-t TARGET_IP]
```

- `-m`/`--mac`：必选参数，指定目标 MAC
- `-t`/`--target`：可选参数，指定目标 IP（默认 `255.255.255.255` 广播）
- 输出差异化：广播模式和单播模式打印不同信息（包含实际发送的 IP）

### Build 注意

- goreleaser `binary: wake-on-lan` 在 Windows 构建时自动加 `.exe` 后缀

### Read 工具限制

- 如遇 `Read` 工具只读 line 1 的 hook 问题，用 `cat -n file` 作为备选

### 发布流程

GitHub Actions（`.github/workflows/release.yml`）在推送 `v*` tag 时触发，使用 goreleaser 构建 Linux/Windows amd64 二进制并发布到 GitHub Releases。
