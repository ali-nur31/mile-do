import { useEffect, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useStore } from '../../store/useUIStore';
import { goalsApi } from '../../api/goals';
import type { Goal, Task } from '../../types';
import { X, Trash2, Calendar, CheckCircle2, Hash, Clock, ArrowRight, Minus, Plus } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { extractDateStr, extractTimeStr, combineToBackend, addMinutesToTime, calculateDuration } from '../../utils/date';
import { useTasks } from '../../hooks/useTasks';
import { Modal } from '../../components/ui/Modal';
import { Button } from '../../components/ui/Button';

export const RightPanel = () => {
  const { selectedTaskId, selectTask } = useStore();
  const { updateTask, deleteTask, tasks } = useTasks();

  const [title, setTitle] = useState('');
  const [editDate, setEditDate] = useState('');
  const [editTime, setEditTime] = useState('');
  const [duration, setDuration] = useState<number>(15);
  const [endTime, setEndTime] = useState('');
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);

  const task = tasks.find((t: Task) => t.id === selectedTaskId);
  const { data: goals } = useQuery({ queryKey: ['goals'], queryFn: goalsApi.getAll });

  const assignedGoal = goals?.find((g: Goal) => g.id === task?.goal_id);
  const displayGoalColor = assignedGoal?.color || "#71717a";

  useEffect(() => {
    if (task) {
      setTitle(task.title);
      const dateStr = extractDateStr(task.scheduled_date);
      const timeStr = extractTimeStr(task.scheduled_time);
      const taskDuration = task.duration_minutes || 15;

      setEditDate(dateStr);
      setEditTime(timeStr || "09:00");
      setDuration(taskDuration);
      if (timeStr) setEndTime(addMinutesToTime(timeStr, taskDuration));
      else setEndTime(addMinutesToTime("09:00", taskDuration));
    }
  }, [task]);

  const handleTitleBlur = () => {
    if (task && title !== task.title && title.trim()) {
      updateTask.mutate({ id: task.id, data: { title } });
    } else {
      setTitle(task?.title || '');
    }
  };

  const handleDateChange = (newDate: string) => {
    setEditDate(newDate);
    if (task) {
      const finalTime = editTime && editTime.trim() ? editTime : '09:00';
      const scheduledDateTime = newDate ? combineToBackend(newDate, finalTime) : undefined;
      if (scheduledDateTime) {
        updateTask.mutate({ id: task.id, data: { scheduled_date_time: scheduledDateTime } });
      }
    }
  };

  const handleTimeChange = (newTime: string) => {
    setEditTime(newTime);
    if (task && editDate) {
      setEndTime(addMinutesToTime(newTime, duration));
      updateTask.mutate({ id: task.id, data: { scheduled_date_time: combineToBackend(editDate, newTime) } });
    }
  };



  const handleEndTimeChange = (newEndTime: string) => {
    setEndTime(newEndTime);
    if (editTime && newEndTime) {
      const newDuration = calculateDuration(editTime, newEndTime);
      if (newDuration >= 15) {
        setDuration(newDuration);
        if (task) updateTask.mutate({ id: task.id, data: { duration_minutes: newDuration } });
      }
    }
  };

  const handleDelete = () => {
    if (selectedTaskId) {
      deleteTask.mutate(selectedTaskId);
      selectTask(null);
      setIsDeleteModalOpen(false);
    }
  };

  if (!selectedTaskId || !task) return null;

  return (
    <>
      <AnimatePresence>
        <motion.aside
          initial={{ x: 20, opacity: 0 }}
          animate={{ x: 0, opacity: 1 }}
          exit={{ x: 20, opacity: 0 }}
          className="w-full h-full flex flex-col bg-white dark:bg-zinc-900 border-l border-zinc-200 dark:border-zinc-800"
        >
          <div className="h-14 border-b border-zinc-100 dark:border-zinc-800 flex items-center justify-between px-4">
            <div className="flex items-center gap-2 text-zinc-400 text-xs font-medium">
              <CheckCircle2 size={14} className={task.is_done ? 'text-blue-600' : ''} />
              {task.is_done ? 'Completed' : 'In Progress'}
            </div>
            <div className="flex gap-1">
              <button
                onClick={() => setIsDeleteModalOpen(true)}
                className="p-2 hover:bg-red-50 dark:hover:bg-red-900/20 text-zinc-400 dark:text-zinc-500 hover:text-red-600 dark:hover:text-red-400 rounded-md transition-colors"
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

          <div className="flex-1 overflow-y-auto p-6 custom-scrollbar space-y-6">
            <textarea
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              onBlur={handleTitleBlur}
              className="w-full text-xl font-bold bg-transparent resize-none outline-none text-zinc-900 dark:text-zinc-100 placeholder:text-zinc-300"
              rows={2}
            />
            <div className="space-y-4">
              <div className="flex items-center gap-3 text-sm text-zinc-600 dark:text-zinc-400">
                <div className="w-8 h-8 rounded-md bg-zinc-50 dark:bg-zinc-800 flex items-center justify-center">
                  <Hash size={16} style={{ color: displayGoalColor }} />
                </div>
                <div className="flex flex-col flex-1 min-w-0">
                  <span className="text-[10px] uppercase tracking-wide opacity-70">List</span>
                  <select
                    value={task.goal_id}
                    onChange={(e) => updateTask.mutate({ id: task.id, data: { goal_id: parseInt(e.target.value) } })}
                    className="w-full bg-transparent border-none p-0 text-sm font-medium text-zinc-900 dark:text-zinc-200 outline-none truncate cursor-pointer"
                  >
                    {goals?.filter(g => g.title !== 'Inbox').map(goal => (
                      <option key={goal.id} value={goal.id} className="dark:bg-zinc-900">
                        {goal.title}
                      </option>
                    ))}
                  </select>
                </div>
              </div>
              <div className="flex items-center gap-3 text-sm text-zinc-600 dark:text-zinc-400">
                <div className="w-8 h-8 rounded-md bg-zinc-50 dark:bg-zinc-800 flex items-center justify-center">
                  <Calendar size={16} />
                </div>
                <div className="flex-1 flex flex-col">
                  <span className="text-[10px] uppercase tracking-wide opacity-70">Due Date</span>
                  <input
                    type="date"
                    value={editDate}
                    onChange={(e) => handleDateChange(e.target.value)}
                    className="bg-transparent border-none p-0 text-sm font-medium text-zinc-900 dark:text-zinc-200 outline-none cursor-pointer"
                  />
                </div>
              </div>
              <div className="flex flex-col gap-4">
                <div className="space-y-2">
                  <span className="text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">Time</span>
                  <div className="flex items-center gap-2 p-2 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl border border-zinc-100 dark:border-zinc-800">
                    <div className="text-zinc-400 dark:text-zinc-500">
                      <Clock size={16} />
                    </div>
                    <div className="flex items-center gap-2 flex-1 justify-center">
                      <input
                        type="time"
                        value={editTime}
                        onChange={(e) => handleTimeChange(e.target.value)}
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

                <div className="space-y-2">
                  <span className="text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">Duration</span>
                  <div className="flex items-center justify-between p-1 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl border border-zinc-100 dark:border-zinc-800 h-[38px] w-full">
                    <button
                      onClick={() => {
                        const newDuration = Math.max(15, duration - 15);
                        setDuration(newDuration);
                        if (editTime) setEndTime(addMinutesToTime(editTime, newDuration));
                        if (task) updateTask.mutate({ id: task.id, data: { duration_minutes: newDuration } });
                      }}
                      className="w-10 h-full flex items-center justify-center text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300 hover:bg-white dark:hover:bg-zinc-700 rounded-lg transition-all"
                    >
                      <Minus size={14} />
                    </button>
                    <div className="flex-1 text-center text-sm font-semibold text-zinc-900 dark:text-zinc-200 tabular-nums">
                      {duration}m
                    </div>
                    <button
                      onClick={() => {
                        const newDuration = duration + 15;
                        setDuration(newDuration);
                        if (editTime) setEndTime(addMinutesToTime(editTime, newDuration));
                        if (task) updateTask.mutate({ id: task.id, data: { duration_minutes: newDuration } });
                      }}
                      className="w-10 h-full flex items-center justify-center text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300 hover:bg-white dark:hover:bg-zinc-700 rounded-lg transition-all"
                    >
                      <Plus size={14} />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </motion.aside>
      </AnimatePresence>

      <Modal isOpen={isDeleteModalOpen} onClose={() => setIsDeleteModalOpen(false)} title="Delete Task">
        <div className="space-y-4">
          <p className="text-sm text-zinc-600 dark:text-zinc-400">
            Are you sure you want to delete <strong>{task.title}</strong>?
          </p>
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="secondary" onClick={() => setIsDeleteModalOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleDelete} className="bg-red-600 hover:bg-red-700 text-white">
              <Trash2 size={16} className="mr-2" />
              Delete Task
            </Button>
          </div>
        </div>
      </Modal>
    </>
  );
};