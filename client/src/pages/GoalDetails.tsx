import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../api/axios';
import { goalsApi } from '../api/goals';
import type { Task, CreateTaskRequest, UpdateTaskRequest } from '../types';
import { TaskItem } from '../features/tasks/TaskItem';
import { Loader2, Plus, Hash } from 'lucide-react';

export const GoalDetails = () => {
  const { id } = useParams<{ id: string }>();
  const goalId = parseInt(id || '0');
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
      const today = new Date().toISOString().split('T')[0];
      const payload: CreateTaskRequest = {
        title,
        goal_id: goalId,
        scheduled_date_time: today 
      };
      return api.post('/tasks/', payload);
    },
    onSuccess: () => {
      setInputValue('');
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

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputValue.trim()) return;
    createTask.mutate(inputValue);
  };

  if (isGoalLoading || isTasksLoading) {
    return <div className="flex justify-center h-64 items-center text-zinc-400"><Loader2 className="animate-spin" /></div>;
  }

  const activeTasks = tasks?.filter(t => !t.is_done) || [];
  const completedTasks = tasks?.filter(t => t.is_done) || [];

  return (
    <div>
      <header className="mb-6 flex items-center gap-3">
        <div className="p-2 bg-zinc-100 dark:bg-zinc-900 rounded-lg text-zinc-500 dark:text-zinc-400">
           <Hash size={24} style={{ color: goal?.color }} />
        </div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">{goal?.title}</h1>
      </header>

      <form onSubmit={handleSubmit} className="mb-6">
        <div className="relative group shadow-sm rounded-lg bg-zinc-50 dark:bg-zinc-900 focus-within:bg-white dark:focus-within:bg-zinc-950 focus-within:ring-2 focus-within:ring-blue-100 dark:focus-within:ring-blue-900 transition-all border border-transparent focus-within:border-blue-200 dark:focus-within:border-blue-800">
          <div className="absolute left-3 top-1/2 -translate-y-1/2 text-blue-500 dark:text-blue-400">
            <Plus size={20} />
          </div>
          <input
            type="text"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            placeholder={`Add to ${goal?.title}...`}
            className="w-full pl-10 pr-4 py-3 bg-transparent border-none outline-none placeholder:text-zinc-400 dark:placeholder:text-zinc-600 text-[15px] text-zinc-900 dark:text-zinc-100"
          />
        </div>
      </form>

      <div className="flex flex-col gap-8">
        <div className="flex flex-col gap-2">
          {activeTasks.map((task) => (
            <TaskItem key={task.id} task={task} onToggle={(id, isDone) => toggleTask.mutate({ id, isDone })} />
          ))}
          {activeTasks.length === 0 && completedTasks.length === 0 && (
            <div className="text-center py-10 text-zinc-400 dark:text-zinc-600 text-sm">Empty list</div>
          )}
        </div>

        {completedTasks.length > 0 && (
          <div>
            <div className="flex items-center gap-2 mb-3 px-2">
               <span className="text-xs font-semibold text-zinc-400 dark:text-zinc-500 bg-zinc-100 dark:bg-zinc-900 px-2 py-0.5 rounded">Completed</span>
               <div className="h-[1px] flex-1 bg-zinc-100 dark:bg-zinc-800"></div>
            </div>
            <div className="opacity-60 flex flex-col gap-2">
              {completedTasks.map((task) => (
                <TaskItem key={task.id} task={task} onToggle={(id, isDone) => toggleTask.mutate({ id, isDone })} />
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
