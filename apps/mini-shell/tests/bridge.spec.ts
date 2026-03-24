import { buildGameURL, buildMiniLoginPayload, buildMiniPayPayload } from '../utils/bridge'

test('buildGameURL adds channel and token', () => {
  expect(buildGameURL('https://game.xxx.com', 'abc')).toContain('channel=wxmini')
})

test('build mini bridge payloads', () => {
  expect(buildMiniLoginPayload('code-1')).toEqual({ code: 'code-1', channel: 'wxmini' })
  expect(buildMiniPayPayload('ord-1')).toEqual({ order_no: 'ord-1', channel: 'wxmini' })
})
