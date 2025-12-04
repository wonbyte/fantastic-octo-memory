import { useQuery } from '@tanstack/react-query';
import { jobsApi } from '../api/jobs';
import { JOB_POLL_INTERVAL } from '../utils/constants';

export const useJob = (jobId: string, enabled = true) => {
  return useQuery({
    queryKey: ['jobs', jobId],
    queryFn: () => jobsApi.getById(jobId),
    enabled: enabled && !!jobId,
    refetchInterval: (query) => {
      const data = query.state.data;
      // Stop polling if no data or job is completed/failed
      if (!data || data.status === 'completed' || data.status === 'failed') {
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
