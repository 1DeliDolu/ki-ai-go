export interface WikiResult {
  pageId: string;
  title: string;
  url: string;
  description?: string;
  extract?: string;
  thumbnail?: string;
  relevanceScore?: number;
}

export interface WikiSearchResponse {
  results: WikiResult[];
  totalResults: number;
  searchTime: number;
}