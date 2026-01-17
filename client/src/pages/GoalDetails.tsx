import { useState } from 'react';
import { useParams, Navigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../api/axios';
import { goalsApi } from '../api/goals';
import type { Task, CreateTaskRequest, UpdateTaskRequest } from '../types';
import { TaskItem } from '../features/tasks/TaskItem';
import { Loader2, Plus, Hash } from 'lucide-react';
import { showToast } from '../utils/toast';

export const GoalDetails = () => {
  const { id } = useParams<{ id: string }>();
  const goalId = id ? parseInt(id, 10) : 0;

  if (!goalId || isNaN(goalId)) {
    return <Navigate to="/goals" replace />;
  }

  const queryClient = useQueryClient();
  const [inputValue, setInputValue] = useState('');

  const { data: goal, isLoading: isGoalLoading } = useQuery({
    queryKey: ['goals', goalId],
    queryFn: () => goalsApi.getById(goalId),
    enabled: !!goalId
  });

  const { data: tasks, isLoading: isTasksLoading } = useQuery({
    queryKey: ['tasks', 'goal', goalId],
    queryFn: async () => {
      const res = await api.get<{ task_data: Task[] }>(`/goals/${goalId}/tasks`);
      return res.data.task_data || [];
    },
    enabled: !!goalId
  });

  const createTask = useMutation({
    mutationFn: async (title: string) => {
      const payload: CreateTaskRequest = {
        title: title.trim(),
        goal_id: goalId,
      };
      return api.post('/tasks/', payload);
    },
    onSuccess: () => {
      setInputValue('');
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      showToast('success', 'Task added to list');
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

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputValue.trim()) return;

    if (inputValue.trim().length < 3) {
      showToast('delete', 'Task title too short (min 3 chars)');
      return;
    }

    createTask.mutate(inputValue);
  };

  if (isGoalLoading || isTasksLoading) {
    return (
      <div className="flex justify-center h-64 items-center text-zinc-400">
        <Loader2 className="animate-spin" size={24} />
      </div>
    );
  }

  const activeTasks = tasks?.filter(t => !t.is_done) || [];
  const completedTasks = tasks?.filter(t => t.is_done) || [];

  return (
    <div className="max-w-3xl mx-auto">
      <header className="mb-6 pb-6 border-b border-zinc-100 dark:border-zinc-800 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <div className="w-12 h-12 rounded-xl bg-zinc-50 dark:bg-zinc-900 flex items-center justify-center border border-zinc-100 dark:border-zinc-800 shadow-sm text-zinc-500 dark:text-zinc-400">
            <Hash size={24} style={{ color: goal?.color }} />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">{goal?.title}</h1>
            <p className="text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wide mt-1">
              {activeTasks.length} tasks remaining
            </p>
          </div>
        </div>
      </header>
      <form onSubmit={handleSubmit} className="mb-8 relative group">
        <div className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-zinc-500 group-focus-within:text-blue-500 dark:group-focus-within:text-blue-400 transition-colors">
          <Plus size={20} />
        </div>
        <input
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          placeholder={`Add a task to ${goal?.title}...`}
          className="w-full pl-12 pr-4 py-3.5 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl shadow-sm outline-none 
                     focus:ring-2 focus:ring-blue-100 dark:focus:ring-blue-900 focus:border-blue-500 dark:focus:border-blue-700 transition-all placeholder:text-zinc-400 dark:placeholder:text-zinc-600 text-zinc-900 dark:text-zinc-100"
        />
      </form>

      <div className="space-y-8">
        <div className="flex flex-col gap-2">
          {activeTasks.map((task) => (
            <TaskItem
              key={task.id}
              task={task}
              onToggle={(id, isDone) => toggleTask.mutate({ id, isDone })}
            />
          ))}
          {activeTasks.length === 0 && completedTasks.length === 0 && (
            <div className="text-zinc-400 dark:text-zinc-600 text-center py-10 text-sm italic">
              This list is empty. Add a task above.
            </div>
          )}
        </div>

        {completedTasks.length > 0 && (
          <div>
            <h2 className="text-xs font-bold text-zinc-400 dark:text-zinc-500 uppercase tracking-widest mb-3 px-1">
              Completed
            </h2>
            <div className="flex flex-col gap-2 opacity-60 hover:opacity-100 transition-opacity">
              {completedTasks.map((task) => (
                <TaskItem
                  key={task.id}
                  task={task}
                  onToggle={(id, isDone) => toggleTask.mutate({ id, isDone })}
                />
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};