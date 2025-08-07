# 轻量级 Web 框架

本框架是基于 Go 语言开发的轻量级 Web 工具，专注于简洁性与扩展性，提供路由管理、中间件机制、模板渲染、静态文件服务等核心功能，适合快速构建中小型 Web 应用。

## 框架特性

### 1. 高效路由系统

- 基于 Trie 树实现路由匹配，支持静态路由、动态路由（如`:param`路径参数）及通配符路由（如`*wildcard`），匹配效率高。
- 支持 GET、POST 等标准 HTTP 方法，提供直观的路由注册接口，便于快速绑定处理函数。
- 支持路由分组功能，可按业务模块对路由归类，方便管理分组内的中间件和路由规则。

### 2. 灵活的中间件机制

- 支持全局中间件与分组中间件，可按需扩展功能（如日志记录、权限校验、跨域处理等）。
- 内置两款实用中间件：
    - `Logger`：自动记录请求的 HTTP 方法、路径、状态码及处理耗时，便于调试和监控。
    - `Recovery`：捕获请求处理过程中的 panic 异常，自动恢复服务并返回 500 错误响应，避免服务崩溃。

### 3. 上下文管理

- 封装

    ```
    Context
    ```

    结构体，整合

    ```
    http.Request
    ```

    与

    ```
    http.ResponseWriter
    ```

    ，提供便捷的参数获取与响应构造能力：

    - 快速获取路径参数（`Param`）、查询参数（`Query`）、表单数据（`PostForm`）。
    - 支持多种响应类型：HTML 模板渲染、JSON 数据、纯文本（String）、二进制数据（Data）等。
    - 内置`Fail`方法统一错误响应格式，简化异常处理流程。

### 4. 模板渲染引擎

- 集成 Go 标准库`html/template`，支持批量加载模板文件并缓存，提升渲染性能。
- 允许注册自定义模板函数（如日期格式化、字符串处理等），增强模板的动态渲染能力。
- 与上下文深度集成，可直接在处理函数中传递数据并指定模板名称渲染。

### 5. 静态文件服务

- 提供静态资源（CSS、JS、图片等）访问能力，支持将本地目录映射为 URL 路径，简化前端资源部署。
- 自动处理静态文件的 MIME 类型，无需手动配置即可正确响应资源类型。

### 6. 错误处理

- 完善的异常捕获机制，配合`Recovery`中间件可优雅处理运行时错误，保障服务稳定性。
- 支持自定义错误响应内容，便于统一 API 的错误格式规范。

## 快速开始

### 环境要求

- Go 1.16+（推荐使用 Go 1.21 + 以支持模块特性及最新标准库功能）。

### 安装框架

```bash
go get github.com/LudensCS/Web/web@v0.0.0-20250804093508-f50e59a4ec6f
```

### 基础使用示例

```go
package main

import (
    "net/http"
    "github.com/LudensCS/Web/web"
)

func main() {
    // 初始化框架实例
    engine := web.New()
    
    // 注册全局中间件（日志与错误恢复）
    engine.Use(web.Logger(), web.Recovery())
    
    // 注册路由
    engine.GET("/", func(c *web.Context) {
        c.HTML(http.StatusOK, "index.tmpl", web.H{
            "title": "轻量级Web框架",
        })
    })
    
    // 启动服务（监听9999端口）
    if err := engine.Run(":9999"); err != nil {
        panic("服务启动失败: " + err.Error())
    }
}
```

## 核心功能使用指南

### 1. 路由管理

#### 基本路由注册

通过`GET`、`POST`等方法直接注册路由：

```go
// 注册GET路由
engine.GET("/users", func(c *web.Context) {
    c.JSON(http.StatusOK, web.H{"message": "用户列表"})
})

// 注册POST路由
engine.POST("/users", func(c *web.Context) {
    username := c.PostForm("username")
    c.JSON(http.StatusCreated, web.H{"message": "创建用户: " + username})
})
```

#### 动态路由与参数获取

支持`:`定义路径参数，`*`定义通配符：

```go
// 路径参数示例
engine.GET("/users/:id", func(c *web.Context) {
    userID := c.Param("id") // 获取路径参数"id"
    c.String(http.StatusOK, "用户ID: %s", userID)
})

// 通配符示例
engine.GET("/files/*path", func(c *web.Context) {
    filePath := c.Param("path") // 获取通配符匹配的完整路径
    c.String(http.StatusOK, "文件路径: %s", filePath)
})
```

#### 路由分组

通过`Group`方法创建路由分组，便于隔离不同模块：

```go
// 创建"/api"分组
apiGroup := engine.Group("/api")
// 为分组注册专属中间件（仅作用于该分组内的路由）
apiGroup.Use(authMiddleware)

// 分组内的路由（实际访问路径为"/api/users"）
apiGroup.GET("/users", func(c *web.Context) {
    c.JSON(http.StatusOK, web.H{"data": []string{"user1", "user2"}})
})
```

### 2. 中间件使用

#### 内置中间件

框架自带`Logger`和`Recovery`，直接通过`Use`注册：

```go
// 注册全局中间件（作用于所有路由）
engine.Use(web.Logger(), web.Recovery())
```

#### 自定义中间件

中间件为`HandleFunc`类型，可通过嵌套逻辑扩展功能：

```go
// 自定义鉴权中间件
func authMiddleware(c *web.Context) {
    // 检查Token
    token := c.Query("token")
    if token == "" {
        c.Fail(http.StatusUnauthorized, "缺少Token")
        return
    }
    // 验证通过，继续执行后续中间件或处理函数
    c.Next()
}

// 注册到路由或分组
engine.Use(authMiddleware)
```

### 3. 模板渲染

#### 配置模板

```go
// 注册自定义模板函数
engine.SetFuncMap(template.FuncMap{
    "FormatAsDate": func(t time.Time) string {
        // 自定义日期格式化逻辑
        return t.Format("2006-01-02")
    },
})

// 加载模板文件（支持通配符）
engine.LoadHTMLGlob("templates/*")
```

#### 渲染模板

在处理函数中通过`HTML`方法渲染：

```go
engine.GET("/articles/:id", func(c *web.Context) {
    c.HTML(http.StatusOK, "article.tmpl", web.H{
        "title":   "文章详情",
        "content": "这是一篇测试文章",
        "date":    time.Now(),
    })
})
```

### 4. 静态文件服务

通过`Static`方法映射静态资源目录：

```go
// 将本地"./static"目录映射到URL路径"/assets"
engine.Static("/assets", "./static")
// 访问"./static/css/style.css"时，URL为"/assets/css/style.css"
```