import axios from "axios";
import type { AxiosInstance, AxiosError, AxiosRequestConfig } from "axios";
import type { DocumentUploadResponse } from "../types/document";
import type { QueryResponse } from "../types/chatMessage";

// Import our comprehensive type definitions
import type {
  ApiConfig,
  RetryableAxiosRequestConfig,
  IApiService,
  HealthCheckResponse,
  ModelListResponse,
  ModelOperationResponse,
  DocumentListResponse,
  DocumentDeleteResponse,
  WikiSearchResponse,
  FileValidationOptions,
  FileValidationResult,
} from "../types/api_ts_type";

// Import values (not types) separately
import {
  ApiError,
  DEFAULT_FILE_VALIDATION,
  API_CONSTANTS,
} from "../types/api_ts_type";

// API Configuration with imported types
const API_CONFIG: ApiConfig = {
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8082/api/v1",
  timeout: API_CONSTANTS.DEFAULT_TIMEOUT,
  retries: API_CONSTANTS.MAX_RETRIES,
  retryDelay: API_CONSTANTS.RETRY_DELAY,
} as const;

// File validation utility with enhanced error handling
const validateFile = (
  file: File,
  options: FileValidationOptions = DEFAULT_FILE_VALIDATION
): FileValidationResult => {
  const errors: string[] = [];

  // Null/undefined check
  if (!file) {
    errors.push("No file provided");
    return { isValid: false, errors };
  }

  // Size validation
  if (file.size > options.maxSize) {
    errors.push(
      `File size (${(file.size / 1024 / 1024).toFixed(1)}MB) exceeds ${(
        options.maxSize /
        1024 /
        1024
      ).toFixed(1)}MB limit`
    );
  }

  // Type validation with better error handling
  const fileName = file.name || "unknown";
  const fileExtension = "." + (fileName.split(".").pop()?.toLowerCase() || "");

  if (!options.allowedTypes.includes(fileExtension)) {
    errors.push(
      `File type ${fileExtension} not supported. Allowed: ${options.allowedTypes.join(
        ", "
      )}`
    );
  }

  // MIME type validation if available
  if (options.allowedMimeTypes && file.type) {
    if (!options.allowedMimeTypes.includes(file.type)) {
      errors.push(
        `MIME type ${
          file.type
        } not supported. Allowed: ${options.allowedMimeTypes.join(", ")}`
      );
    }
  }

  return {
    isValid: errors.length === 0,
    errors,
  };
};

// Create axios instance with modern configuration
const createApiClient = (): AxiosInstance => {
  const client = axios.create({
    baseURL: API_CONFIG.baseURL,
    timeout: API_CONFIG.timeout,
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    withCredentials: false,
  });

  // Request interceptor with proper logging
  client.interceptors.request.use(
    (config) => {
      const method = config.method?.toUpperCase() || "REQUEST";
      const url = config.url || "unknown";
      console.log(`🚀 ${method} ${url}`);
      return config;
    },
    (error: AxiosError) => {
      console.error("❌ Request error:", error.message);
      return Promise.reject(
        new ApiError("Request failed", undefined, error.code)
      );
    }
  );

  // Response interceptor with enhanced retry logic
  client.interceptors.response.use(
    (response) => {
      const status = response.status;
      const url = response.config.url || "unknown";
      console.log(`✅ ${status} ${url}`);
      return response;
    },
    async (error: AxiosError) => {
      const status = error.response?.status;
      const message = error.message;
      console.error("❌ Response error:", status, message);

      // Type-safe retry logic
      const config = error.config as RetryableAxiosRequestConfig;

      if (!config) {
        return Promise.reject(
          new ApiError("No config available for retry", status, error.code)
        );
      }

      // Retry logic for network errors and timeouts
      const shouldRetry =
        error.code === "NETWORK_ERROR" ||
        error.code === "ECONNABORTED" ||
        error.code === "ERR_NETWORK" ||
        (status && status >= 500);

      if (shouldRetry) {
        config.__retryCount = (config.__retryCount || 0) + 1;

        if (config.__retryCount <= API_CONFIG.retries) {
          console.log(
            `🔄 Retrying request (${config.__retryCount}/${
              API_CONFIG.retries
            }) for ${config.url || "unknown"}`
          );

          await new Promise((resolve) =>
            setTimeout(resolve, API_CONFIG.retryDelay * config.__retryCount!)
          );

          return client.request(config as AxiosRequestConfig);
        }
      }

      // Create meaningful error messages
      let errorMessage = "Request failed";
      if (status === 403) {
        errorMessage = "Access forbidden - check CORS configuration";
      } else if (status === 404) {
        errorMessage = "Endpoint not found";
      } else if (status === 500) {
        errorMessage = "Server error";
      } else if (error.code === "NETWORK_ERROR") {
        errorMessage = "Network connection failed";
      } else if (error.code === "ECONNABORTED") {
        errorMessage = "Request timeout";
      }

      return Promise.reject(new ApiError(errorMessage, status, error.code));
    }
  );

  return client;
};

// Modern API Service class implementing the interface
export class ApiService implements IApiService {
  private readonly api: AxiosInstance;

  constructor() {
    this.api = createApiClient();
  }

  // Generic error handler
  private handleError(error: unknown, context: string): never {
    if (error instanceof ApiError) {
      throw error;
    }

    if (axios.isAxiosError(error)) {
      const status = error.response?.status;
      const message = `${context}: ${error.message}`;
      throw new ApiError(message, status, error.code);
    }

    throw new ApiError(`${context}: Unknown error occurred`);
  }

  // Health check with proper typing
  async healthCheck(): Promise<HealthCheckResponse> {
    try {
      const response = await this.api.get("/health");
      return response.data;
    } catch (error) {
      console.error("🏥 Health check failed:", error);
      this.handleError(error, "Health check failed");
    }
  }

  // Model management with proper return types
  async getModels(): Promise<ModelListResponse> {
    try {
      const response = await this.api.get("/models");
      return response.data;
    } catch (error) {
      console.error("🎯 Failed to fetch models:", error);
      this.handleError(error, "Failed to fetch models");
    }
  }

  async downloadModel(
    name: string,
    url: string
  ): Promise<ModelOperationResponse> {
    if (!name || !url) {
      throw new ApiError("Model name and URL are required");
    }

    try {
      const response = await this.api.post("/models/download", { name, url });
      return response.data;
    } catch (error) {
      console.error("⬇️ Model download failed:", error);
      this.handleError(error, `Failed to download model: ${name}`);
    }
  }

  async loadModel(name: string): Promise<ModelOperationResponse> {
    if (!name) {
      throw new ApiError("Model name is required");
    }

    try {
      const response = await this.api.post("/models/load", { name });
      return response.data;
    } catch (error) {
      console.error("🔄 Model load failed:", error);
      this.handleError(error, `Failed to load model: ${name}`);
    }
  }

  async deleteModel(name: string): Promise<ModelOperationResponse> {
    if (!name) {
      throw new ApiError("Model name is required");
    }

    try {
      const response = await this.api.delete(
        `/models/${encodeURIComponent(name)}`
      );
      return response.data;
    } catch (error) {
      console.error("🗑️ Model deletion failed:", error);
      this.handleError(error, `Failed to delete model: ${name}`);
    }
  }

  // Document management with enhanced validation
  async getDocuments(): Promise<DocumentListResponse> {
    try {
      const response = await this.api.get("/documents");
      return response.data;
    } catch (error) {
      console.error("📁 Failed to fetch documents:", error);
      this.handleError(error, "Failed to fetch documents");
    }
  }

  async uploadDocument(file: File): Promise<DocumentUploadResponse> {
    if (!file) {
      throw new ApiError("File is required for upload");
    }

    // Use the validation utility
    const validation = validateFile(file);
    if (!validation.isValid) {
      throw new ApiError(validation.errors.join("; "));
    }

    try {
      const formData = new FormData();
      formData.append("file", file);

      const response = await this.api.post("/documents/upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
        timeout: API_CONSTANTS.FILE_UPLOAD_TIMEOUT,
      });
      return response.data;
    } catch (error) {
      console.error("📤 Document upload failed:", error);
      this.handleError(error, `Failed to upload document: ${file.name}`);
    }
  }

  async deleteDocument(id: string): Promise<DocumentDeleteResponse> {
    try {
      const response = await this.api.delete(`/documents/${id}`);
      return {
        success: true,
        message: response.data.message || "Document deleted successfully",
        deletedId: id,
      };
    } catch (error) {
      this.handleError(error, "Failed to delete document");
    }
  }

  // Wiki search with proper typing
  async searchWiki(query: string): Promise<WikiSearchResponse> {
    if (!query || query.trim().length === 0) {
      throw new ApiError("Search query is required");
    }

    try {
      const encodedQuery = encodeURIComponent(query.trim());
      const response = await this.api.get(`/wiki/search?q=${encodedQuery}`);
      return response.data;
    } catch (error) {
      console.error("🌐 Wiki search failed:", error);
      this.handleError(error, `Failed to search Wikipedia for: ${query}`);
    }
  }

  // AI Query with enhanced typing
  async query(params: {
    query: string;
    include_wiki: boolean;
    model_name: string;
  }): Promise<QueryResponse> {
    if (!params.query || params.query.trim().length === 0) {
      throw new ApiError("Query text is required");
    }

    if (!params.model_name || params.model_name.trim().length === 0) {
      throw new ApiError("Model name is required");
    }

    try {
      const response = await this.api.post(
        "/query",
        {
          query: params.query.trim(),
          include_wiki: Boolean(params.include_wiki),
          model_name: params.model_name.trim(),
        },
        {
          timeout: API_CONSTANTS.AI_QUERY_TIMEOUT,
        }
      );
      return response.data;
    } catch (error) {
      console.error("🤖 AI query failed:", error);
      this.handleError(error, `Failed to process query: ${params.query}`);
    }
  }

  // Connection test method
  async testConnection(): Promise<boolean> {
    try {
      await this.healthCheck();
      return true;
    } catch {
      return false;
    }
  }
}

// Export singleton instance with error boundary
export default new ApiService();

// Export the error class for use in components
export { ApiError };
