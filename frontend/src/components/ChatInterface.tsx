// frontend/src/components/ChatInterface.tsx
import React, { useState, useRef, useEffect } from "react";
import {
    Box,
    TextField,
    Button,
    Typography,
    Paper,
    CircularProgress,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    IconButton,
    Chip,
    Alert,
} from "@mui/material";
import {
    Send,
    AttachFile,
    Close,
    PlayArrow,
    ModelTraining,
} from "@mui/icons-material";
import type { QueryResponse } from "../types/chatMessage";
import type { Model } from "../types/model";

interface ChatMessage {
    id: string;
    type: "user" | "assistant";
    content: string;
    timestamp: Date;
    sources?: {
        documents: Array<{ title?: string }>;
        wiki: Array<{ title?: string; url?: string }>;
    };
}

interface ChatInterfaceProps {
    onQuery: (query: string, includeWiki?: boolean) => Promise<void>;
    onModelLoad?: (modelName: string) => Promise<void>;
    onDocumentUpload?: (file: File) => Promise<void>;
    models?: Model[] | { models: Model[] };
    loading: boolean;
    response: QueryResponse | null;
    selectedModel: string;
    isConnected: boolean;
}

const ChatInterface: React.FC<ChatInterfaceProps> = ({
    onQuery,
    onModelLoad,
    onDocumentUpload,
    models = [],
    loading,
    response,
    selectedModel,
    isConnected,
}) => {
    const [query, setQuery] = useState("");
    const [selectedModelForChat, setSelectedModelForChat] =
        useState(selectedModel);
    const [attachedFiles, setAttachedFiles] = useState<File[]>([]);
    const [chatHistory, setChatHistory] = useState<ChatMessage[]>([]);
    const fileInputRef = useRef<HTMLInputElement>(null);
    const messagesEndRef = useRef<HTMLDivElement>(null);

    // Auto-scroll to bottom when new messages are added
    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    };

    useEffect(() => {
        scrollToBottom();
    }, [chatHistory]);

    // Update chat history when response changes
    useEffect(() => {
        if (response && response.response) {
            // Check if this response is already in history (to avoid duplicates)
            const lastMessage = chatHistory[chatHistory.length - 1];
            if (
                !lastMessage ||
                lastMessage.type !== "assistant" ||
                lastMessage.content !== response.response
            ) {
                const assistantMessage: ChatMessage = {
                    id: Date.now().toString() + "_assistant",
                    type: "assistant",
                    content: response.response,
                    timestamp: new Date(),
                    sources: response.sources,
                };
                setChatHistory((prev) => [...prev, assistantMessage]);
            }
        }
    }, [response, chatHistory]);

    // Extract models from the prop structure
    const rawModels = Array.isArray(models)
        ? models
        : (models as { models: Model[] })?.models || [];
    const safeModels = Array.isArray(rawModels) ? rawModels : [];
    const availableModels = safeModels.filter(
        (model) => model.status === "available"
    );

    // Helper function to get model icon
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

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!query.trim() || loading || !isConnected) return;

        // Add user message to history
        const userMessage: ChatMessage = {
            id: Date.now().toString() + "_user",
            type: "user",
            content: query,
            timestamp: new Date(),
        };
        setChatHistory((prev) => [...prev, userMessage]);

        // Upload attached files first
        if (attachedFiles.length > 0 && onDocumentUpload) {
            for (const file of attachedFiles) {
                try {
                    await onDocumentUpload(file);
                } catch (error) {
                    console.error("Failed to upload file:", file.name, error);
                }
            }
            setAttachedFiles([]); // Clear files after upload
        }

        // Load model if different from current
        if (selectedModelForChat !== selectedModel && onModelLoad) {
            try {
                await onModelLoad(selectedModelForChat);
            } catch (error) {
                console.error("Failed to load model:", error);
            }
        }

        const currentQuery = query;
        setQuery("");
        await onQuery(currentQuery);
    };

    const clearHistory = () => {
        setChatHistory([]);
    };

    const handleFileSelect = () => {
        fileInputRef.current?.click();
    };

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const files = Array.from(e.target.files || []);
        if (files.length > 0) {
            setAttachedFiles((prev) => [...prev, ...files]);
        }
        // Reset input
        if (fileInputRef.current) {
            fileInputRef.current.value = "";
        }
    };

    const removeAttachedFile = (index: number) => {
        setAttachedFiles((prev) => prev.filter((_, i) => i !== index));
    };

    const formatFileSize = (bytes: number): string => {
        const sizes = ["Bytes", "KB", "MB", "GB"];
        if (bytes === 0) return "0 Bytes";
        const i = Math.floor(Math.log(bytes) / Math.log(1024));
        return Math.round((bytes / Math.pow(1024, i)) * 100) / 100 + " " + sizes[i];
    };

    const formatTime = (date: Date) => {
        return date.toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    return (
        <Box
            sx={{
                height: "100%",
                display: "flex",
                flexDirection: "column",
                gap: 2,
            }}
        >
            {/* Header with model selection */}
            <Box
                sx={{
                    display: "flex",
                    alignItems: "center",
                    gap: 2,
                    p: 2,
                    bgcolor: "background.paper",
                    borderRadius: 1,
                    border: 1,
                    borderColor: "divider",
                }}
            >
                <ModelTraining color="primary" />
                <FormControl size="small" sx={{ minWidth: 200, flexGrow: 1 }}>
                    <InputLabel>Select Model</InputLabel>
                    <Select
                        value={selectedModelForChat}
                        onChange={(e) => setSelectedModelForChat(e.target.value)}
                        disabled={!isConnected || availableModels.length === 0}
                        label="Select Model"
                    >
                        {availableModels.map((model) => (
                            <MenuItem key={model.id} value={model.id}>
                                <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                                    <Typography sx={{ fontSize: "1.1em" }}>
                                        {getModelIcon(model.name)}
                                    </Typography>
                                    <Box>
                                        <Typography variant="body2" sx={{ fontWeight: "medium" }}>
                                            {model.name}
                                        </Typography>
                                        <Typography variant="caption" color="text.secondary">
                                            {model.size}
                                        </Typography>
                                    </Box>
                                </Box>
                            </MenuItem>
                        ))}
                    </Select>
                </FormControl>

                {selectedModel && (
                    <Chip
                        icon={<PlayArrow />}
                        label={`Active: ${availableModels.find((m) => m.id === selectedModel)?.name ||
                            selectedModel
                            }`}
                        color="primary"
                        size="small"
                    />
                )}

                {chatHistory.length > 0 && (
                    <Button
                        size="small"
                        onClick={clearHistory}
                        variant="outlined"
                        color="secondary"
                    >
                        Clear History
                    </Button>
                )}
            </Box>

            {!isConnected && (
                <Alert severity="warning">
                    Not connected to backend. Chat functionality is disabled.
                </Alert>
            )}

            {/* Chat Messages Area */}
            <Paper
                sx={{
                    flex: 1,
                    display: "flex",
                    flexDirection: "column",
                    overflow: "hidden",
                    bgcolor: "background.default",
                }}
            >
                <Box
                    sx={{
                        flex: 1,
                        overflow: "auto",
                        p: 2,
                        display: "flex",
                        flexDirection: "column",
                        gap: 2,
                    }}
                >
                    {chatHistory.length === 0 ? (
                        <Box
                            sx={{
                                display: "flex",
                                alignItems: "center",
                                justifyContent: "center",
                                height: "100%",
                                color: "text.secondary",
                                textAlign: "center",
                            }}
                        >
                            <Typography>
                                Start a conversation by typing your question below...
                            </Typography>
                        </Box>
                    ) : (
                        chatHistory.map((message) => (
                            <Box
                                key={message.id}
                                sx={{
                                    display: "flex",
                                    flexDirection:
                                        message.type === "user" ? "row-reverse" : "row",
                                    gap: 1,
                                }}
                            >
                                <Box
                                    sx={{
                                        maxWidth: "70%",
                                        p: 2,
                                        borderRadius: 2,
                                        bgcolor:
                                            message.type === "user"
                                                ? "primary.main"
                                                : "background.paper",
                                        color:
                                            message.type === "user"
                                                ? "primary.contrastText"
                                                : "text.primary",
                                        border: message.type === "assistant" ? 1 : 0,
                                        borderColor: "divider",
                                    }}
                                >
                                    <Typography
                                        variant="body1"
                                        sx={{ whiteSpace: "pre-wrap", mb: 1 }}
                                    >
                                        {message.content}
                                    </Typography>

                                    {message.sources && (
                                        <Box sx={{ mt: 2 }}>
                                            {message.sources.documents &&
                                                message.sources.documents.length > 0 && (
                                                    <Box sx={{ mb: 1 }}>
                                                        <Typography
                                                            variant="caption"
                                                            color="primary"
                                                            display="block"
                                                        >
                                                            📄 Document Sources (
                                                            {message.sources.documents.length})
                                                        </Typography>
                                                        <Box
                                                            sx={{
                                                                display: "flex",
                                                                flexWrap: "wrap",
                                                                gap: 0.5,
                                                                mt: 0.5,
                                                            }}
                                                        >
                                                            {message.sources.documents.map((doc, index) => (
                                                                <Chip
                                                                    key={index}
                                                                    label={doc.title || `Document ${index + 1}`}
                                                                    size="small"
                                                                    variant="outlined"
                                                                    sx={{ fontSize: "0.7rem", height: "auto" }}
                                                                />
                                                            ))}
                                                        </Box>
                                                    </Box>
                                                )}

                                            {message.sources.wiki &&
                                                message.sources.wiki.length > 0 && (
                                                    <Box>
                                                        <Typography
                                                            variant="caption"
                                                            color="secondary"
                                                            display="block"
                                                        >
                                                            🌐 Wiki Sources ({message.sources.wiki.length})
                                                        </Typography>
                                                        <Box
                                                            sx={{
                                                                display: "flex",
                                                                flexWrap: "wrap",
                                                                gap: 0.5,
                                                                mt: 0.5,
                                                            }}
                                                        >
                                                            {message.sources.wiki.map((wiki, index) => (
                                                                <Chip
                                                                    key={index}
                                                                    label={wiki.title || `Wiki ${index + 1}`}
                                                                    size="small"
                                                                    variant="outlined"
                                                                    color="secondary"
                                                                    sx={{ fontSize: "0.7rem", height: "auto" }}
                                                                />
                                                            ))}
                                                        </Box>
                                                    </Box>
                                                )}
                                        </Box>
                                    )}

                                    <Typography
                                        variant="caption"
                                        sx={{
                                            display: "block",
                                            mt: 1,
                                            opacity: 0.7,
                                            textAlign: message.type === "user" ? "right" : "left",
                                        }}
                                    >
                                        {formatTime(message.timestamp)}
                                    </Typography>
                                </Box>
                            </Box>
                        ))
                    )}

                    {loading && (
                        <Box
                            sx={{
                                display: "flex",
                                alignItems: "center",
                                gap: 1,
                                p: 2,
                            }}
                        >
                            <CircularProgress size={20} />
                            <Typography color="text.secondary">AI is thinking...</Typography>
                        </Box>
                    )}

                    <div ref={messagesEndRef} />
                </Box>
            </Paper>

            {/* Attached files display */}
            {attachedFiles.length > 0 && (
                <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
                    {attachedFiles.map((file, index) => (
                        <Chip
                            key={index}
                            label={`${file.name} (${formatFileSize(file.size)})`}
                            onDelete={() => removeAttachedFile(index)}
                            deleteIcon={<Close />}
                            color="primary"
                            variant="outlined"
                        />
                    ))}
                </Box>
            )}

            {/* Input area */}
            <Box component="form" onSubmit={handleSubmit}>
                <Box sx={{ display: "flex", gap: 1, alignItems: "flex-end" }}>
                    <input
                        type="file"
                        ref={fileInputRef}
                        onChange={handleFileChange}
                        multiple
                        accept=".pdf,.txt,.docx,.doc,.md"
                        style={{ display: "none" }}
                    />

                    <IconButton
                        onClick={handleFileSelect}
                        disabled={!isConnected}
                        color="primary"
                        sx={{ mb: 1 }}
                    >
                        <AttachFile />
                    </IconButton>

                    <TextField
                        fullWidth
                        multiline
                        maxRows={4}
                        value={query}
                        onChange={(e) => setQuery(e.target.value)}
                        placeholder={
                            isConnected
                                ? "Ask a question about your documents..."
                                : "Connect to backend to start chatting..."
                        }
                        disabled={loading || !isConnected}
                        variant="outlined"
                        onKeyDown={(e) => {
                            if (e.key === "Enter" && !e.shiftKey) {
                                e.preventDefault();
                                handleSubmit(e);
                            }
                        }}
                    />

                    <Button
                        type="submit"
                        variant="contained"
                        disabled={!query.trim() || loading || !isConnected}
                        startIcon={loading ? <CircularProgress size={16} /> : <Send />}
                        sx={{ mb: 1, minWidth: 100 }}
                    >
                        {loading ? "Sending..." : "Send"}
                    </Button>
                </Box>

                <Typography
                    variant="caption"
                    color="text.secondary"
                    sx={{ mt: 1, display: "block" }}
                >
                    {attachedFiles.length > 0 &&
                        `${attachedFiles.length} file(s) will be uploaded. `}
                    Press Enter to send, Shift+Enter for new line. Your conversation
                    history is preserved.
                </Typography>
            </Box>
        </Box>
    );
};

export default ChatInterface;
