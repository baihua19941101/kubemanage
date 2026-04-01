# TASKS

## 当前阶段

- 阶段：第三阶段（P306 已完成）
- 更新时间：2026-04-01

## 任务列表

| ID | 任务 | 状态 | 说明 |
| --- | --- | --- | --- |
| T001 | README 拆分为“长期蓝图 + 第一阶段 MVP” | 已完成 | 已完成需求分层、MVP 边界、验收标准、任务计划与状态更新 |
| T002 | 项目骨架初始化（前后端工程） | 已完成 | 已完成前后端目录、基础路由、配置与构建脚手架 |
| T003 | 集群管理 MVP | 已完成 | 已实现集群列表、当前集群查询、集群切换 |
| T004 | 名称空间管理 MVP | 已完成 | 已实现列表、创建、删除、YAML 查看/下载 |
| T005 | Deployment/Pod MVP | 已完成 | 已实现列表、详情、日志、YAML 编辑 |
| T006 | Service/ConfigMap/Secret MVP | 已完成 | 已实现列表、详情、Secret 脱敏 |
| T007 | 权限与审计 MVP | 已完成 | 已实现最小 RBAC、关键写操作审计 |
| T008 | MVP 联调与验收测试 | 已完成 | 已完成联调脚本与核心流程验收检查 |
| T009 | 启动命令整理并写入 README | 已完成 | 已补充前后端安装、启动、联调与验证命令 |
| T010 | 后端接入 MySQL/Redis 启动配置 | 已完成 | 已增加配置项、启动连接校验与文档说明 |
| T011 | T003 前置：数据库备份与功能分支创建 | 已完成 | 已备份 kubemanage，已创建 feature/t003-cluster-mvp |
| T012 | T004 前置：数据库备份与功能分支创建 | 已完成 | 已备份 kubemanage，已创建 feature/t004-namespace-mvp |
| T013 | T005 前置：数据库备份与功能分支创建 | 已完成 | 已备份 kubemanage，已创建 feature/t005-workload-mvp |
| T014 | T006 前置：数据库备份与功能分支创建 | 已完成 | 已备份 kubemanage，已创建 feature/t006-service-config-secret-mvp |
| T015 | T007 前置：数据库备份与功能分支创建 | 已完成 | 已备份 kubemanage，已创建 feature/t007-auth-audit-mvp |
| T016 | T008 前置：数据库备份与功能分支创建 | 已完成 | 已备份 kubemanage，已创建 feature/t008-mvp-integration-qa |
| P201 | 第二阶段任务1：K8s真实集群数据接入 | 已转入第三阶段 | 已由第三阶段 P301/P302 接续承接 |
| P202 | 第二阶段任务2：工作负载扩展 | 已完成 | 已完成 P202-A~P202-D（后端/前端/联调与冒烟） |
| P203 | 第二阶段任务3：服务发现扩展 | 已完成 | 已完成 Ingress/HPA 列表、详情与关联查询 |
| P204 | 第二阶段任务4：存储管理扩展 | 已完成 | 已完成 PV/PVC/StorageClass/ConfigMap/Secret 资源能力与菜单对齐 |
| P205 | 第二阶段任务5：日志与调试增强 | 已完成 | 已完成 MVP、第二轮细化与权限模型兼容回归 |
| P206 | 第二阶段任务6：权限与审计增强 | 已完成 | 已完成第一版命名空间级授权与审计筛选 |
| P207 | 第二阶段任务7：联调与验收 | 已完成 | 已完成第二阶段统一验收脚本与整体验收 |
| P301 | 第三阶段任务1：真实集群接入与导入 | 已完成（基础闭环） | 已完成基础闭环与验收，后续深化拆分到 P302/P303 |
| P302 | 第三阶段任务2：k8s client adapter 与连接配置存储 | 已完成 | 已完成 Adapter 分层、敏感字段加密存储、mock/live/auto 模式归一化与回归验收 |
| P303 | 第三阶段任务3：真实资源读链路切换 | 已完成 | 已完成 Cluster/Namespace/Workloads/Service Discovery/Storage 主读接口 live 优先切换 |
| P304 | 第三阶段任务4：真实写操作安全化 | 已完成 | 已完成写操作确认、失败回显与审计增强，并通过 p304 冒烟验证 |
| P305 | 第三阶段任务5：真实日志与终端能力 | 已完成 | 已完成多容器日志与终端能力占位对接，并通过 p305 冒烟验证 |
| P306 | 第三阶段任务6：第三阶段联调与验收 | 已完成 | 已完成第三阶段统一验收脚本并通过回归验证 |
| R101 | R1：壳层重构（左侧菜单+顶栏+路由骨架） | 已完成 | 已完成 ShellLayout、左侧菜单、顶部栏与路由切换 |
| R102 | R2：通用资源页框架重构 | 已完成 | 已完成通用页面组件并迁移 Cluster 页示范 |
| R103 | R3：模块迁移到新框架 | 已完成 | 已完成 Namespace/Workload/Resource/AuthAudit 迁移 |
| R104 | R4：视觉与交互细节对齐 | 已完成 | 已完成菜单分组、面包屑、页面动线与视觉统一 |
| R105 | R5：重构后联调与验收 | 已完成 | 已完成统一验收脚本与回归验证 |
| T018 | R1 前置：数据库备份与功能分支创建 | 已完成 | 已备份 kubemanage，已创建 feature/r1-rancher-shell |
| T019 | R2 前置：数据库备份与功能分支创建 | 已完成 | 已备份 kubemanage，已创建 feature/r2-resource-framework |
| T020 | R3 前置：数据库备份与功能分支创建 | 已完成 | 已备份 kubemanage，已创建 feature/r3-module-migration |
| T021 | R4 前置：数据库备份与功能分支创建 | 已完成 | 已备份 kubemanage，已创建 feature/r4-rancher-polish |
| T022 | R5 前置：数据库备份与功能分支创建 | 已完成 | 已备份 kubemanage，已创建 feature/r5-rebuild-qa |
| T023 | P202 前置：数据库备份与功能分支创建 | 已完成 | 已备份 kubemanage，已创建 feature/p202-workload-extensions |
| T024 | P202 重启：范围澄清与验收标准固化 | 已完成 | 已在 README 固化范围、拆分、执行版验收标准 |
| T025 | P202-A：后端 service 扩展 | 已完成 | 已完成 4 类资源模型、列表/详情、YAML 读写与统一校验函数 |
| T026 | P202-B：后端 handler/router 扩展 | 已完成 | 已新增 4 类资源读写路由与 handler，并更新路由测试 |
| T027 | P202-C：前端工作负载页重建 | 已完成 | 已支持 6 类资源切换、详情抽屉、YAML 编辑与 Pod 日志 |
| T028 | P202-D：联调与回归验证 | 已完成 | 已完成后端测试、前端构建、P202 冒烟脚本验证 |
| T029 | P202-E：Workloads 下拉导航改造 | 已完成 | 已支持点击“工作负载”下拉展示 6 类资源并子路由访问 |
| T030 | P202-F：导航层级扁平化（去三级） | 已完成 | 已改为一级菜单点击展开二级子菜单（Cluster/Workloads/Configuration/Security） |
| T031 | P203 前置：数据库备份与功能分支创建 | 已完成 | 已备份 kubemanage，已创建 feature/p203-service-discovery-extensions |
| T032 | P203-A：后端 Ingress/HPA 与关联查询接口 | 已完成 | 已新增 ingresses/hpas 列表详情及关联查询接口并补测试 |
| T033 | P203-B：前端服务发现扩展页面 | 已完成 | 已在资源页新增 Ingress/HPA 视图与关联信息展示 |
| T034 | P203-C：联调与验收 | 已完成 | 已通过 go test、前端构建、scripts/p203_smoke_test.sh |
| T035 | P204-A：后端存储资源接口扩展 | 已完成 | 已新增 PV/PVC/StorageClass 列表与详情接口并补测试 |
| T036 | P204-B：前端 Storage 页面扩展 | 已完成 | 已新增 StoragePage，支持 PV/PVC/StorageClass/ConfigMap/Secret |
| T037 | 导航结构对齐需求（Service Discovery/Storage 一级） | 已完成 | 已拆分 Service Discovery、Storage 一级菜单与子资源路由 |
| T038 | 布局一致性对齐（向 Workloads 看齐） | 已完成 | 子资源页面统一为筛选栏+资源表+详情抽屉布局 |
| T039 | Storage 子资源命名规范对齐 | 已完成 | 已统一为 PersistentVolumes/StorageClasses/ConfigMaps/PersistentVolumeClaims/Secrets |
| T040 | P205 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p205-logs-debug-enhancements，已备份 kubemanage-20260401-124825.sql |
| T041 | P205-A：后端日志增强与终端占位接口 | 已完成 | 已支持日志查询参数兼容，新增 terminal capabilities/session 占位接口 |
| T042 | P205-B：前端日志弹窗增强 | 已完成 | 已支持筛选、导出、终端入口与占位提示 |
| T043 | P205-C：联调与验收 | 已完成 | 已通过 go test、前端构建、scripts/p205_smoke_test.sh |
| T044 | P205-D：日志跟随刷新与参数联动 | 已完成 | 已支持 follow 刷新、前后端筛选参数联动与动态日志标记 |
| T045 | P205-E：日志交互细节完善 | 已完成 | 已支持复制、统计信息、清空筛选与手动刷新控制 |
| T046 | P205-F：第二轮联调与验收 | 已完成 | 已通过 go test、前端构建与 p205 冒烟验证 |
| T047 | P206 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p206-auth-audit-enhancements，已备份 kubemanage-20260401-131519.sql |
| T048 | P206-A：后端命名空间级授权增强 | 已完成 | 已支持 operator 仅写 dev、admin 全量，越权写入返回 403 |
| T049 | P206-B：后端审计筛选能力 | 已完成 | 已支持 user/role/method/path/statusCode/limit 过滤 |
| T050 | P206-C：前端权限与审计页增强 | 已完成 | 已支持授权说明、审计筛选栏与结果展示 |
| T051 | P206-D：联调与验收 | 已完成 | 已通过 go test、前端构建、scripts/p206_smoke_test.sh |
| T052 | P207 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p207-stage2-qa，已备份 kubemanage-20260401-132420.sql |
| T053 | P207-A：第二阶段统一验收脚本 | 已完成 | 已新增 scripts/p207_stage2_qa.sh，串联 p202/p203/p204/p205/p206 与基础构建测试 |
| T054 | P207-B：回归修正与结果沉淀 | 已完成 | 已修正 p202/p205 脚本权限模型兼容与 p207 脚本清理/断言问题 |
| T055 | P301 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p301-cluster-import-live-data，已备份 kubemanage-20260401-133753.sql |
| T056 | P301-A：第三阶段规划落盘 | 已完成 | 已在 README/TASKS 固化第三阶段扩展方案与 P301 范围 |
| T057 | P301-B：真实集群接入技术方案细化 | 已完成 | 已细化配置模型、导入协议、连接测试与 adapter 设计 |
| T058 | P301-C：后端连接配置与 live 读链路基础实现 | 已完成 | 已完成连接模型、持久化、导入/测试/激活接口与 live clusters/namespaces |
| T059 | P301-D：前端集群导入与连接测试页面 | 已完成 | 已完成导入表单、连接测试、激活与 live 数据展示 |
| T060 | P301-E：P301 冒烟脚本与基础验收 | 已完成 | 已新增 scripts/p301_smoke_test.sh，并通过前端构建与 P301 冒烟验证 |
| T061 | 文档状态对齐（README/TASKS） | 已完成 | 已统一第三阶段当前状态：P301 基础闭环完成，P302 待开始 |
| T062 | P302 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p302-adapter-storage-deepen，已通过 docker exec mysql8 完成备份 kubemanage-20260401-152041-p302.sql |
| T063 | P302-A：adapter 分层重构 | 已完成 | 已拆分 fake/live adapter 与 rest config 构建到独立文件，并补充调用超时控制 |
| T064 | P302-B：连接配置存储深化 | 已完成 | 已新增 KM_SECRET_KEY 透明加解密与历史明文兼容读取 |
| T065 | P302-C：数据源切换策略固化 | 已完成 | 已增加 mock/live/auto 模式归一化解析与路由层统一选择 |
| T066 | P302-D：回归与验收 | 已完成 | 已通过 go test ./...、npm run build、scripts/p302_smoke_test.sh |
| T067 | P303 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p303-live-read-switch，已通过 docker exec mysql8 完成备份 kubemanage-20260401-154202-p303.sql |
| T068 | P303-A：Cluster/Namespace 主链路切换 | 已完成 | 已将 /clusters 与 /namespaces 主接口切换为 live 优先读取，并补无激活连接错误语义测试 |
| T069 | P303-B：Workloads/Service Discovery/Storage 读链路扩展 | 已完成 | 已完成 Workloads 全量列表/详情与 Service Discovery/Storage 主读接口 live 优先切换 |
| T070 | P303-C：回归与验收 | 已完成 | 已通过 go test ./... 与 npm run build（real-only 读链路） |
| T071 | real-only 需求变更：停止 mock 数据链路 | 已完成 | 已移除主读接口的 mock 回退逻辑，并统一为真实数据依赖 |
| T072 | real-only 需求变更：文档与验收口径同步 | 已完成 | README/TASKS 已同步 real-only 约束与当前进展 |
| T073 | P303 二次优化：真实连接管理页去 mock 展示 | 已完成 | 已清理示例集群与 Live 概览，页面仅保留真实连接与真实集群列表 |
| T074 | P303 二次优化：集群列表字段扩展 | 已完成 | 已展示 State/Name/Provider/Distro/Kubernetes Version/Architecture/CPU/Memory/Pods |
| T075 | P303 二次优化：清理非 test1 历史连接数据 | 已完成 | 已清理 cluster_connections 中 test1 之外的历史连接记录，仅保留 test1 |
| T076 | P303 二次优化：集群列表字段组合收敛 | 已完成 | 已收敛为 Provider/Distro 与 Kubernetes Version/Architecture 两个组合字段 |
| T077 | P304 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p304-write-safety，已备份 kubemanage-20260401-162722-p304.sql |
| T078 | P304-A：写操作确认机制 | 已完成 | 关键写接口已接入 X-Action-Confirm 校验 |
| T079 | P304-B：失败回显增强 | 已完成 | 权限/确认失败响应已增加 requestId/code/hint |
| T080 | P304-C：审计字段增强 | 已完成 | 审计记录新增 requestId/namespace/error 字段 |
| T081 | P304-D：回归与验收 | 已完成 | 已通过 go test ./...、npm run build、scripts/p304_smoke_test.sh |
| T082 | P305 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p305-live-logs-terminal，已备份 kubemanage-20260401-164536-p305.sql |
| T083 | P305-A：真实 Pod 日志能力 | 已完成 | 已接入真实 Pod 日志读取，支持 container 参数与前端多容器选择 |
| T084 | P305-B：真实终端能力对接 | 已完成 | 已返回容器能力并打通终端会话容器参数，占位返回明确提示 |
| T085 | P305-C：回归与验收 | 已完成 | 已通过 go test ./...、npm run build、scripts/p305_smoke_test.sh |
| T086 | P306 前置：数据库备份与功能分支确认 | 已完成 | 已在 feature/p305-live-logs-terminal 上备份 kubemanage-20260401-165353-p306.sql |
| T087 | P306-A：第三阶段统一验收脚本 | 已完成 | 已新增 scripts/p306_stage3_qa.sh 并串联 stage3 核心验证 |
| T088 | P306-B：回归结果沉淀与文档同步 | 已完成 | 已完成 stage3 验证结果沉淀并同步 README/TASKS |

## 完成记录

- 2026-03-31：完成 T001（README 拆分与项目阶段化定义）
- 2026-03-31：完成 T002（项目骨架初始化，后端测试与前端构建通过）
- 2026-03-31：完成 T009（启动命令整理并写入 README）
- 2026-03-31：完成 T010（后端接入 MySQL/Redis 启动配置）
- 2026-03-31：完成 T011（T003 前置：数据库备份与功能分支创建）
- 2026-03-31：完成 T003（集群管理 MVP，含后端接口与前端切换交互）
- 2026-03-31：完成 T012（T004 前置：数据库备份与功能分支创建）
- 2026-03-31：完成 T004（名称空间管理 MVP，含后端接口与前端页面）
- 2026-03-31：完成 T013（T005 前置：数据库备份与功能分支创建）
- 2026-03-31：完成 T005（Deployment/Pod MVP，含后端接口与前端页面）
- 2026-03-31：完成 T014（T006 前置：数据库备份与功能分支创建）
- 2026-03-31：完成 T006（Service/ConfigMap/Secret MVP，含后端接口与前端页面）
- 2026-03-31：完成 T015（T007 前置：数据库备份与功能分支创建）
- 2026-03-31：完成 T007（权限与审计 MVP，含最小 RBAC 与写操作审计）
- 2026-03-31：完成 T016（T008 前置：数据库备份与功能分支创建）
- 2026-03-31：完成 T008（MVP 联调与验收测试，烟雾脚本验证通过）
- 2026-03-31：完成第二阶段规划（目标、范围、任务清单、状态迁移）
- 2026-04-01：确认第二阶段方向调整为 Rancher 风格重构路线
- 2026-04-01：完成 T018（R1 前置：数据库备份与功能分支创建）
- 2026-04-01：完成 R101 开发（Rancher 风格壳层重构，待验收）
- 2026-04-01：完成 T019（R2 前置：数据库备份与功能分支创建）
- 2026-04-01：完成 R102 开发（通用资源页框架重构，待验收）
- 2026-04-01：完成 T020（R3 前置：数据库备份与功能分支创建）
- 2026-04-01：完成 R103 开发（模块迁移到新框架，待验收）
- 2026-04-01：完成 T021（R4 前置：数据库备份与功能分支创建）
- 2026-04-01：完成 R104 开发（视觉与交互细节对齐，待验收）
- 2026-04-01：完成 T022（R5 前置：数据库备份与功能分支创建）
- 2026-04-01：完成 R105 开发（重构后联调与验收，待验收）
- 2026-04-01：完成 R1-R5（Rancher 风格重构路线全量完成）
- 2026-04-01：完成 T023（P202 前置：数据库备份与功能分支创建）
- 2026-04-01：按会话恢复决策回退 P202 未完成代码，任务状态恢复为“待开始（前置已完成）”
- 2026-04-01：完成 T024（P202 重启规划：范围、拆分、执行版验收标准固化）
- 2026-04-01：完成数据库备份（backups/kubemanage-20260401-095824.sql），作为 P202 开发前置
- 2026-04-01：完成 T025（P202-A：后端 service 扩展，含 4 类资源与 YAML 读写）
- 2026-04-01：完成 T026（P202-B：后端 handler/router 扩展，含路由测试更新）
- 2026-04-01：完成 T027（P202-C：前端工作负载页重建，前端构建通过）
- 2026-04-01：完成 T028（P202-D：联调与回归验证，含 scripts/p202_smoke_test.sh 冒烟通过）
- 2026-04-01：完成 P202（工作负载扩展全链路交付）
- 2026-04-01：完成数据库备份（backups/kubemanage-20260401-105855.sql），用于 Workloads 导航改造前置
- 2026-04-01：完成 T029（Workloads 下拉子菜单导航改造，前端构建通过）
- 2026-04-01：按需求调整 Workloads 导航为点击“工作负载”下拉展示 6 类资源
- 2026-04-01：完成数据库备份（backups/kubemanage-20260401-112423.sql），用于导航层级扁平化改造前置
- 2026-04-01：完成 T030（导航改为两级结构，前端构建通过）
- 2026-04-01：完成数据库备份（backups/kubemanage-20260401-121425.sql），并创建分支 feature/p203-service-discovery-extensions
- 2026-04-01：完成 T032（后端 Ingress/HPA 与关联查询接口，含路由测试）
- 2026-04-01：完成 T033（前端资源页扩展 Ingress/HPA 视图与关联展示）
- 2026-04-01：完成 T034（P203 联调回归通过，含 scripts/p203_smoke_test.sh）
- 2026-04-01：完成 P203（服务发现扩展）
- 2026-04-01：完成数据库备份（backups/kubemanage-20260401-123611.sql），用于菜单与存储能力对齐前置
- 2026-04-01：完成 T035（后端扩展 PV/PVC/StorageClass 接口并补测试）
- 2026-04-01：完成 T036（新增 Storage 页面并接入 5 类存储资源）
- 2026-04-01：完成 T037（Service Discovery/Storage 拆分为一级菜单）
- 2026-04-01：完成 T038（子资源页面布局统一对齐 Workloads）
- 2026-04-01：完成 P204（存储管理扩展与导航对齐）
- 2026-04-01：完成 T039（Storage 子资源标签命名统一为全称）
- 2026-04-01：完成 T040（P205 前置：数据库备份与功能分支创建，使用 backups/kubemanage-20260401-124825.sql）
- 2026-04-01：恢复会话后确认 P205 主体开发尚未开始，当前状态同步为“前置已完成，主任务待开始”
- 2026-04-01：确认 P205 采用 MVP 路线，范围为日志筛选/导出与终端接口预留，任务状态切换为“进行中”
- 2026-04-01：完成 T041（P205-A：后端日志查询参数兼容，新增终端占位接口）
- 2026-04-01：完成 T042（P205-B：前端日志弹窗支持筛选、导出与终端入口）
- 2026-04-01：完成 T043（P205-C：go test、npm run build、scripts/p205_smoke_test.sh 验证通过）
- 2026-04-01：完成数据库备份（backups/kubemanage-20260401-130722.sql），用于 P205 第二轮完善前置
- 2026-04-01：确认 P205 第二轮完善范围为日志跟随刷新、复制统计与前后端筛选联动
- 2026-04-01：完成 T044（P205-D：支持 follow 刷新、前后端筛选参数联动与动态日志标记）
- 2026-04-01：完成 T045（P205-E：支持复制、统计信息、清空筛选与手动刷新控制）
- 2026-04-01：完成 T046（P205-F：go test、npm run build、独立端口 p205 冒烟验证通过）
- 2026-04-01：完成 T047（P206 前置：数据库备份与功能分支创建，使用 backups/kubemanage-20260401-131519.sql）
- 2026-04-01：确认 P206 第一版范围为命名空间级授权与审计筛选，任务状态切换为“进行中”
- 2026-04-01：完成 T048（P206-A：后端命名空间级授权增强）
- 2026-04-01：完成 T049（P206-B：后端审计筛选能力）
- 2026-04-01：完成 T050（P206-C：前端权限与审计页增强）
- 2026-04-01：完成 T051（P206-D：go test、npm run build、scripts/p206_smoke_test.sh 验证通过）
- 2026-04-01：完成 T052（P207 前置：数据库备份与功能分支创建，使用 backups/kubemanage-20260401-132420.sql）
- 2026-04-01：启动 P207（第二阶段统一验收脚本与整体验收）
- 2026-04-01：完成 T053（P207-A：新增 scripts/p207_stage2_qa.sh，统一执行第二阶段回归验收）
- 2026-04-01：完成 T054（P207-B：修正 p202/p205 脚本权限兼容与 p207 脚本清理/断言问题）
- 2026-04-01：完成 P207（第二阶段统一验收通过，含 backend go test、frontend build、P202/P203/P204/P205/P206 验证）
- 2026-04-01：完成 T055（P301 前置：数据库备份与功能分支创建，使用 backups/kubemanage-20260401-133753.sql）
- 2026-04-01：完成 T056（第三阶段扩展方案与 P301 范围落盘，确认真实集群接入为第一优先级）
- 2026-04-01：完成 T057（P301-B：已在 README 固化连接配置模型、导入协议、接口与 adapter 设计）
- 2026-04-01：完成 T058（P301-C：完成连接模型、持久化、导入/测试/激活接口与 live clusters/namespaces）
- 2026-04-01：完成 T059（P301-D：完成前端导入表单、连接测试、激活与 live 数据展示，前端构建通过）
- 2026-04-01：完成 T060（P301-E：新增 scripts/p301_smoke_test.sh，并通过独立端口 P301 冒烟验证）
- 2026-04-01：按当前会话收口要求停止继续开发，已将 P301 最新基础闭环进度同步到 README/TASKS，供下个会话直接接续
- 2026-04-01：完成 T061（对齐 README/TASKS 状态口径：P301 基础闭环完成，P302 待开始）
- 2026-04-01：完成 T062（创建 feature/p302-adapter-storage-deepen，并通过 docker exec mysql8 生成 backups/kubemanage-20260401-152041-p302.sql）
- 2026-04-01：启动 T063（P302-A：adapter 分层重构）
- 2026-04-01：完成 T063（adapter 逻辑拆分到 cluster_connection_adapter，并通过 go test）
- 2026-04-01：完成 T064（新增 KM_SECRET_KEY 透明加密存储，兼容历史明文读取）
- 2026-04-01：完成 T065（新增 mock/live/auto adapter 模式归一化解析）
- 2026-04-01：完成 T066（通过 go test ./...、npm run build、scripts/p302_smoke_test.sh）
- 2026-04-01：完成 P302（k8s adapter 与连接配置存储深化）
- 2026-04-01：完成 T067（创建 feature/p303-live-read-switch，并通过 docker exec mysql8 生成 backups/kubemanage-20260401-154202-p303.sql）
- 2026-04-01：启动 T068（P303-A：Cluster/Namespace 主链路切换）
- 2026-04-01：完成 T068（/clusters、/namespaces 主接口切换为 live 优先读取，并通过 go test）
- 2026-04-01：启动 T069（Workloads/Service Discovery/Storage 读链路扩展）
- 2026-04-01：完成 T069 第一批（/deployments、/pods 主接口切换为 live 优先读取，并通过 go test）
- 2026-04-01：完成 real-only 变更前置备份（backups/kubemanage-20260401-155246-real-only.sql）
- 2026-04-01：启动 T071（停止 mock 数据链路，统一真实数据依赖）
- 2026-04-01：完成 T072（README/TASKS 同步 real-only 约束与状态）
- 2026-04-01：完成 T069（补齐 Workloads 其余资源与 Service Discovery/Storage 主读接口 live 优先切换）
- 2026-04-01：完成 T071（移除主读接口 mock 回退逻辑，默认模式收敛为 live）
- 2026-04-01：完成 T070（go test ./... 与 npm run build 通过）
- 2026-04-01：完成 P303（真实资源读链路切换）
- 2026-04-01：完成 T073（真实连接管理页去 mock 展示，移除 Live 真实数据概览）
- 2026-04-01：完成 T074（集群列表字段扩展为 State/Name/Provider/Distro/Kubernetes Version/Architecture/CPU/Memory/Pods）
- 2026-04-01：完成 P303 二次优化（真实连接管理页 real-only 收口）
- 2026-04-01：完成需求变更前置备份（backups/kubemanage-20260401-162126-cleanup-non-test1.sql）
- 2026-04-01：启动 T075（清理非 test1 历史连接数据）
- 2026-04-01：启动 T076（集群列表字段组合收敛）
- 2026-04-01：完成 T075（数据库仅保留 test1 连接，并保持为默认激活）
- 2026-04-01：完成 T076（前端集群列表字段组合收敛）
- 2026-04-01：完成 T077（创建 feature/p304-write-safety，并备份 kubemanage-20260401-162722-p304.sql）
- 2026-04-01：启动 T078（写操作确认机制）
- 2026-04-01：完成 T078（关键写路由确认头校验接入）
- 2026-04-01：完成 T079（失败响应增加 requestId/code/hint）
- 2026-04-01：完成 T080（审计记录增加 requestId/namespace/error）
- 2026-04-01：完成 T081（go test ./...、npm run build、scripts/p304_smoke_test.sh 通过）
- 2026-04-01：完成 P304（真实写操作安全化）
- 2026-04-01：完成 T082（创建 feature/p305-live-logs-terminal，并备份 kubemanage-20260401-164536-p305.sql）
- 2026-04-01：启动 T083（真实 Pod 日志能力）
- 2026-04-01：完成 T083（真实 Pod 日志读取 + 多容器选择，且修复前端写操作确认头）
- 2026-04-01：完成 T084（终端能力返回容器列表并透传容器参数）
- 2026-04-01：完成 T085（go test ./...、npm run build、scripts/p305_smoke_test.sh 通过）
- 2026-04-01：完成 P305（真实日志与终端能力）
- 2026-04-01：完成 T086（备份 kubemanage-20260401-165353-p306.sql）
- 2026-04-01：完成 T087（新增 scripts/p306_stage3_qa.sh 并执行通过）
- 2026-04-01：完成 T088（第三阶段回归结果沉淀与文档同步）
- 2026-04-01：完成 P306（第三阶段联调与验收）
- 2026-04-01：完成 T086（备份 kubemanage-20260401-165353-p306.sql）
- 2026-04-01：启动 T087（第三阶段统一验收脚本）
- 2026-04-01：完成需求变更前置备份（backups/kubemanage-20260401-161043-real-cluster-fields.sql）
- 2026-04-01：启动 T073（真实连接管理页去 mock 展示）
