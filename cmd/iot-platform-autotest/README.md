# IoT Platform 自动化测试框架

ThingsPanel IoT 平台的自动化测试框架,用于验证 MQTT 设备接入、数据上报、指令下发等核心功能。

## 功能特性

- ✅ **直连设备模拟** - 支持 MQTT 直连设备完整测试
- ✅ **遥测数据** - 上报与验证(历史数据 + 当前数据)
- ✅ **属性数据** - 上报与验证(覆盖更新逻辑)
- ✅ **事件数据** - 上报与验证(method + params)
- ✅ **控制指令** - 平台下发遥测控制测试
- ✅ **属性设置** - 平台下发属性设置与响应
- ✅ **命令下发** - 平台下发命令与设备响应
- ✅ **数据库验证** - 自动验证数据正确入库
- ✅ **消息响应** - MQTT 消息接收与匹配验证
- 🚧 **网关设备** - 网关及多级拓扑测试(规划中)

## 快速开始

### 1. 环境要求

- Go 1.21+
- PostgreSQL 数据库访问权限
- MQTT Broker 访问权限
- ThingsPanel 平台 API Key

### 2. 安装依赖

```bash
go mod download
```

## 项目架构

### 目录结构

```
iot-platform-autotest/
├── cmd/autotest/              # 命令行工具
├── internal/
│   ├── config/                # 配置管理
│   ├── device/                # 设备层
│   │   ├── device.go          # 设备接口定义
│   │   ├── direct_device.go   # 直连设备实现
│   │   ├── gateway_device.go  # 网关设备实现(待实现)
│   │   └── factory.go         # 设备工厂
│   ├── protocol/              # 协议层
│   │   ├── message_builder.go # 消息构建器接口
│   │   ├── direct_builder.go  # 直连设备消息构建
│   │   └── gateway_builder.go # 网关设备消息构建(待实现)
│   ├── platform/              # 平台交互层
│   │   ├── api_client.go      # HTTP API 客户端
│   │   └── db_client.go       # 数据库客户端
│   └── utils/                 # 工具函数
├── tests/
│   ├── direct/                # 直连设备测试
│   ├── gateway/               # 网关设备测试(待添加)
│   └── common/                # 公共测试工具
├── testdata/                  # 测试数据样本
└── docs/                      # 文档
```

### 设计原则

1. **接口驱动** - 通过 `Device` 接口统一直连设备和网关设备
2. **配置驱动** - 通过 `device_type` 配置自动选择设备实现
3. **分层架构** - 设备层、协议层、平台交互层职责清晰
4. **易于扩展** - 新增设备类型只需实现接口

### 4. 运行测试

```bash
go test ./tests/direct/... -v
go test ./tests/gateway/... -v
```

**运行指定测试**:

```bash
# 遥测数据测试
go test ./tests/direct/telemetry_test.go -v
go test ./tests/gateway/telemetry_test.go -v

# 属性数据测试
go test ./tests/direct/attribute_test.go -v
go test ./tests/gateway/attribute_test.go -v

# 事件数据测试
go test ./tests/direct/event_test.go -v
go test ./tests/gateway/event_test.go -v

# 命令测试
go test ./tests/direct/command_test.go -v
go test ./tests/gateway/command_test.go -v
```
