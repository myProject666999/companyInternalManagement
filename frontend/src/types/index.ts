export interface User {
  id: number;
  username: string;
  name: string;
  gender?: string;
  birth_date?: string;
  phone?: string;
  email?: string;
  address?: string;
  role: Role;
  department_id?: number;
  department?: Department;
  position?: string;
  join_date?: string;
  status: string;
  created_at: string;
}

export type Role = 'general_manager' | 'department_manager' | 'employee';

export interface Department {
  id: number;
  name: string;
  phone?: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface Attendance {
  id: number;
  user_id: number;
  user?: User;
  type: AttendanceType;
  date: string;
  time?: string;
  status: AttendanceStatus;
  reason?: string;
  remark?: string;
  approver_id?: number;
  created_at: string;
  updated_at: string;
}

export type AttendanceType = 'clock_in' | 'clock_out' | 'leave' | 'business_trip' | 'overtime';
export type AttendanceStatus = 'pending' | 'approved' | 'rejected';

export interface WorkTask {
  id: number;
  title: string;
  description?: string;
  assignee_id: number;
  assignee?: User;
  assignor_id: number;
  assignor?: User;
  start_date?: string;
  end_date?: string;
  priority: string;
  status: string;
  progress: number;
  evaluation?: string;
  evaluation_score?: number;
  created_at: string;
  updated_at: string;
}

export interface Information {
  id: number;
  title: string;
  content?: string;
  type: string;
  author_id: number;
  author?: User;
  is_public: boolean;
  department_id?: number;
  view_count: number;
  created_at: string;
  updated_at: string;
}

export interface Document {
  id: number;
  title: string;
  description?: string;
  file_path: string;
  file_size: number;
  file_type?: string;
  uploader_id: number;
  uploader?: User;
  is_public: boolean;
  department_id?: number;
  download_count: number;
  created_at: string;
  updated_at: string;
}

export interface Message {
  id: number;
  sender_id: number;
  sender?: User;
  receiver_id: number;
  receiver?: User;
  subject: string;
  content?: string;
  is_read: boolean;
  reply_to_id?: number;
  parent_id?: number;
  created_at: string;
  updated_at: string;
}

export interface WorkLog {
  id: number;
  user_id: number;
  user?: User;
  date: string;
  work_summary?: string;
  work_content?: string;
  tomorrow_plan?: string;
  problems?: string;
  task_id?: number;
  created_at: string;
  updated_at: string;
}

export interface WorkReport {
  id: number;
  user_id: number;
  user?: User;
  title: string;
  content?: string;
  type?: string;
  task_id?: number;
  created_at: string;
  updated_at: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface ApiResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}
