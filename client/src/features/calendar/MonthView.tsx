import { startOfMonth, endOfMonth, startOfWeek, endOfWeek, eachDayOfInterval, format, isSameDay, isSameMonth } from 'date-fns';
import { useDroppable } from '@dnd-kit/core';
import { clsx } from 'clsx';
import { Plus } from 'lucide-react';
import { getLocalISOString } from '../../utils/date';
import type { Task } from '../../types';
import { DraggableTask } from './DraggableTask';

interface MonthViewProps {
  currentDate: Date;
  tasks: Task[];
  onAddClick: (date: Date) => void;
  onTaskClick: (task: Task) => void;
}

const CalendarDay = ({
  date,
  tasks,
  isCurrentMonth,
  onAddClick,
  onTaskClick
}: {
  date: Date;
  tasks: Task[];
  isCurrentMonth: boolean;
  onAddClick: (e: React.MouseEvent) => void;
  onTaskClick: (task: Task) => void;
}) => {
  const { setNodeRef, isOver } = useDroppable({
    id: getLocalISOString(date),
    data: { date }
  });

  const isToday = isSameDay(date, new Date());

  return (
    <div
      ref={setNodeRef}
      className={clsx(
        "flex-1 min-h-0 border-r border-b border-zinc-200 dark:border-zinc-800 transition-colors relative group flex flex-col h-full",
        !isCurrentMonth && "bg-zinc-50/50 dark:bg-zinc-900/30 text-zinc-400 dark:text-zinc-600",
        isOver && "bg-blue-50 dark:bg-blue-900/20 ring-2 ring-inset ring-blue-400"
      )}
    >
      <div className="flex items-center justify-between p-1 sm:p-1.5">
        <span className={clsx(
          "text-[10px] sm:text-xs font-semibold rounded-full w-5 h-5 sm:w-6 sm:h-6 flex items-center justify-center",
          isToday ? "bg-blue-600 text-white" : "text-zinc-500 dark:text-zinc-400"
        )}>
          {format(date, 'd')}
        </span>
        <button
          onClick={onAddClick}
          className="opacity-0 group-hover:opacity-100 p-0.5 rounded hover:bg-zinc-200 dark:hover:bg-zinc-700 text-zinc-400 hover:text-blue-500 transition-all touch-manipulation"
        >
          <Plus size={12} className="sm:hidden" />
          <Plus size={14} className="hidden sm:block" />
        </button>
      </div>

      <div className="flex-1 overflow-y-auto custom-scrollbar px-0.5 sm:px-1 pb-0.5 sm:pb-1 space-y-0.5 sm:space-y-1">
        {tasks.map(task => (
          <DraggableTask
            key={task.id}
            task={task}
            onClick={onTaskClick}
          />
        ))}
      </div>
    </div>
  );
};

export const MonthView = ({ currentDate, tasks, onAddClick, onTaskClick }: MonthViewProps) => {
  const monthStart = startOfMonth(currentDate);
  const monthEnd = endOfMonth(monthStart);
  const startDate = startOfWeek(monthStart);
  const endDate = endOfWeek(monthEnd);
  const days = eachDayOfInterval({ start: startDate, end: endDate });

  return (
    <div className="flex flex-col h-full">
      <div className="grid grid-cols-7 border-b border-zinc-200 dark:border-zinc-800">
        {['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'].map(day => (
          <div key={day} className="py-1.5 sm:py-2 text-center text-[10px] sm:text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
            <span className="hidden sm:inline">{day}</span>
            <span className="sm:hidden">{day[0]}</span>
          </div>
        ))}
      </div>

      <div className="flex-1 grid grid-cols-7 grid-rows-5 lg:grid-rows-6">
        {days.map(day => (
          <CalendarDay
            key={day.toISOString()}
            date={day}
            tasks={tasks.filter(t => isSameDay(new Date(t.scheduled_date), day))}
            isCurrentMonth={isSameMonth(day, monthStart)}
            onAddClick={(e) => { e.stopPropagation(); onAddClick(day); }}
            onTaskClick={onTaskClick}
          />
        ))}
      </div>
    </div>
  );
};
