import { useState, useEffect, useCallback } from 'react';
import { bidsApi } from '../api/bids';
import { Bid, GenerateBidRequest } from '../types';

export const useBids = (projectId: string) => {
  const [data, setData] = useState<Bid[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchBids = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const bids = await bidsApi.getProjectBids(projectId);
      setData(bids);
    } catch (err) {
      setError(err as Error);
    } finally {
      setIsLoading(false);
    }
  }, [projectId]);

  useEffect(() => {
    if (projectId) {
      fetchBids();
    }
  }, [projectId, fetchBids]);

  return {
    data,
    isLoading,
    error,
    refetch: fetchBids,
  };
};

export const useBid = (bidId: string) => {
  const [data, setData] = useState<Bid | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchBid = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const bid = await bidsApi.getBid(bidId);
      setData(bid);
    } catch (err) {
      setError(err as Error);
    } finally {
      setIsLoading(false);
    }
  }, [bidId]);

  useEffect(() => {
    if (bidId) {
      fetchBid();
    }
  }, [bidId, fetchBid]);

  return {
    data,
    isLoading,
    error,
    refetch: fetchBid,
  };
};

export const useGenerateBid = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const generateBid = async (projectId: string, request: GenerateBidRequest): Promise<Bid | null> => {
    try {
      setIsLoading(true);
      setError(null);
      const bid = await bidsApi.generateBid(projectId, request);
      return bid;
    } catch (err) {
      setError(err as Error);
      return null;
    } finally {
      setIsLoading(false);
    }
  };

  return {
    generateBid,
    isLoading,
    error,
  };
};
