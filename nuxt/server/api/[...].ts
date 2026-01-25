import { auth } from '~/lib/auth'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const goApiUrl = config.goApiUrl || process.env.GO_API_URL || 'http://api:8080'

  // Get the path after /api/
  const path = event.path

  // Build the target URL
  const targetUrl = `${goApiUrl}${path}`

  // Get the request method
  const method = event.method

  // Get request body for non-GET requests
  let body: string | undefined
  if (method !== 'GET' && method !== 'HEAD') {
    body = await readRawBody(event)
  }

  try {
    const { token } = await auth.api.getToken({
      headers: event.headers
    })

    // Make the request to the Go API
    const response = await fetch(targetUrl, {
      method,
      headers: {
        ...event.headers,
        Authorization: token ? `Bearer ${token}` : ''
      },
      body,
      credentials: 'include'
    })

    // Forward response headers
    const responseHeaders: Record<string, string> = {}
    response.headers.forEach((value, key) => {
      // Don't forward certain headers that Nuxt handles
      if (!['transfer-encoding', 'connection', 'keep-alive', 'content-length'].includes(key.toLowerCase())) {
        responseHeaders[key] = value
      }
    })

    // Set response headers
    for (const [key, value] of Object.entries(responseHeaders)) {
      setHeader(event, key, value)
    }

    // Set status code
    setResponseStatus(event, response.status)

    // Return the response body
    const contentType = response.headers.get('content-type') || ''
    if (contentType.includes('application/json')) {
      const obj = await response.json()
      return obj
    }

    return response.text()
  } catch (error) {
    console.error('Proxy error:', error)
    throw createError({
      statusCode: 502,
      statusMessage: 'Bad Gateway',
      message: 'Failed to proxy request to Go API'
    })
  }
})
