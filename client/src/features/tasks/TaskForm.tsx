import React, { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { goalsApi } from '../../api/goals';
import { Modal } from '../../components/ui/Modal';
import { Button } from '../../components/ui/Button';
import type { CreateTaskRequest } from '../../types';
import { Calendar as CalendarIcon, Clock, Hash } from 'lucide-react';
import { getLocalISOString, combineToBackend } from '../../utils/date';
import { showToast } from '../../utils/toast';
import { useTasks } from '../../hooks/useTasks';

interface TaskFormProps {
  isOpen: boolean;
  onClose: () => void;
  initialDate?: Date;
}

export const TaskForm: React.FC<TaskFormProps> = ({ isOpen, onClose, initialDate }) => {
  const { createTask } = useTasks();
  const [title, setTitle] = useState('');
  const [selectedGoalId, setSelectedGoalId] = useState<number>(0);
  const [date, setDate] = useState('');
  const [time, setTime] = useState('09:00');

  const { data: goals } = useQuery({
    queryKey: ['goals'],
    queryFn: goalsApi.getAll
  });

  useEffect(() => {
    if (isOpen) {
      setTitle('');
      setSelectedGoalId(0);
      setTime('09:00');
      
      if (initialDate) {
        setDate(getLocalISOString(initialDate));
      } else {
        setDate('');
      }
    }
  }, [isOpen, initialDate]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    const payload: CreateTaskRequest = {
      title: title.trim(),
      goal_id: selectedGoalId
    };

    if (date && date.trim() !== '') {
      const finalTime = time && time.trim() !== '' ? time : '09:00';
      payload.scheduled_date_time = combineToBackend(date, finalTime);
    }

    createTask.mutate(payload, {
      onSuccess: () => {
        showToast('success', 'Task Created');
        onClose();
      },
      onError: (error) => {
        console.error('Failed to create task:', error);
        showToast('delete', 'Failed to create task');
      }
    });
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="New Task">
      <form onSubmit={handleSubmit} className="space-y-5">
        <div>
          <input
            placeholder="What needs to be done?"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            autoFocus
            className="w-full text-lg font-medium bg-transparent border-none outline-none placeholder:text-zinc-400 text-zinc-900 dark:text-zinc-100"
          />
        </div>

        <div className="space-y-4">
          <div className="flex flex-col gap-2">
            <label className="text-xs font-semibold text-zinc-500 uppercase tracking-wider flex items-center gap-1.5">
              <Hash size={12} /> Goal
            </label>
            <div className="flex flex-wrap gap-2">
              <button
                type="button"
                onClick={() => setSelectedGoalId(0)}
                className={`
                  px-3 py-1.5 rounded-lg text-xs font-medium border transition-all
                  ${selectedGoalId === 0 
                    ? 'bg-zinc-900 text-white border-zinc-900 dark:bg-zinc-100 dark:text-zinc-900 dark:border-zinc-100' 
                    : 'bg-white dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 border-zinc-200 dark:border-zinc-700 hover:border-zinc-400 dark:hover:border-zinc-500'}
                `}
              >
                Inbox
              </button>
              {goals?.map((goal) => (
                <button
                  key={goal.id}
                  type="button"
                  onClick={() => setSelectedGoalId(goal.id)}
                  className={`
                    px-3 py-1.5 rounded-lg text-xs font-medium border transition-all
                    ${selectedGoalId === goal.id 
                      ? 'bg-zinc-900 text-white border-zinc-900 dark:bg-zinc-100 dark:text-zinc-900 dark:border-zinc-100' 
                      : 'bg-white dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 border-zinc-200 dark:border-zinc-700 hover:border-zinc-400 dark:hover:border-zinc-500'}
                  `}
                >
                  {goal.title}
                </button>
              ))}
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="flex flex-col gap-2">
              <label className="text-xs font-semibold text-zinc-500 uppercase tracking-wider flex items-center gap-1.5">
                <CalendarIcon size={12} /> Date
              </label>
              <input
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
                className="w-full px-3 py-2 bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-lg text-sm text-zinc-900 dark:text-zinc-100 focus:ring-2 focus:ring-blue-500 outline-none transition-all"
              />
            </div>
            <div className="flex flex-col gap-2">
              <label className="text-xs font-semibold text-zinc-500 uppercase tracking-wider flex items-center gap-1.5">
                <Clock size={12} /> Time
              </label>
              <input
                type="time"
                value={time}
                onChange={(e) => setTime(e.target.value)}
                className="w-full px-3 py-2 bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-lg text-sm text-zinc-900 dark:text-zinc-100 focus:ring-2 focus:ring-blue-500 outline-none transition-all"
              />
            </div>
          </div>
        </div>

        <div className="flex justify-end pt-2">
          <Button type="submit" disabled={!title.trim() || createTask.isPending}>
            {createTask.isPending ? 'Saving...' : 'Save Task'}
          </Button>
        </div>
      </form>
    </Modal>
  );
};