import apiClient from './client';
import {
  Blueprint,
  UploadUrlRequest,
  UploadUrlResponse,
  CompleteUploadRequest,
  TriggerAnalysisResponse,
  AnalysisResult,
  TakeoffSummary,
} from '../types';
import axios from 'axios';

export const blueprintsApi = {
  getByProjectId: async (projectId: string): Promise<Blueprint[]> => {
    const response = await apiClient.get<Blueprint[]>(`/projects/${projectId}/blueprints`);
    return response.data;
  },

  getById: async (id: string): Promise<Blueprint> => {
    const response = await apiClient.get<Blueprint>(`/blueprints/${id}`);
    return response.data;
  },

  requestUploadUrl: async (
    projectId: string,
    data: UploadUrlRequest
  ): Promise<UploadUrlResponse> => {
    const response = await apiClient.post<UploadUrlResponse>(
      `/projects/${projectId}/blueprints/upload-url`,
      data
    );
    return response.data;
  },

  uploadToS3: async (uploadUrl: string, file: Blob | File): Promise<void> => {
    await axios.put(uploadUrl, file, {
      headers: {
        'Content-Type': file.type,
      },
    });
  },

  completeUpload: async (blueprintId: string, data: CompleteUploadRequest): Promise<void> => {
    await apiClient.post(`/blueprints/${blueprintId}/complete-upload`, data);
  },

  triggerAnalysis: async (blueprintId: string): Promise<TriggerAnalysisResponse> => {
    const response = await apiClient.post<TriggerAnalysisResponse>(
      `/blueprints/${blueprintId}/analyze`
    );
    return response.data;
  },

  getAnalysis: async (blueprintId: string): Promise<AnalysisResult> => {
    const response = await apiClient.get<AnalysisResult>(
      `/blueprints/${blueprintId}/analysis`
    );
    return response.data;
  },

  getTakeoffSummary: async (blueprintId: string): Promise<TakeoffSummary> => {
    const response = await apiClient.get<TakeoffSummary>(
      `/blueprints/${blueprintId}/takeoff-summary`
    );
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/blueprints/${id}`);
  },
};
