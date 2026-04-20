# wake-on-lan 设计文档

**日期**：2026-04-20
**状态**：已批准

## 概述

用 Go 编写的命令行工具，通过 Wake-on-LAN (WOL) 魔术包唤醒局域网内的主机。

**模块**：`github.com/lpx0312/wake-on-lan`
**仓库**：`git@github.com:lpx0312/wake-on-lan.git`

## CLI 接口

```bash
wake-on-lan <MAC_ADDRESS>
```

- **参数**：MAC 地址
- **支持的格式**：
  - `00:11:22:33:44:55`
  - `00-11-22-33-44-55`
  - `001122334455`
- **成功输出**：`Wake-on-LAN packet sent to <MAC>`
- **失败输出**：错误信息到 stderr，exit code 1
- **无参数**：显示帮助信息

## 项目结构

```
wake-on-lan/
├── main.go              # 入口文件
├── go.mod               # Go 模块定义
├── .goreleaser.yml      # GoReleaser 构建配置
└── README.md            # 使用说明
```

## 核心函数

| 函数 | 用途 |
|---|---|
| `parseMAC` | 解析 MAC 地址字符串为字节数组 |
| `parseHexByte` | 解析 2 字符十六进制字符串 |
| `createMagicPacket` | 构建魔术包（6×0xFF + 16×MAC） |
| `sendWOL` | UDP 广播发送到 255.255.255.255:9 |

## 构建配置

### GoReleaser 构建目标

| OS | Arch | 输出文件名 |
|---|---|---|
| Windows | amd64 | `wake-on-lan_windows_amd64.exe` |
| Linux | amd64 | `wake-on-lan_linux_amd64` |

### Release 流程

1. 推送 tag（例如 `v1.0.0`）
2. GoReleaser 自动构建多平台二进制
3. 上传二进制到 GitHub Release

## 错误处理

- 无效 MAC 格式 → `invalid MAC address format`
- 解析失败 → `invalid MAC address`
- UDP 连接失败 → `failed to create connection`
- 发送失败 → `failed to send packet`
