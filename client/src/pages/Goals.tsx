import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useGoals } from '../hooks/useGoals';
import { goalsApi } from '../api/goals';
import { tasksApi } from '../api/tasks';
import { api } from '../api/axios';
import type { Task } from '../types';
import { Loader2, Plus, CheckCircle2, Trash2 } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { Modal } from '../components/ui/Modal';
import { useNavigate } from 'react-router-dom';
import { useStore } from '../store/useUIStore';
import { showToast } from '../utils/toast';

export const Goals = () => {
  const { goals, isLoading, createGoal } = useGoals();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteAllOpen, setIsDeleteAllOpen] = useState(false);
  const [newGoalTitle, setNewGoalTitle] = useState('');
  const [newGoalColor, setNewGoalColor] = useState('#2563eb');
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { openContextMenu } = useStore();

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (newGoalTitle.trim().length < 3) {
        showToast('delete', 'List name too short (min 3 chars)');
        return;
    }
    if (createGoal.isPending) return;
    
    await createGoal.mutateAsync({ title: newGoalTitle, color: newGoalColor, category_type: 'growth' });
    showToast('success', 'List Created');
    setNewGoalTitle('');
    setIsModalOpen(false);
  };

  const completeGoal = useMutation({
    mutationFn: async (id: number) => {
      const res = await api.get<{ task_data: Task[] }>(`/goals/${id}/tasks`);
      const tasks = res.data.task_data || [];
      await Promise.all(tasks.map(t => tasksApi.delete(t.id)));
      return goalsApi.delete(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      showToast('success', 'Goal Completed!');
    },
    onError: () => {
        showToast('delete', 'Cannot delete default list');
    }
  });

  const deleteAllGoals = useMutation({
    mutationFn: async () => {
      if (!goals) return 0;
      const deletableGoals = goals.filter(g => g.title !== 'Other' && g.title !== 'Routine');
      
      for (const goal of deletableGoals) {
          const res = await api.get<{ task_data: Task[] }>(`/goals/${goal.id}/tasks`);
          const tasks = res.data.task_data || [];
          await Promise.all(tasks.map(t => tasksApi.delete(t.id)));
          await goalsApi.delete(goal.id);
      }
      return deletableGoals.length;
    },
    onSuccess: (count) => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      setIsDeleteAllOpen(false);
      
      if (count > 0) {
          showToast('success', `Deleted ${count} lists`);
      } else {
          showToast('delete', 'Default lists were kept');
      }
    }
  });

  if (isLoading) return <div className="flex justify-center h-64 items-center text-zinc-400"><Loader2 className="animate-spin" /></div>;

  return (
    <>
      <div className="flex items-end justify-between gap-4 mb-8">
        <div><h1 className="text-3xl font-bold text-zinc-900 dark:text-zinc-100">Lists & Goals</h1><p className="text-zinc-500 dark:text-zinc-400 mt-1">Manage your projects.</p></div>
        <div className="flex gap-2">
          {goals && goals.length > 0 && (<Button variant="ghost" onClick={() => setIsDeleteAllOpen(true)} className="text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20"><Trash2 size={16} /></Button>)}
          <Button onClick={() => setIsModalOpen(true)}><Plus size={16} className="mr-2" /> New List</Button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {goals?.map((goal) => {
            const isProtected = goal.title === 'Other' || goal.title === 'Routine';
            return (
              <div key={goal.id} className="group bg-white dark:bg-zinc-900 p-5 rounded-xl border border-zinc-200 dark:border-zinc-800 shadow-sm hover:shadow-md cursor-pointer transition-all relative" onClick={() => navigate(`/goals/${goal.id}`)} onContextMenu={(e) => openContextMenu(e, 'list', goal.id, { title: goal.title })}>
                <div className="flex items-center justify-between mb-3">
                  <div className="w-3 h-3 rounded-full" style={{ backgroundColor: goal.color || '#2563eb' }} />
                  <button onClick={(e) => { e.stopPropagation(); if(isProtected) { showToast('delete', 'Cannot delete default list'); return; } completeGoal.mutate(goal.id); }} className={`p-1.5 rounded-full transition-colors opacity-0 group-hover:opacity-100 ${isProtected ? "cursor-not-allowed text-zinc-300" : "hover:bg-green-50 text-zinc-300 hover:text-green-600"}`}><CheckCircle2 size={18} /></button>
                </div>
                <h3 className="font-semibold text-zinc-900 dark:text-zinc-100">{goal.title}</h3>
                <p className="text-xs text-zinc-400 dark:text-zinc-500 mt-2">Created {new Date(goal.created_at).toLocaleDateString()}</p>
              </div>
            );
        })}
      </div>

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title="Create New List">
        <form onSubmit={handleCreate} className="space-y-4">
          <input placeholder="List Name" value={newGoalTitle} onChange={(e) => setNewGoalTitle(e.target.value)} autoFocus className="w-full px-3 py-2 bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-lg text-sm text-zinc-900 dark:text-zinc-100" />
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Color</label>
            <div className="flex gap-2">{['#2563eb', '#dc2626', '#16a34a', '#d97706', '#9333ea'].map((c) => (<button key={c} type="button" onClick={() => setNewGoalColor(c)} className={`w-6 h-6 rounded-full border-2 ${newGoalColor === c ? 'border-zinc-900 dark:border-zinc-100 scale-110' : 'border-transparent'}`} style={{ backgroundColor: c }} />))}</div>
          </div>
          <div className="flex justify-end gap-2 pt-4"><Button type="button" variant="secondary" onClick={() => setIsModalOpen(false)}>Cancel</Button><Button type="submit" disabled={!newGoalTitle || createGoal.isPending}>{createGoal.isPending ? 'Creating...' : 'Create'}</Button></div>
        </form>
      </Modal>

      <Modal isOpen={isDeleteAllOpen} onClose={() => setIsDeleteAllOpen(false)} title="Delete All Lists?">
        <div className="space-y-4"><p className="text-sm text-zinc-600 dark:text-zinc-400">This will delete all lists (except Default ones) and their associated tasks.</p><div className="flex justify-end gap-2 pt-2"><Button variant="secondary" onClick={() => setIsDeleteAllOpen(false)}>Cancel</Button><Button onClick={() => deleteAllGoals.mutate()} className="bg-red-600 hover:bg-red-700 text-white">Delete All</Button></div></div>
      </Modal>
    </>
  );
};
