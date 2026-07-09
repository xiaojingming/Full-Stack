type Envelope<T = unknown> = {
  data: T
  error?: string
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })

  const body: Envelope<T> = await res.json()

  if (body.error) {
    throw new Error(body.error)
  }

  return body.data as T
}

export function query<T = unknown>(endpoint: string): Promise<T> {
  return request<T>(`/api/${endpoint}`)
}

export function exec(endpoint: string, params?: Record<string, unknown>): Promise<number> {
  return request<number>(`/api/${endpoint}`, {
    method: 'POST',
    body: params ? JSON.stringify(params) : undefined,
  })
}
