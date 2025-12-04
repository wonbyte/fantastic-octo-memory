import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { blueprintsApi } from '../api/blueprints';
import { UploadUrlRequest } from '../types';

export const useBlueprints = (projectId: string) => {
  return useQuery({
    queryKey: ['blueprints', projectId],
    queryFn: () => blueprintsApi.getByProjectId(projectId),
    enabled: !!projectId,
  });
};

export const useBlueprint = (id: string) => {
  return useQuery({
    queryKey: ['blueprints', 'detail', id],
    queryFn: () => blueprintsApi.getById(id),
    enabled: !!id,
  });
};

export const useRequestUploadUrl = () => {
  return useMutation({
    mutationFn: ({ projectId, data }: { projectId: string; data: UploadUrlRequest }) =>
      blueprintsApi.requestUploadUrl(projectId, data),
  });
};

export const useCompleteUpload = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ blueprintId, success, error_message }: {
      blueprintId: string;
      success: boolean;
      error_message?: string;
    }) => blueprintsApi.completeUpload(blueprintId, { success, error_message }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['blueprints'] });
    },
  });
};

export const useTriggerAnalysis = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (blueprintId: string) => blueprintsApi.triggerAnalysis(blueprintId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['blueprints'] });
    },
  });
};
