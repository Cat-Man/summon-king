import { runtimeConfig } from '@/config/runtime'

export interface APIResponse<T> {
  code: number
  message: string
  data: T
  trace_id: string
}

function resolveRequestUrl(input: string): string {
  if (/^https?:\/\//.test(input) || input.startsWith('/api/')) {
    return input
  }

  if (input.startsWith('/')) {
    return `${runtimeConfig.apiBaseUrl}${input}`
  }

  return input
}

export async function request<T>(input: string, init: RequestInit = {}): Promise<T> {
  const { headers, ...rest } = init
  const response = await fetch(resolveRequestUrl(input), {
    ...rest,
    headers: {
      'Content-Type': 'application/json',
      ...(headers ?? {})
    }
  })

  if (!response.ok) {
    throw new Error(`request failed with status ${response.status}`)
  }

  const payload = (await response.json()) as APIResponse<T>
  if (payload.code !== 0) {
    throw new Error(payload.message || 'request failed')
  }

  return payload.data
}
