import { spawnSync } from 'node:child_process'

const cmd = process.platform === 'win32' ? 'pnpm.cmd' : 'pnpm'
const result = spawnSync(cmd, ['exec', 'vitest', 'run'], {
  stdio: 'inherit'
})

process.exit(result.status ?? 1)
