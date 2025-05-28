import type { AxiosRequestConfig } from "axios";
import type { Model } from "./model";
import type { Document, DocumentUploadResponse } from "./document";
import type { WikiResult } from "./wikiResult";
import type { QueryResponse } from "./chatMessage";

// ==================== API Configuration Types ====================

export interface ApiConfig {
  readonly baseURL: string;
  readonly timeout: number;
  readonly retries: number;
  readonly retryDelay: number;
}

export interface RetryableAxiosRequestConfig extends AxiosRequestConfig {
  __retryCount?: number;
  url?: string; // Add url property for better type safety
}

// ==================== Error Types ====================

export interface ApiErrorOptions {
  message: string;
  status?: number;
  code?: string;
  context?: string;
}

export class ApiError extends Error {
  public readonly status?: number;
  public readonly code?: string;
  public readonly context?: string;

  constructor(
    message: string,
    status?: number,
    code?: string,
    context?: string
  ) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
    this.context = context;
  }

  toJSON() {
    return {
      name: this.name,
      message: this.message,
      status: this.status,
      code: this.code,
      context: this.context,
    };
  }
}

// ==================== HTTP Status Types ====================

export type HttpStatus =
  | 200
  | 201
  | 204 // Success
  | 400
  | 401
  | 403
  | 404
  | 409
  | 422 // Client Error
  | 500
  | 502
  | 503
  | 504; // Server Error

export interface HttpErrorMap {
  400: "Bad Request";
  401: "Unauthorized";
  403: "Forbidden";
  404: "Not Found";
  409: "Conflict";
  422: "Unprocessable Entity";
  500: "Internal Server Error";
  502: "Bad Gateway";
  503: "Service Unavailable";
  504: "Gateway Timeout";
}

// ==================== API Response Types ====================

export interface ApiResponse<T = unknown> {
  data: T;
  message?: string;
  success: boolean;
  timestamp?: number;
}

export interface PaginatedResponse<T> extends ApiResponse<T[]> {
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export interface HealthCheckResponse {
  status: "ok" | "error";
  message: string;
  timestamp: number;
  uptime?: number;
  version?: string;
}

// ==================== Model API Types ====================

export interface ModelListResponse {
  models: Model[];
  total?: number;
}

export interface ModelDownloadRequest {
  name: string;
  url: string;
  version?: string;
  force?: boolean;
}

export interface ModelLoadRequest {
  name: string;
  options?: {
    temperature?: number;
    max_tokens?: number;
    context_length?: number;
  };
}

export interface ModelOperationResponse {
  message: string;
  success: boolean;
  model?: Partial<Model>;
}

// ==================== Document API Types ====================

export interface DocumentListResponse {
  documents: Document[];
  total?: number;
}

export interface DocumentUploadOptions {
  extractText?: boolean;
  generateEmbeddings?: boolean;
  chunkSize?: number;
  overlapSize?: number;
}

export interface DocumentDeleteResponse {
  message: string;
  success: boolean;
  deletedId: string;
}

// ==================== Wiki API Types ====================

export interface WikiSearchRequest {
  query: string;
  limit?: number;
  language?: string;
  includeImages?: boolean;
}

export interface WikiSearchResponse {
  results: WikiResult[];
  totalResults: number;
  searchTime: number;
  language: string;
}

// ==================== AI Query Types ====================

export interface QueryRequest {
  query: string;
  model_name: string;
  include_wiki?: boolean;
  include_documents?: boolean;
  max_tokens?: number;
  temperature?: number;
  stream?: boolean;
  context?: string[];
}

export interface QueryOptions {
  timeout?: number;
  signal?: AbortSignal;
  onProgress?: (chunk: string) => void;
}

// ==================== File Validation Types ====================

export interface FileValidationOptions {
  maxSize: number;
  allowedTypes: string[];
  allowedMimeTypes?: string[];
}

export interface FileValidationResult {
  isValid: boolean;
  errors: string[];
  warnings?: string[];
}

// Export as const value (not type)
export const DEFAULT_FILE_VALIDATION: FileValidationOptions = {
  maxSize: 10 * 1024 * 1024, // 10MB
  allowedTypes: [".pdf", ".txt", ".docx", ".md"],
  allowedMimeTypes: [
    "application/pdf",
    "text/plain",
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "text/markdown",
  ],
};

// ==================== API Service Interface ====================

export interface IApiService {
  // Health
  healthCheck(): Promise<HealthCheckResponse>;
  testConnection(): Promise<boolean>;

  // Models
  getModels(): Promise<ModelListResponse>;
  downloadModel(name: string, url: string): Promise<ModelOperationResponse>;
  loadModel(name: string): Promise<ModelOperationResponse>;
  deleteModel(name: string): Promise<ModelOperationResponse>;

  // Documents
  getDocuments(): Promise<DocumentListResponse>;
  uploadDocument(
    file: File,
    options?: DocumentUploadOptions
  ): Promise<DocumentUploadResponse>;
  deleteDocument(id: string): Promise<DocumentDeleteResponse>;

  // Wiki
  searchWiki(
    query: string,
    options?: Partial<WikiSearchRequest>
  ): Promise<WikiSearchResponse>;

  // AI Query
  query(params: QueryRequest, options?: QueryOptions): Promise<QueryResponse>;
}

// ==================== Request/Response Interceptor Types ====================

export interface RequestInterceptor {
  onRequest?: (
    config: AxiosRequestConfig
  ) => AxiosRequestConfig | Promise<AxiosRequestConfig>;
  onRequestError?: (
    error: Error | ApiError
  ) => Error | ApiError | Promise<Error | ApiError>;
}

export interface ResponseInterceptor {
  onResponse?: (response: unknown) => unknown;
  onResponseError?: (
    error: Error | ApiError
  ) => Error | ApiError | Promise<Error | ApiError>;
}

// ==================== Retry Strategy Types ====================

export interface RetryConfig {
  retries: number;
  retryDelay: number;
  retryCondition?: (error: Error | ApiError) => boolean;
  shouldResetTimeout?: boolean;
}

export type RetryStrategy = "fixed" | "exponential" | "linear";

// ==================== API Client Configuration ====================

export interface ApiClientConfig extends ApiConfig {
  retryConfig?: RetryConfig;
  retryStrategy?: RetryStrategy;
  requestInterceptors?: RequestInterceptor[];
  responseInterceptors?: ResponseInterceptor[];
  validateStatus?: (status: number) => boolean;
  enableLogging?: boolean;
  logLevel?: "debug" | "info" | "warn" | "error";
}

// ==================== Utility Types ====================

export type ApiMethod =
  | "GET"
  | "POST"
  | "PUT"
  | "DELETE"
  | "PATCH"
  | "HEAD"
  | "OPTIONS";

export interface RequestOptions {
  timeout?: number;
  signal?: AbortSignal;
  retries?: number;
  validateStatus?: (status: number) => boolean;
}

export interface LogEntry {
  timestamp: number;
  level: "debug" | "info" | "warn" | "error";
  method: string;
  url: string;
  status?: number;
  duration?: number;
  error?: string;
}

// ==================== Constants ====================

// Export as const value (not type)
export const API_CONSTANTS = {
  DEFAULT_TIMEOUT: 30000,
  FILE_UPLOAD_TIMEOUT: 60000,
  AI_QUERY_TIMEOUT: 120000,
  MAX_RETRIES: 3,
  RETRY_DELAY: 1000,
  MAX_FILE_SIZE: 10 * 1024 * 1024, // 10MB
} as const;

export const HTTP_STATUS_MESSAGES: Record<number, string> = {
  200: "OK",
  201: "Created",
  204: "No Content",
  400: "Bad Request",
  401: "Unauthorized",
  403: "Forbidden",
  404: "Not Found",
  409: "Conflict",
  422: "Unprocessable Entity",
  500: "Internal Server Error",
  502: "Bad Gateway",
  503: "Service Unavailable",
  504: "Gateway Timeout",
} as const;

export const ERROR_CODES = {
  NETWORK_ERROR: "NETWORK_ERROR",
  TIMEOUT: "ECONNABORTED",
  CONNECTION_REFUSED: "ECONNREFUSED",
  UNKNOWN: "UNKNOWN_ERROR",
} as const;
