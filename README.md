# summon-king

《召唤之王》项目仓库。

当前仓库为首版 monorepo 工程骨架，包含：
- Go 后端服务骨架
- H5 游戏前端骨架
- 管理后台骨架
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

当前为第一版初始化提交，重点完成：
- monorepo 基础目录创建
- 根目录构建脚本与工作区配置
- 后端 `bootstrap` 配置加载最小实现
- 基础测试命令、开发命令占位

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
make lint
make dev-api
```

## 文档入口

- 详细设计文档：`docs/plans/2026-03-21-召唤之王详细设计文档.md`
- 实施计划：`docs/plans/2026-03-21-召唤之王实施计划.md`

## 说明

原始需求资料文件体积较大，仅保留在本地工作目录中，不提交到仓库。
