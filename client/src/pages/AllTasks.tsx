import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../api/axios';
import { tasksApi } from '../api/tasks';
import type { ListTasksResponse, UpdateTaskRequest } from '../types';
import { TaskItem } from '../features/tasks/TaskItem';
import { Loader2, Layers, Search, Trash2 } from 'lucide-react';
import { Modal } from '../components/ui/Modal';
import { Button } from '../components/ui/Button';
import { showToast } from '../utils/toast';

export const AllTasks = () => {
  const [search, setSearch] = useState('');
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['tasks', 'all'],
    queryFn: async () => {
      const res = await api.get<ListTasksResponse>('/tasks/');
      return res.data;
    }
  });

  const toggleTask = useMutation({
    mutationFn: async ({ id, isDone }: { id: number; isDone: boolean }) => {
      const payload: UpdateTaskRequest = { is_done: !isDone };
      return api.patch(`/tasks/${id}`, payload);
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tasks'] })
  });

  const deleteAll = useMutation({
    mutationFn: async (ids: number[]) => {
      await Promise.all(ids.map(id => tasksApi.delete(id)));
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      setIsDeleteModalOpen(false);
      showToast('delete', 'All Tasks Deleted', 'Your system is now empty.');
    }
  });

  if (isLoading) return <div className="flex justify-center h-64 items-center text-zinc-400"><Loader2 className="animate-spin" /></div>;

  const allTasks = data?.task_data || [];
  const filteredTasks = allTasks.filter(t => t.title.toLowerCase().includes(search.toLowerCase()));

  const scheduledTasks = filteredTasks.filter(t => t.scheduled_date && !t.scheduled_date.startsWith('0001') && !t.is_done);
  const backlogTasks = filteredTasks.filter(t => (!t.scheduled_date || t.scheduled_date.startsWith('0001')) && !t.is_done);
  const completedTasks = filteredTasks.filter(t => t.is_done);

  return (
    <div className="max-w-3xl mx-auto">
      <header className="mb-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-zinc-100 dark:bg-zinc-800 rounded-xl flex items-center justify-center text-zinc-500">
              <Layers size={20} />
            </div>
            <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">All Tasks</h1>
          </div>
          {allTasks.length > 0 && (
            <button 
              onClick={() => setIsDeleteModalOpen(true)}
              className="text-xs text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 px-3 py-2 rounded-lg transition-colors flex items-center gap-2 font-medium"
            >
              <Trash2 size={16} /> Delete All
            </button>
          )}
        </div>

        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400" size={16} />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search tasks..."
            className="w-full pl-10 pr-4 py-2.5 bg-zinc-50 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-700 rounded-lg text-sm text-zinc-900 dark:text-zinc-100 focus:ring-2 focus:ring-blue-500 outline-none transition-all"
          />
        </div>
      </header>

      <div className="space-y-10">
        {scheduledTasks.length > 0 && (
          <section>
            <h2 className="text-xs font-bold text-zinc-400 dark:text-zinc-500 uppercase tracking-widest mb-3 px-1">Scheduled</h2>
            <div className="flex flex-col gap-2">
              {scheduledTasks.map((task) => (
                <TaskItem key={task.id} task={task} onToggle={(id, isDone) => toggleTask.mutate({ id, isDone })} />
              ))}
            </div>
          </section>
        )}

        {backlogTasks.length > 0 && (
          <section>
            <h2 className="text-xs font-bold text-zinc-400 dark:text-zinc-500 uppercase tracking-widest mb-3 px-1">Backlog</h2>
            <div className="flex flex-col gap-2">
              {backlogTasks.map((task) => (
                <TaskItem key={task.id} task={task} onToggle={(id, isDone) => toggleTask.mutate({ id, isDone })} />
              ))}
            </div>
          </section>
        )}

        {completedTasks.length > 0 && (
          <section>
            <h2 className="text-xs font-bold text-zinc-400 dark:text-zinc-500 uppercase tracking-widest mb-3 px-1">Completed</h2>
            <div className="flex flex-col gap-2 opacity-60">
              {completedTasks.map((task) => (
                <TaskItem key={task.id} task={task} onToggle={(id, isDone) => toggleTask.mutate({ id, isDone })} />
              ))}
            </div>
          </section>
        )}
      </div>

      <Modal isOpen={isDeleteModalOpen} onClose={() => setIsDeleteModalOpen(false)} title="Delete Everything?">
        <div className="space-y-4">
          <p className="text-sm text-zinc-600 dark:text-zinc-400">
            This will permanently delete <strong>{allTasks.length}</strong> tasks. This action cannot be undone.
          </p>
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="secondary" onClick={() => setIsDeleteModalOpen(false)}>Cancel</Button>
            <Button onClick={() => deleteAll.mutate(allTasks.map(t => t.id))} className="bg-red-600 hover:bg-red-700 text-white">Delete All</Button>
          </div>
        </div>
      </Modal>
    </div>
  );
};
