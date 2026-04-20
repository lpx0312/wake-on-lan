# wake-on-lan

命令行工具，通过 Wake-on-LAN (WOL) 魔术包唤醒局域网主机。

## 安装

### 从 Release 下载

前往 [Releases](https://github.com/lpx0312/wake-on-lan/releases) 下载对应平台的二进制文件。

### 从源码构建

```bash
go install
```

## 使用

```bash
wake-on-lan <MAC_ADDRESS>
```

### 支持的 MAC 地址格式

- `00:11:22:33:44:55`
- `00-11-22-33:44:55`
- `001122334455`

### 示例

```bash
wake-on-lan 00:11:22:33:44:55
```

## 构建

使用 GoReleaser 构建多平台二进制：

```bash
goreleaser release --snapshot --clean
```
