import { useEffect, useRef } from 'react';
import { useStore } from '../../store/useUIStore';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { tasksApi } from '../../api/tasks';
import { goalsApi } from '../../api/goals';
import { api } from '../../api/axios';
import { Trash2, Edit2, CheckSquare } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { showToast } from '../../utils/toast';
import type { Task } from '../../types';

export const ContextMenu = () => {
  const { contextMenu, closeContextMenu, selectTask } = useStore();
  const menuRef = useRef<HTMLDivElement>(null);
  const queryClient = useQueryClient();

  useEffect(() => {
    const handleClick = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        closeContextMenu();
      }
    };
    window.addEventListener('click', handleClick);
    return () => window.removeEventListener('click', handleClick);
  }, [closeContextMenu]);

  const deleteTask = useMutation({
    mutationFn: (id: number) => tasksApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      showToast('delete', 'Task Deleted', 'Permanently removed.');
      closeContextMenu();
    }
  });

  const deleteList = useMutation({
    mutationFn: async (id: number) => {
      const res = await api.get<{ task_data: Task[] }>(`/goals/${id}/tasks`);
      const tasks = res.data.task_data || [];
      await Promise.all(tasks.map(t => tasksApi.delete(t.id)));
      return goalsApi.delete(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      showToast('delete', 'List Deleted', 'List and its tasks removed.');
      closeContextMenu();
    }
  });

  const toggleTaskStatus = useMutation({
    mutationFn: async (id: number) => {
      const task = await tasksApi.getById(id);
      return api.patch(`/tasks/${id}`, { is_done: !task.is_done });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      closeContextMenu();
    }
  });

  if (!contextMenu.isOpen) return null;

  return (
    <AnimatePresence>
      <motion.div
        ref={menuRef}
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.95 }}
        transition={{ duration: 0.1 }}
        style={{ top: contextMenu.y, left: contextMenu.x }}
        className="fixed z-50 w-52 bg-white dark:bg-zinc-900 rounded-lg shadow-xl border border-zinc-100 dark:border-zinc-800 py-1 overflow-hidden"
      >
        {contextMenu.type === 'task' && contextMenu.targetId && (
          <>
            <button
              onClick={() => {
                selectTask(contextMenu.targetId);
                closeContextMenu();
              }}
              className="w-full text-left px-4 py-2.5 text-sm text-zinc-700 dark:text-zinc-200 hover:bg-zinc-50 dark:hover:bg-zinc-800 flex items-center gap-2"
            >
              <Edit2 size={14} /> Edit Task
            </button>
            <button
              onClick={() => toggleTaskStatus.mutate(contextMenu.targetId!)}
              className="w-full text-left px-4 py-2.5 text-sm text-zinc-700 dark:text-zinc-200 hover:bg-zinc-50 dark:hover:bg-zinc-800 flex items-center gap-2"
            >
              <CheckSquare size={14} /> Toggle Status
            </button>
            <div className="h-[1px] bg-zinc-100 dark:bg-zinc-800 my-1" />
            <button
              onClick={() => deleteTask.mutate(contextMenu.targetId!)}
              className="w-full text-left px-4 py-2.5 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 flex items-center gap-2"
            >
              <Trash2 size={14} /> Delete
            </button>
          </>
        )}

        {contextMenu.type === 'list' && contextMenu.targetId && (
          <>
            <div className="px-4 py-2 text-xs font-semibold text-zinc-400 dark:text-zinc-500 uppercase tracking-wider border-b border-zinc-50 dark:border-zinc-800 mb-1">
              {contextMenu.data?.title || 'List Options'}
            </div>

            {!['routine', 'other'].includes(contextMenu.data?.title?.toLowerCase()) && (
              <button
                onClick={() => deleteList.mutate(contextMenu.targetId!)}
                className="w-full text-left px-4 py-2.5 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 flex items-center gap-2"
              >
                <Trash2 size={14} /> Delete List
              </button>
            )}
          </>
        )}
      </motion.div>
    </AnimatePresence>
  );
};
