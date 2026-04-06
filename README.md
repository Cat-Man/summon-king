# summon-king

《召唤之王》项目仓库。

当前仓库已经从纯骨架推进到“后端可跑 + H5 首条闭环可联调”的阶段，包含：
- Go 后端 API（健康检查、访客登录、首页聚合、地图/副本/修行链路）
- H5 游戏前端（Vue 3 + Vite + Router + Pinia）
- 管理后台占位工程
- 微信小程序壳工程骨架
- 详细设计文档与实施计划

## 技术栈

- 后端：Go 1.22+ / Gin / MySQL / Redis
- 玩家前端：Vue 3 / TypeScript / Vite / Pinia / Vant
- 后台前端：Vue 3 / TypeScript / Element Plus
- 小程序：mini-shell + web-view
- 部署：Docker / Nginx

## 目录结构

```text
apps/
  backend/      Go 后端
  game-web/     玩家 H5
  admin-web/    运营/管理后台
  mini-shell/   微信小程序壳

docs/
  plans/        详细设计与实施计划
```

## 当前状态

当前主线重点完成：
- 后端统一响应、中间件、模块路由与 `home/overview`
- 访客登录与前端会话链路
- H5 正式工程壳、路由壳和首条“首页 -> 地图/副本 -> 修行”可联调闭环
- `game-web` 真实单测、类型检查与生产构建

## 本地启动

### 1. 安装依赖

请先确保本机安装：
- Go
- Node.js
- pnpm

### 2. 常用命令

```bash
make test-backend
make test-frontend
make build-frontend
make lint
make smoke-backend
make dev-api
```

说明：
- `make test-frontend` 当前会执行 `apps/game-web` 的真实单测和生产构建。
- `apps/admin-web` 仍是占位工程，因此未再纳入仓库级验证，避免出现“占位脚本通过但被误认为前端已验完”的假通过。
- `make smoke-backend` 只做后端入口编译烟测，不会像直接运行服务那样阻塞 CI。

## 文档入口

- 详细设计文档：`docs/plans/2026-03-21-召唤之王详细设计文档.md`
- 当前阶段实施计划：`docs/plans/2026-04-02-底座与首条闭环实施计划.md`

## 说明

原始需求资料文件体积较大，仅保留在本地工作目录中，不提交到仓库。
