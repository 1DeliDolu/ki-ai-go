export interface Model {
  id: string;
  name: string;
  size: string;
  status: 'available' | 'downloading' | 'loading' | 'loaded' | 'error';
  downloadProgress?: number;
  description?: string;
  modelType: 'chat' | 'embedding' | 'multimodal';
}

export interface ModelManagerProps {
  models: Model[];
  selectedModel: string;
  onDownload: (name: string, url: string) => Promise<void>;
  onLoad: (modelName: string) => Promise<void>;
  onRefresh: () => Promise<void>;
  isConnected: boolean;
}
