import { useEffect, useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useStore } from '../../store/useUIStore';
import { tasksApi } from '../../api/tasks';
import { goalsApi } from '../../api/goals';
import type { Goal } from '../../types';
import { X, Trash2, Calendar, CheckCircle2, Hash, Clock } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { toast } from 'sonner';
import { extractDateStr, extractTimeStr, combineToBackend } from '../../utils/date';
import { useTasks } from '../../hooks/useTasks';

export const RightPanel = () => {
  const { selectedTaskId, selectTask } = useStore();
  const { updateTask, tasks } = useTasks();
  
  const queryClient = useQueryClient();
  const [title, setTitle] = useState('');
  const [editDate, setEditDate] = useState('');
  const [editTime, setEditTime] = useState('');

  const task = tasks.find(t => t.id === selectedTaskId);

  const { data: goals } = useQuery({
    queryKey: ['goals'],
    queryFn: goalsApi.getAll
  });

  const assignedGoal = goals?.find((g: Goal) => g.id === task?.goal_id);
  const displayGoalName = assignedGoal?.title || "Inbox";
  const displayGoalColor = assignedGoal?.color || "#71717a";

  useEffect(() => {
    if (task) {
      setTitle(task.title);
      
      const dateStr = extractDateStr(task.scheduled_date);
      const timeStr = extractTimeStr(task.scheduled_time);
      
      console.log('Task loaded in RightPanel:', {
        task_id: task.id,
        raw_scheduled_date: task.scheduled_date,
        raw_scheduled_time: task.scheduled_time,
        extracted_date: dateStr,
        extracted_time: timeStr
      });
      
      setEditDate(dateStr);
      setEditTime(timeStr || "09:00");
    }
  }, [task]);

  const deleteTask = useMutation({
    mutationFn: async () => {
      if (!selectedTaskId) return;
      return tasksApi.delete(selectedTaskId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      selectTask(null);
      toast.success('Task deleted');
    }
  });

  const handleTitleBlur = () => {
    if (task && title !== task.title && title.trim() !== '') {
      updateTask.mutate({ id: task.id, data: { title } });
    } else if (title.trim() === '') {
      setTitle(task?.title || '');
    }
  };

  const handleDateChange = (newDate: string) => {
    setEditDate(newDate);
    
    if (task) {
      if (newDate && newDate.trim() !== '') {
        const finalTime = editTime && editTime.trim() !== '' ? editTime : '09:00';
        const scheduledDateTime = combineToBackend(newDate, finalTime);
        updateTask.mutate({ 
          id: task.id, 
          data: { scheduled_date_time: scheduledDateTime } 
        });
      } else {
        updateTask.mutate({ 
          id: task.id, 
          data: { scheduled_date_time: "0001-01-01 00:00:00" } 
        });
      }
    }
  };

  const handleTimeChange = (newTime: string) => {
    setEditTime(newTime);
    
    if (task && editDate && editDate.trim() !== '') {
      const finalTime = newTime && newTime.trim() !== '' ? newTime : '09:00';
      const scheduledDateTime = combineToBackend(editDate, finalTime);
      updateTask.mutate({ 
        id: task.id, 
        data: { scheduled_date_time: scheduledDateTime } 
      });
    }
  };

  const handleDelete = () => {
    deleteTask.mutate();
  };

  if (!selectedTaskId || !task) return null;

  return (
    <AnimatePresence>
      <motion.aside
        initial={{ x: 20, opacity: 0 }}
        animate={{ x: 0, opacity: 1 }}
        exit={{ x: 20, opacity: 0 }}
        transition={{ duration: 0.2, ease: "easeOut" }}
        className="w-full h-full flex flex-col bg-white dark:bg-zinc-900 transition-colors duration-200"
      >
        <div className="h-14 border-b border-zinc-100 dark:border-zinc-800 flex items-center justify-between px-4 flex-shrink-0">
          <div className="flex items-center gap-2 text-zinc-400 dark:text-zinc-500 text-xs font-medium">
            <CheckCircle2 size={14} className={task.is_done ? 'text-blue-600 dark:text-blue-500' : ''} />
            {task.is_done ? 'Completed' : 'In Progress'}
          </div>
          <div className="flex items-center gap-1">
            <button 
              onClick={handleDelete}
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

        <div className="flex-1 overflow-y-auto p-6 custom-scrollbar">
          <textarea
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            onBlur={handleTitleBlur}
            className="w-full text-xl font-bold text-zinc-900 dark:text-zinc-100 bg-transparent resize-none outline-none placeholder:text-zinc-300 dark:placeholder:text-zinc-700 min-h-[3rem]"
            placeholder="Task title"
            rows={2}
          />

          <div className="mt-4 space-y-4">
            <div className="flex items-center gap-3 text-sm text-zinc-600 dark:text-zinc-400">
              <div className="w-8 h-8 rounded-md bg-zinc-50 dark:bg-zinc-800 flex items-center justify-center">
                <Hash size={16} style={{ color: displayGoalColor }} />
              </div>
              <div className="flex flex-col">
                <span className="text-[10px] text-zinc-400 dark:text-zinc-500 uppercase tracking-wide">List</span>
                <span className="font-medium text-zinc-900 dark:text-zinc-200">{displayGoalName}</span>
              </div>
            </div>

            <div className="flex items-center gap-3 text-sm text-zinc-600 dark:text-zinc-400">
              <div className="w-8 h-8 rounded-md bg-zinc-50 dark:bg-zinc-800 flex items-center justify-center text-zinc-400 dark:text-zinc-500">
                <Calendar size={16} />
              </div>
              <div className="flex-1 flex flex-col">
                <span className="text-[10px] text-zinc-400 dark:text-zinc-500 uppercase tracking-wide">Due Date</span>
                <input 
                  type="date" 
                  value={editDate}
                  onChange={(e) => handleDateChange(e.target.value)}
                  className="bg-transparent border-none p-0 text-sm font-medium text-zinc-900 dark:text-zinc-200 focus:ring-0 cursor-pointer w-full"
                />
              </div>
            </div>

            <div className="flex items-center gap-3 text-sm text-zinc-600 dark:text-zinc-400">
              <div className="w-8 h-8 rounded-md bg-zinc-50 dark:bg-zinc-800 flex items-center justify-center text-zinc-400 dark:text-zinc-500">
                <Clock size={16} />
              </div>
              <div className="flex-1 flex flex-col">
                <span className="text-[10px] text-zinc-400 dark:text-zinc-500 uppercase tracking-wide">Time</span>
                <input 
                  type="time" 
                  value={editTime}
                  onChange={(e) => handleTimeChange(e.target.value)}
                  className="bg-transparent border-none p-0 text-sm font-medium text-zinc-900 dark:text-zinc-200 focus:ring-0 cursor-pointer w-full"
                />
              </div>
            </div>
          </div>
        </div>
      </motion.aside>
    </AnimatePresence>
  );
};