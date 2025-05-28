export interface UseApiDataOptions {
  pollInterval?: number;
  enabled?: boolean;
  cacheTime?: number;
}

export interface ApiCacheItem<T = unknown> {
  data: T;
  timestamp: number;
  expires: number;
}

export interface ApiCache {
  [key: string]: ApiCacheItem<unknown>;
}

export interface UseApiDataReturn<T> {
  data: T | null;
  loading: boolean;
  error: Error | null;
  refresh: () => Promise<T | undefined>;
  refetch: () => Promise<T | undefined>;
}

export interface ApiDataHook {
  <T>(
    key: string,
    fetcher: () => Promise<T>,
    options?: UseApiDataOptions
  ): UseApiDataReturn<T>;
}
