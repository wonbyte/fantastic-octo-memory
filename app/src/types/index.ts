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
  company_name?: string;
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
  version: number;
  parent_blueprint_id?: string;
  is_latest: boolean;
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

// Bid Types
export type BidStatus = 'draft' | 'sent' | 'accepted' | 'rejected';

export interface Bid {
  id: string;
  project_id: string;
  job_id?: string;
  name?: string;
  total_cost?: number;
  labor_cost?: number;
  material_cost?: number;
  markup_percentage?: number;
  final_price?: number;
  status: BidStatus;
  bid_data?: string; // JSONB stored as string
  pdf_url?: string;
  pdf_s3_key?: string;
  version: number;
  parent_bid_id?: string;
  is_latest: boolean;
  created_at: string;
  updated_at: string;
}

export interface LineItem {
  description: string;
  trade: string;
  quantity: number;
  unit: string;
  unit_cost: number;
  total: number;
}

export interface BidData {
  bid_id: string;
  project_id: string;
  status: string;
  scope_of_work: string;
  line_items: LineItem[];
  labor_cost: number;
  material_cost: number;
  subtotal: number;
  markup_amount: number;
  total_price: number;
  exclusions: string[];
  inclusions: string[];
  schedule: Record<string, string>;
  payment_terms: string;
  warranty_terms: string;
  closing_statement: string;
}

export interface GenerateBidRequest {
  blueprint_id: string;
  markup_percentage?: number;
  company_name?: string;
  bid_name?: string;
}

export interface PricingConfig {
  material_prices: Record<string, number>;
  labor_rates: Record<string, number>;
  overhead_rate: number;
  profit_margin: number;
}

export interface PricingSummary {
  line_items: LineItem[];
  labor_cost: number;
  material_cost: number;
  subtotal: number;
  overhead_amount: number;
  markup_amount: number;
  total_price: number;
  costs_by_trade: Record<string, number>;
}

// Revision Types
export type ChangeType = 'added' | 'removed' | 'modified';

export interface BlueprintRevision {
  id: string;
  blueprint_id: string;
  version: number;
  filename: string;
  s3_key: string;
  file_size?: number;
  mime_type?: string;
  analysis_data?: string;
  changes_summary?: string;
  created_by?: string;
  created_at: string;
}

export interface BidRevision {
  id: string;
  bid_id: string;
  version: number;
  name?: string;
  total_cost?: number;
  labor_cost?: number;
  material_cost?: number;
  markup_percentage?: number;
  final_price?: number;
  status: BidStatus;
  bid_data?: string;
  changes_summary?: string;
  created_by?: string;
  created_at: string;
}

export interface BlueprintChange {
  change_type: ChangeType;
  category: string; // room, opening, fixture, measurement, material
  description: string;
  old_value?: unknown;
  new_value?: unknown;
  impact?: string; // High, Medium, Low
}

export interface BidChange {
  change_type: ChangeType;
  category: string; // cost, quantity, scope, timeline, terms, line_item
  trade?: string;
  description: string;
  old_value?: unknown;
  new_value?: unknown;
  impact?: string; // High, Medium, Low
}

export interface ComparisonSummary {
  total_changes: number;
  added_count: number;
  removed_count: number;
  modified_count: number;
  high_impact_count: number;
  changes_by_category: Record<string, number>;
}

export interface BlueprintComparison {
  from_version: number;
  to_version: number;
  changes: BlueprintChange[];
  summary: ComparisonSummary;
}

export interface BidComparison {
  from_version: number;
  to_version: number;
  changes: BidChange[];
  summary: ComparisonSummary;
}
