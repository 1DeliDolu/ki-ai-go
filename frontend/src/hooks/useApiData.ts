import { useState, useEffect, useCallback, useRef } from "react";
import type {
  UseApiDataOptions,
  ApiCache,
  UseApiDataReturn,
} from "../types/useApiData_types";

const apiCache: ApiCache = {};
const ongoingRequests = new Map<string, Promise<unknown>>();

export function useApiData<T>(
  key: string,
  fetcher: () => Promise<T>,
  options: UseApiDataOptions = {}
): UseApiDataReturn<T> {
  const {
    pollInterval = 30000, // 30 seconds default
    enabled = true,
    cacheTime = 10000, // 10 seconds cache
  } = options;

  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const intervalRef = useRef<NodeJS.Timeout | undefined>(undefined);
  const mountedRef = useRef(true);

  const fetchData = useCallback(
    async (force = false): Promise<T | undefined> => {
      if (!enabled) return undefined;

      // Check cache first
      const cached = apiCache[key];
      const now = Date.now();

      if (!force && cached && now < cached.expires) {
        setData(cached.data as T);
        return cached.data as T;
      }

      // Check if there's an ongoing request
      const ongoing = ongoingRequests.get(key);
      if (ongoing) {
        try {
          const result = await ongoing;
          if (mountedRef.current) {
            setData(result as T);
          }
          return result as T;
        } catch (err) {
          if (mountedRef.current) {
            setError(err as Error);
          }
          throw err;
        }
      }

      // Create new request
      setLoading(true);
      setError(null);

      const request = fetcher();
      ongoingRequests.set(key, request);

      try {
        const result = await request;

        // Cache the result
        apiCache[key] = {
          data: result,
          timestamp: now,
          expires: now + cacheTime,
        };

        if (mountedRef.current) {
          setData(result);
          setLoading(false);
        }

        return result;
      } catch (err) {
        if (mountedRef.current) {
          setError(err as Error);
          setLoading(false);
        }
        throw err;
      } finally {
        ongoingRequests.delete(key);
      }
    },
    [key, fetcher, enabled, cacheTime]
  );

  const refresh = useCallback(() => {
    return fetchData(true);
  }, [fetchData]);

  useEffect(() => {
    mountedRef.current = true;

    if (enabled) {
      // Initial fetch
      fetchData();

      // Set up polling if interval is specified
      if (pollInterval > 0) {
        intervalRef.current = setInterval(() => {
          fetchData();
        }, pollInterval);
      }
    }

    return () => {
      mountedRef.current = false;
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [fetchData, pollInterval, enabled]);

  useEffect(() => {
    return () => {
      mountedRef.current = false;
    };
  }, []);

  return {
    data,
    loading,
    error,
    refresh,
    refetch: refresh,
  };
}

// Utility to clear cache
export function clearApiCache(key?: string) {
  if (key) {
    delete apiCache[key];
  } else {
    Object.keys(apiCache).forEach((k) => delete apiCache[k]);
  }
}
