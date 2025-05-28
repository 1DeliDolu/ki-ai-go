export interface Document {
  id: string;
  name: string;
  type: "pdf" | "docx" | "txt" | "md";
  size: number;
  uploadDate: string;
  status: "processing" | "ready" | "error";
  chunks?: number;
  embeddings?: boolean;
  path?: string; // Add path field for backend compatibility
}

export interface DocumentUploadResponse {
  success: boolean;
  document?: Document;
  message: string;
}

export interface DocumentUploaderProps {
  onUpload: (file: File) => Promise<void>;
  isConnected: boolean;
}

export interface DocumentListProps {
  documents: Document[];
  onRefresh: () => Promise<void>;
  onDelete: (id: string) => Promise<void>;
  isConnected: boolean;
}
