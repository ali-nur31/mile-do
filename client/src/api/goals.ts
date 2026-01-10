import { api } from './axios';
import type { CreateGoalRequest, ListGoalsResponse, Goal, Task } from '../types';

export interface UpdateGoalRequest extends Partial<CreateGoalRequest> {
  id: number;
  is_archived?: boolean;
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

  update: async (data: UpdateGoalRequest) => {
    const response = await api.patch('/goals/', data);
    return response.data;
  },

  delete: async (id: number) => {
    await api.delete(`/goals/${id}`);
  }
};
