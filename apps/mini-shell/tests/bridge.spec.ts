import { buildGameURL } from '../utils/bridge'

test('buildGameURL adds channel and token', () => {
  expect(buildGameURL('https://game.xxx.com', 'abc')).toContain('channel=wxmini')
})
