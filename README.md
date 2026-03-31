# k8s 管理平台

一个基于 Web 的 Kubernetes 可视化管理平台，用于多集群下的资源管理、运维操作与审计追踪。

## 项目目标

- 贴合 Rancher 简洁实用风格，可融合适度赛博科技风
- 所有能力围绕 k8s 高频运维操作设计，优先提升排障和变更效率
- 支持桌面端访问
- Secrets、Token 等敏感信息默认脱敏，权限严格遵循 RBAC，操作可审计可追溯
- 支持多语言切换
- 页面主题支持 `Light / Auto / Dark`

## 技术栈

- 前端：React + TypeScript + Material-UI (MUI) + Zustand
- 后端：Go + Gin + GORM，对接 Kubernetes API / 自定义后端服务
- 数据库：MySQL + Redis

## 开发运行

### 环境要求

- Node.js 20+
- npm 10+
- Go 1.23+
- MySQL 8+（默认连接：`localhost:3306`）
- Redis 7+（默认连接：`localhost:6379`）

### 首次安装依赖

前端（已配置国内源 `frontend/.npmrc`）：

```bash
cd frontend
npm install
```

后端（使用国内 Go 代理）：

```bash
cd backend
GOPROXY=https://goproxy.cn,direct go mod tidy
```

### 启动后端

启动前请确认 MySQL 与 Redis 已运行。当前默认配置：

- MySQL DSN：`root:123456@tcp(localhost:3306)/kubemanage?charset=utf8mb4&parseTime=True&loc=Local`
- Redis 地址：`localhost:6379`
- Redis 密码：空
- Redis DB：`0`
- MySQL 库名：`kubemanage`（需提前创建）

支持环境变量覆盖：

- `KM_LISTEN_ADDR`
- `KM_MYSQL_DSN`
- `KM_REDIS_ADDR`
- `KM_REDIS_PASS`
- `KM_REDIS_DB`

```bash
cd backend
go run ./cmd/server
```

默认监听：`http://localhost:8080`

### 启动前端

新开一个终端：

```bash
cd frontend
npm run dev
```

默认地址：`http://localhost:5173`

说明：前端已在 `vite.config.ts` 配置 `/api` 代理到 `http://localhost:8080`。

### 可选：使用 Makefile

在仓库根目录执行：

```bash
make backend-run
make frontend-dev
```

### 基础验证

后端健康检查：

```bash
curl http://localhost:8080/api/v1/healthz
```

集群列表接口：

```bash
curl http://localhost:8080/api/v1/clusters
```

当前集群接口：

```bash
curl http://localhost:8080/api/v1/clusters/current
```

切换当前集群：

```bash
curl -X POST http://localhost:8080/api/v1/clusters/switch \
  -H "Content-Type: application/json" \
  -d '{"name":"staging-cluster"}'
```

名称空间列表：

```bash
curl http://localhost:8080/api/v1/namespaces
```

创建名称空间：

```bash
curl -X POST http://localhost:8080/api/v1/namespaces \
  -H "Content-Type: application/json" \
  -d '{"name":"qa"}'
```

查看名称空间 YAML：

```bash
curl http://localhost:8080/api/v1/namespaces/qa/yaml
```

下载名称空间 YAML：

```bash
curl -OJ http://localhost:8080/api/v1/namespaces/qa/yaml/download
```

删除名称空间：

```bash
curl -X DELETE http://localhost:8080/api/v1/namespaces/qa
```

Deployment 列表：

```bash
curl http://localhost:8080/api/v1/deployments
```

查看 Deployment YAML：

```bash
curl http://localhost:8080/api/v1/deployments/web-api/yaml
```

更新 Deployment YAML：

```bash
curl -X PUT http://localhost:8080/api/v1/deployments/web-api/yaml \
  -H "Content-Type: application/json" \
  -d '{"yaml":"apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: web-api\n"}'
```

Pod 列表：

```bash
curl http://localhost:8080/api/v1/pods
```

查看 Pod 日志：

```bash
curl http://localhost:8080/api/v1/pods/web-api-7bf59f6f9c-abcde/logs
```

## 开发原则

- Go 实现高并发处理，统一连接池管理
- Redis 缓存热点数据，必要场景使用分布式锁
- 功能贴合 Kubernetes 官方资源标准，避免自定义非标准语义
- 关键操作具备审计能力，逐步接入监控与告警

## 功能需求拆分

### 长期蓝图（Vision）

#### 1. 集群管理
- 集群列表展示（名称、版本、状态、节点数）
- 集群详情查看
- 多集群导入、切换
- 集群事件查看

#### 2. 节点管理
- 节点列表、角色、版本、IP、状态、资源使用率（系统类型/CPU/内存/磁盘/POD数量/存在时间）
- 节点标签与污点查看、编辑
- 节点详情与事件查看
- 节点配置下载、编辑（YAML/Config）
- 节点删除

#### 3. 名称空间管理
- 名称空间查看、创建、删除、修改配置（YAML/Config）
- 名称空间配置下载

#### 4. 工作负载
- Deployment、StatefulSet、DaemonSet、CronJob、Job、Pod 全量管理
- 通用动作：Show Configuration、Edit Config/YAML、Delete、下载 YAML
- 增强动作：Execute Shell、Redeploy、Rollback、View Logs、手动触发/暂停任务、查看关联资源

#### 5. 服务发现
- Service、Ingress、HPA 的列表、详情、配置编辑、删除、下载 YAML
- 关联服务/Pod 查询，连通性测试与控制器状态（可选）

#### 6. 存储管理
- PV、PVC、StorageClass、ConfigMap、Secret 的列表、详情、配置编辑、删除、下载 YAML
- 关联资源查询
- Secret 脱敏展示与受控查看（可选）

#### 7. Policy 管理
- LimitRange、NetworkPolicy、ResourceQuota、Audit Log Policy 的列表、详情、配置编辑、删除、下载 YAML
- 策略生效状态与审计样例查看（可选）

#### 8. 日志与调试
- 全局日志概览、Pod/容器实时日志、事件日志、组件日志
- 条件筛选（命名空间/资源类型/时间范围/日志级别）
- 日志搜索、导出、下载、Web Terminal、故障排查指引

#### 9. 权限与用户
- 用户/用户组/角色/绑定/Token 管理
- 自定义角色与权限范围配置
- RBAC 权限校验（可选）
- 用户操作日志追踪

### 第一阶段 MVP（当前执行范围）

#### MVP 目标
在不引入过高复杂度的前提下，先交付一个可用的基础运维台，覆盖“查看-排障-变更”核心闭环。

#### MVP 范围（必须完成）
- 多集群：集群列表 + 集群切换 + 基础状态查看
- 名称空间：列表、创建、删除、YAML 查看/下载
- 工作负载：Deployment、Pod 列表与详情
- 运维动作：Pod 日志查看、Deployment YAML 编辑并应用、资源删除（二次确认）
- 服务发现：Service 列表与详情
- 配置资源：ConfigMap、Secret 列表与详情（Secret 脱敏）
- 基础审计：记录关键写操作（编辑、删除、触发动作）
- 基础权限：基于角色的页面/操作可见性控制（最小可用 RBAC）

#### MVP 暂不纳入（后续阶段）
- StatefulSet、DaemonSet、CronJob、Job 的完整操作能力
- Ingress、HPA、PV/PVC/StorageClass 全量能力
- NetworkPolicy、ResourceQuota、LimitRange、Audit Log Policy
- 用户组与 Token 全生命周期管理
- Web Terminal、高级日志检索、多语言、完整主题系统

#### MVP 验收标准
- 核心资源页可稳定完成“列表-详情-YAML 查看/下载”
- 至少 2 类写操作可用且有审计记录（例如 Deployment 编辑、Pod 删除）
- Secret 敏感字段不以明文展示
- 同一账号越权操作被拒绝并有日志

## 设计方案

### 架构设计（MVP）

- 前端：`React + MUI + Zustand`，按“资源模块 + 通用资源表格 + YAML 编辑器”组织
- 后端：`Gin` 提供统一 REST API，`service` 层封装 k8s client 与权限校验
- 数据层：MySQL 存平台数据（用户、审计、集群配置），Redis 缓存热点查询和会话
- 安全：接口级权限校验 + 关键操作审计 + Secret 脱敏策略

### 模块划分（MVP）

- `cluster`：集群注册、切换、健康状态
- `namespace`：基础 CRUD 与 YAML 能力
- `workload`：Deployment / Pod 查询与基础运维动作
- `service-discovery`：Service 查询
- `config`：ConfigMap / Secret 查询与脱敏
- `audit`：操作日志写入与查询
- `authz`：角色校验与接口拦截

## 任务计划

### 第一阶段 MVP 任务清单

1. 项目骨架初始化（前后端工程、配置、路由、基础页面）
2. 集群管理 MVP（列表、切换、状态）
3. 名称空间 MVP（列表、创建、删除、YAML 查看/下载）
4. Deployment / Pod MVP（列表、详情、日志、YAML 编辑）
5. Service / ConfigMap / Secret MVP（列表、详情、Secret 脱敏）
6. 权限与审计 MVP（关键写操作审计、最小 RBAC）
7. 联调与验收测试（核心流程、权限校验、审计验证）

## 任务状态

- 当前阶段：`第一阶段 MVP`
- 当前任务：`Deployment / Pod MVP（任务 4）`
- 当前状态：`待验收`
- 下一任务：`Service / ConfigMap / Secret MVP（任务 5）`
