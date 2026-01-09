import { useState } from 'react';
import { useGoals } from '../hooks/useGoals';
import { Loader2, Plus } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { Modal } from '../components/ui/Modal';
import { Input } from '../components/ui/Input';
import { useNavigate } from 'react-router-dom';

export const Goals = () => {
  const { goals, isLoading, createGoal } = useGoals();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newGoalTitle, setNewGoalTitle] = useState('');
  const [newGoalColor, setNewGoalColor] = useState('#2563eb');
  const navigate = useNavigate();

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newGoalTitle.trim()) return;
    
    await createGoal.mutateAsync({
      title: newGoalTitle,
      color: newGoalColor,
      category_type: 'growth'
    });
    setNewGoalTitle('');
    setIsModalOpen(false);
  };

  if (isLoading) {
    return <div className="flex justify-center h-64 items-center text-zinc-400"><Loader2 className="animate-spin" /></div>;
  }

  return (
    <>
      <div className="flex items-end justify-between gap-4 mb-8">
        <div>
          <h1 className="text-3xl font-bold text-zinc-900 dark:text-zinc-100 tracking-tight">Lists & Goals</h1>
          <p className="text-zinc-500 dark:text-zinc-400 mt-1">Manage your projects.</p>
        </div>
        <Button onClick={() => setIsModalOpen(true)}>
          <Plus size={16} className="mr-2" />
          New List
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {goals?.map((goal) => (
          <div 
            key={goal.id} 
            onClick={() => navigate(`/goals/${goal.id}`)}
            className="group bg-white dark:bg-zinc-900 p-5 rounded-xl border border-zinc-200 dark:border-zinc-800 shadow-sm hover:shadow-md hover:border-blue-200 dark:hover:border-blue-900/50 cursor-pointer transition-all"
          >
            <div className="flex items-center justify-between mb-3">
              <div 
                className="w-2 h-2 rounded-full"
                style={{ backgroundColor: goal.color || '#2563eb' }}
              />
            </div>
            <h3 className="font-semibold text-zinc-900 dark:text-zinc-100 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">{goal.title}</h3>
            <p className="text-xs text-zinc-400 dark:text-zinc-500 mt-2">
               Created {new Date(goal.created_at).toLocaleDateString()}
            </p>
          </div>
        ))}
      </div>

      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Create New List"
      >
        <form onSubmit={handleCreate} className="space-y-4">
          <Input 
            placeholder="List Name" 
            value={newGoalTitle}
            onChange={(e) => setNewGoalTitle(e.target.value)}
            autoFocus
            className="dark:bg-zinc-950 dark:border-zinc-700 dark:text-zinc-100 dark:placeholder:text-zinc-600"
          />
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">Color</label>
            <div className="flex gap-2">
              {['#2563eb', '#dc2626', '#16a34a', '#d97706', '#9333ea'].map((c) => (
                <button
                  key={c}
                  type="button"
                  onClick={() => setNewGoalColor(c)}
                  className={`w-6 h-6 rounded-full border-2 transition-all ${newGoalColor === c ? 'border-zinc-900 dark:border-zinc-100 scale-110' : 'border-transparent'}`}
                  style={{ backgroundColor: c }}
                />
              ))}
            </div>
          </div>
          <div className="flex justify-end gap-2 pt-4">
            <Button type="button" variant="secondary" onClick={() => setIsModalOpen(false)}>Cancel</Button>
            <Button type="submit" disabled={!newGoalTitle}>Create</Button>
          </div>
        </form>
      </Modal>
    </>
  );
};
