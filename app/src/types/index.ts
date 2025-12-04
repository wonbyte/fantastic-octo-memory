// API Response Types
export interface ApiResponse<T = unknown> {
  data?: T;
  error?: string;
  message?: string;
}

// User Types
export interface User {
  id: string;
  email: string;
  name?: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

// Project Types
export type ProjectStatus = 'active' | 'completed' | 'archived';

export interface Project {
  id: string;
  name: string;
  description?: string;
  status: ProjectStatus;
  created_at: string;
  updated_at: string;
  user_id: string;
}

export interface CreateProjectRequest {
  name: string;
  description?: string;
}

// Blueprint Types
export type BlueprintUploadStatus = 'pending' | 'uploaded' | 'failed';
export type BlueprintAnalysisStatus = 'not_started' | 'queued' | 'processing' | 'completed' | 'failed';

export interface Blueprint {
  id: string;
  project_id: string;
  filename: string;
  file_size: number;
  content_type: string;
  upload_status: BlueprintUploadStatus;
  analysis_status: BlueprintAnalysisStatus;
  s3_key?: string;
  thumbnail_url?: string;
  created_at: string;
  updated_at: string;
}

export interface UploadUrlRequest {
  filename: string;
  content_type: string;
}

export interface UploadUrlResponse {
  blueprint_id: string;
  upload_url: string;
  expires_at: string;
}

export interface CompleteUploadRequest {
  success: boolean;
  error_message?: string;
}

// Job Types
export type JobStatus = 'queued' | 'processing' | 'completed' | 'failed';

export interface Job {
  id: string;
  blueprint_id: string;
  status: JobStatus;
  progress?: number;
  error_message?: string;
  result?: AnalysisResult;
  created_at: string;
  updated_at: string;
  started_at?: string;
  completed_at?: string;
}

export interface AnalysisResult {
  rooms?: Room[];
  openings?: Opening[];
  fixtures?: Fixture[];
  measurements?: Measurement[];
  materials?: Material[];
  summary?: {
    total_rooms: number;
    total_openings: number;
    total_fixtures: number;
    total_area?: number;
  };
}

export interface Room {
  name: string;
  dimensions: string;
  area: number;
  room_type?: string;
}

export interface Opening {
  opening_type: string;
  count: number;
  size: string;
  details?: string;
}

export interface Fixture {
  fixture_type: string;
  category: string;
  count: number;
  details?: string;
}

export interface Measurement {
  measurement_type: string;
  value: number;
  unit: string;
  location?: string;
}

export interface Material {
  material_name: string;
  quantity: number;
  unit: string;
  specifications?: string;
}

export interface TakeoffSummary {
  total_area: number;
  total_perimeter: number;
  opening_counts: Record<string, number>;
  fixture_counts: Record<string, number>;
  room_count: number;
  room_breakdown: RoomSummary[];
  opening_breakdown: OpeningSummary[];
  fixture_breakdown: FixtureSummary[];
}

export interface RoomSummary {
  name: string;
  room_type?: string;
  area: number;
  dimensions: string;
}

export interface OpeningSummary {
  opening_type: string;
  count: number;
  size: string;
}

export interface FixtureSummary {
  fixture_type: string;
  category: string;
  count: number;
}

export interface Coordinate {
  x: number;
  y: number;
}

// Trigger Analysis
export interface TriggerAnalysisResponse {
  job_id: string;
  status: JobStatus;
}
