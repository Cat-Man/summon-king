// @vitest-environment node

import { expect, test } from 'vitest'

import config from '../../../vite.config'

test('proxies api requests to local backend in dev', async () => {
  const resolved = (
    typeof config === 'function'
      ? await (config as (env: {
          command: 'serve',
          mode: 'development',
          isPreview: false,
          isSsrBuild: false
        }) => unknown)({
          command: 'serve',
          mode: 'development',
          isPreview: false,
          isSsrBuild: false
        })
      : config
  ) as {
    server?: {
      proxy?: Record<string, unknown>
    }
  }

  expect(resolved.server?.proxy?.['/api']).toMatchObject({
    target: 'http://localhost:8080',
    changeOrigin: true
  })
})
