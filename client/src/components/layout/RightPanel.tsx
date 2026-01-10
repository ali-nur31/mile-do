import { useEffect, useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useStore } from '../../store/useUIStore';
import { api } from '../../api/axios';
import { tasksApi } from '../../api/tasks';
import type { UpdateTaskRequest } from '../../types';
import { X, Trash2, Calendar, CheckCircle2 } from 'lucide-react';
import { toast } from 'sonner';

export const RightPanel = () => {
  const { selectedTaskId, selectTask } = useStore();
  const queryClient = useQueryClient();
  const [title, setTitle] = useState('');

  const { data: task, isLoading } = useQuery({
    queryKey: ['task', selectedTaskId],
    queryFn: () => tasksApi.getById(selectedTaskId!),
    enabled: !!selectedTaskId,
  });

  useEffect(() => {
    if (task) setTitle(task.title);
  }, [task]);

  const updateTask = useMutation({
    mutationFn: async (data: UpdateTaskRequest) => {
      if (!selectedTaskId) return;
      return api.patch(`/tasks/${selectedTaskId}`, data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    },
    onError: () => toast.error('Failed to update task')
  });

  const deleteTask = useMutation({
    mutationFn: async () => {
      if (!selectedTaskId) return;
      return tasksApi.delete(selectedTaskId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      selectTask(null);
      toast.custom(() => (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-100 dark:border-red-900/50 text-red-600 dark:text-red-400 px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 w-full">
          <Trash2 size={18} />
          <span className="font-medium text-sm">Task permanently deleted</span>
        </div>
      ), { duration: 3000 });
    }
  });

  const handleTitleBlur = () => {
    if (task && title !== task.title && title.trim() !== '') {
      updateTask.mutate({ title });
    } else if (title.trim() === '') {
      setTitle(task?.title || '');
    }
  };

  if (!selectedTaskId) return null;

  return (
    <div className="w-full h-full flex flex-col bg-white dark:bg-zinc-900 transition-colors duration-200">
      <div className="h-14 border-b border-zinc-100 dark:border-zinc-800 flex items-center justify-between px-6 flex-shrink-0">
        <div className="flex items-center gap-2 text-zinc-400 dark:text-zinc-500 text-xs font-medium">
          <CheckCircle2 size={14} className={task?.is_done ? 'text-blue-600 dark:text-blue-500' : ''} />
          {task?.is_done ? 'Completed' : 'In Progress'}
        </div>
        <div className="flex items-center gap-1">
          <button 
            onClick={() => deleteTask.mutate()}
            className="p-2 hover:bg-red-50 dark:hover:bg-red-900/20 text-zinc-400 dark:text-zinc-500 hover:text-red-600 dark:hover:text-red-400 rounded-md transition-colors"
            title="Delete Task"
          >
            <Trash2 size={18} />
          </button>
          <button 
            onClick={() => selectTask(null)}
            className="p-2 hover:bg-zinc-100 dark:hover:bg-zinc-800 text-zinc-400 dark:text-zinc-500 hover:text-zinc-600 dark:hover:text-zinc-300 rounded-md transition-colors"
          >
            <X size={20} />
          </button>
        </div>
      </div>

      {isLoading ? (
        <div className="p-8 text-zinc-400 dark:text-zinc-500 text-sm">Loading details...</div>
      ) : (
        <div className="flex-1 overflow-y-auto p-6 custom-scrollbar">
          <textarea
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            onBlur={handleTitleBlur}
            className="w-full text-xl font-bold text-zinc-900 dark:text-zinc-100 bg-transparent resize-none outline-none placeholder:text-zinc-300 dark:placeholder:text-zinc-700 min-h-[3rem]"
            placeholder="Task title"
            rows={2}
          />

          <div className="mt-6 space-y-4">
            <div className="flex items-center gap-3 text-sm text-zinc-600 dark:text-zinc-400">
              <div className="w-8 h-8 rounded-md bg-zinc-50 dark:bg-zinc-800 flex items-center justify-center text-zinc-400 dark:text-zinc-500">
                <Calendar size={16} />
              </div>
              <div className="flex flex-col">
                <span className="text-xs text-zinc-400 dark:text-zinc-500">Due Date</span>
                <span className="font-medium text-zinc-900 dark:text-zinc-200">
                  {task?.scheduled_date ? new Date(task.scheduled_date).toLocaleDateString() : 'No Date'}
                </span>
              </div>
            </div>
            
            <div className="p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-lg border border-zinc-100 dark:border-zinc-800 mt-8">
              <p className="text-xs text-zinc-400 dark:text-zinc-500 italic">
                Description and Priority fields are hidden.
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
