import { useState, useEffect, useCallback, useRef } from "react";
import type {
  UseConnectionStatusReturn,
  ConnectionStatusOptions,
} from "../types/useConnectionStatus_type";

const HEALTH_CHECK_INTERVAL = 30000; // 30 seconds
const RETRY_INTERVAL = 5000; // 5 seconds when disconnected

export function useConnectionStatus(
  baseURL: string,
  options: ConnectionStatusOptions = {}
): UseConnectionStatusReturn {
  const {
    healthCheckInterval = HEALTH_CHECK_INTERVAL,
    retryInterval = RETRY_INTERVAL,
    timeout = 5000,
  } = options;

  const [isConnected, setIsConnected] = useState(false);
  const [lastCheck, setLastCheck] = useState<Date | null>(null);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  const retryTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const mountedRef = useRef(true);

  const checkConnection = useCallback(async () => {
    try {
      const response = await fetch(`${baseURL}/api/v1/health`, {
        method: "GET",
        headers: {
          "Cache-Control": "no-cache",
        },
        signal: AbortSignal.timeout(timeout),
      });

      const connected = response.ok;

      if (mountedRef.current) {
        setIsConnected(connected);
        setLastCheck(new Date());
      }

      return connected;
    } catch {
      if (mountedRef.current) {
        setIsConnected(false);
        setLastCheck(new Date());
      }
      return false;
    }
  }, [baseURL, timeout]);

  const startPolling = useCallback(() => {
    // Clear existing intervals
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
    }
    if (retryTimeoutRef.current) {
      clearTimeout(retryTimeoutRef.current);
    }

    const poll = async () => {
      const connected = await checkConnection();

      if (mountedRef.current) {
        if (connected) {
          // If connected, use normal interval
          intervalRef.current = setTimeout(poll, healthCheckInterval);
        } else {
          // If disconnected, retry more frequently
          retryTimeoutRef.current = setTimeout(poll, retryInterval);
        }
      }
    };

    // Start polling
    poll();
  }, [checkConnection, healthCheckInterval, retryInterval]);

  useEffect(() => {
    mountedRef.current = true;
    startPolling();

    return () => {
      mountedRef.current = false;
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
      if (retryTimeoutRef.current) {
        clearTimeout(retryTimeoutRef.current);
      }
    };
  }, [startPolling]);

  const forceCheck = useCallback(() => {
    return checkConnection();
  }, [checkConnection]);

  return {
    isConnected,
    lastCheck,
    checkConnection: forceCheck,
  };
}
