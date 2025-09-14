# go-protos

A Go DDD-based gRPC template project with Protobuf.  
一个基于 Go DDD 架构的 gRPC 模板项目。

---

## 📌 Introduction | 项目简介

This project is a **DDD (Domain-Driven Design)** template for building Go microservices with **gRPC** and **Protobuf**, designed to provide a clean, extensible architecture.  

本项目是一个基于 **DDD（领域驱动设计）** 的 Go 微服务模板，使用 **gRPC** 和 **Protobuf** 进行服务通信，旨在提供简洁、可扩展的架构。

---

## ✨ Features | 功能特性

- 📂 **Strict DDD layered architecture** （严格的 DDD 分层架构）  
- ⚡ **gRPC + Protobuf integration** （内置 gRPC 与 Protobuf 集成）  
- 🛠 **Viper for configuration management** （使用 Viper 管理配置）  
- 💾 **GORM for persistence** （使用 GORM 作为持久化层）  
- 🧪 **In-memory & database repository implementations** （支持内存与数据库仓储实现）  
- 🚀 **Ready for microservices** （开箱即用，适合微服务场景）  

---

## 📂 Project Structure | 项目结构

```go
go-protos/
├── cmd/                # Application entrypoints
├── config/             # Configurations (default, dev, prod)
├── internal/           # Core business logic
│   ├── domain/         # Domain layer (entities, services, repositories)
│   ├── application/    # Application layer (use cases, commands, queries)
│   ├── infrastructure/ # Infrastructure (DB, event bus, persistence)
│   └── interfaces/     # Interfaces (gRPC, HTTP, MQ)
├── proto/              # Protobuf definitions
├── sql/                # Database migrations
├── pkg/                # Shared utilities (config, logger, db, etc.)
└── README.md
````

---

## 🚀 Getting Started | 快速开始

### 1. Clone the repo
```bash
git clone https://github.com/yourname/go-protos.git
cd go-protos
````

### 2. Generate protobuf

```bash
protoc --go_out=. --go-grpc_out=. proto/user.proto
```

### 3. Run the service

```bash
go run cmd/service-name/main.go
```

---

## 🧰 Tech Stack | 技术栈

* **Go 1.22+**
* **gRPC + Protobuf**
* **Viper** (configuration)
* **GORM** (database persistence)
* **SQLite/MariaDB/other DBs**
* **DDD layered design**

---

## 📖 License

MIT License. Free for personal and commercial use.
MIT 协议，个人与商业项目均可自由使用。

### ☕ Support
<p>
  <a href="https://buymeacoffee.com/xiaomai">
    <img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" height="50" width="210" alt="buy me a coffee" />
  </a>
</p>
