# wake-on-lan 实现计划

**Goal:** 构建一个 Go CLI 工具，通过 Wake-on-LAN 魔术包唤醒局域网主机

**Architecture:** 简单单文件 CLI，直接使用 Go 标准库 `net` 包发送 UDP 广播魔术包

**Tech Stack:** Go 1.21+, GoReleaser, GitHub Actions

---

## 文件结构

```
wake-on-lan/
├── main.go              # 入口文件，包含所有逻辑
├── go.mod               # Go 模块定义
├── .goreleaser.yml      # GoReleaser 构建配置
└── README.md            # 使用说明
```

---

### Task 1: 初始化 Go 模块

**Files:**
- Create: `go.mod`
- Create: `main.go` (基础框架)

- [ ] **Step 1: 创建 go.mod**

```bash
cd D:/Users/Desktop/wake-on-lan
go mod init github.com/lpx0312/wake-on-lan
```

- [ ] **Step 2: 创建 main.go 框架**

```go
package main

import "fmt"

func main() {
    fmt.Println("wake-on-lan v1.0.0")
}
```

- [ ] **Step 3: 提交**

```bash
git add go.mod main.go
git commit -m "init: project scaffold"
```

---

### Task 2: 实现 MAC 解析函数

**Files:**
- Modify: `main.go`

- [ ] **Step 1: 编写测试用例**

```go
func TestParseMAC(t *testing.T) {
    tests := []struct {
        input    string
        expected []byte
        wantErr  bool
    }{
        {"00:11:22:33:44:55", []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}, false},
        {"00-11-22-33-44-55", []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}, false},
        {"001122334455", []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}, false},
        {"invalid", nil, true},
        {"00:11:22:33:44", nil, true},
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            result, err := parseMAC(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("parseMAC() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && string(result) != string(tt.expected) {
                t.Errorf("parseMAC() = %v, expected %v", result, tt.expected)
            }
        })
    }
}
```

- [ ] **Step 2: 运行测试验证失败**

```bash
go test -v -run TestParseMAC
```

Expected: FAIL - undefined function parseMAC

- [ ] **Step 3: 实现 parseMAC 和 parseHexByte**

```go
// parseHexByte parses a 2-character hex string into a byte.
func parseHexByte(s string) (byte, error) {
    var result byte
    for i := range 2 {
        c := s[i]
        var val byte
        switch {
        case c >= '0' && c <= '9':
            val = c - '0'
        case c >= 'a' && c <= 'f':
            val = c - 'a' + 10
        case c >= 'A' && c <= 'F':
            val = c - 'A' + 10
        default:
            return 0, fmt.Errorf("invalid hex character: %c", c)
        }
        result = result<<4 | val
    }
    return result, nil
}

// parseMAC parses a MAC address string and returns its byte representation.
// Supports formats: "00:11:22:33:44:55", "00-11-22-33-44-55", "001122334455"
func parseMAC(macStr string) ([]byte, error) {
    cleaned := strings.ReplaceAll(macStr, ":", "")
    cleaned = strings.ReplaceAll(cleaned, "-", "")

    if len(cleaned) != 12 {
        return nil, fmt.Errorf("invalid MAC address format: %s", macStr)
    }

    mac := make([]byte, 6)
    for i := range 6 {
        byteStr := cleaned[i*2 : i*2+2]
        b, err := parseHexByte(byteStr)
        if err != nil {
            return nil, fmt.Errorf("invalid MAC address: %s", macStr)
        }
        mac[i] = b
    }

    return mac, nil
}
```

- [ ] **Step 4: 运行测试验证通过**

```bash
go test -v -run TestParseMAC
```

Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add main.go
git commit -m "feat: implement MAC address parsing"
```

---

### Task 3: 实现魔术包创建和发送

**Files:**
- Modify: `main.go`

- [ ] **Step 1: 编写测试用例**

```go
func TestCreateMagicPacket(t *testing.T) {
    mac := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
    packet := createMagicPacket(mac)

    // 检查长度
    if len(packet) != 102 {
        t.Errorf("packet length = %d, expected 102", len(packet))
    }

    // 检查前6字节是 0xFF
    for i := 0; i < 6; i++ {
        if packet[i] != 0xFF {
            t.Errorf("packet[%d] = %x, expected 0xFF", i, packet[i])
        }
    }

    // 检查后96字节是16次MAC重复
    for i := 0; i < 16; i++ {
        for j := 0; j < 6; j++ {
            if packet[6+i*6+j] != mac[j] {
                t.Errorf("packet[%d] = %x, expected %x", 6+i*6+j, packet[6+i*6+j], mac[j])
            }
        }
    }
}
```

- [ ] **Step 2: 运行测试验证失败**

```bash
go test -v -run TestCreateMagicPacket
```

Expected: FAIL - undefined function createMagicPacket

- [ ] **Step 3: 实现 createMagicPacket 和 sendWOL**

```go
// createMagicPacket creates a Wake-on-LAN magic packet for the given MAC address.
// Format: 6 bytes of 0xFF followed by 16 repetitions of the MAC address.
func createMagicPacket(mac []byte) []byte {
    packet := make([]byte, 102)
    for i := range 6 {
        packet[i] = 0xFF
    }
    for i := range 16 {
        copy(packet[6+i*6:], mac)
    }
    return packet
}

// sendWOL sends a Wake-on-LAN magic packet to the specified MAC address.
func sendWOL(mac []byte) error {
    packet := createMagicPacket(mac)

    conn, err := net.Dial("udp", "255.255.255.255:9")
    if err != nil {
        return fmt.Errorf("failed to create connection: %w", err)
    }
    defer func() {
        _ = conn.Close()
    }()

    udpConn := conn.(*net.UDPConn)
    if _, err := udpConn.Write(packet); err != nil {
        return fmt.Errorf("failed to send packet: %w", err)
    }

    return nil
}
```

- [ ] **Step 4: 运行测试验证通过**

```bash
go test -v -run TestCreateMagicPacket
```

Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add main.go
git commit -m "feat: implement magic packet creation and WOL send"
```

---

### Task 4: 实现 main 函数

**Files:**
- Modify: `main.go`

- [ ] **Step 1: 编写 main 函数逻辑**

```go
func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s <MAC_ADDRESS>\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "Example: %s 00:11:22:33:44:55\n", os.Args[0])
        os.Exit(1)
    }

    macStr := os.Args[1]

    mac, err := parseMAC(macStr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    if err := sendWOL(mac); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Wake-on-LAN packet sent to %s\n", macStr)
}
```

- [ ] **Step 2: 添加必要的 import**

确保 import 包含:
```go
import (
    "fmt"
    "net"
    "os"
    "strings"
)
```

- [ ] **Step 3: 本地测试**

```bash
go build -o wake-on-lan.exe .
./wake-on-lan.exe
```

Expected: 显示帮助信息

```bash
./wake-on-lan.exe 00:11:22:33:44:55
```

Expected: 尝试发送（可能报错，如果没有目标机器）

- [ ] **Step 4: 提交**

```bash
git add main.go
git commit -m "feat: implement CLI entry point"
```

---

### Task 5: 配置 GoReleaser

**Files:**
- Create: `.goreleaser.yml`

- [ ] **Step 1: 创建 .goreleaser.yml**

```yaml
before:
  hooks:
    - go mod download

builds:
  - id: wake-on-lan
    dir: .
    main: ./main.go
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
    goarch:
      - amd64
    binary: wake-on-lan

archives:
  - id: default
    format: binary
    format_overrides:
      - goos: windows
        format: zip
```

- [ ] **Step 2: 本地测试构建**

```bash
goreleaser build --snapshot --clean
```

- [ ] **Step 3: 提交**

```bash
git add .goreleaser.yml
git commit -m "ci: add GoReleaser configuration"
```

---

### Task 6: 创建 README

**Files:**
- Create: `README.md`

- [ ] **Step 1: 创建 README.md**

```markdown
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
- `00-11-22-33-44-55`
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
```

- [ ] **Step 2: 提交**

```bash
git add README.md
git commit -m "docs: add README"
```

---

### Task 7: 推送代码并创建 Release

**Files:**
- Modify: 远程仓库

- [ ] **Step 1: 添加远程仓库并推送**

```bash
git remote add origin git@github.com:lpx0312/wake-on-lan.git
git branch -M main
git push -u origin main
```

- [ ] **Step 2: 创建 v1.0.0 tag 并推送**

```bash
git tag v1.0.0
git push origin v1.0.0
```

---

## 自检清单

- [x] spec 覆盖：CLI 接口、MAC 解析、魔术包构建、GoReleaser 配置
- [x] 无 placeholder：所有步骤都有完整代码
- [x] 类型一致性：函数签名在所有 task 中保持一致
