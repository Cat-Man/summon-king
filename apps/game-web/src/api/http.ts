interface APIResponse<T> {
  code: number
  message: string
  data: T
  trace_id: string
}

export async function request<T>(input: string, init: RequestInit = {}): Promise<T> {
  const { headers, ...rest } = init
  const response = await fetch(input, {
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
