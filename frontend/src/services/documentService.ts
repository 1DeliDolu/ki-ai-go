import type { Document, DocumentUploadResponse } from "../types/document";

const API_BASE_URL = "http://localhost:8082/api/v1";

// Interface for backend document response
interface BackendDocument {
  id: string | number;
  name: string;
  type: string;
  size: number;
  uploadDate?: string;
  upload_date?: string;
  status?: string;
  chunks?: number;
  embeddings?: boolean;
  path?: string;
}

interface BackendDocumentsResponse {
  documents: BackendDocument[];
}

export class DocumentService {
  private normalizeFileType(type: string): "pdf" | "docx" | "txt" | "md" {
    const cleanType = type.replace(".", "").toLowerCase();
    switch (cleanType) {
      case "pdf":
        return "pdf";
      case "docx":
      case "doc":
        return "docx";
      case "txt":
        return "txt";
      case "md":
      case "markdown":
        return "md";
      default:
        return "txt"; // fallback to txt for unknown types
    }
  }

  private normalizeStatus(status?: string): "processing" | "ready" | "error" {
    switch (status?.toLowerCase()) {
      case "processing":
        return "processing";
      case "error":
      case "failed":
        return "error";
      case "ready":
      case "completed":
      case "done":
      default:
        return "ready";
    }
  }

  async listDocuments(): Promise<Document[]> {
    try {
      const response = await fetch(`${API_BASE_URL}/documents`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data: BackendDocumentsResponse | BackendDocument[] =
        await response.json();
      console.log("Documents from backend:", data); // Debug log

      // Extract documents array from the response
      let documents: BackendDocument[] = [];
      if (Array.isArray(data)) {
        documents = data;
      } else if (data.documents && Array.isArray(data.documents)) {
        documents = data.documents;
      } else {
        console.warn("Unexpected response format:", data);
        return [];
      }

      // Normalize the document data to match frontend expectations
      const normalizedDocs = documents.map(
        (doc: BackendDocument): Document => ({
          id: doc.id?.toString() || "",
          name: doc.name || "",
          type: this.normalizeFileType(doc.type),
          size: doc.size || 0,
          uploadDate:
            doc.uploadDate || doc.upload_date || new Date().toISOString(),
          status: this.normalizeStatus(doc.status),
          chunks: doc.chunks,
          embeddings: doc.embeddings,
          path: doc.path,
        })
      );

      console.log("Normalized documents:", normalizedDocs); // Debug log
      console.log("Returning array:", Array.isArray(normalizedDocs)); // Debug log

      return normalizedDocs;
    } catch (error) {
      console.error("Failed to fetch documents:", error);
      return [];
    }
  }

  async uploadDocument(file: File): Promise<DocumentUploadResponse> {
    try {
      const formData = new FormData();
      formData.append("file", file);

      const response = await fetch(`${API_BASE_URL}/documents/upload`, {
        method: "POST",
        body: formData,
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(
          data.message || `HTTP error! status: ${response.status}`
        );
      }

      return data;
    } catch (error) {
      console.error("Failed to upload document:", error);
      throw error;
    }
  }

  async deleteDocument(id: string): Promise<void> {
    try {
      const response = await fetch(`${API_BASE_URL}/documents/${id}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(
          data.message || `HTTP error! status: ${response.status}`
        );
      }
    } catch (error) {
      console.error("Failed to delete document:", error);
      throw error;
    }
  }

  async checkConnection(): Promise<boolean> {
    try {
      const response = await fetch(`${API_BASE_URL}/health`, {
        method: "HEAD",
      });
      return response.ok;
    } catch (error) {
      console.error("Backend connection check failed:", error);
      return false;
    }
  }
}

export const documentService = new DocumentService();
