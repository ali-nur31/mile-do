import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../api/axios';
import { goalsApi } from '../api/goals';
import { tasksApi } from '../api/tasks';
import type { CreateTaskRequest, UpdateTaskRequest, ListTasksResponse, Goal } from '../types';
import { extractDateStr, extractTimeStr, combineToBackend, getLocalISOString } from '../utils/date';
import { showToast } from '../utils/toast';

export const useTasks = () => {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['tasks', 'all'],
    queryFn: async () => {
      try {
        const res = await api.get<ListTasksResponse>('/tasks/');
        return res.data;
      } catch (error) {
        return { task_data: [] };
      }
    },
    retry: 1
  });

  const allTasks = data?.task_data || [];

  const createTask = useMutation({
    mutationFn: async (payload: CreateTaskRequest) => {
      if (!payload.title || payload.title.trim().length < 3) {
        throw new Error("Task title must be at least 3 characters.");
      }

      let finalPayload = { ...payload };

      if (!finalPayload.goal_id || finalPayload.goal_id === 0) {
        try {
          let goals = queryClient.getQueryData<Goal[]>(['goals']);
          if (!goals || goals.length === 0) {
            goals = await goalsApi.getAll();
            queryClient.setQueryData(['goals'], goals);
          }

          const otherGoal = goals?.find(g => g.title === 'Other');
          const routineGoal = goals?.find(g => g.title === 'Routine');

          if (!routineGoal) await goalsApi.create({ title: "Routine", color: "#73260A", category_type: "maintenance" });
          let targetGoal = otherGoal;
          if (!targetGoal) targetGoal = await goalsApi.create({ title: "Other", color: "#3b82f6", category_type: "other" });

          if (!otherGoal || !routineGoal) queryClient.invalidateQueries({ queryKey: ['goals'] });

          if (targetGoal) finalPayload.goal_id = targetGoal.id;
        } catch (e) {
          console.error(e);
        }
      }

      if (finalPayload.scheduled_date_time && finalPayload.duration_minutes) {
        const duration = Math.max(15, finalPayload.duration_minutes);
        try {
          const [datePart, timePart] = finalPayload.scheduled_date_time.split(' ');
          if (datePart && timePart) {
            const [hours, minutes] = timePart.split(':').map(Number);
            const startDate = new Date();
            startDate.setFullYear(
              parseInt(datePart.split('-')[0]),
              parseInt(datePart.split('-')[1]) - 1,
              parseInt(datePart.split('-')[2])
            );
            startDate.setHours(hours || 0, minutes || 0, 0, 0);

            const endDate = new Date(startDate);
            endDate.setMinutes(endDate.getMinutes() + duration);

            const endDateStr = getLocalISOString(endDate);
            const endTimeStr = `${String(endDate.getHours()).padStart(2, '0')}:${String(endDate.getMinutes()).padStart(2, '0')}:00`;
            finalPayload.scheduled_end_date_time = `${endDateStr} ${endTimeStr}`;
          }
        } catch (e) {
          console.warn('Failed to calculate end date time for create:', e);
        }
      }

      const { duration_minutes, ...payloadToSend } = finalPayload;
      (payloadToSend as any).duration_minutes = duration_minutes;

      const response = await api.post('/tasks/', payloadToSend);
      return response.data;
    },
    onSuccess: (newTask) => {
      queryClient.setQueryData(['tasks', 'all'], (old: ListTasksResponse | undefined) => {
        if (!old) return { user_id: 0, task_data: [newTask] };
        return { ...old, task_data: [...old.task_data, newTask] };
      });
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    },
    onError: (error: any) => {
      const msg = error.message || error.response?.data?.message || 'Failed to create task';
      showToast('delete', msg);
    }
  });

  const updateTask = useMutation({
    mutationFn: async ({ id, data }: { id: number; data: UpdateTaskRequest }) => {
      const currentTask = allTasks.find(t => t.id === id);
      if (!currentTask) throw new Error("Task not found");

      if (data.title !== undefined && data.title.trim().length < 3) {
        throw new Error("Task title must be at least 3 characters.");
      }

      let scheduledDateTime: string;

      if (data.scheduled_date_time !== undefined) {
        scheduledDateTime = data.scheduled_date_time;
      } else {
        const existingDate = extractDateStr(currentTask.scheduled_date);
        const existingTime = extractTimeStr(currentTask.scheduled_time);
        if (existingDate && !existingDate.startsWith('0001')) {
          scheduledDateTime = combineToBackend(existingDate, existingTime || "09:00") || "";
        } else {
          const today = new Date();
          const todayStr = getLocalISOString(today);
          scheduledDateTime = combineToBackend(todayStr, "09:00") || "";
        }
      }

      let scheduledEndDateTime: string | undefined = undefined;
      const duration = data.duration_minutes ?? currentTask.duration_minutes ?? 15;

      if (scheduledDateTime && scheduledDateTime !== "" && !scheduledDateTime.startsWith('0001')) {
        try {
          const [datePart, timePart] = scheduledDateTime.split(' ');
          if (datePart && timePart) {
            const [hours, minutes] = timePart.split(':').map(Number);
            const endDate = new Date();
            endDate.setFullYear(parseInt(datePart.split('-')[0]), parseInt(datePart.split('-')[1]) - 1, parseInt(datePart.split('-')[2]));
            endDate.setHours(hours || 0, minutes || 0, 0, 0);
            endDate.setMinutes(endDate.getMinutes() + duration);

            const endDateStr = getLocalISOString(endDate);
            const endTimeStr = `${String(endDate.getHours()).padStart(2, '0')}:${String(endDate.getMinutes()).padStart(2, '0')}:00`;
            scheduledEndDateTime = `${endDateStr} ${endTimeStr}`;
          }
        } catch (e) {
          console.warn('Failed to calculate end date time:', e);
        }
      }

      const goalIdValue = data.goal_id ?? currentTask.goal_id;
      const titleValue = data.title ?? currentTask.title;
      const isDoneValue = data.is_done ?? currentTask.is_done;

      const payload: any = {
        goal_id: goalIdValue,
        title: titleValue,
        is_done: isDoneValue,
        scheduled_date_time: scheduledDateTime,
        duration_minutes: duration
      };

      if (scheduledEndDateTime) {
        payload.scheduled_end_date_time = scheduledEndDateTime;
      }

      return api.patch(`/tasks/${id}`, payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    },
    onError: (err: any) => {
      const msg = err.message || err.response?.data?.message || 'Failed to update task';
      showToast('delete', msg);
    }
  });

  const completeTask = useMutation({
    mutationFn: async (id: number) => {
      return tasksApi.complete(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    },
    onError: () => {
      showToast('delete', 'Failed to complete task');
    }
  });

  const deleteTask = useMutation({
    mutationFn: async (id: number) => {
      return api.delete(`/tasks/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      showToast('success', 'Task deleted');
    },
    onError: () => {
      showToast('delete', 'Failed to delete task');
    }
  });

  return {
    tasks: allTasks,
    isLoading,
    createTask,
    updateTask,
    completeTask,
    deleteTask
  };
};
