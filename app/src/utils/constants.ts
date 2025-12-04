export const APP_NAME = 'Construction Estimator';

export const API_TIMEOUT = 30000;

export const MAX_FILE_SIZE = 50 * 1024 * 1024; // 50MB

export const SUPPORTED_FILE_TYPES = {
  'application/pdf': ['.pdf'],
  'image/png': ['.png'],
  'image/jpeg': ['.jpg', '.jpeg'],
};

export const JOB_POLL_INTERVAL = 5000; // 5 seconds

export const COLORS = {
  primary: '#3B82F6',
  secondary: '#10B981',
  error: '#EF4444',
  warning: '#F59E0B',
  success: '#10B981',
  text: {
    primary: '#111827',
    secondary: '#6B7280',
    light: '#9CA3AF',
  },
  background: {
    primary: '#FFFFFF',
    secondary: '#F3F4F6',
  },
  border: '#E5E7EB',
};

export const STATUS_COLORS = {
  queued: '#F59E0B',
  processing: '#3B82F6',
  completed: '#10B981',
  failed: '#EF4444',
  pending: '#9CA3AF',
  uploaded: '#10B981',
  not_started: '#9CA3AF',
  active: '#10B981',
  archived: '#6B7280',
};
