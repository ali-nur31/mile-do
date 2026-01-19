import { api } from './axios';
import type { CreateTaskRequest, UpdateTaskRequest, ListTasksResponse, Task } from '../types';

export const tasksApi = {
  getInbox: async () => {
    const response = await api.get<ListTasksResponse>('/tasks/inbox');
    return response.data.task_data || [];
  },
  getAll: async () => {
    const response = await api.get<ListTasksResponse>('/tasks/');
    return response.data.task_data || [];
  },
  getByPeriod: async (period: string) => {
    const response = await api.get<ListTasksResponse>(`/tasks/period?p=${period}`);
    return response.data.task_data || [];
  },
  getById: async (id: number) => {
    const response = await api.get<Task>(`/tasks/${id}`);
    return response.data;
  },
  create: async (data: CreateTaskRequest) => {
    const response = await api.post<Task>('/tasks/', data);
    return response.data;
  },
  update: async (id: number, data: UpdateTaskRequest) => {
    const response = await api.patch<Task>(`/tasks/${id}`, data);
    return response.data;
  },
  complete: async (id: number) => {
    const response = await api.patch<Task>(`/tasks/${id}/complete`);
    return response.data;
  },
  delete: async (id: number) => {
    await api.delete(`/tasks/${id}`);
  },
  analyzeToday: async () => {
    const response = await api.get('/tasks/analyze');
    return response.data;
  }
};
