import apiClient from './client';
import { Job } from '../types';

export const jobsApi = {
  getById: async (id: string): Promise<Job> => {
    const response = await apiClient.get<Job>(`/jobs/${id}`);
    return response.data;
  },

  getByBlueprintId: async (blueprintId: string): Promise<Job[]> => {
    const response = await apiClient.get<Job[]>(`/blueprints/${blueprintId}/jobs`);
    return response.data;
  },
};
