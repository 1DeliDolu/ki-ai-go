import React, { useState, useEffect, useCallback } from "react";
import { Box, Typography, Alert } from "@mui/material";
import DocumentList from "./DocumentList";
import type { Document } from "../types/document";
import { documentService } from "../services/documentService";

const DocumentManager: React.FC = () => {
  const [documents, setDocuments] = useState<Document[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const checkConnection = useCallback(async () => {
    const connected = await documentService.checkConnection();
    setIsConnected(connected);
    return connected;
  }, []);

  const fetchDocuments = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const connected = await checkConnection();
      if (connected) {
        const docs = await documentService.listDocuments();
        console.log("Fetched documents from service:", docs);
        console.log("Is docs an array?", Array.isArray(docs));

        // Ensure we only set an array
        if (Array.isArray(docs)) {
          setDocuments(docs);
        } else {
          console.error("DocumentService returned non-array:", docs);
          // If it's an object with documents property, extract it
          if (docs && typeof docs === "object" && "documents" in docs) {
            //eslint-disable-next-line @typescript-eslint/no-explicit-any
            const extractedDocs = (docs as any).documents;
            if (Array.isArray(extractedDocs)) {
              setDocuments(extractedDocs);
            } else {
              setDocuments([]);
            }
          } else {
            setDocuments([]);
          }
        }
      } else {
        setError("Cannot connect to backend server");
        setDocuments([]);
      }
    } catch (err) {
      console.error("Error fetching documents:", err);
      setError(
        err instanceof Error ? err.message : "Failed to fetch documents"
      );
      setDocuments([]);
    } finally {
      setLoading(false);
    }
  }, [checkConnection]);

  const handleRefresh = useCallback(async () => {
    await fetchDocuments();
  }, [fetchDocuments]);

  const handleDelete = useCallback(
    async (id: string) => {
      try {
        await documentService.deleteDocument(id);
        setDocuments((prev) => prev.filter((doc) => doc.id !== id));
      } catch (err) {
        console.error("Error deleting document:", err);
        throw err;
      }
    },
    //eslint-disable-next-line react-hooks/exhaustive-deps
    [documentService]
  );

  useEffect(() => {
    fetchDocuments();
    const interval = setInterval(checkConnection, 30000);
    return () => clearInterval(interval);
  }, [fetchDocuments, checkConnection]);

  if (loading && documents.length === 0) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography>Loading documents...</Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ height: "100%", display: "flex", flexDirection: "column" }}>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <DocumentList
        documents={documents}
        onRefresh={handleRefresh}
        onDelete={handleDelete}
        isConnected={isConnected}
      />
    </Box>
  );
};

export default DocumentManager;
