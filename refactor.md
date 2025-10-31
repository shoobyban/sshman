# Refactoring Plan: Standardized JSON Error Responses

## 1. Objective

To refactor the backend API and frontend client to use a standardized JSON format for error responses. This will allow detailed, actionable error messages from the backend to be displayed directly to the user in the frontend, improving debuggability and user experience.

## 2. The Problem

Currently, the backend API uses `http.Error`, which sends errors as a plain `text/plain` response.

**Backend Code (Example):**
```go
// api/hosts.go
http.Error(w, err.Error(), http.StatusInternalServerError)
```

The frontend API client catches non-2xx responses but can only throw a generic, hardcoded error message because it cannot parse the plain text response for details.

**Frontend Code (Example):**
```typescript
// frontend/src/lib/api.ts
export const createHost = async (host: Partial<Host>): Promise<Host> => {
  const response = await fetch(`${API_BASE}/hosts`, { /* ... */ });
  // This check is correct, but the error thrown is generic.
  if (!response.ok) throw new Error("Failed to create host");
  return response.json();
};

// frontend/src/pages/HostForm.tsx
toast({
  title: "Error creating host",
  description: "An unknown error occurred. Please try again.", // The user sees this.
  variant: "destructive",
});
```

This leads to a poor user experience, where a specific backend error like `dial tcp: address homeassistant.local: missing port in address` is hidden behind a generic message like "Failed to create host. Please try again."

## 3. Proposed Solution: JSON Error Structure

We will implement a consistent JSON structure for all error responses.

**Proposed JSON Error Structure:**
```json
{
  "error": {
    "message": "A user-friendly summary of the error.",
    "details": "The detailed, technical error message from the backend.",
    "code": 400
  }
}
```
*   `message`: A high-level, safe-to-display error message.
*   `details`: The raw `error.Error()` string from the Go backend. This is invaluable for debugging.
*   `code`: The HTTP status code.

---

## 4. Backend Refactoring Plan

### Step 4.1: Create a JSON Error Helper

In the `api` package, we will create a new file `api/errors.go` (or add to `api/main.go`) with a helper function to standardize sending JSON errors.

**File: `api/errors.go`**
```go
package api

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the structure for our standard JSON error.
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Details string `json:"details,omitempty"`
		Code    int    `json:"code"`
	} `json:"error"`
}

// JSONError sends a structured JSON error response.
func JSONError(w http.ResponseWriter, message string, details string, code int) {
	response := ErrorResponse{}
	response.Error.Message = message
	response.Error.Details = details
	response.Error.Code = code

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// As a fallback, write the plain text error
		http.Error(w, message, code)
	}
}
```

### Step 4.2: Update API Handlers

Replace all instances of `http.Error()` in `api/hosts.go`, `api/users.go`, and `api/groups.go` with our new `JSONError` helper.

**Before:**
```go
// api/hosts.go
if err != nil {
    cfg.Log().Errorf("Can't decode host %s", err)
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
}
```

**After:**
```go
// api/hosts.go
if err != nil {
    details := err.Error()
    cfg.Log().Errorf("Can't decode host: %s", details)
    JSONError(w, "Invalid request body.", details, http.StatusBadRequest)
    return
}
```
This change needs to be applied to every error path in the API handlers.

---

## 5. Frontend Refactoring Plan

### Step 5.1: Create a Custom Error Type

To better handle structured errors, we'll define a custom error type.

**File: `frontend/src/lib/errors.ts` (New File)**
```typescript
export interface ApiErrorPayload {
  error: {
    message: string;
    details?: string;
    code: number;
  };
}

export class ApiError extends Error {
  public readonly details?: string;
  public readonly statusCode: number;

  constructor(payload: ApiErrorPayload) {
    super(payload.error.message);
    this.name = 'ApiError';
    this.details = payload.error.details;
    this.statusCode = payload.error.code;
  }
}
```

### Step 5.2: Update API Client to Parse JSON Errors

Modify the API functions in `frontend/src/lib/api.ts` to parse the JSON error body and throw our new `ApiError`.

**Before:**
```typescript
// frontend/src/lib/api.ts
export const createHost = async (host: Partial<Host>): Promise<Host> => {
  const response = await fetch(`${API_BASE}/hosts`, { /* ... */ });
  if (!response.ok) throw new Error("Failed to create host");
  return response.json();
};
```

**After:**
```typescript
// frontend/src/lib/api.ts
import { ApiError, ApiErrorPayload } from './errors'; // Import new types

export const createHost = async (host: Partial<Host>): Promise<Host> => {
  const response = await fetch(`${API_BASE}/hosts`, { /* ... */ });
  if (!response.ok) {
    const errorPayload: ApiErrorPayload = await response.json();
    throw new ApiError(errorPayload);
  }
  return response.json();
};
```
This pattern should be applied to **all** API request functions (`getHosts`, `updateUser`, `deleteGroup`, etc.).

### Step 5.3: Update UI Components to Display Detailed Errors

Finally, update the `onError` callbacks in our React Query mutations to display the detailed error message.

**Before:**
```tsx
// frontend/src/pages/HostForm.tsx
const mutation = useMutation({
  mutationFn: createHost,
  onError: () => {
    toast({
      title: "Error creating host",
      description: "An unknown error occurred. Please try again.",
      variant: "destructive",
    });
  },
});
```

**After:**
```tsx
// frontend/src/pages/HostForm.tsx
import { ApiError } from '@/lib/errors'; // Import custom error

const mutation = useMutation({
  mutationFn: createHost,
  onError: (error) => {
    let description = "An unknown error occurred. Please try again.";
    if (error instanceof ApiError && error.details) {
      description = error.details; // Use the detailed error from the backend!
    } else if (error instanceof Error) {
      description = error.message;
    }
    
    toast({
      title: "Error Creating Host",
      description: description, // Display the detailed message
      variant: "destructive",
    });
  },
});
```
This pattern should be applied to all `useMutation` hooks throughout the application. The toast component may need styling adjustments to better display potentially long, technical error messages (e.g., using a `<pre>` tag or allowing more vertical space).

## 6. Benefits

1.  **Improved User Experience:** Users see specific, meaningful errors, helping them resolve issues on their own (e.g., fixing a malformed hostname).
2.  **Faster Debugging:** Developers can immediately see the root cause of an issue from frontend screenshots or user reports.
3.  **Consistent API:** The API becomes more robust and predictable.
4.  **Extensibility:** The structured error format can be extended in the future with more context, like error codes or trace IDs.
