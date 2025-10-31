export interface ApiErrorPayload {
  error: {
    message: string;
    details?: string;
    code: number;
    id?: string;
    meta?: Record<string, unknown>;
  };
}

export class ApiError extends Error {
  public readonly details?: string;
  public readonly statusCode: number;
  public readonly id?: string;
  public readonly meta?: Record<string, unknown>;

  constructor(payload: ApiErrorPayload) {
    super(payload.error.message);
    this.name = 'ApiError';
    this.details = payload.error.details;
    this.statusCode = payload.error.code;
    this.id = payload.error.id;
    this.meta = payload.error.meta;
  }
}
