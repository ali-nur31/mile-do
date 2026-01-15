import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../api/axios';
import { tasksApi } from '../api/tasks';
import type { CreateTaskRequest, UpdateTaskRequest, UpdateTaskPayload, ListTasksResponse } from '../types';
import { extractDateStr, extractTimeStr, combineToBackend } from '../utils/date';

export const useTasks = () => {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['tasks', 'all'],
    queryFn: async () => {
      const res = await api.get<ListTasksResponse>('/tasks/');
      return res.data;
    }
  });

  const allTasks = data?.task_data || [];

  const createTask = useMutation({
    mutationFn: async (payload: CreateTaskRequest) => {
      console.log('🚀 Creating task with payload:', JSON.stringify(payload, null, 2));
      
      const response = await api.post('/tasks/', payload);
      
      console.log('Backend response:', JSON.stringify(response.data, null, 2));
      
      return response;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    }
  });

  const updateTask = useMutation({
    mutationFn: async ({ id, data }: { id: number; data: UpdateTaskRequest }) => {
      const currentTask = allTasks.find(t => t.id === id);
      
      if (!currentTask) throw new Error("Task not found in cache");

      let scheduledDateTime = "";
      
      if (data.scheduled_date_time !== undefined) {
        scheduledDateTime = data.scheduled_date_time;
      } else {
        const existingDate = extractDateStr(currentTask.scheduled_date);
        const existingTime = extractTimeStr(currentTask.scheduled_time);
        
        if (existingDate) {
          scheduledDateTime = combineToBackend(existingDate, existingTime || "09:00");
        } else {
          scheduledDateTime = "0001-01-01 00:00:00";
        }
      }

      const payload: UpdateTaskPayload = {
        title: data.title ?? currentTask.title,
        goal_id: data.goal_id ?? currentTask.goal_id,
        is_done: data.is_done ?? currentTask.is_done,
        scheduled_date_time: scheduledDateTime
      };

      console.log('Updating task:', id, 'with payload:', JSON.stringify(payload, null, 2));

      const response = await api.patch(`/tasks/${id}`, payload);
      
      console.log('Update response:', JSON.stringify(response.data, null, 2));
      
      return response;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    }
  });

  return {
    tasks: allTasks,
    isLoading,
    createTask,
    updateTask,
    deleteTask: tasksApi.delete
  };
};