import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../api/axios';
import type { ListTasksResponse, UpdateTaskRequest, CreateTaskRequest } from '../types';
import { TaskItem } from '../features/tasks/TaskItem';
import { Loader2, Plus, Calendar } from 'lucide-react';
import { useTasks } from '../hooks/useTasks';
import { getTodayStr, combineToBackend } from '../utils/date';
import { showToast } from '../utils/toast';

export const Today = () => {
  const [inputValue, setInputValue] = useState('');
  const queryClient = useQueryClient();
  const { createTask } = useTasks();

  const { data, isLoading } = useQuery({
    queryKey: ['tasks', 'today'],
    queryFn: async () => {
      const todayStr = getTodayStr();
      const d = new Date();
      d.setDate(d.getDate() + 1);
      const year = d.getFullYear();
      const month = String(d.getMonth() + 1).padStart(2, '0');
      const day = String(d.getDate()).padStart(2, '0');
      const tomorrowStr = `${year}-${month}-${day}`;

      const res = await api.get<ListTasksResponse>(`/tasks/period?after_date=${todayStr}&before_date=${tomorrowStr}`);
      return res.data;
    }
  });

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputValue.trim()) return;

    const today = getTodayStr();
    const scheduledDateTime = combineToBackend(today, "09:00");
    
    const payload: CreateTaskRequest = {
      title: inputValue.trim(),
      goal_id: 0,
      scheduled_date_time: scheduledDateTime
    };
    
    createTask.mutate(payload, {
      onSuccess: () => {
        setInputValue('');
        showToast('success', 'Task added for today');
      }
    });
  };

  const toggleTask = useMutation({
    mutationFn: async ({ id, isDone }: { id: number; isDone: boolean }) => {
      const payload: UpdateTaskRequest = { is_done: !isDone };
      return api.patch(`/tasks/${id}`, payload);
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tasks'] })
  });

  const tasks = data?.task_data || [];
  const todoTasks = tasks.filter(t => !t.is_done);
  const doneTasks = tasks.filter(t => t.is_done);

  if (isLoading) return <div className="flex justify-center h-64 items-center text-zinc-400"><Loader2 className="animate-spin" /></div>;

  return (
    <div>
      <header className="mb-6 flex items-center gap-3">
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">Today</h1>
        <span className="text-zinc-400 dark:text-zinc-500 font-normal text-lg">{new Date().toLocaleDateString(undefined, { weekday: 'short', day: 'numeric' })}</span>
      </header>

      <form onSubmit={handleCreate} className="mb-6">
        <div className="relative group shadow-sm rounded-lg bg-zinc-50 dark:bg-zinc-900 focus-within:bg-white dark:focus-within:bg-zinc-950 focus-within:ring-2 focus-within:ring-blue-100 dark:focus-within:ring-blue-900 transition-all border border-transparent focus-within:border-blue-200 dark:focus-within:border-blue-800">
          <div className="absolute left-3 top-1/2 -translate-y-1/2 text-blue-500 dark:text-blue-400">
            <Plus size={20} />
          </div>
          <input
            type="text"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            placeholder="Add a task for today..."
            className="w-full pl-10 pr-4 py-3 bg-transparent border-none outline-none placeholder:text-zinc-400 dark:placeholder:text-zinc-600 text-[15px] text-zinc-900 dark:text-zinc-100"
          />
        </div>
      </form>

      <div className="flex flex-col gap-8">
        <div className="flex flex-col gap-2">
          {todoTasks.map((task) => (
            <TaskItem key={task.id} task={task} onToggle={(id, isDone) => toggleTask.mutate({ id, isDone })} />
          ))}
          {tasks.length === 0 && (
             <div className="text-center py-16">
                <div className="w-16 h-16 mx-auto bg-zinc-100 dark:bg-zinc-900 rounded-full flex items-center justify-center mb-4 text-zinc-300 dark:text-zinc-700">
                    <Calendar size={32} />
                </div>
                <p className="text-zinc-400 dark:text-zinc-600 text-sm">No tasks scheduled for today</p>
             </div>
          )}
        </div>

        {doneTasks.length > 0 && (
          <div>
             <div className="flex items-center gap-2 mb-3 px-2">
                <span className="text-xs font-semibold text-zinc-400 dark:text-zinc-500 bg-zinc-100 dark:bg-zinc-900 px-2 py-0.5 rounded">Completed</span>
                <div className="h-[1px] flex-1 bg-zinc-100 dark:bg-zinc-800"></div>
             </div>
             <div className="opacity-60 flex flex-col gap-2">
               {doneTasks.map((task) => (
                 <TaskItem key={task.id} task={task} onToggle={(id, isDone) => toggleTask.mutate({ id, isDone })} />
               ))}
             </div>
          </div>
        )}
      </div>
    </div>
  );
};