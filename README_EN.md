# Katool-Go

<div align="center">

<img src="logo.png" alt="Katool-Go Logo" width="400">

<h1>🛠️ Katool-Go</h1>

<p>
  <a href="https://pkg.go.dev/github.com/karosown/katool-go"><img src="https://pkg.go.dev/badge/github.com/karosown/katool-go.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/karosown/katool-go"><img src="https://goreportcard.com/badge/github.com/karosown/katool-go" alt="Go Report Card"></a>
  <a href="https://github.com/karosown/katool-go/releases"><img src="https://img.shields.io/github/v/release/karosown/katool-go" alt="GitHub release"></a>
  <a href="https://github.com/karosown/katool-go/blob/main/LICENSE"><img src="https://img.shields.io/github/license/karosown/katool-go" alt="License"></a>
  <a href="https://golang.org/dl/"><img src="https://img.shields.io/github/go-mod/go-version/karosown/katool-go" alt="Go Version"></a>
</p>

<b><i>A comprehensive Go utility library inspired by Java ecosystem best practices, providing full-spectrum development support</i></b>

<p>
  <a href="README.md">🇨🇳 中文</a> |
  <a href="README_EN.md">🇺🇸 English</a>
</p>

</div>

<hr>

## 📋 Table of Contents

- [📝 Introduction](#introduction)
- [✨ Features](#features)
- [📦 Installation](#installation)
- [🚀 Quick Start](#quick-start)
- [🔧 Core Modules](#core-modules)
- [💡 Best Practices](#best-practices)
- [👥 Contributing](#contributing)
- [📄 License](#license)

<hr>

## 📝 Introduction

**Katool-Go** is a modern, comprehensive Go utility library designed to enhance development efficiency and code quality. It draws inspiration from mature design patterns in the Java ecosystem while fully leveraging Go's modern features like generics and goroutines, providing developers with a complete toolkit solution.

### 🎯 Design Goals

- **Type Safety**: Fully utilize Go 1.18+ generics for type-safe APIs
- **High Performance**: Built-in concurrency optimizations leveraging Go's performance advantages
- **Easy to Use**: Provide Java Stream API-like chaining operations to reduce learning curve
- **Production Ready**: Complete error handling, logging system, and test coverage

<hr>

## ✨ Features

<table>
  <tr>
    <td><b>🌊 Stream Processing</b></td>
    <td>Java 8 Stream API-like chaining operations with parallel processing, complete map/filter/reduce/collect operations</td>
  </tr>
  <tr>
    <td><b>📚 Containers & Collections</b></td>
    <td>Enhanced collection types: Map, SafeMap, SortedMap, HashBasedMap, Optional, all with generics support</td>
  </tr>
  <tr>
    <td><b>💉 Dependency Injection</b></td>
    <td>Lightweight IOC container supporting component registration, retrieval, and lifecycle management</td>
  </tr>
  <tr>
    <td><b>🔒 Concurrency Control</b></td>
    <td>LockSupport (Java-like park/unpark), synchronization lock wrappers for goroutine control</td>
  </tr>
  <tr>
    <td><b>🔄 Data Conversion</b></td>
    <td>Struct property copying, type conversion, file export (CSV/JSON), serialization for comprehensive data processing</td>
  </tr>
  <tr>
    <td><b>⚡ Rule Engine</b></td>
    <td>Flexible business rule processing, rule chain building, middleware support for enterprise-grade rule management</td>
  </tr>
  <tr>
    <td><b>🌐 Network & More</b></td>
    <td>HTTP client, web crawler, database tools, logging system, algorithms, text processing, and utility functions</td>
  </tr>
</table>

<hr>

## 📦 Installation

Install the latest version using `go get`:

```bash
go get -u github.com/karosown/katool-go
```

> ⚠️ **System Requirements**
> - Go version >= 1.23.1
> - Generics support required

<hr>

## 🚀 Quick Start

### 🌊 Stream Processing Example

```go
package main

import (
    "fmt"
    "github.com/karosown/katool-go/container/stream"
)

type User struct {
    Name  string
    Age   int
    Money int
}

func main() {
    users := []User{
        {Name: "Alice", Age: 25, Money: 1000},
        {Name: "Bob", Age: 30, Money: 1500},
        {Name: "Charlie", Age: 35, Money: 2000},
    }
    
    // Parallel stream processing
    adultUsers := stream.ToStream(&users).
        Parallel().
        Filter(func(u User) bool { return u.Age >= 30 }).
        Sort(func(a, b User) bool { return a.Money > b.Money }).
        ToList()
    
    fmt.Printf("Adults sorted by income: %+v\n", adultUsers)
}
```

### 📚 Optional Container Example

```go
package main

import (
    "fmt"
    "strings"
    "github.com/karosown/katool-go/container/optional"
)

func main() {
    // Safe null handling
    nameOpt := optional.Of("John Doe")
    nameOpt.IfPresent(func(name string) {
        fmt.Printf("Username: %s\n", name)
    })
    
    // Chaining operations
    result := optional.MapTyped(optional.Of("  hello  "), strings.TrimSpace).
        Filter(func(s string) bool { return len(s) > 0 }).
        OrElse("empty")
    
    fmt.Printf("Processed: %s\n", result)
}
```

### ⚡ Rule Engine Example

```go
package main

import (
    "fmt"
    "github.com/karosown/katool-go/ruleengine"
)

type User struct {
    Name string
    Age  int
}

func main() {
    engine := ruleengine.NewRuleEngine[User]()
    
    // Register rules
    engine.RegisterRule("validate_age",
        func(user User, _ any) bool { return user.Age > 0 },
        func(user User, _ any) (User, any, error) {
            if user.Age < 18 {
                return user, "Minor user", ruleengine.FALLTHROUGH
            }
            return user, "Adult user", nil
        },
    )
    
    // Build and execute
    engine.NewBuilder("user_check").AddRule("validate_age").Build()
    result := engine.Execute("user_check", User{Name: "John", Age: 25})
    fmt.Printf("Result: %v\n", result.Result)
}
```

<hr>

## 🔧 Core Modules

### 📚 Optional Container

Safe null value handling inspired by Java's Optional:

- **Creation**: `Of()`, `Empty()`, `OfNullable()`
- **Checking**: `IsPresent()`, `IsEmpty()`
- **Retrieval**: `Get()`, `OrElse()`, `OrElseGet()`
- **Functional**: `Filter()`, `Map()`, `MapTyped()`, `IfPresent()`

### 🌊 Stream Processing

Java 8 Stream API-like operations with Go generics:

- **Creation**: `ToStream()`, `Parallel()`
- **Intermediate**: `Filter()`, `Map()`, `Sort()`, `Distinct()`
- **Terminal**: `ToList()`, `Reduce()`, `ForEach()`, `GroupBy()`

### ⚡ Rule Engine

Enterprise-grade business rule processing:

- **Flow Control**: `EOF` (terminate), `FALLTHROUGH` (skip)
- **Middleware**: Logging, monitoring, custom processing
- **Execution**: Rule chains, rule trees, batch processing

### 🔄 Data Conversion

Comprehensive data transformation tools:

- **Struct Operations**: Property copying, type conversion
- **File Export**: CSV, JSON serialization
- **Type Utilities**: Slice conversion, type casting

### 🔒 Concurrency Control

Java-inspired concurrency utilities:

- **LockSupport**: `Park()` and `Unpark()` operations
- **Synchronization**: `Synchronized()` code blocks
- **Thread Safety**: Concurrent collections and utilities

<hr>

## 💡 Best Practices

### 🌊 Stream Processing

```go
// ✅ Filter first, then process (efficient)
stream.ToStream(&data).
    Filter(func(item Item) bool { return item.IsValid() }).
    Map(func(item Item) Result { return item.Process() }).
    ToList()

// ❌ Process all, then filter (inefficient)
stream.ToStream(&data).
    Map(func(item Item) Result { return item.Process() }).
    Filter(func(result Result) bool { return result.IsValid() }).
    ToList()
```

### 📚 Optional Container

```go
// ✅ Use MapTyped for type safety
result := optional.MapTyped(opt, strings.TrimSpace).
    Filter(func(s string) bool { return len(s) > 0 }).
    OrElse("default")

// ❌ Avoid type assertions
result := opt.Map(func(s any) any { 
    return strings.TrimSpace(s.(string)) // Risky
}).OrElse("default")
```

### ⚡ Rule Engine

```go
// ✅ Single responsibility rules
engine.RegisterRule("validate_email",
    func(user User, _ any) bool { return user.Email != "" },
    func(user User, _ any) (User, any, error) {
        if !isValidEmail(user.Email) {
            return user, "Invalid email", ruleengine.EOF
        }
        return user, "Valid email", nil
    },
)

// ❌ Overly complex rules
engine.RegisterRule("validate_everything", // Too broad
    func(user User, _ any) bool { return true },
    func(user User, _ any) (User, any, error) {
        // Multiple validations in one rule - hard to maintain
    },
)
```

<hr>

## 👥 Contributing

We welcome contributions! Here's how to get started:

### 🚀 Quick Contribution Guide

1. **Fork & Clone**
   ```bash
   git clone https://github.com/your-username/katool-go.git
   cd katool-go
   ```

2. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Develop & Test**
   ```bash
   go test ./...
   go fmt ./...
   go vet ./...
   ```

4. **Commit & Push**
   ```bash
   git commit -m "feat: add your feature description"
   git push origin feature/your-feature-name
   ```

5. **Create Pull Request**

### ✅ Code Standards

- Pass all tests (`go test ./...`)
- Follow Go conventions (`go fmt`, `go vet`)
- Add documentation for public APIs
- Include test cases for new features
- Consider performance implications

### 🐛 Report Issues

Found a bug or have a suggestion? Please [create an issue](https://github.com/karosown/katool-go/issues).

<hr>

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

### 📜 License Summary

- ✅ Commercial use allowed
- ✅ Modification allowed
- ✅ Distribution allowed
- ✅ Private use allowed
- ❗ No warranty provided
- ❗ Authors not liable

### 🤝 Acknowledgments

Special thanks to:
- [Go Team](https://golang.org/) - For the excellent language
- [resty](https://github.com/go-resty/resty) - HTTP client library
- [rod](https://github.com/go-rod/rod) - Chrome automation
- [jieba](https://github.com/yanyiwu/gojieba) - Chinese word segmentation
- All contributors and users of this project

---

<div align="center">
  <sub>Built with ❤️ by <a href="https://github.com/karosown">Karosown Team</a></sub>
  <br>
  <sub>⭐ If this project helps you, please give us a Star!</sub>
</div> 