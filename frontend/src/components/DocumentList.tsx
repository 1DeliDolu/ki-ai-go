// frontend/src/components/DocumentList.tsx
import React, { useState, useMemo } from "react";
import {
  Box,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  Chip,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  CircularProgress,
} from "@mui/material";
import {
  Refresh,
  PictureAsPdf,
  Description,
  TextSnippet,
  ArticleOutlined,
  Delete,
} from "@mui/icons-material";
import type { Document } from "../types/document";

interface DocumentListProps {
  documents: Document[] | { documents: Document[] };
  onRefresh: () => Promise<void>;
  onDelete: (id: string) => Promise<void>;
  isConnected: boolean;
}

const DocumentList: React.FC<DocumentListProps> = ({
  documents = [],
  onRefresh,
  onDelete,
  isConnected,
}) => {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [documentToDelete, setDocumentToDelete] = useState<Document | null>(
    null
  );
  const [deleting, setDeleting] = useState(false);

  // Extract documents from the prop structure - it comes as {documents: [...]} or Document[]
  const safeDocuments = useMemo(() => {
    const rawDocuments = Array.isArray(documents)
      ? documents
      : (documents as { documents: Document[] })?.documents || [];
    return Array.isArray(rawDocuments) ? rawDocuments : [];
  }, [documents]);

  // Enhanced debug logging
  React.useEffect(() => {

    if (!Array.isArray(documents) && !documents?.documents) {
      console.error("DocumentList received invalid documents:", documents);
    }
  }, [documents, safeDocuments]);

  const formatFileSize = (bytes: number): string => {
    const sizes = ["Bytes", "KB", "MB", "GB"];
    if (bytes === 0) return "0 Bytes";
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return Math.round((bytes / Math.pow(1024, i)) * 100) / 100 + " " + sizes[i];
  };

  const formatDate = (dateString: string): string => {
    try {
      return new Date(dateString).toLocaleDateString();
    } catch {
      console.warn("Invalid date string:", dateString);
      return "Invalid date";
    }
  };

  const getFileIcon = (type: string) => {
    // Handle both with and without dots
    const cleanType = type.replace(".", "").toLowerCase();
    switch (cleanType) {
      case "pdf":
        return <PictureAsPdf color="error" />;
      case "txt":
        return <Description color="primary" />;
      case "docx":
      case "doc":
        return <TextSnippet color="info" />;
      case "md":
        return <ArticleOutlined color="success" />;
      default:
        return <Description />;
    }
  };

  const getStatusColor = (
    status: string
  ): "success" | "warning" | "error" | "default" => {
    switch (status) {
      case "ready":
        return "success";
      case "processing":
        return "warning";
      case "error":
        return "error";
      default:
        return "default";
    }
  };

  const handleDeleteClick = (doc: Document) => {
    setDocumentToDelete(doc);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = async () => {
    if (!documentToDelete) return;

    setDeleting(true);
    try {
      await onDelete(documentToDelete.id);
      setDeleteDialogOpen(false);
      setDocumentToDelete(null);
      await onRefresh(); // Refresh the list after deletion
    } catch (error) {
      console.error("Delete failed:", error);
    } finally {
      setDeleting(false);
    }
  };

  const handleDeleteCancel = () => {
    setDeleteDialogOpen(false);
    setDocumentToDelete(null);
  };

  return (
    <Box
      sx={{
        width: "100%",
        height: "100%",
        display: "flex",
        flexDirection: "column",
        overflow: "hidden",
      }}
    >
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
          flexWrap: "wrap",
          gap: 1,
          flexShrink: 0,
        }}
      >
        <Typography variant="h6">
          Document List ({safeDocuments.length})
        </Typography>
        <Button
          startIcon={<Refresh />}
          variant="outlined"
          onClick={onRefresh}
          disabled={!isConnected}
        >
          Refresh
        </Button>
      </Box>

      {!isConnected && (
        <Box
          sx={{
            mb: 2,
            p: 2,
            bgcolor: "warning.light",
            borderRadius: 1,
            flexShrink: 0,
          }}
        >
          <Typography color="warning.dark">
            ⚠️ Not connected to backend. Document list may be outdated.
          </Typography>
        </Box>
      )}

      <Box sx={{ flex: 1, overflow: "auto", width: "100%" }}>
        {safeDocuments.length === 0 ? (
          <Box
            sx={{
              textAlign: "center",
              py: 4,
              bgcolor: "background.paper",
              borderRadius: 1,
              height: "100%",
              display: "flex",
              flexDirection: "column",
              justifyContent: "center",
              alignItems: "center",
            }}
          >
            <Typography color="text.secondary">
              No documents uploaded yet
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
              Upload some documents to get started!
            </Typography>
          </Box>
        ) : (
          <Grid container spacing={2} sx={{ width: "100%", margin: 0 }}>
            {safeDocuments.map((doc) => (
              <Grid
                item
                xs={12}
                sm={6}
                lg={4}
                key={doc.id}
                sx={{ padding: "8px !important" }}
              >
                <Card variant="outlined" sx={{ height: "100%", width: "100%" }}>
                  <CardContent>
                    <Box
                      sx={{
                        display: "flex",
                        alignItems: "flex-start",
                        justifyContent: "space-between",
                      }}
                    >
                      <Box
                        sx={{
                          display: "flex",
                          alignItems: "flex-start",
                          flexGrow: 1,
                        }}
                      >
                        <Box sx={{ mr: 1 }}>{getFileIcon(doc.type)}</Box>
                        <Box sx={{ flexGrow: 1, minWidth: 0 }}>
                          <Typography
                            variant="subtitle1"
                            noWrap
                            title={doc.name}
                          >
                            {doc.name}
                          </Typography>
                          <Box
                            sx={{
                              display: "flex",
                              flexWrap: "wrap",
                              gap: 1,
                              mt: 1,
                            }}
                          >
                            <Chip
                              label={doc.type.toUpperCase()}
                              size="small"
                              variant="outlined"
                            />
                            <Chip
                              label={formatFileSize(doc.size)}
                              size="small"
                              variant="outlined"
                            />
                            <Chip
                              label={doc.status}
                              size="small"
                              color={getStatusColor(doc.status)}
                            />
                          </Box>
                          <Typography
                            variant="caption"
                            color="text.secondary"
                            display="block"
                            sx={{ mt: 1 }}
                          >
                            Uploaded: {formatDate(doc.uploadDate)}
                          </Typography>
                        </Box>
                      </Box>
                      <IconButton
                        size="small"
                        color="error"
                        onClick={() => handleDeleteClick(doc)}
                        disabled={!isConnected}
                        sx={{ ml: 1 }}
                      >
                        <Delete fontSize="small" />
                      </IconButton>
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        )}
      </Box>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onClose={handleDeleteCancel}>
        <DialogTitle>Delete Document</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete "{documentToDelete?.name}"? This
            action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleDeleteCancel} disabled={deleting}>
            Cancel
          </Button>
          <Button
            onClick={handleDeleteConfirm}
            color="error"
            variant="contained"
            disabled={deleting}
            startIcon={deleting ? <CircularProgress size={16} /> : <Delete />}
          >
            {deleting ? "Deleting..." : "Delete"}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default DocumentList;
