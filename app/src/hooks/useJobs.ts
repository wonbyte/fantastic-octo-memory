import { useQuery } from '@tanstack/react-query';
import { jobsApi } from '../api/jobs';
import { JOB_POLL_INTERVAL } from '../utils/constants';

export const useJob = (jobId: string, enabled = true) => {
  return useQuery({
    queryKey: ['jobs', jobId],
    queryFn: () => jobsApi.getById(jobId),
    enabled: enabled && !!jobId,
    refetchInterval: (query) => {
      const status = query.state.data?.status;
      // Stop polling when job is completed or failed
      if (status === 'completed' || status === 'failed') {
        return false;
      }
      return JOB_POLL_INTERVAL;
    },
  });
};

export const useJobsByBlueprint = (blueprintId: string) => {
  return useQuery({
    queryKey: ['jobs', 'blueprint', blueprintId],
    queryFn: () => jobsApi.getByBlueprintId(blueprintId),
    enabled: !!blueprintId,
  });
};
