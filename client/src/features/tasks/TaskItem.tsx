import { memo } from 'react';
import type { Task } from '../../types';
import { Check, Calendar, MoreVertical } from 'lucide-react';
import { motion } from 'framer-motion';
import { useStore } from '../../store/useUIStore';
import { extractTimeFromBackend } from '../../utils/date';

interface TaskItemProps {
  task: Task;
  onToggle: (id: number, isDone: boolean) => void;
}

export const TaskItem = memo(({ task, onToggle }: TaskItemProps) => {
  const { selectTask, selectedTaskId, openContextMenu } = useStore();
  
  const isDateValid = task.scheduled_date && !task.scheduled_date.startsWith('0001');
  const dateObj = isDateValid ? new Date(task.scheduled_date) : null;
  const isSelected = selectedTaskId === task.id;

  const getDateLabel = (date: Date) => {
    const today = new Date();
    const isToday = date.toDateString() === today.toDateString();
    const tomorrow = new Date(today);
    tomorrow.setDate(tomorrow.getDate() + 1);
    const isTomorrow = date.toDateString() === tomorrow.toDateString();
    const baseDate = date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });

    if (isToday) return 'Today';
    if (isTomorrow) return 'Tomorrow';

    const timeStr = extractTimeFromBackend(task.scheduled_time);
    if (timeStr && timeStr !== "00:00") {
      return `${baseDate} • ${timeStr}`;
    }
    return baseDate;
  };

  const dateLabel = dateObj ? getDateLabel(dateObj) : '';
  const isOverdue = dateObj && dateObj < new Date() && !task.is_done && dateLabel !== 'Today';

  return (
    <div 
      onClick={() => selectTask(task.id)}
      onContextMenu={(e) => openContextMenu(e, 'task', task.id)}
      className={`
        group flex items-start gap-3 py-3.5 px-4 rounded-xl border transition-all cursor-pointer w-full
        ${isSelected 
          ? 'bg-blue-50 dark:bg-blue-500/10 border-blue-200 dark:border-blue-500/20' 
          : 'bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 hover:border-zinc-300 dark:hover:border-zinc-700 shadow-sm'}
      `}
    >
      <button
        onClick={(e) => {
          e.stopPropagation();
          onToggle(task.id, task.is_done);
        }}
        className={`
          flex-shrink-0 mt-0.5 w-5 h-5 rounded-[6px] border flex items-center justify-center transition-all duration-200
          ${task.is_done 
            ? 'bg-zinc-300 dark:bg-zinc-600 border-zinc-300 dark:border-zinc-600' 
            : `bg-transparent border-zinc-400 dark:border-zinc-500 hover:border-blue-500 dark:hover:border-blue-400`} 
        `}
      >
        <motion.div
          initial={false}
          animate={{ scale: task.is_done ? 1 : 0 }}
          transition={{ type: "spring", stiffness: 500, damping: 30 }}
        >
          <Check size={12} className="text-white dark:text-zinc-950 stroke-[3px]" />
        </motion.div>
      </button>
      
      <div className="flex-1 min-w-0 flex flex-col">
        <span className={`text-[15px] leading-snug transition-all ${task.is_done ? 'text-zinc-400 dark:text-zinc-600 line-through' : 'text-zinc-800 dark:text-zinc-100 font-medium'}`}>
          {task.title || "Untitled Task"}
        </span>
        
        {isDateValid && (
          <div className="flex items-center gap-2 mt-1.5 h-4">
            <span className={`
              text-[11px] flex items-center gap-1 font-medium rounded px-1.5 py-0.5 -ml-1.5
              ${task.is_done 
                ? 'text-zinc-400 dark:text-zinc-600' 
                : isOverdue 
                  ? 'text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20' 
                  : 'text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20'}
            `}>
              <Calendar size={11} className="mb-[1px]" />
              {dateLabel}
            </span>
          </div>
        )}
      </div>

      <button 
        className="md:hidden p-1 text-zinc-300 dark:text-zinc-600 hover:text-zinc-600 dark:hover:text-zinc-300"
        onClick={(e) => {
          e.stopPropagation();
          openContextMenu(e, 'task', task.id);
        }}
      >
        <MoreVertical size={16} />
      </button>
    </div>
  );
});