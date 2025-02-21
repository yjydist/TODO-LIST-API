# Todo List API

在此项目中, 你需要开发一个 RESTful API, 以允许用户管理其待办事项列表(Todo List)

先前的后端项目仅关注 CRUD 操作, 但是该项目还需要你实施用户身份验证

## 目标
你将从该项目中学到的技能包括:  
- 用户身份验证
- 模式设计和数据库
- REATful API 设计
- CRUD 操作
- 错误处理
- 安全

## 要求
你需要开发带有以下接口的 RESTful API:  
- 用户注册以创建新用户
- 登录接口, 用于用户认证并生成令牌
- 用于管理待办事项清单的 CRUD 操作
- 实现用户身份验证，仅允许授权用户访问待办事项列表
- 实施错误处理和安全措施
- 使用数据库存储用户和待办事项列表数据(可以使用任何数据库)
- 实施适当的数据验证
- 实现待办事项列表的分页和筛选功能

以下列出接口和请求与响应的详细信息:

### 用户注册
使用一下请求注册用户
```json
POST /register
{
  "name": "John Doe",
  "email": "john@doe.com",
  "password": "password"
}
```

这将验证给定的详细信息, 确保电子邮件是唯一的, 并将用户详细信息存储在数据库中

在将密码存储在数据库之前, 请确保把密码哈希

如果注册成功, 返回带有可以用于用户认证的令牌的响应

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
}
```

该令牌可以是 JWT 或随机字符串, 用于身份验证. 具体实现细节由自己决定

### 用户登录
使用以下请求对用户进行身份验证:  
```json
POST /login
{
  "email": "john@doe.com",
  "password": "password"
}
```

这将验证提供的邮箱和密码, 若认证成功, 则返回一个令牌
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
}
```

### 创建待办事项
使用以下请求创建一个新的待办事项:
```json
POST /todos
{
  "title": "Buy groceries",
  "description": "Buy milk, bread, eggs, and butter"
}
```

用户需将登录端点获取的令牌放入请求头中验证身份

请使用 `"Authorization"` 头, 其值为令牌

若令牌缺失或无效, 则返回错误和状态码 401  

```json
{
  "message": "Unauthorized"
}
```

创建待办事项成功后, 请回复该事项的详细信息。
```json
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Buy milk, eggs, and bread"
}
```

### 更新待办事项
使用以下请求更新现有的待办事项:
```json
PUT /todos/1
{
  "title": "Buy groceries",
  "description": "Buy milk, eggs, bread, and cheese"
}
```

就像创建待办事项一样, 用户必须发送收到的令牌

另外, 请确保验证用户有权更新待办事项, 即用户是他们正在更新的 TODO 项目的创建者

如果用户未授权更新项目, 则使用错误和状态代码 403 响应

```json
{
  "message": "Forbidden"
}
```

成功更新待办事项后, 请响应该项目的更新详细信息

```json
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Buy milk, eggs, bread, and cheese"
}
```

### 删除待办事项

使用以下请求删除现有的待办事项:

```
DELETE /todos/1
```

必须对用户进行身份验证并授权删除待办事项. 成功删除后, 使用状态代码 `204` 响应

### 获取待办事项

使用以下请求获取待办事项列表:

```
GET /todos?page=1&limit=10
```

必须对用户进行身份验证以访问任务, 并且应分页响应, 响应待办事项清单以及分页详细信息

### 额外任务

实施待办事项列表的过滤和排序

实施 API 的单元测试

实施限制 API 的费率和节流

实施身份验证的刷新令牌机制



