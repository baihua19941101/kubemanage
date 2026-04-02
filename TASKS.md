# TASKS

## 当前阶段

- 阶段：第二十阶段（P1401 已完成）
- 更新时间：2026-04-02

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
| P401 | 第四阶段任务1：前端错误回显与可观测性增强 | 已完成 | 已完成统一错误解析并展示 requestId/code/hint |
| P402 | 第四阶段任务2：真实终端 exec 链路打通 | 已完成 | 已完成会话管理、WebSocket/exec 桥接、安全校验与真实集群 e2e 验收 |
| P501 | 第五阶段任务1：账号认证与用户体系 | 已完成 | 已完成账号密码认证闭环、管理员创建用户、Bearer 优先鉴权与前端登录态接入 |
| P601 | 第六阶段任务1：用户管理增强 | 已完成 | 已完成用户列表、用户启停、重置密码、前端页面改造与冒烟验收 |
| P701 | 第七阶段任务1：细粒度授权基础能力 | 已完成 | 已完成用户角色与授权范围在线编辑、前端交互与冒烟验收 |
| P801 | 第八阶段任务1：认证源管理基础能力（local/ldap） | 已完成 | 已完成认证源模型、管理接口、登录 provider 语义与前端登录源选择 |
| P901 | 第九阶段任务1：真实 LDAP Bind 最小闭环 | 已完成 | 已完成 provider=ldap 真实登录链路、账号映射策略与验收脚本 |
| P1001 | 第十阶段任务1：节点管理第一版能力 | 已完成 | 已完成节点列表/详情/YAML 读能力、前端节点页面与冒烟验收 |
| P1101 | 第十一阶段任务1：权限体系深化第一版（Token 生命周期） | 已完成 | 已完成 token 会话列表、撤销与登出全部能力并接入前端权限页 |
| P1201 | 第十二阶段任务1：Policy 管理第一版能力 | 已完成 | 已完成 LimitRange/ResourceQuota/NetworkPolicy 列表详情与 YAML 查看/下载、前端页面路由与冒烟验收 |
| P1202 | 第十二阶段任务2：Policy YAML 写能力第一版 | 已完成 | 已完成 Policy YAML 保存接口、前端编辑保存交互与写入验收 |
| P1203 | 第十二阶段任务3：Policy 删除能力第一版 | 已完成 | 已完成 3 类 Policy 删除接口、前端删除交互与联调验收 |
| P1204 | 第十二阶段任务4：Policy 新建能力第一版 | 已完成 | 已完成 3 类 Policy 创建接口、前端新建交互与联调验收 |
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
| T089 | P401 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p401-terminal-exec-mvp，已备份 kubemanage-20260401-171126-p401.sql |
| T090 | P401-A：前端统一错误解析 | 已完成 | 新增 parseApiError 与 ApiRequestError，统一解析 error/code/hint/requestId |
| T091 | P401-B：前端关键写操作错误展示优化 | 已完成 | Cluster/Namespace/Workload 关键写操作统一展示 requestId 错误信息 |
| T092 | P401-C：回归与验收 | 已完成 | 已通过 npm run build 与关键写操作路径回归 |
| T093 | P402-A：终端会话管理 | 已完成 | 已完成会话创建、TTL 与 sessionId 一次性消费校验，并补充单元测试 |
| T094 | P402-B：WebSocket exec 通道 | 已完成 | 已完成会话校验分支测试、WebSocket 通道路由接入与 exec 桥接链路落地 |
| T095 | P402-C：回归与验收 | 已完成 | 已完成 go test、frontend build、p402 冒烟与真实集群 e2e exec 收发验收 |
| T096 | P402 文档细化：执行拆分与验收标准固化 | 已完成 | 已在 README/TASKS 固化 P402 范围、设计要点、拆分与验收口径 |
| T097 | P402 前置：数据库备份与开发基线确认 | 已完成 | 已备份 backups/kubemanage-20260401-175333-p402.sql 并完成基线验证 |
| T098 | P402-B 加固：会话创建者绑定校验 | 已完成 | 已增加 ws attach 的 user/role 绑定校验并补充分支测试 |
| T099 | P402-C 文档化：真实集群 e2e 验收说明与记录模板 | 已完成 | 已在 README 增加执行步骤、失败分支校验项与验收记录模板 |
| T100 | P402-B 修正：live 模式 Pod 命名空间解析 | 已完成 | 已修正 terminal session/live Pod 写路由在 live 模式下的 namespace 解析链路 |
| T101 | P402 运维增强：终端会话 TTL 可配置化 | 已完成 | 已支持 KM_TERMINAL_SESSION_TTL_SECONDS 并补解析单测与文档说明 |
| T102 | P402 可用性增强：会话过期信息前后端联动 | 已完成 | 会话创建响应新增 ttlSeconds/expiresAt，前端终端提示增加 TTL/过期时间展示 |
| T103 | P501 前置：数据库备份与功能分支基线确认 | 已完成 | 已备份 backups/kubemanage-20260401-202511-p501.sql，当前分支 feature/p401-terminal-exec-mvp 继续开发 |
| T104 | P501-A：用户/令牌数据模型与默认管理员初始化 | 已完成 | 已新增 users/refresh_tokens 模型并接入 AutoMigrate，首次启动自动创建 admin/123456 |
| T105 | P501-B：认证接口（登录/刷新/登出） | 已完成 | 已新增 /auth/login /auth/refresh /auth/logout 账号密码认证闭环 |
| T106 | P501-C：管理员创建用户与角色授权 | 已完成 | 已新增 /auth/users（仅 admin），支持 standard-user/readonly 与授权命名空间 |
| T107 | P501-D：鉴权中间件升级与旧头兼容过渡 | 已完成 | 已实现 Bearer Token 优先解析，兼容 X-User/X-User-Role 过渡 1 阶段 |
| T108 | P501-E：联调与验收 | 已完成 | 已通过 go test、frontend build、scripts/p501_auth_smoke_test.sh |
| T109 | P501-F：前端登录态接入与页面切换 | 已完成 | 已完成登录页、Token 存储、路由守卫与顶部用户态展示/退出 |
| T110 | P601 前置：数据库备份与功能分支创建 | 已完成 | 当前分支 feature/p601-user-management，已备份 backups/kubemanage-20260401-205935-p601.sql |
| T111 | P601-A：后端用户管理接口 | 已完成 | 已新增用户列表、用户启停、密码重置接口与权限/确认头接入 |
| T112 | P601-B：前端用户管理页面改造 | 已完成 | AuthAuditPage 已升级为“用户管理 + 审计日志”双区块页面 |
| T113 | P601 收口：前端构建链路修复 | 已完成 | 已恢复 AuthAuditPage 路由链路，frontend `npm run build` 通过 |
| T114 | P601-C：联调与验收 | 已完成 | 已通过 go test、frontend build、scripts/p601_user_management_smoke_test.sh |
| T115 | P701 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p701-fine-grained-auth，已备份 backups/kubemanage-20260401-214449-p701.sql |
| T116 | P701-A：README/TASKS 计划与状态同步 | 已完成 | 已同步第七阶段范围、拆分与当前状态（进行中） |
| T117 | P701-B：后端用户授权编辑接口 | 已完成 | 已新增 PATCH /auth/users/:username，支持角色与授权范围更新 |
| T118 | P701-C：前端用户授权编辑交互 | 已完成 | 已新增编辑授权弹窗并接入用户列表刷新 |
| T119 | P701-D：联调与验收 | 已完成 | 已通过 go test、frontend build、scripts/p701_fine_grained_auth_smoke_test.sh |
| T120 | P801 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p801-auth-provider-foundation，已备份 backups/kubemanage-20260401-215918-p801.sql |
| T121 | P801-A：README/TASKS 计划与状态同步 | 已完成 | 已同步第八阶段范围、拆分与当前状态（进行中） |
| T122 | P801-B：后端认证源模型与管理接口 | 已完成 | 已新增 auth_providers 模型、默认初始化与认证源管理接口 |
| T123 | P801-C：登录 provider 语义与前端登录页扩展 | 已完成 | 登录接口支持 provider，登录页新增认证源选择并读取公开认证源列表 |
| T124 | P801-D：联调与验收 | 已完成 | 已通过 go test、frontend build、scripts/p801_auth_provider_smoke_test.sh |
| T125 | P901 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p901-ldap-bind-mvp，已备份 backups/kubemanage-20260401-221003-p901.sql |
| T126 | P901-A：README/TASKS 计划与状态同步 | 已完成 | 已同步第九阶段范围、拆分与当前状态（进行中） |
| T127 | P901-B：后端 LDAP Bind 登录链路 | 已完成 | 已完成 provider=ldap 真实 Bind 登录与本地用户映射签发 |
| T128 | P901-C：联调与验收 | 已完成 | 已通过 go test、frontend build、scripts/p901_ldap_bind_smoke_test.sh |
| T129 | 会话恢复后任务进度核对与文档对齐 | 已完成 | 已核对代码与文档偏差，修正 P801/P901 状态口径并补充 2026-04-02 复测记录 |
| T130 | P1001 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p1001-node-management-mvp，已备份 backups/kubemanage-20260402-085330-p1001.sql |
| T131 | P1001-A：README/TASKS 计划与状态同步 | 已完成 | 已同步第十阶段节点管理范围、拆分与当前状态（进行中） |
| T132 | P1001-B：后端节点读链路实现 | 已完成 | 已新增 nodes 列表/详情/YAML 查看下载接口与 live 适配 |
| T133 | P1001-C：前端节点管理页面与路由接入 | 已完成 | 已新增 Node 页面、Cluster 菜单入口、详情抽屉与 YAML 查看下载 |
| T134 | P1001-D：联调与验收 | 已完成 | 已通过 go test、frontend build、scripts/p1001_node_management_smoke_test.sh |
| T135 | P1101 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p1101-authz-lifecycle，已备份 backups/kubemanage-20260402-091110-p1101.sql |
| T136 | P1101-A：README/TASKS 计划与状态同步 | 已完成 | 已同步第十一阶段权限体系深化范围、拆分与当前状态（进行中） |
| T137 | P1101-B：后端 token 生命周期接口 | 已完成 | 已新增 token 会话列表、按 token 撤销、按用户批量撤销接口 |
| T138 | P1101-C：前端会话令牌管理能力 | 已完成 | 已在权限页新增 token 列表筛选、撤销与批量撤销操作 |
| T139 | P1101-D：联调与验收 | 已完成 | 已通过 go test、frontend build、scripts/p1101_token_lifecycle_smoke_test.sh |
| T140 | P1201 前置：数据库备份与功能分支创建 | 已完成 | 已创建 feature/p1201-policy-management，已备份 backups/kubemanage-20260402-095515-p1201.sql |
| T141 | P1201-A：README/TASKS 计划与状态同步 | 已完成 | 已同步第十二阶段 Policy 管理范围、拆分与当前状态（进行中） |
| T142 | P1201-B：后端 Policy 读链路实现 | 已完成 | 已新增 LimitRange/ResourceQuota/NetworkPolicy 列表详情与 YAML 查看/下载接口并接入路由测试 |
| T143 | P1201-C：前端 Policy 页面与路由接入 | 已完成 | 已新增 Policy 页面、菜单入口、资源切换与 YAML 查看/下载交互 |
| T144 | P1201-D：联调与验收 | 已完成 | 已通过 go test ./...、frontend build、scripts/p1201_policy_smoke_test.sh |
| T145 | P1202 前置：数据库备份与任务启动 | 已完成 | 已备份 backups/kubemanage-20260402-102151-p1202.sql，进入 P1202 开发 |
| T146 | P1202-A：README/TASKS 计划与状态同步 | 已完成 | 已同步 P1202 范围、拆分与当前状态（进行中） |
| T147 | P1202-B：后端 Policy YAML 写接口 | 已完成 | 已新增 3 类 Policy YAML 保存接口并接入路由、权限、确认头与作用域校验 |
| T148 | P1202-C：前端 Policy YAML 编辑保存交互 | 已完成 | PolicyPage 已支持 YAML 查看/编辑/保存，按角色控制保存能力 |
| T149 | P1202-D：联调与验收 | 已完成 | 已通过 go test ./...、frontend build、scripts/p1202_policy_write_smoke_test.sh |
| T150 | P1203 前置：数据库备份与任务启动 | 已完成 | 已备份 backups/kubemanage-20260402-111032-p1203.sql，进入 P1203 开发 |
| T151 | P1203-A：README/TASKS 计划与状态同步 | 已完成 | 已同步 P1203 范围、拆分与当前状态（进行中） |
| T152 | P1203-B：后端 Policy 删除接口 | 已完成 | 已新增 3 类 Policy 删除接口并接入作用域权限、确认头与路由测试 |
| T153 | P1203-C：前端 Policy 删除交互 | 已完成 | PolicyPage 已新增删除按钮、确认弹窗与删除后刷新 |
| T154 | P1203-D：联调与验收 | 已完成 | 已通过 go test ./...、frontend build、scripts/p1203_policy_delete_smoke_test.sh |
| T155 | 导航调整：Policy 与 Security 菜单分组修正 | 已完成 | 已将 LimitRange/ResourceQuota/NetworkPolicy 归入 Policy 菜单，Security 仅保留权限与审计 |
| T156 | P1204 前置：数据库备份与任务启动 | 已完成 | 已备份 backups/kubemanage-20260402-120216-p1204.sql，进入 P1204 开发 |
| T157 | P1204-A：README/TASKS 计划与状态同步 | 已完成 | 已同步 P1204 范围、拆分与当前状态（进行中） |
| T158 | P1204-B：后端 Policy 创建接口 | 已完成 | 已新增 3 类 Policy 创建接口并接入作用域权限、确认头与路由测试 |
| T159 | P1204-C：前端 Policy 新建交互 | 已完成 | PolicyPage 已新增新建 YAML 弹窗、创建请求与创建后刷新 |
| T160 | P1204-D：联调与验收 | 已完成 | 已通过 go test ./...、frontend build、scripts/p1204_policy_create_smoke_test.sh |
| T161 | P1204 收口优化：Create 文案与命名空间下拉选择 | 已完成 | 已将按钮改为 Create，并支持关联全部 namespace 下拉选择创建 |
| T162 | 会话收口：暂停新功能开发并同步进度文档 | 已完成 | 已按指令停止继续开发，完成 README/TASKS 最新进度核对与收口记录 |
| T163 | P1301-A：README real-only 口径纠偏 | 已完成 | 已标注基础验证中 live/mock 能力边界，修正易误导命令说明 |
| T164 | P1301-B：`/clusters/switch` real-only 兼容 | 已完成 | 已支持 live 模式按连接名切换激活集群，并补充路由测试 |
| T165 | P1301-C：回归与验收 | 已完成 | 已通过 go test ./...、npm run build，并完成 README/TASKS 状态同步 |
| T166 | P1302-A：第二轮对齐前置备份与计划同步 | 已完成 | 已完成数据库备份 backups/kubemanage-20260402-130753-p1302.sql，并同步 README/TASKS 范围 |
| T167 | P1302-B：Namespace real-only 写链路对齐 | 已完成 | 已支持 live 模式 namespace 创建/删除/YAML 查看下载 |
| T168 | P1302-C：Deployment real-only YAML 写链路对齐 | 已完成 | 已支持 live 模式 deployment YAML 查看/保存，并补名称/命名空间一致性校验 |
| T169 | P1302-C 补充：Deployment 写权限作用域解析对齐 | 已完成 | Deployment YAML 写入路由作用域解析已改为 live 优先命名空间解析 |
| T170 | P1302-D：回归与验收 | 已完成 | 已通过 go test ./...、npm run build，并完成 README/TASKS 状态同步 |
| T171 | Namespace 创建权限前端对齐修复 | 已完成 | 已按授权命名空间范围控制创建/删除按钮状态，并展示可写范围提示 |
| T172 | P1303-A：YAML 弹窗折叠样式统一改造 | 已完成 | 已将共享 YamlDialog 改为可折叠样式并支持展开/收起 |
| T173 | P1303-B：YAML 样式统一接入验证 | 已完成 | Namespace/Workload/Node/Policy 的 YAML 弹窗已统一生效 |
| T174 | P1303-C：回归与验收 | 已完成 | 已通过 frontend npm run build，并完成 README/TASKS 同步 |
| T175 | P1304-A：字段级折叠改造前置备份与依赖准备 | 已完成 | 已完成数据库备份 backups/kubemanage-20260402-132840-p1304.sql，并新增 frontend `yaml` 依赖 |
| T176 | P1304-B：YamlDialog 字段级折叠视图改造 | 已完成 | 已新增“结构视图/源码视图”，支持 key 级折叠与全量展开/折叠 |
| T177 | P1304-C：回归与验收 | 已完成 | 已通过 frontend npm run build，并完成 README/TASKS 同步 |
| T178 | P1305-A：Rancher 风格编辑器依赖接入 | 已完成 | 已新增 `react-ace` 与 `ace-builds` 并使用国内源安装 |
| T179 | P1305-B：YamlDialog Rancher 风格改造 | 已完成 | 已改为 YAML 编辑器主视图（行号+gutter 折叠）并保留结构视图 |
| T180 | P1305-C：回归与验收 | 已完成 | 已通过 frontend npm run build，并完成 README/TASKS 同步 |
| T181 | P1306-A：YAML 工具栏增强前置备份 | 已完成 | 已完成数据库备份 backups/kubemanage-20260402-133851-p1306.sql |
| T182 | P1306-B：YamlDialog 工具栏增强 | 已完成 | 已新增导入文件、下载、还原、变更预览操作并保留 key 折叠编辑 |
| T183 | P1306-C：回归与验收 | 已完成 | 已通过 frontend npm run build，并完成 README/TASKS 同步 |
| T184 | Deployment YAML 按钮无反馈修复 | 已完成 | WorkloadPage 已补 YAML 打开/保存异常提示与加载态，避免点击无反馈 |
| T185 | P1307-A：5 类工作负载对齐前置备份 | 已完成 | 已完成数据库备份 backups/kubemanage-20260402-143130-p1307.sql |
| T186 | P1307-B：live reader 补齐 5 类 YAML 读写 | 已完成 | 已支持 Pod/StatefulSet/DaemonSet/Job/CronJob YAML 查看与保存 |
| T187 | P1307-C：live 写入作用域解析对齐 | 已完成 | 已将 5 类写路由命名空间解析统一改为 live 优先 |
| T188 | P1307-D：回归与验收 | 已完成 | 已通过 go test ./...、npm run build，并完成 README/TASKS 同步 |
| T189 | Workload YAML 保存反馈增强 | 已完成 | 已补保存成功提示、失败精确错误（含 requestId）与保存后自动刷新列表 |
| T190 | Workload YAML 保存后留在编辑态增强 | 已完成 | 保存后保留弹窗、显示 requestId 成功提示并刷新基线内容 |
| T191 | YAML 保存元信息展示增强 | 已完成 | YamlDialog 已显示最近保存时间、最近 requestId 与最近保存历史 |
| T192 | YAML 保存历史按资源隔离 | 已完成 | 已按“资源类型+名称”独立记录保存历史，切换资源不互相覆盖 |
| T193 | 会话收口：停止新功能开发并同步进度 | 已完成 | 已停止继续开发，仅更新 README/TASKS 最新状态并等待下一轮优先级 |
| P1401 | 第二十阶段任务1：前端 Pod 终端实操化（xterm） | 已完成 | 已完成 xterm 终端接入、页面集成与构建回归 |
| T194 | P1401-A：前置备份与 README/TASKS 计划同步 | 已完成 | 已创建 feature/p1401-terminal-ui 并完成备份 backups/kubemanage-20260402-154417-p1401.sql |
| T195 | P1401-B：前端终端组件接入 | 已完成 | 已新增 TerminalDialog，完成 xterm + ws 输入输出与自适应 |
| T196 | P1401-C：页面集成与交互完善 | 已完成 | 已完成日志弹窗接入、自动连接、状态提示与手动重连 |
| T197 | P1401-D：联调与验收 | 已完成 | 已通过 go test ./... 与 frontend npm run build |

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
- 2026-04-01：完成 T089（创建 feature/p401-terminal-exec-mvp，并备份 kubemanage-20260401-171126-p401.sql）
- 2026-04-01：启动 T090（前端统一错误解析）
- 2026-04-01：完成 T090（前端统一错误解析与错误对象标准化）
- 2026-04-01：完成 T091（关键写操作错误展示 requestId/code/hint）
- 2026-04-01：完成 T092（npm run build 通过）
- 2026-04-01：完成 P401（前端错误回显与可观测性增强）
- 2026-04-01：启动 T093（终端会话管理）
- 2026-04-01：启动 T094（WebSocket exec 通道）
- 2026-04-01：完成 T096（P402 执行拆分与验收标准在 README/TASKS 固化）
- 2026-04-01：完成 T093（终端会话管理与 session store 单元测试，go test ./... 与 npm run build 通过）
- 2026-04-01：完成 T097（P402 开发前数据库备份：backups/kubemanage-20260401-175333-p402.sql）
- 2026-04-01：完成 T094（新增 terminal ws handler 测试与 session 过期语义修复）
- 2026-04-01：启动 T095（已通过 go test ./...、npm run build、scripts/p402_smoke_test.sh）
- 2026-04-01：完成 T098（终端 ws 会话创建者绑定校验，go test ./... 与 p402 冒烟通过）
- 2026-04-01：完成 T099（补充 P402-C 真实集群 e2e 验收执行说明与记录模板）
- 2026-04-01：完成 T100（修复 live 模式 terminal session 鉴权命名空间解析，真实 Pod 可创建会话）
- 2026-04-01：完成 T095（真实集群 e2e 通过：会话创建 201、WebSocket 收发回显成功、重复消费 401、TTL 过期 401、owner mismatch 403）
- 2026-04-01：完成 P402（真实终端 exec 链路打通并通过真实集群验收）
- 2026-04-01：完成 T101（终端会话 TTL 支持环境变量配置，默认 120 秒并补充解析单测）
- 2026-04-01：完成 T102（会话创建响应新增 ttlSeconds/expiresAt，前端终端提示增加会话时效展示）
- 2026-04-01：完成 T103（P501 开发前数据库备份：backups/kubemanage-20260401-202511-p501.sql）
- 2026-04-01：启动 P501（账号认证与用户体系）
- 2026-04-01：启动 T104（用户/令牌数据模型与默认管理员初始化）
- 2026-04-01：完成 T104（users/refresh_tokens 模型落库并自动初始化默认管理员 admin/123456）
- 2026-04-01：完成 T105（完成账号密码登录/刷新/登出接口）
- 2026-04-01：完成 T106（完成管理员创建用户接口与角色范围校验）
- 2026-04-01：完成 T107（完成 Bearer Token 优先鉴权与旧头兼容过渡）
- 2026-04-01：启动 T108（认证联调与验收，待补脚本与前端登录态验证）
- 2026-04-01：完成 T108（go test ./...、npm run build、scripts/p501_auth_smoke_test.sh 通过）
- 2026-04-01：启动 T109（前端登录态接入与页面切换）
- 2026-04-01：完成 T109（登录页、Token 本地存储、路由守卫、用户态展示与退出）
- 2026-04-01：完成 P501（账号认证与用户体系首版交付，含认证冒烟脚本）
- 2026-04-01：完成 T110（P601 前置：确认分支 feature/p601-user-management，并完成数据库备份 backups/kubemanage-20260401-205935-p601.sql）
- 2026-04-01：完成 T111（P601-A：后端新增用户列表/启停/重置密码接口）
- 2026-04-01：完成 T112（P601-B：前端 AuthAuditPage 升级为“用户管理 + 审计日志”页面）
- 2026-04-01：完成 T113（恢复 frontend/src/pages/AuthAuditPage.tsx，前端构建恢复通过）
- 2026-04-01：完成 T114（P601-C：go test ./...、npm run build、scripts/p601_user_management_smoke_test.sh 通过）
- 2026-04-01：完成 P601（用户管理增强交付并验收通过）
- 2026-04-01：完成 T115（P701 前置：创建 feature/p701-fine-grained-auth，并完成数据库备份 backups/kubemanage-20260401-214449-p701.sql）
- 2026-04-01：完成 T116（P701：README/TASKS 同步第七阶段范围、拆分与状态）
- 2026-04-01：启动 T117（P701-B：后端用户授权编辑接口）
- 2026-04-01：完成 T117（新增 PATCH /api/v1/auth/users/:username，支持角色与授权范围编辑）
- 2026-04-01：完成 T118（前端用户管理页新增“编辑授权”弹窗与提交流程）
- 2026-04-01：完成 T119（go test ./...、npm run build、scripts/p701_fine_grained_auth_smoke_test.sh 通过）
- 2026-04-01：完成 P701（细粒度授权基础能力交付并验收通过）
- 2026-04-01：完成 T120（P801 前置：创建 feature/p801-auth-provider-foundation，并完成数据库备份 backups/kubemanage-20260401-215918-p801.sql）
- 2026-04-01：完成 T121（P801：README/TASKS 同步第八阶段范围、拆分与状态）
- 2026-04-01：完成 T122（新增 auth_providers 模型、默认 provider 初始化与认证源管理接口）
- 2026-04-01：完成 T123（登录接口支持 provider，前端登录页新增认证源选择）
- 2026-04-01：完成 T124（go test ./...、npm run build、scripts/p801_auth_provider_smoke_test.sh 通过）
- 2026-04-01：完成 P801（认证源管理基础能力交付并验收通过）
- 2026-04-01：完成 T125（P901 前置：创建 feature/p901-ldap-bind-mvp，并完成数据库备份 backups/kubemanage-20260401-221003-p901.sql）
- 2026-04-01：完成 T126（P901：README/TASKS 同步第九阶段范围、拆分与状态）
- 2026-04-01：启动 T127（接入 provider=ldap 真实登录链路与 LDAP 配置解析）
- 2026-04-01：完成阶段验证（go test ./...、npm run build、scripts/p801_auth_provider_smoke_test.sh 通过）
- 2026-04-01：T127 进入阻塞（scripts/p901_ldap_bind_smoke_test.sh 登录返回 502，LDAP 连接被重置，待专项排查）
- 2026-04-02：完成会话恢复后复核（go test ./...、npm run build 通过；scripts/p901_ldap_bind_smoke_test.sh 仍返回 502）
- 2026-04-02：完成 T129（修正 README 的 P801 502 口径，更新 T128 为“进行中（阻塞）”并同步 TASKS 更新时间）
- 2026-04-02：完成 P901 冒烟脚本修正（修复 LDAP_USER_FILTER 默认占位符截断；默认端口改为 31389 映射容器 10389；默认测试账号调整为 fry/fry）
- 2026-04-02：完成 T127/T128（go test ./...、npm run build、scripts/p901_ldap_bind_smoke_test.sh 全部通过）
- 2026-04-02：完成 P901（真实 LDAP Bind 最小闭环交付并验收通过）
- 2026-04-02：完成 T130（创建 feature/p1001-node-management-mvp，并完成数据库备份 backups/kubemanage-20260402-085330-p1001.sql）
- 2026-04-02：完成 T131（同步 README/TASKS：第十阶段节点管理范围、拆分与状态）
- 2026-04-02：完成 T132（后端新增 nodes 列表/详情/YAML 查看下载接口，接入 live reader 与路由测试）
- 2026-04-02：完成 T133（前端新增 Node 页面并接入 Cluster 菜单与路由）
- 2026-04-02：完成 T134（go test ./...、npm run build、scripts/p1001_node_management_smoke_test.sh 通过）
- 2026-04-02：完成 P1001（节点管理第一版能力交付并验收通过）
- 2026-04-02：完成 T135（创建 feature/p1101-authz-lifecycle，并完成数据库备份 backups/kubemanage-20260402-091110-p1101.sql）
- 2026-04-02：完成 T136（同步 README/TASKS：第十一阶段权限体系深化范围、拆分与状态）
- 2026-04-02：完成 T137（后端新增 /auth/tokens 列表、/auth/tokens/:id/revoke、/auth/tokens/revoke-all 接口）
- 2026-04-02：完成 T138（前端 AuthAuditPage 新增会话令牌管理区块与筛选/撤销能力）
- 2026-04-02：完成 T139（go test ./...、npm run build、scripts/p1101_token_lifecycle_smoke_test.sh 通过）
- 2026-04-02：完成 P1101（权限体系深化第一版 Token 生命周期交付并验收通过）
- 2026-04-02：完成 T140（创建 feature/p1201-policy-management，并完成数据库备份 backups/kubemanage-20260402-095515-p1201.sql）
- 2026-04-02：完成 T141（同步 README/TASKS：第十二阶段 Policy 管理范围、拆分与状态）
- 2026-04-02：完成 T142（后端 Policy 读链路实现，新增 LimitRange/ResourceQuota/NetworkPolicy 列表详情与 YAML 查看/下载接口，并通过 go test ./...）
- 2026-04-02：完成 T143（前端新增 PolicyPage 与菜单/路由接入，支持 LimitRange/ResourceQuota/NetworkPolicy 切换与 YAML 查看下载）
- 2026-04-02：完成 T144（完成联调验收：go test ./...、npm run build、scripts/p1201_policy_smoke_test.sh 全部通过）
- 2026-04-02：完成 P1201（Policy 管理第一版能力交付完成）
- 2026-04-02：完成 T145（P1202 开发前数据库备份：backups/kubemanage-20260402-102151-p1202.sql）
- 2026-04-02：完成 T146（同步 README/TASKS：P1202 范围、拆分与状态）
- 2026-04-02：完成 T147（后端新增 Policy YAML 保存接口，接入命名空间作用域权限与确认头）
- 2026-04-02：完成 T148（前端 PolicyPage 新增 YAML 编辑保存交互，readonly 仅查看）
- 2026-04-02：完成 T149（联调验收通过：go test ./...、npm run build、scripts/p1202_policy_write_smoke_test.sh）
- 2026-04-02：完成 P1202（Policy YAML 写能力第一版交付完成）
- 2026-04-02：完成 T150（P1203 开发前数据库备份：backups/kubemanage-20260402-111032-p1203.sql）
- 2026-04-02：完成 T151（同步 README/TASKS：P1203 范围、拆分与状态）
- 2026-04-02：完成 T152（后端新增 Policy 删除接口并接入权限/确认头与路由测试）
- 2026-04-02：完成 T153（前端 PolicyPage 新增删除交互与删除后状态刷新）
- 2026-04-02：完成 T154（联调验收通过：go test ./...、npm run build、scripts/p1203_policy_delete_smoke_test.sh）
- 2026-04-02：完成 P1203（Policy 删除能力第一版交付完成）
- 2026-04-02：完成 T155（菜单分组修正：Policy 独立承载三类策略，Security 仅保留权限与审计）
- 2026-04-02：完成 T156（P1204 开发前数据库备份：backups/kubemanage-20260402-120216-p1204.sql）
- 2026-04-02：完成 T157（同步 README/TASKS：P1204 范围、拆分与状态）
- 2026-04-02：完成 T158（后端新增 Policy 创建接口并接入权限/确认头与路由测试）
- 2026-04-02：完成 T159（前端 PolicyPage 新增新建 YAML 交互与创建后刷新）
- 2026-04-02：完成 T160（联调验收通过：go test ./...、npm run build、scripts/p1204_policy_create_smoke_test.sh）
- 2026-04-02：完成 P1204（Policy 新建能力第一版交付完成）
- 2026-04-02：完成 T161（需求调整：Create 按钮文案与命名空间全量下拉选择）
- 2026-04-02：完成 T162（按会话指令暂停新功能开发，并完成任务文档进度同步收口）
- 2026-04-02：完成 T163（README 基础验证口径纠偏，明确 live/mock 能力边界）
- 2026-04-02：完成 T164（`/clusters/switch` 支持 real-only 按连接名激活，并补充路由测试）
- 2026-04-02：完成 T165（回归通过：go test ./...、npm run build，完成阶段状态更新）
- 2026-04-02：完成 T166（第二轮对齐前置数据库备份：backups/kubemanage-20260402-130753-p1302.sql）
- 2026-04-02：完成 T167（Namespace real-only 写链路打通：创建/删除/YAML 查看下载）
- 2026-04-02：完成 T168（Deployment real-only YAML 查看/保存打通并增加一致性校验）
- 2026-04-02：完成 T169（Deployment YAML 写入作用域解析调整为 live 优先）
- 2026-04-02：完成 T170（第二轮回归通过：go test ./...、npm run build，README/TASKS 已同步）
- 2026-04-02：完成 T171（Namespace 前端权限对齐修复：创建/删除按命名空间授权范围控制并补提示）
- 2026-04-02：完成 T172（共享 YamlDialog 重构为可折叠样式，支持展开/收起）
- 2026-04-02：完成 T173（Namespace/Workload/Node/Policy YAML 交互统一生效）
- 2026-04-02：完成 T174（前端构建通过并完成文档与任务状态同步）
- 2026-04-02：完成 T175（字段级折叠改造前置备份完成，并新增 `yaml` 解析依赖）
- 2026-04-02：完成 T176（YAML 结构视图支持 key 级折叠/展开与源码双视图切换）
- 2026-04-02：完成 T177（前端构建通过并完成 README/TASKS 状态同步）
- 2026-04-02：完成 T178（接入 `react-ace/ace-builds`，用于 YAML 编辑器折叠能力）
- 2026-04-02：完成 T179（YAML 编辑默认切换为 Rancher 风格代码编辑器，支持 key 折叠）
- 2026-04-02：完成 T180（前端构建通过并完成 README/TASKS 状态同步）
- 2026-04-02：完成 T181（YAML 工具栏增强前置数据库备份：backups/kubemanage-20260402-133851-p1306.sql）
- 2026-04-02：完成 T182（YAML 编辑器新增导入/下载/还原/变更预览操作）
- 2026-04-02：完成 T183（前端构建通过并完成 README/TASKS 状态同步）
- 2026-04-02：完成 T184（修复 Deployment YAML 按钮点击无反馈，新增错误提示与加载状态）
- 2026-04-02：完成 T185（5 类工作负载 real-only 对齐前置数据库备份：backups/kubemanage-20260402-143130-p1307.sql）
- 2026-04-02：完成 T186（live reader 补齐 Pod/StatefulSet/DaemonSet/Job/CronJob YAML 读写）
- 2026-04-02：完成 T187（5 类写路由作用域解析切换为 live 优先命名空间判定）
- 2026-04-02：完成 T188（后端测试与前端构建通过，并完成 README/TASKS 同步）
- 2026-04-02：完成 T189（Workload YAML 保存反馈增强：成功提示、错误透出、保存后自动刷新）
- 2026-04-02：完成 T190（YAML 保存后保持编辑弹窗并展示 requestId 成功提示）
- 2026-04-02：完成 T191（YAML 弹窗新增保存元信息展示：最近保存时间/requestId/历史）
- 2026-04-02：完成 T192（YAML 保存历史改为按资源隔离记录，避免跨资源覆盖）
- 2026-04-02：完成 T193（按指令停止新功能开发，并完成 README/TASKS 收口同步）
- 2026-04-02：完成 T194（P1401 前置：origin/main 新分支与数据库备份，并同步 README/TASKS 计划）
- 2026-04-02：完成 T195（接入 xterm 终端组件，打通 WebSocket 输入输出与自适应）
- 2026-04-02：完成 T196（Workload 日志弹窗集成真实终端，支持自动连接与手动重连）
- 2026-04-02：完成 T197（P1401 联调验收：go test ./...、npm run build 通过）
