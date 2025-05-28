import React, { useState, useEffect, useCallback } from "react";
import {
  ThemeProvider,
  createTheme,
  CssBaseline,
  AppBar,
  Toolbar,
  Typography,
  Box,
  Tabs,
  Tab,
  Paper,
  Chip,
  useMediaQuery,
  Alert,
  Snackbar,
} from "@mui/material";
import {
  SmartToy,
  ModelTraining,
  Article,
  Chat as ChatIcon,
} from "@mui/icons-material";
import ModelManager from "./components/ModelManager";
import DocumentUploader from "./components/DocumentUploader";
import ChatInterface from "./components/ChatInterface";
import WikiResults from "./components/WikiResults";
import DocumentList from "./components/DocumentList";
import { useApiData } from "./hooks/useApiData";
import { useConnectionStatus } from "./hooks/useConnectionStatus";
import type { Model } from "./types/model";
import type { Document } from "./types/document";
import type { QueryResponse } from "./types/chatMessage";

// Custom hook for error handling
const useErrorHandler = () => {
  const [error, setError] = useState<string | null>(null);

  const handleError = useCallback((error: unknown, fallbackMessage: string) => {
    console.error("Error occurred:", error);

    if (error instanceof Error) {
      setError(error.message);
    } else {
      setError(fallbackMessage);
    }
  }, []);

  const clearError = useCallback(() => setError(null), []);

  return { error, handleError, clearError };
};

const App: React.FC = () => {
  const [activeTab, setActiveTab] = useState<number>(0);
  const [selectedModel, setSelectedModel] = useState("");
  const [response, setResponse] = useState<QueryResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const prefersDarkMode = useMediaQuery("(prefers-color-scheme: dark)");

  const { isConnected } = useConnectionStatus("http://localhost:8082");

  const { data: models, refresh: refreshModels } = useApiData<Model[]>(
    "models",
    () =>
      fetch(`http://localhost:8082/api/v1/models`).then((res) => res.json()),
    {
      enabled: isConnected,
      pollInterval: 60000, // 1 minute
      cacheTime: 30000, // 30 seconds cache
    }
  );

  const { data: documents, refresh: refreshDocuments } = useApiData<Document[]>(
    "documents",
    () =>
      fetch(`http://localhost:8082/api/v1/documents`).then((res) => res.json()),
    {
      enabled: isConnected,
      pollInterval: 60000, // 1 minute
      cacheTime: 30000, // 30 seconds cache
    }
  );

  // Custom hook for error handling
  const { error, clearError } = useErrorHandler();

  const theme = createTheme({
    palette: {
      mode: prefersDarkMode ? "dark" : "light",
      primary: {
        main: "#646cff",
      },
      secondary: {
        main: "#535bf2",
      },
      background: {
        default: prefersDarkMode ? "#242424" : "#f5f5f5",
        paper: prefersDarkMode ? "#1a1a1a" : "#ffffff",
      },
    },
    typography: {
      fontFamily: "system-ui, Avenir, Helvetica, Arial, sans-serif",
    },
    components: {
      MuiButton: {
        styleOverrides: {
          root: {
            borderRadius: 8,
            textTransform: "none",
          },
        },
      },
    },
  });

  // Initialize data on component mount
  useEffect(() => {
    const initialize = async () => {
      await Promise.allSettled([refreshModels(), refreshDocuments()]);
    };

    initialize();
  }, [refreshModels, refreshDocuments]);

  const handleQuery = useCallback(
    async (query: string, includeWiki: boolean = true) => {
      if (!isConnected) return;

      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8082/api/v1/query`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            query,
            include_wiki: includeWiki,
            include_documents: true,
            model_name: selectedModel || "default",
          }),
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(
            errorData.error ||
              `Query failed: ${response.status} ${response.statusText}`
          );
        }

        const data = await response.json();

        // Ensure response has proper structure
        const normalizedResponse: QueryResponse = {
          response: data.response || "",
          sources: {
            documents: data.sources?.documents || [],
            wiki: data.sources?.wiki || [],
          },
        };

        setResponse(normalizedResponse);
      } catch (error) {
        console.error("Query failed:", error);
        // Set error response
        setResponse({
          response:
            error instanceof Error
              ? error.message
              : "Query failed. Please try again.",
          sources: {
            documents: [],
            wiki: [],
          },
        });
      } finally {
        setLoading(false);
      }
    },
    [selectedModel, isConnected]
  );

  const handleModelDownload = useCallback(
    async (name: string, url: string) => {
      if (!isConnected) return;

      try {
        await fetch(`http://localhost:8082/api/v1/models/download`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ name, url }),
        });
        await refreshModels();
      } catch (error) {
        console.error("Model download failed:", error);
      }
    },
    [refreshModels, isConnected]
  );

  const handleModelLoad = useCallback(
    async (modelName: string) => {
      if (!isConnected) return;

      try {
        const response = await fetch(
          `http://localhost:8082/api/v1/models/load`,
          {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ name: modelName }),
          }
        );

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.error || "Failed to load model");
        }

        setSelectedModel(modelName);
        await refreshModels();
      } catch (error) {
        console.error("Model load failed:", error);
        // You could show an error message to the user here
      }
    },
    [refreshModels, isConnected]
  );

  const handleDocumentUpload = useCallback(
    async (file: File) => {
      if (!isConnected) return;

      const formData = new FormData();
      formData.append("file", file);

      try {
        await fetch(`http://localhost:8082/api/v1/documents/upload`, {
          method: "POST",
          body: formData,
        });
        await refreshDocuments();
      } catch (error) {
        console.error("Document upload failed:", error);
      }
    },
    [refreshDocuments, isConnected]
  );

  const handleTabChange = useCallback(
    (_: React.SyntheticEvent, newValue: number) => {
      setActiveTab(newValue);
    },
    []
  );

  const handleDeleteDocument = async (id: string) => {
    try {
      await fetch(`http://localhost:8082/api/v1/documents/${id}`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
      });
      await refreshDocuments();
    } catch (error) {
      console.error("Delete document error:", error);
    }
  };

  // Safe document count with null check
  const documentCount = documents?.length ?? 0;
  const modelCount = models?.length ?? 0;

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box
        sx={{
          flexGrow: 1,
          display: "flex",
          flexDirection: "column",
          minHeight: "100vh",
          width: "100vw",
          margin: 0,
          padding: 0,
        }}
      >
        <AppBar position="static" elevation={1} color="default">
          <Toolbar>
            <SmartToy sx={{ mr: 2 }} />
            <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
              Local AI Assistant
            </Typography>
            <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
              <Chip
                label={isConnected ? "Connected" : "Disconnected"}
                color={isConnected ? "success" : "error"}
                size="small"
                variant="outlined"
              />
              <Chip
                label={
                  selectedModel
                    ? `Model: ${selectedModel}`
                    : "No model selected"
                }
                color={selectedModel ? "primary" : "default"}
                size="small"
              />
              <Chip
                label={`Documents: ${documentCount}`}
                color="primary"
                size="small"
                variant="outlined"
              />
              <Chip
                label={`Models: ${modelCount}`}
                color="secondary"
                size="small"
                variant="outlined"
              />
            </Box>
          </Toolbar>
          <Tabs
            value={activeTab}
            onChange={handleTabChange}
            indicatorColor="primary"
            textColor="primary"
            centered
          >
            <Tab icon={<ChatIcon />} label="Chat" />
            <Tab icon={<ModelTraining />} label="Models" />
            <Tab icon={<Article />} label="Documents" />
          </Tabs>
        </AppBar>

        <Box sx={{ flexGrow: 1, width: "100%", overflow: "hidden" }}>
          {activeTab === 0 && (
            <Box
              sx={{
                display: "flex",
                flexDirection: { xs: "column", md: "row" },
                gap: 2,
                height: "calc(100vh - 200px)",
                p: 2,
              }}
            >
              <Paper
                elevation={2}
                sx={{
                  flex: 3,
                  p: 2,
                  borderRadius: 2,
                  display: "flex",
                  flexDirection: "column",
                }}
              >
                <ChatInterface
                  onQuery={handleQuery}
                  onModelLoad={handleModelLoad}
                  onDocumentUpload={handleDocumentUpload}
                  models={models || []}
                  loading={loading}
                  response={response}
                  selectedModel={selectedModel}
                  isConnected={isConnected}
                />
              </Paper>

              {response?.sources?.wiki && response.sources.wiki.length > 0 && (
                <Paper
                  elevation={2}
                  sx={{
                    flex: 2,
                    p: 2,
                    borderRadius: 2,
                    display: "flex",
                    flexDirection: "column",
                  }}
                >
                  <WikiResults
                    results={response.sources.wiki.map((source) => ({
                      pageId: source.url || "",
                      title: source.title || "",
                      url: source.url || "",
                      description: source.content,
                      extract: source.content,
                      relevanceScore: source.relevanceScore,
                    }))}
                    isLoading={false}
                  />
                </Paper>
              )}
            </Box>
          )}

          {activeTab === 1 && (
            <Box
              sx={{
                p: 2,
                height: "calc(100vh - 200px)",
                display: "flex",
                flexDirection: "column",
              }}
            >
              <Paper
                elevation={2}
                sx={{
                  p: 3,
                  borderRadius: 2,
                  height: "100%",
                  display: "flex",
                  flexDirection: "column",
                }}
              >
                <ModelManager
                  models={models || []}
                  selectedModel={selectedModel}
                  onDownload={handleModelDownload}
                  onLoad={handleModelLoad}
                  onRefresh={async () => {
                    await refreshModels();
                  }}
                  isConnected={isConnected}
                />
              </Paper>
            </Box>
          )}

          {activeTab === 2 && (
            <Box
              sx={{
                display: "flex",
                flexDirection: "column",
                gap: 2,
                height: "calc(100vh - 200px)",
                p: 2,
              }}
            >
              <Paper
                elevation={2}
                sx={{
                  p: 3,
                  borderRadius: 2,
                  flex: 1,
                  display: "flex",
                  flexDirection: "column",
                }}
              >
                <DocumentUploader
                  onUpload={handleDocumentUpload}
                  isConnected={isConnected}
                />
              </Paper>
              <Paper
                elevation={2}
                sx={{
                  p: 3,
                  borderRadius: 2,
                  flex: 1,
                  display: "flex",
                  flexDirection: "column",
                }}
              >
                <DocumentList
                  documents={documents || []}
                  onRefresh={async () => {
                    await refreshDocuments();
                  }}
                  onDelete={handleDeleteDocument}
                  isConnected={isConnected}
                />
              </Paper>
            </Box>
          )}
        </Box>

        {/* Error notification */}
        <Snackbar
          open={!!error}
          autoHideDuration={6000}
          onClose={clearError}
          anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
        >
          <Alert onClose={clearError} severity="error" sx={{ width: "100%" }}>
            {error}
          </Alert>
        </Snackbar>
      </Box>
    </ThemeProvider>
  );
};

export default App;
