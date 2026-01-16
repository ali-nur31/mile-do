import { api } from './axios';
import type { CreateGoalRequest, ListGoalsResponse, Goal, Task } from '../types';

export interface UpdateGoalRequest {
  id: number;
  title: string;
  color: string;
  category_type: string;
  is_archived: boolean;
}

export const goalsApi = {
  getAll: async () => {
    const response = await api.get<ListGoalsResponse>('/goals/');
    return response.data.data || [];
  },

  getById: async (id: number) => {
    const response = await api.get<Goal>(`/goals/${id}`);
    return response.data;
  },

  getTasksByGoal: async (id: number) => {
    const response = await api.get<{ task_data: Task[] }>(`/goals/${id}/tasks`);
    return response.data.task_data || [];
  },

  create: async (data: CreateGoalRequest) => {
    const response = await api.post<Goal>('/goals/', data);
    return response.data;
  },

  update: async (id: number, current: Goal, updates: Partial<UpdateGoalRequest>) => {
    const payload: UpdateGoalRequest = {
      id: id,
      title: updates.title ?? current.title,
      color: updates.color ?? current.color,
      category_type: updates.category_type ?? current.category_type,
      is_archived: updates.is_archived ?? current.is_archived
    };
    
    const response = await api.patch('/goals/', payload);
    return response.data;
  },

  delete: async (id: number) => {
    await api.delete(`/goals/${id}`);
  }
};