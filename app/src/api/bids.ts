import apiClient from './client';
import { Bid, GenerateBidRequest, PricingSummary } from '../types';

export const bidsApi = {
  getProjectBids: async (projectId: string): Promise<Bid[]> => {
    const response = await apiClient.get<Bid[]>(`/projects/${projectId}/bids`);
    return response.data;
  },

  getBid: async (bidId: string): Promise<Bid> => {
    const response = await apiClient.get<Bid>(`/bids/${bidId}`);
    return response.data;
  },

  generateBid: async (projectId: string, data: GenerateBidRequest): Promise<Bid> => {
    const response = await apiClient.post<Bid>(`/projects/${projectId}/generate-bid`, data);
    return response.data;
  },

  getBidPDF: async (bidId: string): Promise<{ pdf_url: string }> => {
    const response = await apiClient.get<{ pdf_url: string }>(`/bids/${bidId}/pdf`);
    return response.data;
  },

  getPricingSummary: async (projectId: string, blueprintId: string): Promise<PricingSummary> => {
    const response = await apiClient.get<PricingSummary>(
      `/projects/${projectId}/pricing-summary?blueprint_id=${blueprintId}`
    );
    return response.data;
  },
};
