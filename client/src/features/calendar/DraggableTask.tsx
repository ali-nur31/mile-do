import { useDraggable } from '@dnd-kit/core';
import { GripVertical } from 'lucide-react';
import { clsx } from 'clsx';
import { useStore } from '../../store/useUIStore';
import type { Task } from '../../types';

interface DraggableTaskProps {
  task: Task;
  isOverlay?: boolean;
  onClick?: (task: Task) => void;
}

export const DraggableTask = ({ task, isOverlay = false, onClick }: DraggableTaskProps) => {
  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({
    id: `task-${task.id}`,
    data: { task }
  });

  const { selectTask } = useStore();

  const handleClick = () => {
    if (isOverlay) return;
    if (onClick) {
      onClick(task);
    } else {
      selectTask(task.id);
    }
  };

  return (
    <div
      ref={setNodeRef}
      {...listeners}
      {...attributes}
      onClick={handleClick}
      className={clsx(
        "p-1.5 rounded border shadow-sm text-[10px] cursor-grab active:cursor-grabbing mb-1 flex items-center gap-1.5 transition-all select-none touch-none",
        isDragging ? "opacity-30" : "opacity-100",
        task.is_done
          ? "bg-zinc-100 dark:bg-zinc-800/50 border-zinc-200 dark:border-zinc-700/50 opacity-60 text-zinc-400"
          : "bg-white dark:bg-zinc-800 border-zinc-200 dark:border-zinc-700 text-zinc-700 dark:text-zinc-200 hover:border-blue-400 dark:hover:border-blue-500",
        isOverlay && "bg-blue-600 text-white border-blue-700 shadow-xl scale-105 z-50 cursor-grabbing !opacity-100"
      )}
    >
      {!task.is_done && <GripVertical size={10} className={clsx(isOverlay ? "text-blue-200" : "text-zinc-300 dark:text-zinc-600")} />}
      <span className={clsx("truncate font-medium", task.is_done && "line-through")}>
        {task.title || "Untitled"}
      </span>
    </div>
  );
};
