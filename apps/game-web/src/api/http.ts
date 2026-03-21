export async function request<T>(input: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(input, {
    headers: {
      'Content-Type': 'application/json',
      ...(init.headers ?? {})
    },
    ...init
  })

  if (!response.ok) {
    throw new Error(`request failed with status ${response.status}`)
  }

  return response.json() as Promise<T>
}
