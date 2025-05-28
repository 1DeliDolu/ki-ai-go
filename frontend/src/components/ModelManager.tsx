// frontend/src/components/ModelManager.tsx
import React, { useState } from "react";
import {
  Typography,
  Button,
  Box,
  Grid,
  Card,
  CardContent,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Chip,
  Divider,
} from "@mui/material";
import {
  Refresh,
  Check,
  Storage,
  ModelTraining,
  PlayArrow,
} from "@mui/icons-material";
import type { Model } from "../types/model";

interface ModelManagerProps {
  models: Model[] | { models: Model[] };
  selectedModel: string;
  onLoad: (modelName: string) => Promise<void>;
  onDownload: (name: string, url: string) => Promise<void>;
  onRefresh: () => Promise<void>;
  isConnected: boolean;
}

const ModelManager: React.FC<ModelManagerProps> = ({
  models = [],
  selectedModel,
  onLoad,
  onRefresh,
  isConnected,
}) => {
  const [selectedModelForLoading, setSelectedModelForLoading] =
    useState<string>("");

  const handleLoadSelectedModel = () => {
    if (selectedModelForLoading) {
      onLoad(selectedModelForLoading);
      setSelectedModelForLoading(""); // Clear selection after loading
    }
  };

  // Extract models from the prop structure
  const rawModels = Array.isArray(models)
    ? models
    : (models as { models: Model[] })?.models || [];
  const safeModels = Array.isArray(rawModels) ? rawModels : [];

  // Helper function to get model icon based on name
  const getModelIcon = (modelName: string): string => {
    const name = modelName.toLowerCase();
    if (name.includes("nemotron")) return "🚀";
    if (name.includes("neural-chat")) return "🧠";
    if (name.includes("openchat")) return "💬";
    if (name.includes("llama")) return "🦙";
    if (name.includes("phi")) return "🔬";
    if (name.includes("mistral")) return "⚡";
    return "🤖";
  };

  // Helper function to get model type
  const getModelType = (modelName: string): string => {
    const name = modelName.toLowerCase();
    if (name.includes("chat") || name.includes("openchat")) return "Chat";
    if (name.includes("vision")) return "Vision";
    return "General";
  };

  // All available models are now dynamic from backend
  const availableModels = safeModels.filter(
    (model) => model.status === "available"
  );

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
        <Typography
          variant="h5"
          component="h2"
          sx={{ display: "flex", alignItems: "center", flexGrow: 1 }}
        >
          <ModelTraining sx={{ mr: 1 }} /> AI Model Selection
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
            ⚠️ Not connected to backend. Model operations are disabled.
          </Typography>
        </Box>
      )}

      {/* Model Selection Section */}
      <Box
        sx={{
          mb: 3,
          p: 3,
          bgcolor: "background.paper",
          borderRadius: 2,
          border: 1,
          borderColor: "divider",
          flexShrink: 0,
        }}
      >
        <Typography
          variant="h6"
          gutterBottom
          sx={{ display: "flex", alignItems: "center" }}
        >
          <PlayArrow sx={{ mr: 1, color: "primary.main" }} />
          Select and Load Model
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          Choose one of your Ollama models to start conversations
        </Typography>

        <Box
          sx={{
            display: "flex",
            gap: 2,
            alignItems: "center",
            flexWrap: "wrap",
          }}
        >
          <FormControl sx={{ minWidth: 350, flexGrow: 1 }}>
            <InputLabel>Choose a Model</InputLabel>
            <Select
              value={selectedModelForLoading}
              onChange={(e) => setSelectedModelForLoading(e.target.value)}
              disabled={!isConnected || availableModels.length === 0}
              label="Choose a Model"
            >
              {availableModels.map((model) => (
                <MenuItem key={model.id} value={model.id}>
                  <Box
                    sx={{
                      display: "flex",
                      alignItems: "center",
                      width: "100%",
                    }}
                  >
                    <Typography sx={{ mr: 1, fontSize: "1.2em" }}>
                      {getModelIcon(model.name)}
                    </Typography>
                    <Box sx={{ flexGrow: 1 }}>
                      <Typography variant="body1" sx={{ fontWeight: "medium" }}>
                        {model.name}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        {model.size} • {getModelType(model.name)}
                      </Typography>
                    </Box>
                  </Box>
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <Button
            variant="contained"
            startIcon={<PlayArrow />}
            onClick={handleLoadSelectedModel}
            disabled={!selectedModelForLoading || !isConnected}
            size="large"
            sx={{ minWidth: 120 }}
          >
            Load Model
          </Button>
        </Box>

        {selectedModel && (
          <Box sx={{ mt: 2 }}>
            <Chip
              icon={<Check />}
              label={`Currently loaded: ${
                availableModels.find((m) => m.id === selectedModel)?.name ||
                selectedModel
              }`}
              color="primary"
              variant="filled"
              size="medium"
            />
          </Box>
        )}
      </Box>

      <Divider sx={{ my: 2, flexShrink: 0 }} />

      {/* Dynamic Models from Ollama */}
      <Box sx={{ flex: 1, overflow: "auto" }}>
        <Typography variant="h6" gutterBottom>
          Your Ollama Models ({availableModels.length} available)
        </Typography>

        <Grid container spacing={2}>
          {safeModels.map((model) => {
            const isSelected = selectedModel === model.id;
            const isAvailable = model.status === "available";

            return (
              <Grid item xs={12} sm={6} key={model.id}>
                <Card
                  sx={{
                    border: isSelected ? 2 : 1,
                    borderColor: isSelected ? "primary.main" : "divider",
                    opacity: isAvailable ? 1 : 0.6,
                    cursor: isAvailable && !isSelected ? "pointer" : "default",
                    transition: "all 0.2s",
                    "&:hover":
                      isAvailable && !isSelected
                        ? {
                            borderColor: "primary.light",
                            transform: "translateY(-2px)",
                            boxShadow: 2,
                          }
                        : {},
                  }}
                  variant="outlined"
                  onClick={() => {
                    if (isAvailable && !isSelected) {
                      setSelectedModelForLoading(model.id);
                    }
                  }}
                >
                  <CardContent>
                    <Box sx={{ display: "flex", alignItems: "flex-start" }}>
                      <Typography sx={{ mr: 2, mt: 0.5, fontSize: "2em" }}>
                        {getModelIcon(model.name)}
                      </Typography>
                      <Box sx={{ flexGrow: 1 }}>
                        <Typography variant="h6" component="div">
                          {model.name}
                        </Typography>
                        <Box
                          sx={{
                            display: "flex",
                            gap: 1,
                            flexWrap: "wrap",
                            mt: 1,
                          }}
                        >
                          <Chip
                            label={getModelType(model.name)}
                            size="small"
                            variant="outlined"
                            color="primary"
                          />
                          <Chip
                            label={model.size}
                            size="small"
                            variant="outlined"
                          />
                          <Chip
                            label={
                              model.status === "available"
                                ? "Ready"
                                : "Not Ready"
                            }
                            size="small"
                            color={
                              model.status === "available"
                                ? "success"
                                : "warning"
                            }
                          />
                          {isSelected && (
                            <Chip
                              label="Active"
                              size="small"
                              color="primary"
                              icon={<Check />}
                            />
                          )}
                        </Box>
                        <Typography
                          variant="caption"
                          color="text.secondary"
                          sx={{ display: "block", mt: 1 }}
                        >
                          Model ID: {model.id}
                        </Typography>
                      </Box>
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
            );
          })}
        </Grid>

        {availableModels.length === 0 && (
          <Box
            sx={{
              textAlign: "center",
              py: 4,
              bgcolor: "background.paper",
              borderRadius: 1,
              mt: 2,
            }}
          >
            <Storage sx={{ fontSize: 48, color: "text.secondary", mb: 2 }} />
            <Typography color="text.secondary" variant="h6">
              No models available
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
              {safeModels.length > 0
                ? "Models detected but not ready. Check Ollama status."
                : "No models detected. Ensure Ollama is running with models installed."}
            </Typography>
          </Box>
        )}
      </Box>
    </Box>
  );
};

export default ModelManager;
