import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../api/axios';
import type { CreateTaskRequest, UpdateTaskRequest, ListTasksResponse } from '../types';

export const useTasks = () => {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['tasks', 'inbox'],
    queryFn: async () => {
      const res = await api.get<ListTasksResponse>('/tasks/inbox');
      return res.data;
    }
  });

  const createTask = useMutation({
    mutationFn: async (title: string) => {
      const payload: CreateTaskRequest = {
        title,
        goal_id: 0,
      };
      return api.post('/tasks/', payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    }
  });

  const toggleTask = useMutation({
    mutationFn: async ({ id, isDone }: { id: number; isDone: boolean }) => {
      const payload: UpdateTaskRequest = { is_done: !isDone };
      return api.patch(`/tasks/${id}`, payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    }
  });

  return {
    tasks: data?.task_data || [],
    isLoading,
    createTask,
    toggleTask
  };
};
