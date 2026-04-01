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
- `KM_K8S_ADAPTER_MODE`（建议仅使用 `live` / `auto`，当前研发目标为 real-only）
- `KM_SECRET_KEY`（连接敏感字段加密密钥，留空则保持明文兼容模式）
- `KM_TERMINAL_SESSION_TTL_SECONDS`（终端会话 TTL 秒数，默认 `120`）
- `KM_AUTH_JWT_SECRET`（认证 JWT 密钥）
- `KM_AUTH_ACCESS_TTL_SECONDS`（Access Token 有效期秒数，默认 `3600`）
- `KM_AUTH_REFRESH_TTL_SECONDS`（Refresh Token 有效期秒数，默认 `604800`）
- `KM_AUTH_COMPAT_STAGE_KEEP`（旧头兼容阶段数，当前默认 `1`）
写操作确认头：`X-Action-Confirm: CONFIRM`（关键写操作必填）
失败排查头：`X-Request-Id`（后端响应会回传）

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

Service 列表：

```bash
curl http://localhost:8080/api/v1/services
```

Ingress 列表：

```bash
curl http://localhost:8080/api/v1/ingresses
```

Ingress 关联 Service：

```bash
curl http://localhost:8080/api/v1/ingresses/web-api-ing/services
```

HPA 列表：

```bash
curl http://localhost:8080/api/v1/hpas
```

HPA 目标查询：

```bash
curl http://localhost:8080/api/v1/hpas/web-api-hpa/target
```

PV 列表：

```bash
curl http://localhost:8080/api/v1/pvs
```

PVC 列表：

```bash
curl http://localhost:8080/api/v1/pvcs
```

StorageClass 列表：

```bash
curl http://localhost:8080/api/v1/storageclasses
```

ConfigMap 列表：

```bash
curl http://localhost:8080/api/v1/configmaps
```

Secret 列表（脱敏）：

```bash
curl http://localhost:8080/api/v1/secrets
```

查看当前角色信息：

```bash
curl -H "X-User-Role: admin" http://localhost:8080/api/v1/auth/me
```

默认管理员登录（首次启动自动初始化 `admin/123456`）：

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
```

管理员创建用户（需携带管理员登录获取的 Access Token）：

```bash
curl -X POST http://localhost:8080/api/v1/auth/users \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "X-Action-Confirm: CONFIRM" \
  -H "Content-Type: application/json" \
  -d '{"username":"readonly1","password":"123456","role":"readonly","allowedNamespaces":["dev"]}'
```

管理员查看用户列表：

```bash
curl -H "Authorization: Bearer <ACCESS_TOKEN>" \
  http://localhost:8080/api/v1/auth/users
```

管理员禁用用户：

```bash
curl -X PATCH http://localhost:8080/api/v1/auth/users/readonly1/status \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "X-Action-Confirm: CONFIRM" \
  -H "Content-Type: application/json" \
  -d '{"isActive":false}'
```

管理员重置用户密码：

```bash
curl -X POST http://localhost:8080/api/v1/auth/users/readonly1/reset-password \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "X-Action-Confirm: CONFIRM" \
  -H "Content-Type: application/json" \
  -d '{"password":"654321"}'
```

管理员编辑用户角色与授权范围：

```bash
curl -X PATCH http://localhost:8080/api/v1/auth/users/readonly1 \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "X-Action-Confirm: CONFIRM" \
  -H "Content-Type: application/json" \
  -d '{"role":"standard-user","allowedNamespaces":["dev","qa"]}'
```

查看公开认证源列表（登录页使用）：

```bash
curl http://localhost:8080/api/v1/auth/providers/public
```

管理员创建 LDAP 认证源：

```bash
curl -X POST http://localhost:8080/api/v1/auth/providers \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "X-Action-Confirm: CONFIRM" \
  -H "Content-Type: application/json" \
  -d '{"name":"corp-ldap","type":"ldap","config":{"url":"ldap://127.0.0.1:389","baseDN":"dc=example,dc=com"}}'
```

管理员设置默认认证源：

```bash
curl -X POST http://localhost:8080/api/v1/auth/providers/2/default \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "X-Action-Confirm: CONFIRM"
```

以 viewer 角色尝试删除名称空间（应返回 403）：

```bash
curl -X DELETE -H "X-User-Role: viewer" \
  http://localhost:8080/api/v1/namespaces/default
```

以 admin 角色查看审计日志：

```bash
curl -H "X-User-Role: admin" http://localhost:8080/api/v1/audits
```

一键执行 MVP 联调验收脚本：

```bash
bash scripts/mvp_smoke_test.sh
```

可选自定义后端地址：

```bash
BASE_URL=http://127.0.0.1:8080 bash scripts/mvp_smoke_test.sh
```

重构阶段统一验收脚本（R5）：

```bash
bash scripts/rebuild_qa.sh
```

P601 用户管理冒烟脚本：

```bash
bash scripts/p601_user_management_smoke_test.sh
```

P701 细粒度授权冒烟脚本：

```bash
bash scripts/p701_fine_grained_auth_smoke_test.sh
```

P801 认证源管理冒烟脚本：

```bash
bash scripts/p801_auth_provider_smoke_test.sh
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

### 第二阶段（原规划，已暂停）

#### 阶段目标
在第一阶段基础上补齐“运维深度能力”，聚焦真实 Kubernetes 资源联动和高级排障能力。

#### 阶段范围（计划）
- 工作负载扩展：StatefulSet、DaemonSet、Job、CronJob（列表/详情/YAML/基础操作）
- 服务发现扩展：Ingress、HPA（列表/详情/关联资源）
- 存储扩展：PV/PVC/StorageClass（列表/详情/关联关系）
- 日志与调试扩展：日志过滤、下载、基础 Web Terminal 接口预留
- 权限体系扩展：用户组、角色绑定、命名空间级权限控制
- 审计扩展：按用户、时间范围、资源类型筛选审计日志

#### 第二阶段验收标准（草案）
- 新增资源模块具备“列表-详情-YAML 查看/编辑”闭环
- 至少 3 类关联查询能力可用（如 Service->Pod、PVC->PV、Ingress->Service）
- 日志模块支持基础筛选并可导出
- 权限控制支持命名空间级别差异化授权

### 第二阶段（Rancher 风格重构路线）

#### 阶段目标
将现有“功能验证型后台”重构为“Rancher 风格控制台”，实现整体布局、导航结构、页面交互的一致性，目标体验达到 Rancher 相似度约 80%。

#### 阶段范围（重构优先级）
- `R1`：统一壳层（顶部栏 + 左侧菜单 + 主内容区 + 路由骨架）
- `R2`：通用资源页面框架（筛选栏 + 资源表 + 详情抽屉 + YAML 编辑区）
- `R3`：迁移现有模块到新框架（Cluster/Namespace/Workload/Resource/AuthAudit）
- `R4`：Rancher 风格视觉细节与交互统一（菜单分组、面包屑、页面级操作栏）

#### 第二阶段验收标准（重构）
- 左侧菜单与顶部栏稳定可用，核心模块可通过路由直达
- 所有已实现模块完成新框架迁移，不再使用临时标签页切换
- 页面布局和交互风格与 Rancher 控制台保持高一致性

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

### 架构演进（第二阶段）

- 引入 `k8s client adapter` 层，逐步替换当前示例数据服务
- 在 `service` 层增加资源关联查询能力（ownerReferences / label selector）
- 扩展审计存储结构，支持分页与条件检索
- 为日志与终端能力预留长连接网关接口

### 架构演进（第三阶段）

- 将当前示例数据服务逐步替换为真实 `k8s client adapter`
- 引入集群连接配置存储，支持多集群导入、测试连接与切换
- 第三阶段第一优先级为真实集群接入，不再继续扩展基于示例数据的功能面
- 集群导入首批同时支持 `kubeconfig` 与 `API Server + Token + CA` 两种方式
- 真实读链路优先覆盖 `Cluster / Namespace / Workloads / Service Discovery / Storage`
- 在真实读链路稳定后，再推进真实写操作保护、实时日志与终端网关
- 新增约束（2026-04-01）：停止继续扩展或返回 mock 数据，后续开发与联调统一依赖真实集群数据

### 前端重构设计（Rancher 风格）

- 布局：`ShellLayout`（`TopBar` + `SideNav` + `ContentArea`）
- 导航：按域分组（Cluster / Workloads / Service Discovery / Storage / Security / Audit）
- 页面容器：统一 `PageHeader + Toolbar + TableArea + DetailPanel`
- 状态：路由状态统一由 `react-router` 管理，业务状态保持 `Zustand`
- 左侧导航采用两级结构：一级菜单 `Cluster / Workloads / Service Discovery / Storage / Security`，点击一级菜单展开二级子菜单
- Workloads 二级菜单：`Deployment / Pod / StatefulSet / DaemonSet / Job / CronJob`
- Service Discovery 二级菜单：`Service / Ingress / HPA`
- Storage 二级菜单：`PersistentVolumes / PersistentVolumeClaims / StorageClasses / ConfigMaps / Secrets`
- 所有二级资源页面布局统一对齐 Workloads（筛选栏 + 资源表 + 详情抽屉）

## 任务计划

### 第一阶段 MVP 任务清单

1. 项目骨架初始化（前后端工程、配置、路由、基础页面）
2. 集群管理 MVP（列表、切换、状态）
3. 名称空间 MVP（列表、创建、删除、YAML 查看/下载）
4. Deployment / Pod MVP（列表、详情、日志、YAML 编辑）
5. Service / ConfigMap / Secret MVP（列表、详情、Secret 脱敏）
6. 权限与审计 MVP（关键写操作审计、最小 RBAC）
7. 联调与验收测试（核心流程、权限校验、审计验证）

### 第二阶段任务清单（规划）

1. K8s 真实集群数据接入（替换示例数据）
2. 工作负载扩展（StatefulSet/DaemonSet/Job/CronJob）
3. 服务发现扩展（Ingress/HPA）与关联关系查询
4. 存储管理扩展（PV/PVC/StorageClass）
5. 日志与调试增强（筛选、导出、终端接口预留）
6. 权限与审计增强（命名空间级授权、审计筛选）
7. 第二阶段联调与验收测试

### 第三阶段任务清单（扩展方案）

1. 真实集群接入与导入（`kubeconfig` / `API Server + Token + CA`）
2. `k8s client adapter` 与连接配置存储落地
3. 真实资源读链路切换（Cluster/Namespace/Workloads/Service Discovery/Storage）
4. 真实写操作安全化（确认、失败回显、审计增强）
5. 真实日志与终端能力（多容器、流式日志、exec/terminal）
6. 第三阶段联调与验收

### 第五阶段任务清单（账号认证与用户体系）

1. 账号认证基础（管理员创建账号、账号密码登录、刷新、登出）
2. 用户模型与默认管理员初始化（默认 `admin/123456`，首次启动自动创建）
3. 角色与作用域策略落地（参考 Rancher 分层：`admin`/`standard-user`/`readonly`）
4. 中间件升级（优先 Bearer Token，兼容 `X-User/X-User-Role` 过渡 1 个阶段）
5. 管理员用户管理接口（创建用户、角色分配、授权命名空间）
6. 联调与验收（后端测试、前端构建、认证冒烟脚本）

### P301 计划（2026-04-01）

#### 范围定义
- 集群导入：支持 `kubeconfig` 文本/文件导入
- 集群导入：支持 `API Server`、`Token`、`CA` 方式导入
- 集群连接测试：导入前或导入后执行连通性校验
- 集群切换：当前集群切换真正影响后端数据源
- 真实读链路首批覆盖：`clusters`、`namespaces`
- 本轮不纳入：真实写操作、真实日志流、真实终端

#### 开发拆分
1. 集群连接配置模型与存储结构设计
2. 后端导入/测试连接/切换接口
3. `k8s client adapter` 初版与 Cluster/Namespace 真实读取
4. 前端集群导入与连接测试页面
5. 联调与回归（接口测试、前端构建、P301 冒烟脚本）

#### 技术方案（执行版）

##### 连接配置模型
- 存储实体：`cluster_connections`
- 核心字段：
  - `id`
  - `name`
  - `mode`：`kubeconfig` / `token`
  - `api_server`
  - `kubeconfig_content`
  - `bearer_token`
  - `ca_cert`
  - `skip_tls_verify`
  - `is_default`
  - `status`：`connected` / `failed` / `unknown`
  - `last_checked_at`
  - `last_error`
- 安全要求：
  - `bearer_token`、`kubeconfig_content`、`ca_cert` 不在列表接口明文返回
  - 前端详情页仅展示脱敏摘要

##### 导入协议
- `kubeconfig` 导入：
  - 表单字段：`name`、`kubeconfigContent`
  - 后端自动解析当前 context、server、certificate 信息
- `API Server + Token + CA` 导入：
  - 表单字段：`name`、`apiServer`、`bearerToken`、`caCert`、`skipTLSVerify`
  - 允许 `caCert` 为空且 `skipTLSVerify=true` 的测试模式

##### 首批接口设计
- `GET /api/v1/clusters/connections`
- `POST /api/v1/clusters/connections/import/kubeconfig`
- `POST /api/v1/clusters/connections/import/token`
- `POST /api/v1/clusters/connections/test`
- `POST /api/v1/clusters/connections/:id/activate`
- `GET /api/v1/clusters/live`
- `GET /api/v1/namespaces/live`

##### adapter 分层
- `handlers`
  - 仅负责协议解析与响应格式
- `service`
  - 负责导入校验、连接测试、当前集群切换
- `adapter`
  - 新增 `k8s client adapter`，统一封装 `client-go`
- `infra/store`
  - 承担连接配置持久化

##### 数据源切换策略
- 保留现有示例数据服务作为 fallback
- 当“当前集群”为真实连接且连接测试通过时，`clusters/namespaces` 优先走 live adapter
- 当连接不可用时，接口返回明确错误，不静默回退到 mock 数据

##### P301 首批验收标准
- 能成功导入 `kubeconfig` 集群连接
- 能成功导入 `API Server + Token + CA` 集群连接
- 能执行连接测试并返回成功/失败原因
- 当前集群激活后，`clusters` / `namespaces` 能从真实集群读取
- 敏感字段不在列表接口明文泄漏
- 后端 `go test ./...`、前端 `npm run build`、`P301` 冒烟脚本通过

### P205 MVP 计划（2026-04-01）

#### 范围定义
- 资源范围：`Pod` 日志查看增强
- 能力范围：日志关键字筛选、大小写敏感切换、仅显示匹配行、日志导出、终端接口预留
- 本轮不纳入：实时日志流、时间范围筛选、多容器切换、真实 Web Terminal

#### 开发拆分
1. 后端日志查询参数兼容与终端占位接口
2. 前端日志弹窗增强（筛选、导出、终端入口）
3. 联调与回归（接口测试、前端构建、P205 冒烟脚本）

#### P205 验收标准（执行版）
- Pod 日志弹窗支持关键字筛选，空关键字时展示全部日志
- 支持大小写敏感开关与“仅显示匹配行”
- 支持将当前日志内容导出为 `.log` 文件
- 后端提供终端预留接口，并明确返回“未启用”占位响应
- 前端提供终端入口，但不伪装为真实终端能力
- 后端 `go test ./...` 通过，前端 `npm run build` 通过
- 提供 `scripts/p205_smoke_test.sh` 基础冒烟脚本

### P205 第二轮完善（2026-04-01）

#### 范围定义
- 日志跟随刷新：支持弹窗内周期刷新，模拟 `follow` 调试体验
- 交互补强：支持复制日志、显示行数/匹配数、快速清空筛选
- 参数联动：前端筛选条件回传后端接口，减少前后端行为偏差

#### 开发拆分
1. 后端补充日志数据变化与筛选参数测试
2. 前端日志弹窗增加跟随刷新、复制、统计与刷新控制
3. 联调与回归，补充 P205 冒烟脚本覆盖

### P206 第一版计划（2026-04-01）

#### 范围定义
- 命名空间级授权：`viewer` 只读，`operator` 仅允许写入 `dev`，`admin` 不受命名空间限制
- 审计筛选：支持按 `user`、`role`、`method`、`path`、`statusCode` 过滤，并支持 `limit`
- 前端页面：补充授权说明、审计筛选栏与筛选结果展示

#### 开发拆分
1. 后端鉴权能力增强（命名空间作用域判定）
2. 后端审计查询增强（筛选参数 + 结果过滤）
3. 前端权限与审计页面增强（授权说明 + 审计筛选）
4. 联调与回归（接口测试、前端构建、P206 冒烟脚本）

### P207 计划（2026-04-01）

#### 范围定义
- 汇总第二阶段已交付能力，执行统一回归验收
- 串联 `P202`、`P203`、`P205`、`P206` 冒烟脚本
- 执行后端测试、前端构建，并沉淀统一验收入口
- 同步 README/TASKS 的第二阶段状态

#### 开发拆分
1. 编写第二阶段统一验收脚本
2. 执行回归验证并修正发现的问题
3. 同步任务状态与验收结果

### P202 重启计划（2026-04-01）

#### 范围定义
- 资源范围：`StatefulSet`、`DaemonSet`、`Job`、`CronJob`
- 能力范围：列表、详情、YAML 查看、YAML 编辑保存
- 本轮不纳入：真实集群接入、批量操作、高级筛选、任务手动触发与暂停

#### 开发拆分
1. 后端 `service` 扩展（4 类资源示例数据与 YAML 存取）
2. 后端 `handler/router` 扩展（读写接口 + RBAC + 写操作审计）
3. 前端工作负载页重建（4 类资源表格/详情抽屉/YAML 编辑）
4. 联调与回归（接口测试 + 前端构建 + 冒烟脚本）

#### P202 验收标准（执行版）
- 4 类资源均支持“列表-详情-YAML 查看-YAML 保存”闭环
- `viewer` 无法写入 YAML（403），`operator/admin` 可写入（204）
- 前端 `workloads` 页面可正常访问，`npm run build` 通过
- 后端 `go test ./...` 通过，且新增接口纳入路由测试

### P402 执行计划（2026-04-01）

#### 范围定义
- 目标范围：打通真实终端 `exec` 最小闭环，覆盖会话创建、WebSocket 通道、后端 stream 桥接
- 能力范围：`sessionId` 一次性消费、TTL 失效、容器参数透传、基础异常回传（过期/无效/不可达）
- 本轮不纳入：终端录屏、多人共享会话、历史回放、复杂命令模板

#### 设计要点
- 会话模型：`Create -> Consume -> Expire`，默认短 TTL，消费后立即失效
- 通道模型：`/terminal/sessions` 负责鉴权和会话签发，`/terminal/ws` 负责输入输出转发
- 执行模型：通过 `client-go remotecommand` 连接 `pods/exec`，将 WebSocket 文本流桥接到 stdin/stdout
- 安全模型：保持现有写操作确认头与 RBAC 作用域校验，不放宽权限边界

#### 开发拆分
1. `P402-A`：终端会话管理（会话创建、TTL、sessionId 校验）
2. `P402-B`：WebSocket exec 通道（输入输出桥接、异常回传、连接生命周期）
3. `P402-C`：回归与验收（`go test`、`npm run build`、`p402` 冒烟脚本）

#### P402 验收标准（执行版）
- 创建终端会话返回 `sessionId` 与可连接的 `wsPath`
- 创建终端会话响应包含 `ttlSeconds` 与 `expiresAt`，便于前端提示会话时效
- 同一 `sessionId` 仅允许消费一次，重复使用返回无效会话错误
- 过期 `sessionId` 返回会话过期错误，且不会触发真实 `exec`
- WebSocket 成功连通后可向容器发送输入并接收输出
- 后端 `go test ./...` 通过，前端 `npm run build` 通过，补充 `scripts/p402_smoke_test.sh`

#### P402-C 真实集群 e2e 验收执行说明
1. 验收前置：
`KM_K8S_ADAPTER_MODE=live`，且已有可用激活连接（`/api/v1/clusters/connections/:id/activate` 已成功）。
2. 创建会话：
`POST /api/v1/pods/:name/terminal/sessions`（带 `X-Action-Confirm: CONFIRM`）并记录返回 `sessionId/wsPath`。
3. 建立连接：
使用 `wsPath` 建立 WebSocket，建议附带命令参数 `?command=sh`。
4. 输入输出校验：
发送 `echo P402_E2E_OK`，期望回包中可见 `P402_E2E_OK`。
5. 失败分支校验：
同一 `sessionId` 二次连接应失败；超 TTL 后连接应返回过期错误；跨用户/角色复用应返回 owner mismatch。
6. 结果沉淀：
将本次验收结果写入 `TASKS.md` 完成记录，并更新 `T095` 状态。

#### P402-C 验收记录模板

```md
- 日期：YYYY-MM-DD
- 环境：cluster=<name>, namespace=<ns>, pod=<pod>, container=<container>
- 会话创建：通过/失败（状态码、requestId）
- WebSocket 连接：通过/失败（错误信息）
- 收发验证：通过/失败（关键输出）
- 失败分支：重复消费=通过/失败；TTL过期=通过/失败；owner绑定=通过/失败
- 结论：通过/不通过
- 备注：问题与后续动作
```

#### P402-C 本次验收结果（2026-04-01）
- 环境：`KM_K8S_ADAPTER_MODE=live`，激活连接 `id=1(test1)`，验证 Pod `nginx-smoke-5c66c96d88-cbh7m`
- 会话创建：通过（`201`，返回 `sessionId/wsPath`）
- WebSocket 收发：通过（发送 `echo P402_E2E_OK`，回显包含 `P402_E2E_OK`）
- 失败分支：重复消费=`401 invalid`，TTL 过期=`401 expired`，跨用户复用=`403 owner mismatch`
- 结论：`P402-C` 真实集群 e2e 验收通过

### P501 计划（2026-04-01）

#### 范围定义
- 仅管理员可创建账号（关闭开放注册）
- 首版仅支持账号密码登录（OAuth/LDAP 后续阶段开发）
- `readonly` 仅允许已授权范围访问
- 兼容 `X-User/X-User-Role` 仅保留 1 个阶段，后续切换到 Bearer Token

#### 角色策略（Rancher 风格映射）
- 全局角色：
  - `admin`：平台全权限，具备用户创建与授权管理能力
  - `standard-user`：普通运维账号，按授权范围访问
  - `readonly`：只读账号，仅可访问授权范围
- 兼容映射（过渡期）：
  - `operator` -> `standard-user`
  - `viewer` -> `readonly`

#### 开发拆分
1. 用户与令牌数据模型（users/refresh_tokens）及默认管理员初始化
2. 认证接口（登录、刷新、登出、管理员创建用户）
3. 鉴权中间件升级（Bearer Token 优先 + 旧头兼容）
4. 作用域权限落地（readonly 仅授权范围）
5. 联调与验收（go test、npm build、auth 冒烟）

#### P501 验收标准（执行版）
- 默认管理员 `admin/123456` 可登录并获取访问令牌
- 非管理员无法创建账号；管理员可创建 `standard-user`/`readonly`
- `readonly` 对未授权命名空间写操作返回 `403`
- 旧 `X-User/X-User-Role` 仍可在过渡期使用，Bearer Token 生效优先级更高
- 后端 `go test ./...`、前端 `npm run build`、认证冒烟脚本通过
- 提供 `scripts/p501_auth_smoke_test.sh` 认证冒烟脚本

### P601 计划（2026-04-01）

#### 功能需求（当前会话确认）
- 管理员用户管理增强：用户列表查询、用户启停、密码重置
- 前端安全域页面同时提供用户管理与审计筛选能力

#### 设计要点
- 后端接口已落地：
  - `GET /api/v1/auth/users`
  - `PATCH /api/v1/auth/users/:username/status`
  - `POST /api/v1/auth/users/:username/reset-password`
- 关键写操作继续沿用 `X-Action-Confirm: CONFIRM` 与 `PermUserManage` 权限校验
- 前端 `AuthAuditPage` 已升级为“用户管理 + 审计日志”双区块页面
- 新增 `scripts/p601_user_management_smoke_test.sh`，覆盖用户创建、启停、重置密码与权限拒绝场景

#### 任务计划（执行结果）
1. 后端用户管理接口联调（已完成）
2. 前端用户管理页面改造（已完成）
3. 回归与冒烟验证（已完成）

#### 阶段状态
- 当前处于 `P601 已完成`
- 验证结果：`go test ./...`、`npm run build`、`scripts/p601_user_management_smoke_test.sh` 全部通过

### P701 计划（2026-04-01）

#### 范围定义
- 目标：落地“细粒度授权基础能力”，允许管理员在线调整用户角色与授权命名空间
- 本轮不纳入：OAuth/LDAP 登录接入、企业组织/用户组同步、外部 IdP 回调链路

#### 能力范围
- 后端新增“更新用户角色与授权范围”接口
- 前端用户管理页增加“编辑授权”能力
- 保持现有 `X-Action-Confirm` 写操作确认与 `PermUserManage` 权限校验

#### 开发拆分
1. `P701-A`：后端接口与服务层校验（角色合法性、命名空间范围）
2. `P701-B`：前端用户授权编辑交互（弹窗编辑 + 列表刷新）
3. `P701-C`：联调与验收（go test、frontend build、p701 冒烟）

#### 当前状态
- 状态：`已完成`
- 已交付：
  - `PATCH /api/v1/auth/users/:username` 用户角色与授权范围编辑接口
  - 前端用户管理页“编辑授权”弹窗交互
  - `scripts/p701_fine_grained_auth_smoke_test.sh` 冒烟脚本
- 验证结果：`go test ./...`、`npm run build`、`scripts/p701_fine_grained_auth_smoke_test.sh` 全部通过

### P801 计划（2026-04-01）

#### 范围定义
- 目标：落地 OAuth/LDAP 前的认证源管理基础能力，统一登录入口的 provider 语义
- 本轮聚焦：`local + ldap` 认证源配置与选择逻辑，先不接入真实 LDAP Bind 与 OAuth 回调

#### 能力范围
- 后端新增认证源模型与管理接口（列表、创建、启停、设为默认）
- 登录接口支持 `provider` 入参；若选择 `ldap`，返回明确“未实现”错误，避免静默回退
- 前端登录页增加认证源选择（默认跟随后端返回的默认 provider）

#### 开发拆分
1. `P801-A`：认证源数据模型与后端接口
2. `P801-B`：登录入口 provider 参数扩展与错误语义
3. `P801-C`：前端登录页 provider 选择与联调
4. `P801-D`：回归与冒烟验收

#### 当前状态
- 状态：`已完成`
- 已交付：
  - 认证源模型 `auth_providers`（`local/ldap`）及默认初始化
  - 认证源接口：`GET /auth/providers/public`、`GET /auth/providers`、`POST /auth/providers`、`PATCH /auth/providers/:id/status`、`POST /auth/providers/:id/default`
  - 登录接口支持 `provider` 参数（默认 provider 为 `ldap` 时返回 `501 not implemented`）
  - 前端登录页新增认证源选择器（默认读取后端公开认证源）
  - `scripts/p801_auth_provider_smoke_test.sh` 冒烟脚本
- 验证结果：`go test ./...`、`npm run build`、`scripts/p801_auth_provider_smoke_test.sh` 全部通过

### 第二阶段任务清单（Rancher 风格重构）

1. `R1`：壳层重构（左侧菜单 + 顶栏 + 路由骨架）
2. `R2`：通用资源页面框架重构
3. `R3`：现有模块迁移到新框架
4. `R4`：视觉与交互细节对齐 Rancher
5. `R5`：重构后联调与验收

## 任务状态

- 当前阶段：`第八阶段（P801 已完成）`
- 当前任务：`P801：认证源管理基础能力（local/ldap）（已完成）`
- 当前状态：`已交付认证源模型、管理接口、登录 provider 语义与登录页选择器`
- 下一任务：`第九阶段规划（真实 LDAP Bind 与 OAuth 回调接入）`
