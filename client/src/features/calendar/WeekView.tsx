import { startOfWeek, endOfWeek, eachDayOfInterval, format, isSameDay } from 'date-fns';
import { useDroppable } from '@dnd-kit/core';
import { clsx } from 'clsx';
import { Plus } from 'lucide-react';
import { getLocalISOString } from '../../utils/date';
import type { Task } from '../../types';
import { DraggableTask } from './DraggableTask';

interface WeekViewProps {
  currentDate: Date;
  tasks: Task[];
  onAddClick: (date: Date) => void;
  onTaskClick: (task: Task) => void;
}

const WeekDayColumn = ({
  date,
  tasks,
  onAddClick,
  onTaskClick
}: {
  date: Date;
  tasks: Task[];
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
        "flex-1 min-w-0 flex flex-col h-full border-r last:border-r-0 border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-950/50 transition-colors group",
        isOver && "bg-blue-50 dark:bg-blue-900/20 ring-inset ring-2 ring-blue-400"
      )}
    >
      <div className={clsx(
        "p-2 sm:p-3 border-b border-zinc-200 dark:border-zinc-800 flex flex-col items-center gap-0.5 sm:gap-1",
        isToday ? "bg-blue-50/50 dark:bg-blue-900/10" : ""
      )}>
        <span className="text-[10px] sm:text-xs uppercase font-medium text-zinc-500 dark:text-zinc-400">
          {format(date, 'EEE')}
        </span>
        <div className={clsx(
          "w-6 h-6 sm:w-8 sm:h-8 rounded-full flex items-center justify-center text-xs sm:text-sm font-bold transition-colors",
          isToday ? "bg-blue-600 text-white" : "text-zinc-900 dark:text-zinc-100 group-hover:bg-zinc-100 dark:group-hover:bg-zinc-800"
        )}>
          {format(date, 'd')}
        </div>
      </div>

      <div
        className="flex-1 overflow-y-auto custom-scrollbar p-1.5 sm:p-2 space-y-1.5 sm:space-y-2 relative"
        onClick={() => { }}
      >
        {tasks.map(task => (
          <DraggableTask
            key={task.id}
            task={task}
            onClick={onTaskClick}
          />
        ))}

        <button
          onClick={onAddClick}
          className="w-full py-1.5 sm:py-2 rounded border border-dashed border-zinc-200 dark:border-zinc-700 text-zinc-400 hover:text-blue-500 hover:border-blue-300 hover:bg-blue-50 dark:hover:bg-blue-900/10 transition-all flex items-center justify-center gap-1 opacity-0 group-hover:opacity-100 touch-manipulation"
        >
          <Plus size={12} className="sm:hidden" />
          <Plus size={14} className="hidden sm:block" />
          <span className="text-[10px] sm:text-xs">Add Task</span>
        </button>
      </div>
    </div>
  );
};

export const WeekView = ({ currentDate, tasks, onAddClick, onTaskClick }: WeekViewProps) => {
  const startDate = startOfWeek(currentDate);
  const endDate = endOfWeek(startDate);
  const days = eachDayOfInterval({ start: startDate, end: endDate });

  return (
    <div className="h-full flex flex-col">
      <div className="flex-1 flex overflow-x-auto">
        {days.map(day => (
          <WeekDayColumn
            key={day.toISOString()}
            date={day}
            tasks={tasks.filter(t => isSameDay(new Date(t.scheduled_date), day))}
            onAddClick={(e) => { e.stopPropagation(); onAddClick(day); }}
            onTaskClick={onTaskClick}
          />
        ))}
      </div>
    </div>
  );
};
