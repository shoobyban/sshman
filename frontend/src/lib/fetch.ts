import { ApiError, ApiErrorPayload } from './errors';

/**
 * Wrapper around fetch that handles JSON error responses from the backend.
 * Throws ApiError for non-2xx responses with structured error data.
 */
export async function apiFetch<T = unknown>(
  input: RequestInfo | URL,
  init?: RequestInit
): Promise<T> {
  let response: Response;

  try {
    response = await fetch(input, init);
  } catch (e) {
    // Network error (no response)
    throw new ApiError({
      error: {
        message: 'Network error. Please check your connection.',
        details: e instanceof Error ? e.message : String(e),
        code: 0,
      },
    });
  }

  // Try to read response body
  const text = await response.text();

  // Parse JSON if possible
  let json: unknown = null;
  try {
    json = text ? JSON.parse(text) : null;
  } catch {
    // Not JSON
    if (!response.ok) {
      // Non-JSON error response
      throw new ApiError({
        error: {
          message: response.statusText || 'An error occurred',
          details: text || 'No details available',
          code: response.status,
        },
      });
    }
    // Success but non-JSON body - return text
    return text as T;
  }

  // Handle error responses
  if (!response.ok) {
    const errorPayload = json as ApiErrorPayload;
    if (errorPayload && errorPayload.error) {
      throw new ApiError(errorPayload);
    } else {
      // Unexpected error format
      throw new ApiError({
        error: {
          message: response.statusText || 'An error occurred',
          details: JSON.stringify(json),
          code: response.status,
        },
      });
    }
  }

  return json as T;
}
