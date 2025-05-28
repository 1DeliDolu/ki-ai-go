// frontend/src/components/DocumentUploader.tsx
import React, { useState, useRef } from 'react';
import { 
    Box, 
    Typography, 
    CircularProgress, 
    Paper
} from '@mui/material';
import { CloudUpload } from '@mui/icons-material';

interface DocumentUploaderProps {
    onUpload: (file: File) => Promise<void>;
    isConnected: boolean;
}

const DocumentUploader: React.FC<DocumentUploaderProps> = ({ onUpload, isConnected }) => {
    const [uploading, setUploading] = useState(false);
    const [dragActive, setDragActive] = useState(false);
    const inputRef = useRef<HTMLInputElement>(null);

    const handleDrag = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        if (e.type === "dragenter" || e.type === "dragover") {
            setDragActive(true);
        } else if (e.type === "dragleave") {
            setDragActive(false);
        }
    };

    const handleDrop = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setDragActive(false);

        if (!isConnected) return;

        if (e.dataTransfer.files && e.dataTransfer.files[0]) {
            handleFile(e.dataTransfer.files[0]);
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        e.preventDefault();
        if (!isConnected) return;
        
        if (e.target.files && e.target.files[0]) {
            handleFile(e.target.files[0]);
        }
    };

    const handleFile = async (file: File) => {
        const allowedTypes = ['.pdf', '.txt', '.docx', '.md'];
        const fileExt = '.' + file.name.split('.').pop()?.toLowerCase();

        if (!allowedTypes.includes(fileExt)) {
            alert('Unsupported file type. Please upload PDF, TXT, DOCX, or MD files.');
            return;
        }

        if (file.size > 10 * 1024 * 1024) { // 10MB limit
            alert('File too large. Maximum size is 10MB.');
            return;
        }

        setUploading(true);
        try {
            await onUpload(file);
        } catch (error) {
            console.error('Upload failed:', error);
        } finally {
            setUploading(false);
            if (inputRef.current) {
                inputRef.current.value = '';
            }
        }
    };

    const onButtonClick = () => {
        if (!isConnected) return;
        inputRef.current?.click();
    };

    return (
        <Box sx={{ 
            width: '100%', 
            height: '100%', 
            display: 'flex', 
            flexDirection: 'column' 
        }}>
            <Typography variant="h6" gutterBottom sx={{ flexShrink: 0 }}>
                Upload Document
            </Typography>
            {!isConnected && (
                <Box sx={{ mb: 2, p: 2, bgcolor: 'warning.light', borderRadius: 1, flexShrink: 0 }}>
                    <Typography color="warning.dark">
                        ⚠️ Not connected to backend. Document upload is disabled.
                    </Typography>
                </Box>
            )}
            <Paper
                elevation={0}
                variant="outlined"
                sx={{
                    p: 3,
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    border: '2px dashed',
                    borderColor: dragActive ? 'primary.main' : 'divider',
                    bgcolor: dragActive ? 'action.hover' : 'background.paper',
                    borderRadius: 2,
                    flex: 1,
                    minHeight: { xs: 200, sm: 250 },
                    cursor: isConnected ? 'pointer' : 'not-allowed',
                    opacity: isConnected ? 1 : 0.6,
                    transition: 'all 0.3s ease',
                    width: '100%',
                    '&:hover': {
                        borderColor: isConnected ? 'primary.main' : 'divider',
                        bgcolor: isConnected ? 'action.hover' : 'background.paper',
                    }
                }}
                onClick={onButtonClick}
                onDragEnter={isConnected ? handleDrag : undefined}
                onDragLeave={isConnected ? handleDrag : undefined}
                onDragOver={isConnected ? handleDrag : undefined}
                onDrop={isConnected ? handleDrop : undefined}
            >
                <input
                    ref={inputRef}
                    type="file"
                    accept=".pdf,.txt,.docx,.md"
                    onChange={handleChange}
                    style={{ display: 'none' }}
                />

                {uploading ? (
                    <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                        <CircularProgress size={40} />
                        <Typography sx={{ mt: 2 }}>Uploading document...</Typography>
                    </Box>
                ) : (
                    <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                        <CloudUpload fontSize="large" color={isConnected ? "primary" : "disabled"} sx={{ mb: 2 }} />
                        <Typography variant="subtitle1" align="center" gutterBottom color={isConnected ? "inherit" : "text.disabled"}>
                            {isConnected ? "Drag and drop a file here, or click to select" : "Connection required to upload documents"}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" align="center">
                            Supported: PDF, TXT, DOCX, MD (max 10MB)
                        </Typography>
                    </Box>
                )}
            </Paper>
        </Box>
    );
};

export default DocumentUploader;

