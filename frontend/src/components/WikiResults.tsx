import React from 'react';
import {
    Box,
    Typography,
    Card,
    CardContent,
    CardMedia,
    Link,
    CircularProgress,
    Chip
} from '@mui/material';
import { Language as LanguageIcon } from '@mui/icons-material';
import type { WikiResult } from '../types/wikiResult';

interface WikiResultsProps {
    results: WikiResult[];
    isLoading?: boolean;
}

const WikiResults: React.FC<WikiResultsProps> = ({ results, isLoading = false }) => {
    if (isLoading) {
        return (
            <Box sx={{ p: 2, textAlign: 'center', height: '100%', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
                <CircularProgress size={30} />
                <Typography sx={{ mt: 2 }}>Searching Wikipedia...</Typography>
            </Box>
        );
    }

    if (!results || results.length === 0) {
        return null;
    }

    return (
        <Box sx={{ width: '100%', height: '100%', display: 'flex', flexDirection: 'column' }}>
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2, flexShrink: 0 }}>
                <LanguageIcon color="primary" sx={{ mr: 1 }} />
                <Typography variant="h6">
                    Wikipedia Results ({results.length})
                </Typography>
            </Box>

            <Box sx={{ 
                display: 'flex', 
                flexDirection: 'column', 
                gap: 2, 
                flex: 1, 
                overflow: 'auto',
                width: '100%'
            }}>
                {results.map((result, index) => (
                    <Card key={result.pageId || index} variant="outlined" sx={{ 
                        display: 'flex', 
                        flexDirection: { xs: 'column', sm: 'row' },
                        width: '100%',
                        flexShrink: 0
                    }}>
                        {result.thumbnail && (
                            <CardMedia
                                component="img"
                                image={result.thumbnail}
                                alt={result.title}
                                sx={{ 
                                    height: { xs: 200, sm: 140 },
                                    width: { xs: '100%', sm: 200 },
                                    flexShrink: 0
                                }}
                                onError={(e) => {
                                    e.currentTarget.style.display = 'none';
                                }}
                            />
                        )}
                        <CardContent sx={{ flex: 1, width: '100%' }}>
                            <Typography variant="h6" component="div" gutterBottom>
                                <Link href={result.url} target="_blank" rel="noopener noreferrer" underline="hover">
                                    {result.title}
                                </Link>
                            </Typography>
                            
                            {result.description && (
                                <Typography variant="body2" color="text.secondary" gutterBottom>
                                    {result.description}
                                </Typography>
                            )}
                            
                            {result.extract && (
                                <Typography variant="body2" paragraph>
                                    {result.extract}
                                </Typography>
                            )}
                            
                            {result.relevanceScore && (
                                <Chip 
                                    label={`Relevance: ${(result.relevanceScore * 100).toFixed(1)}%`} 
                                    size="small" 
                                    color="primary" 
                                    variant="outlined"
                                />
                            )}
                        </CardContent>
                    </Card>
                ))}
            </Box>
        </Box>
    );
};

export default WikiResults;