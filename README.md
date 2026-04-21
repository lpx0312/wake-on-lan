# wake-on-lan

命令行工具，通过 Wake-on-LAN (WOL) 魔术包唤醒局域网主机。

## 安装

### 从 Release 下载

前往 [Releases](https://github.com/lpx0312/wake-on-lan/releases) 下载对应平台的二进制文件。

### 从源码构建

```bash
git clone https://github.com/lpx0312/wake-on-lan.git
cd wake-on-lan
go build -o wake-on-lan.exe .
```

## 使用

```bash
wake-on-lan -m <MAC_ADDRESS> [-t TARGET_IP]
wake-on-lan --mac <MAC_ADDRESS> [--target TARGET_IP]
```

### 参数说明

- `-m`, `--mac` (必需) - 要唤醒的目标机器 MAC 地址
- `-t`, `--target` (可选) - 目标机器的 IP 地址。使用单播传输更可靠。如果省略，则使用广播地址 `255.255.255.255`

### 支持的 MAC 地址格式

- `00:11:22:33:44:55`
- `00-11-22-33:44:55`
- `001122334455`

### 示例

```bash
# 使用广播唤醒（默认）
wake-on-lan -m 00:11:22:33:44:55

# 使用单播直接发送到目标 IP（更可靠）
wake-on-lan --mac 00:11:22:33:44:55 --target 192.168.0.198
```

## 构建

使用 GoReleaser 构建多平台二进制：

```bash
goreleaser release --snapshot --clean
```
