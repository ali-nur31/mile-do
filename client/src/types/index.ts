export interface AuthResponse {
  access_token: string;
  refresh_token: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  confirm_password: string;
}

export interface Goal {
  id: number;
  user_id: number;
  title: string;
  color: string;
  category_type: 'growth' | 'maintenance' | 'other';
  is_archived: boolean;
  created_at: string;
}

export interface ListGoalsResponse {
  user_id: number;
  data: Goal[];
}

export interface CreateGoalRequest {
  title: string;
  color: string;
  category_type: string;
}

export interface Task {
  id: number;
  goal_id: number;
  user_id: number;
  title: string;
  is_done: boolean;
  scheduled_date: string;
  has_time: boolean;
  scheduled_time: string;
  duration_minutes: number;
  reschedule_count: number;
  created_at: string;
}

export interface ListTasksResponse {
  user_id: number;
  task_data: Task[];
}

export interface CreateTaskRequest {
  goal_id: number;
  title: string;
  scheduled_date_time?: string;
  scheduled_end_date_time?: string;
}

export interface UpdateTaskRequest {
  goal_id?: number;
  title?: string;
  is_done?: boolean;
  scheduled_date_time?: string;
}
