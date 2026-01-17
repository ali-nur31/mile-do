import React, { useState, useEffect, useCallback } from 'react';
import { useQuery } from '@tanstack/react-query';
import { goalsApi } from '../../api/goals';
import { Modal } from '../../components/ui/Modal';
import { Button } from '../../components/ui/Button';
import type { CreateTaskRequest } from '../../types';
import { Calendar as CalendarIcon, Clock, Hash, ArrowRight, Minus, Plus } from 'lucide-react';
import { getLocalISOString, combineToBackend, addMinutesToTime, calculateDuration } from '../../utils/date';
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
  const [startTime, setStartTime] = useState('09:00');
  const [endTime, setEndTime] = useState('09:30');
  const [duration, setDuration] = useState<number>(30);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const { data: goals } = useQuery({ queryKey: ['goals'], queryFn: goalsApi.getAll });

  useEffect(() => {
    if (isOpen) {
      setTitle('');
      setStartTime('09:00');
      setEndTime('09:30');
      setDuration(30);
      setDate(initialDate ? getLocalISOString(initialDate) : '');
      setIsSubmitting(false);

      if (goals && goals.length > 0) {
        const defaultGoal = goals.find(g => g.title === 'Other') || goals.find(g => g.title !== 'Inbox') || goals[0];
        if (defaultGoal) setSelectedGoalId(defaultGoal.id);
      }
    }
  }, [isOpen, initialDate, goals]);

  const handleStartTimeChange = (newTime: string) => {
    setStartTime(newTime);
    if (newTime) setEndTime(addMinutesToTime(newTime, duration));
  };



  const handleEndTimeChange = (newEndTime: string) => {
    setEndTime(newEndTime);
    if (startTime && newEndTime) {
      const newDur = calculateDuration(startTime, newEndTime);
      if (newDur >= 15) setDuration(newDur);
    }
  };

  const handleSubmit = useCallback((e: React.FormEvent) => {
    e.preventDefault();

    if (isSubmitting || !title.trim()) return;

    if (title.trim().length < 3) {
      showToast('delete', 'Task title too short (min 3 chars)');
      return;
    }

    setIsSubmitting(true);

    const payload: CreateTaskRequest = {
      title: title.trim(),
      goal_id: selectedGoalId,
      duration_minutes: duration
    };

    if (date && date.trim() !== '') {
      payload.scheduled_date_time = combineToBackend(date, startTime);
    }

    createTask.mutate(payload, {
      onSuccess: () => {
        showToast('success', 'Task Created');
        onClose();
      },
      onSettled: () => {
        setIsSubmitting(false);
      }
    });
  }, [isSubmitting, title, selectedGoalId, duration, date, startTime, createTask, onClose]);

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="New Task">
      <form onSubmit={handleSubmit} className="space-y-5">
        <input
          placeholder="What needs to be done?"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          autoFocus
          className="w-full text-lg font-medium bg-transparent border-none outline-none text-zinc-900 dark:text-zinc-100 placeholder:text-zinc-400"
        />

        <div className="space-y-4">
          <div className="flex flex-col gap-2">
            <label className="text-xs font-semibold text-zinc-500 uppercase tracking-wider flex items-center gap-1.5">
              <Hash size={12} /> Goal
            </label>
            <div className="flex flex-wrap gap-2">
              {goals?.filter(g => g.title !== 'Inbox').map((goal) => (
                <button
                  key={goal.id}
                  type="button"
                  onClick={() => setSelectedGoalId(goal.id)}
                  className={`px-3 py-1.5 rounded-lg text-xs font-medium border transition-all ${selectedGoalId === goal.id
                    ? 'bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900'
                    : 'bg-white dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 border-zinc-200 dark:border-zinc-700'
                    }`}
                >
                  {goal.title}
                </button>
              ))}
            </div>
          </div>

          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <label className="text-xs font-semibold text-zinc-500 uppercase tracking-wider flex items-center gap-1.5">
                <CalendarIcon size={12} /> Date
              </label>
              <input
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
                className="w-full px-3 py-2 bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-lg text-sm text-zinc-900 dark:text-zinc-100 focus:ring-2 focus:ring-blue-500 outline-none"
              />
            </div>

            <div className="flex flex-col gap-2">
              <span className="text-xs font-semibold text-zinc-500 uppercase tracking-wider">Time</span>
              <div className="flex items-center gap-2 p-2 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl border border-zinc-100 dark:border-zinc-800">
                <div className="text-zinc-400 dark:text-zinc-500">
                  <Clock size={16} />
                </div>
                <div className="flex items-center gap-2 flex-1 justify-center">
                  <input
                    type="time"
                    value={startTime}
                    onChange={(e) => handleStartTimeChange(e.target.value)}
                    className="bg-transparent border-none p-0 text-sm font-semibold text-zinc-900 dark:text-zinc-200 outline-none cursor-pointer w-[88px] text-center"
                  />
                  <ArrowRight size={12} className="text-zinc-300 dark:text-zinc-600" />
                  <input
                    type="time"
                    value={endTime}
                    onChange={(e) => handleEndTimeChange(e.target.value)}
                    className="bg-transparent border-none p-0 text-sm font-semibold text-zinc-900 dark:text-zinc-200 outline-none cursor-pointer w-[88px] text-center"
                  />
                </div>
              </div>
            </div>

            <div className="flex flex-col gap-2">
              <span className="text-xs font-semibold text-zinc-500 uppercase tracking-wider">Duration</span>
              <div className="flex items-center justify-between p-1 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl border border-zinc-100 dark:border-zinc-800 h-[38px] w-full">
                <button
                  type="button"
                  onClick={() => {
                    const newDuration = Math.max(15, duration - 15);
                    setDuration(newDuration);
                    if (startTime) setEndTime(addMinutesToTime(startTime, newDuration));
                  }}
                  className="w-10 h-full flex items-center justify-center text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300 hover:bg-white dark:hover:bg-zinc-700 rounded-lg transition-all"
                >
                  <Minus size={14} />
                </button>
                <div className="flex-1 text-center text-sm font-semibold text-zinc-900 dark:text-zinc-200 tabular-nums">
                  {duration}m
                </div>
                <button
                  type="button"
                  onClick={() => {
                    const newDuration = duration + 15;
                    setDuration(newDuration);
                    if (startTime) setEndTime(addMinutesToTime(startTime, newDuration));
                  }}
                  className="w-10 h-full flex items-center justify-center text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300 hover:bg-white dark:hover:bg-zinc-700 rounded-lg transition-all"
                >
                  <Plus size={14} />
                </button>
              </div>
            </div>
          </div>
        </div>

        <div className="flex justify-end pt-2">
          <Button type="submit" disabled={!title.trim() || isSubmitting}>
            {isSubmitting ? 'Saving...' : 'Save Task'}
          </Button>
        </div>
      </form>
    </Modal>
  );
};