# Contributing Guide

感谢参与《召唤之王》项目。

## 基本原则

- 先文档、后开发：涉及需求、架构、玩法、接口、数据结构变更时，先更新设计文档或实施计划。
- 小步提交：每次提交尽量聚焦单一主题，避免混入无关改动。
- 可验证：提交前至少运行与改动相关的测试、lint 或启动命令。
- 不提交本地原始需求资料、临时文件、密钥和环境变量。

## 分支建议

推荐使用以下分支命名：

- `feat/<name>`：新功能
- `fix/<name>`：缺陷修复
- `docs/<name>`：文档调整
- `chore/<name>`：工程与工具链维护

## 提交信息建议

推荐采用 Conventional Commits：

- `feat:` 新功能
- `fix:` 修复
- `docs:` 文档
- `chore:` 工程配置
- `refactor:` 重构
- `test:` 测试

示例：

```text
feat: add player session bootstrap
docs: refine battle module design
chore: add github actions ci workflow
```

## 提交前检查

在仓库根目录执行：

```bash
make test-backend
make test-frontend
make lint
make dev-api
```

## 文档入口

- `docs/plans/2026-03-21-召唤之王详细设计文档.md`
- `docs/plans/2026-03-21-召唤之王实施计划.md`

## Pull Request 建议

PR 描述建议包含：

- 背景 / 目的
- 具体改动点
- 验证方式
- 风险点与回滚方式
