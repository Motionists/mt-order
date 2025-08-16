# 外卖订单系统 (MT-Order)

一个基于 Go + Vue 3 + React 的全栈外卖订单管理系统，包含用户端、商家管理后台和后端 API。

## 项目结构

```
mt-order/
├── server/                 # Go 后端服务
│   ├── internal/           
│   │   ├── config/         # 配置管理
│   │   ├── database/       # 数据库连接
│   │   ├── handlers/       # HTTP 处理器
│   │   ├── middleware/     # 中间件
│   │   ├── models/         # 数据模型
│   │   ├── router/         # 路由配置
│   │   └── services/       # 业务逻辑
│   ├── main.go            # 入口文件
│   ├── go.mod             # Go 模块依赖
│   └── .env               # 环境配置
├── web/                    # Vue 3 用户端
│   ├── src/
│   │   ├── components/     # 公共组件
│   │   ├── pages/          # 页面组件
│   │   ├── stores/         # Pinia 状态管理
│   │   └── api/            # API 请求
│   └── package.json
├── admin/                  # React 管理后台
│   ├── src/
│   │   ├── components/     # 组件
│   │   ├── pages/          # 页面
│   │   └── services/       # API 服务
│   └── package.json
└── docker-compose.yml      # Docker 编排
```

## 技术栈

### 后端
- **Go 1.21** - 后端编程语言
- **Gin** - Web 框架
- **GORM** - ORM 数据库操作
- **MySQL** - 关系型数据库
- **Redis** - 缓存和会话存储
- **JWT** - 用户认证
- **CORS** - 跨域支持

### 前端 (用户端)
- **Vue 3** - 前端框架
- **Pinia** - 状态管理
- **Vue Router** - 路由管理
- **Element Plus** - UI 组件库
- **Axios** - HTTP 客户端
- **Vite** - 构建工具

### 管理后台
- **React 18** - 前端框架
- **Ant Design** - UI 组件库
- **Zustand** - 状态管理
- **React Router** - 路由管理
- **Recharts** - 图表组件

## 功能特性

### 用户端功能
- ✅ 用户注册/登录
- ✅ 浏览商家和菜品
- ✅ 购物车管理
- ✅ 下单和支付
- ✅ 订单查询和追踪
- ✅ 用户信息管理

### 管理后台功能
- ✅ 商家管理
- ✅ 菜品管理
- ✅ 订单管理
- ✅ 用户管理
- ✅ 数据统计和分析
- ✅ 系统设置

### 后端 API 功能
- ✅ RESTful API 设计
- ✅ JWT 身份验证
- ✅ 请求参数验证
- ✅ 错误处理
- ✅ 数据库事务
- ✅ 自动数据迁移

## 快速开始

### 环境要求

- **Go** >= 1.21
- **Node.js** >= 18
- **MySQL** >= 8.0
- **Redis** >= 6.0 (可选)

### 1. 克隆项目

```bash
git clone https://github.com/Motionists/mt-order.git
cd mt-order
```

### 2. 后端配置

```bash
cd server

# 安装依赖
go mod tidy

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置数据库连接等信息
```

配置 `.env` 文件：
```env
PORT=8080
DATABASE_URL=root:password@tcp(localhost:3306)/mt_order?charset=utf8mb4&parseTime=True&loc=Local
REDIS_URL=localhost:6379
JWT_SECRET=your-secret-key-here
```

### 3. 数据库配置

```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE mt_order CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
exit;
```

### 4. 启动后端服务

```bash
cd server
go run main.go
```

服务将在 `http://localhost:8080` 启动

### 5. 启动前端服务

#### 用户端 (Vue)
```bash
cd web
npm install
npm run dev
```

访问地址：`http://localhost:5173`

#### 管理后台 (React)
```bash
cd admin
npm install
npm run dev
```

访问地址：`http://localhost:3000`

### 6. 使用 Docker (可选)

```bash
# 启动数据库服务
docker-compose up -d mysql redis

# 构建并启动所有服务
docker-compose up -d
```

## API 文档

### 认证接口
```
POST /api/auth/register    # 用户注册
POST /api/auth/login       # 用户登录
```

### 商家接口
```
GET  /api/merchants        # 获取商家列表
GET  /api/merchants/:id    # 获取商家详情
GET  /api/merchants/:id/dishes  # 获取商家菜品
```

### 购物车接口 (需要认证)
```
GET    /api/cart           # 获取购物车
POST   /api/cart           # 添加商品到购物车
PUT    /api/cart/:id       # 更新购物车商品
DELETE /api/cart/:id       # 删除购物车商品
```

### 订单接口 (需要认证)
```
GET  /api/orders           # 获取订单列表
POST /api/orders           # 创建订单
GET  /api/orders/:id       # 获取订单详情
```

## 开发指南

### 后端开发

#### 添加新的 API 端点
1. 在 `internal/models/` 中定义数据模型
2. 在 `internal/handlers/` 中创建处理器
3. 在 `internal/router/router.go` 中注册路由

#### 数据库迁移
```bash
# GORM 会自动处理数据库迁移
# 修改模型后重启服务即可
go run main.go
```

### 前端开发

#### 添加新页面 (Vue)
1. 在 `web/src/pages/` 中创建页面组件
2. 在 `web/src/router/` 中配置路由
3. 在 `web/src/stores/` 中管理状态

#### 添加新组件 (React)
1. 在 `admin/src/components/` 中创建组件
2. 在 `admin/src/pages/` 中使用组件
3. 在 `admin/src/services/` 中处理 API 调用

## 部署

### 生产环境构建

#### 后端
```bash
cd server
CGO_ENABLED=0 GOOS=linux go build -o mt-order-server .
```

#### 前端
```bash
# 用户端
cd web
npm run build

# 管理后台
cd admin
npm run build
```

### Docker 部署
```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d
```

## 测试

### 后端测试
```bash
cd server
go test ./...
```

### API 测试
```bash
# 使用 curl 测试注册接口
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com", 
    "password": "123456"
  }'
```

## 常见问题

### 1. 数据库连接失败
- 检查 MySQL 服务是否启动
- 验证数据库连接字符串
- 确认数据库用户权限

### 2. 端口冲突
- 修改 `.env` 文件中的 `PORT` 配置
- 检查端口占用：`lsof -i :8080`

### 3. 前端构建失败
- 清除 node_modules：`rm -rf node_modules && npm install`
- 检查 Node.js 版本是否兼容

## 贡献指南

1. Fork 项目
2. 创建功能分支：`git checkout -b feature/new-feature`
3. 提交更改：`git commit -am 'Add new feature'`
4. 推送分支：`git push origin feature/new-feature`
5. 提交 Pull Request

## 许可证

MIT License

## 联系方式

- 项目地址：https://github.com/Motionists/mt-order
- 问题反馈：https://github.com/Motionists/mt-order/issues

---

**开发团队：Motionists**