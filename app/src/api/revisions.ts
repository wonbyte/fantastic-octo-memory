import apiClient from './client';
import {
  BlueprintRevision,
  BlueprintComparison,
  BidRevision,
  BidComparison,
} from '../types';

export const revisionsApi = {
  // Blueprint revision endpoints
  getBlueprintRevisions: async (blueprintId: string): Promise<BlueprintRevision[]> => {
    const response = await apiClient.get<BlueprintRevision[]>(
      `/blueprints/${blueprintId}/revisions`
    );
    return response.data;
  },

  createBlueprintRevision: async (blueprintId: string): Promise<BlueprintRevision> => {
    const response = await apiClient.post<BlueprintRevision>(
      `/blueprints/${blueprintId}/revisions`
    );
    return response.data;
  },

  compareBlueprintRevisions: async (
    blueprintId: string,
    fromVersion: number,
    toVersion: number
  ): Promise<BlueprintComparison> => {
    const response = await apiClient.get<BlueprintComparison>(
      `/blueprints/${blueprintId}/compare?from=${fromVersion}&to=${toVersion}`
    );
    return response.data;
  },

  // Bid revision endpoints
  getBidRevisions: async (bidId: string): Promise<BidRevision[]> => {
    const response = await apiClient.get<BidRevision[]>(`/bids/${bidId}/revisions`);
    return response.data;
  },

  createBidRevision: async (bidId: string): Promise<BidRevision> => {
    const response = await apiClient.post<BidRevision>(`/bids/${bidId}/revisions`);
    return response.data;
  },

  compareBidRevisions: async (
    bidId: string,
    fromVersion: number,
    toVersion: number
  ): Promise<BidComparison> => {
    const response = await apiClient.get<BidComparison>(
      `/bids/${bidId}/compare?from=${fromVersion}&to=${toVersion}`
    );
    return response.data;
  },
};
