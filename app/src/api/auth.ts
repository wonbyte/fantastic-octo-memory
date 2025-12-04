import apiClient, { setAuthToken, clearAuthToken } from './client';
import { LoginRequest, LoginResponse } from '../types';

export const authApi = {
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post<LoginResponse>('/auth/login', credentials);
    if (response.data.token) {
      await setAuthToken(response.data.token);
    }
    return response.data;
  },

  logout: async (): Promise<void> => {
    await clearAuthToken();
  },

  register: async (data: LoginRequest & { name?: string }): Promise<LoginResponse> => {
    const response = await apiClient.post<LoginResponse>('/auth/register', data);
    if (response.data.token) {
      await setAuthToken(response.data.token);
    }
    return response.data;
  },
};
