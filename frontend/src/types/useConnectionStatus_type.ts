export interface UseConnectionStatusReturn {
  isConnected: boolean;
  lastCheck: Date | null;
  checkConnection: () => Promise<boolean>;
}

export interface ConnectionStatusOptions {
  healthCheckInterval?: number;
  retryInterval?: number;
  timeout?: number;
}

export interface ConnectionStatusHook {
  (
    baseURL: string,
    options?: ConnectionStatusOptions
  ): UseConnectionStatusReturn;
}

export interface HealthCheckResponse {
  status: "ok" | "error";
  timestamp?: string;
  uptime?: number;
  version?: string;
}
