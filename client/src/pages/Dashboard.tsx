import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../api/axios';
import { tasksApi } from '../api/tasks';
import type { CreateTaskRequest, UpdateTaskRequest, ListTasksResponse } from '../types';
import { TaskItem } from '../features/tasks/TaskItem';
import { Plus, Loader2, Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { Modal } from '../components/ui/Modal';
import { Button } from '../components/ui/Button';

export const Dashboard = () => {
  const [inputValue, setInputValue] = useState('');
  const [showCompleted, setShowCompleted] = useState(false);
  const [isClearModalOpen, setIsClearModalOpen] = useState(false);
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['tasks', 'all'],
    queryFn: async () => {
      const res = await api.get<ListTasksResponse>('/tasks/');
      return res.data;
    }
  });

  const createTask = useMutation({
    mutationFn: async (title: string) => {
      const today = new Date().toISOString().split('T')[0];
      const payload: CreateTaskRequest = {
        title,
        goal_id: 0,
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
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tasks'] })
  });

  const deleteCompleted = useMutation({
    mutationFn: async (ids: number[]) => {
      await Promise.all(ids.map(id => tasksApi.delete(id)));
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      setIsClearModalOpen(false);
      toast.custom(() => (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-100 dark:border-red-900/50 text-red-600 dark:text-red-400 px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 w-full">
          <Trash2 size={18} />
          <span className="font-medium text-sm">Completed tasks cleared</span>
        </div>
      ), { duration: 3000 });
    }
  });

  const handleBatchDelete = () => {
    const idsToDelete = doneTasks.map(t => t.id);
    deleteCompleted.mutate(idsToDelete);
  };

  const allTasks = data?.task_data || [];
  const inboxTasks = allTasks.filter(t => t.goal_id === 0 || !t.goal_id);
  const todoTasks = inboxTasks.filter(t => !t.is_done);
  const doneTasks = inboxTasks.filter(t => t.is_done);

  if (isLoading) return <div className="flex justify-center h-64 items-center text-zinc-400"><Loader2 className="animate-spin" /></div>;

  return (
    <div>
      <header className="mb-6 flex justify-between items-center">
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">Inbox</h1>
        {doneTasks.length > 0 && (
          <button 
            onClick={() => setIsClearModalOpen(true)}
            className="text-xs text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 px-2 py-1 rounded transition-colors flex items-center gap-1"
          >
            <Trash2 size={12} /> Clear Completed
          </button>
        )}
      </header>

      <form onSubmit={(e) => { e.preventDefault(); if(inputValue.trim()) createTask.mutate(inputValue); }} className="mb-6">
        <div className="relative group shadow-sm rounded-lg bg-zinc-50 dark:bg-zinc-900 focus-within:bg-white dark:focus-within:bg-zinc-950 focus-within:ring-2 focus-within:ring-blue-100 dark:focus-within:ring-blue-900 transition-all border border-transparent focus-within:border-blue-200 dark:focus-within:border-blue-800">
          <div className="absolute left-3 top-1/2 -translate-y-1/2 text-blue-500 dark:text-blue-400">
            <Plus size={20} />
          </div>
          <input
            type="text"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            placeholder="Add a task to Inbox..."
            className="w-full pl-10 pr-4 py-3 bg-transparent border-none outline-none placeholder:text-zinc-400 dark:placeholder:text-zinc-600 text-[15px] text-zinc-900 dark:text-zinc-100"
          />
        </div>
      </form>

      <div className="flex flex-col gap-8">
        <div className="flex flex-col gap-2">
          {todoTasks.map((task) => (
            <TaskItem key={task.id} task={task} onToggle={(id, isDone) => toggleTask.mutate({ id, isDone })} />
          ))}
          {inboxTasks.length === 0 && (
             <div className="text-center py-16">
                <img src="https://cdn-icons-png.flaticon.com/512/7486/7486744.png" alt="Empty" className="w-16 h-16 mx-auto opacity-10 grayscale dark:invert mb-3" />
                <p className="text-zinc-400 dark:text-zinc-600 text-sm">Your inbox is empty</p>
             </div>
          )}
        </div>

        {doneTasks.length > 0 && (
          <div>
             <div className="flex items-center gap-3 mb-3 px-1">
                <button 
                  onClick={() => setShowCompleted(!showCompleted)}
                  className="text-xs font-bold text-zinc-400 dark:text-zinc-500 bg-zinc-100 dark:bg-zinc-900 px-2 py-1 rounded hover:bg-zinc-200 dark:hover:bg-zinc-800 transition-colors uppercase tracking-wide"
                >
                  {showCompleted ? 'Hide' : 'Show'} Completed ({doneTasks.length})
                </button>
                <div className="h-[1px] flex-1 bg-zinc-100 dark:bg-zinc-800"></div>
             </div>
             
             {showCompleted && (
               <div className="opacity-60 flex flex-col gap-2">
                 {doneTasks.slice(0, 50).map((task) => (
                   <TaskItem key={task.id} task={task} onToggle={(id, isDone) => toggleTask.mutate({ id, isDone })} />
                 ))}
                 {doneTasks.length > 50 && (
                   <div className="text-center text-xs text-zinc-400 py-2">
                     + {doneTasks.length - 50} more items hidden
                   </div>
                 )}
               </div>
             )}
          </div>
        )}
      </div>

      <Modal
        isOpen={isClearModalOpen}
        onClose={() => setIsClearModalOpen(false)}
        title="Clear Completed Tasks"
      >
        <div className="space-y-4">
          <p className="text-sm text-zinc-600 dark:text-zinc-400">
            Are you sure you want to delete all <strong>{doneTasks.length}</strong> completed tasks? This action cannot be undone.
          </p>
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="secondary" onClick={() => setIsClearModalOpen(false)}>Cancel</Button>
            <Button 
              onClick={handleBatchDelete} 
              className="bg-red-600 hover:bg-red-700 text-white focus:ring-red-500"
            >
              Clear Completed
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
};
