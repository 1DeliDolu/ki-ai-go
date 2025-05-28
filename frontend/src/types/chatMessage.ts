export interface QueryResponse {
  response: string;
  sources: {
    documents: Source[];
    wiki: Source[];
  };
}

export interface ChatMessage {
  id: string;
  content: string;
  timestamp: Date;
  isUser: boolean;
}

export interface Source {
  type: "document" | "wiki";
  title: string;
  content: string;
  relevanceScore: number;
  documentId?: string;
  url?: string;
}
export interface ChatInterfaceProps {
  onQuery: (query: string, includeWiki?: boolean) => Promise<void>;
  loading: boolean;
  response: QueryResponse | null;
  selectedModel: string;
  isConnected: boolean;
}
