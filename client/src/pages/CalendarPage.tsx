import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  DndContext, 
  useDraggable, 
  useDroppable, 
  useSensor, 
  useSensors, 
  MouseSensor,
  TouchSensor,
  DragOverlay,
  defaultDropAnimationSideEffects
} from '@dnd-kit/core';
import type { DragEndEvent, DragStartEvent, DropAnimation } from '@dnd-kit/core';
import { startOfMonth, endOfMonth, startOfWeek, endOfWeek, eachDayOfInterval, format, isSameDay, isSameMonth } from 'date-fns';
import { api } from '../api/axios';
import type { ListTasksResponse, Task } from '../types';
import { TaskForm } from '../features/tasks/TaskForm';
import { Plus, GripVertical, Loader2, CalendarX, Check } from 'lucide-react';
import { clsx } from 'clsx';
import { safeDate, getLocalISOString, combineToBackend } from '../utils/date';
import { useStore } from '../store/useUIStore';
import { showToast } from '../utils/toast';
import { Button } from '../components/ui/Button';
import { Modal } from '../components/ui/Modal';
import { useTasks } from '../hooks/useTasks';

const DraggableTask = ({ task, isOverlay = false }: { task: Task; isOverlay?: boolean }) => {
  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({
    id: `task-${task.id}`,
    data: { task }
  });
  const { selectTask } = useStore();

  return (
    <div
      ref={setNodeRef}
      {...listeners}
      {...attributes}
      onClick={() => !isOverlay && selectTask(task.id)}
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
      {task.is_done && <Check size={10} className="text-zinc-400" />}
      <span className={clsx("truncate font-medium", task.is_done && "line-through")}>
        {task.title || "Untitled"}
      </span>
    </div>
  );
};

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
      <div className="p-1.5 flex justify-between items-start flex-shrink-0">
        <span className={clsx(
          "text-xs font-semibold rounded-full w-6 h-6 flex items-center justify-center",
          isToday 
            ? "bg-blue-600 text-white" 
            : "text-zinc-500 dark:text-zinc-400"
        )}>
          {format(date, 'd')}
        </span>
        <button 
          onClick={onAddClick}
          className="opacity-0 group-hover:opacity-100 p-0.5 rounded hover:bg-zinc-200 dark:hover:bg-zinc-700 text-zinc-400 hover:text-blue-500 transition-all"
        >
          <Plus size={14} />
        </button>
      </div>
      
      <div className="flex-1 overflow-y-auto px-1 pb-1 min-h-0 custom-scrollbar">
        <div className="flex flex-col gap-0.5">
          {tasks.map(task => (
            <div 
              key={task.id} 
              onClick={(e) => { e.stopPropagation(); onTaskClick(task); }}
              className={clsx(
                "text-[10px] truncate px-1.5 py-0.5 rounded border cursor-pointer hover:opacity-80 transition-opacity",
                task.is_done 
                  ? "bg-zinc-100 dark:bg-zinc-800 border-zinc-200 dark:border-zinc-700 text-zinc-400 line-through" 
                  : "bg-blue-50 dark:bg-blue-900/20 border-blue-100 dark:border-blue-800 text-blue-700 dark:text-blue-300"
              )}
            >
              {task.title || "Untitled"}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export const CalendarPage = () => {
  const { selectTask } = useStore();
  const queryClient = useQueryClient();
  const { updateTask } = useTasks();
  
  const [currentDate] = useState(new Date());
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isClearModalOpen, setIsClearModalOpen] = useState(false);
  const [selectedDate, setSelectedDate] = useState<Date | undefined>();
  const [activeDragTask, setActiveDragTask] = useState<Task | null>(null);

  const sensors = useSensors(
    useSensor(MouseSensor, { activationConstraint: { distance: 5 } }),
    useSensor(TouchSensor, { activationConstraint: { delay: 100, tolerance: 5 } })
  );

  const { data, isLoading } = useQuery({
    queryKey: ['tasks', 'all'],
    queryFn: async () => {
      const res = await api.get<ListTasksResponse>('/tasks/');
      return res.data;
    }
  });

  const clearCalendar = useMutation({
    mutationFn: async (tasksToClear: Task[]) => {
      await Promise.all(tasksToClear.map(t => 
        updateTask.mutateAsync({ 
          id: t.id, 
          data: { scheduled_date_time: "0001-01-01 00:00:00" } 
        })
      ));
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      setIsClearModalOpen(false);
      showToast('success', "Calendar cleared.");
    }
  });

  const handleDragStart = (event: DragStartEvent) => {
    setActiveDragTask(event.active.data.current?.task || null);
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    const task = active.data.current?.task as Task;
    setActiveDragTask(null);

    if (over && task) {
      const newDateStr = over.id as string;
      const scheduledDateTime = combineToBackend(newDateStr, "09:00");
      updateTask.mutate({ 
        id: task.id, 
        data: { scheduled_date_time: scheduledDateTime } 
      });
    }
  };

  const openCreateForm = (date: Date) => {
    setSelectedDate(date);
    setIsFormOpen(true);
  };

  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <Loader2 className="animate-spin text-zinc-400" />
      </div>
    );
  }

  const allTasks = data?.task_data || [];
  const backlogTasks = allTasks.filter(t => !safeDate(t.scheduled_date) && !t.is_done);
  const scheduledTasks = allTasks.filter(t => safeDate(t.scheduled_date));

  const monthStart = startOfMonth(currentDate);
  const monthEnd = endOfMonth(monthStart);
  const startDate = startOfWeek(monthStart);
  const endDate = endOfWeek(monthEnd);
  const days = eachDayOfInterval({ start: startDate, end: endDate });

  const dropAnimation: DropAnimation = {
    sideEffects: defaultDropAnimationSideEffects({ styles: { active: { opacity: '0.5' } } }),
  };

  return (
    <DndContext sensors={sensors} onDragStart={handleDragStart} onDragEnd={handleDragEnd}>
      <div className="flex flex-col lg:flex-row h-full w-full bg-white dark:bg-zinc-950 overflow-hidden">
        <div className="w-full lg:w-64 h-48 lg:h-full border-r border-b lg:border-b-0 border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-950 flex flex-col flex-shrink-0">
          <div className="p-4 border-b border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 flex justify-between items-center sticky top-0 z-10">
            <h2 className="font-semibold text-sm text-zinc-900 dark:text-zinc-100 uppercase tracking-wide">Backlog</h2>
            <span className="text-[10px] bg-zinc-100 dark:bg-zinc-800 text-zinc-500 px-2 py-0.5 rounded-full">{backlogTasks.length}</span>
          </div>
          <div className="flex-1 overflow-y-auto p-3 custom-scrollbar">
            {backlogTasks.map(task => (
              <DraggableTask key={task.id} task={task} />
            ))}
            {backlogTasks.length === 0 && (
              <div className="text-center py-10 text-zinc-400 text-xs italic">All unscheduled tasks done</div>
            )}
          </div>
        </div>

        <div className="flex-1 flex flex-col h-full bg-white dark:bg-zinc-900 overflow-hidden relative">
          <div className="p-4 flex items-center justify-between border-b border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 flex-shrink-0">
            <h2 className="text-xl font-bold text-zinc-900 dark:text-zinc-100">{format(currentDate, 'MMMM yyyy')}</h2>
            {scheduledTasks.length > 0 && (
              <button onClick={() => setIsClearModalOpen(true)} className="text-xs text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 px-2 py-1.5 rounded flex items-center gap-1">
                <CalendarX size={14} /> Clear
              </button>
            )}
          </div>

          <div className="flex-1 flex flex-col overflow-hidden relative">
            <div className="absolute inset-0 flex flex-col">
              <div className="w-full flex-1 flex flex-col h-full">
                <div className="grid grid-cols-7 border-b border-zinc-200 dark:border-zinc-800 bg-zinc-50/50 dark:bg-zinc-900/50 flex-shrink-0">
                  {['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'].map(day => (
                    <div key={day} className="py-2 text-center text-xs font-semibold text-zinc-400 dark:text-zinc-500 uppercase">{day}</div>
                  ))}
                </div>
                <div className="flex-1 grid grid-cols-7 auto-rows-fr bg-white dark:bg-zinc-950 min-h-0">
                  {days.map(day => (
                    <CalendarDay
                      key={day.toISOString()}
                      date={day}
                      tasks={scheduledTasks.filter(t => isSameDay(new Date(t.scheduled_date), day))}
                      isCurrentMonth={isSameMonth(day, monthStart)}
                      onAddClick={(e) => { e.stopPropagation(); openCreateForm(day); }}
                      onTaskClick={(task) => selectTask(task.id)}
                    />
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <DragOverlay dropAnimation={dropAnimation}>
        {activeDragTask ? <DraggableTask task={activeDragTask} isOverlay /> : null}
      </DragOverlay>

      <TaskForm isOpen={isFormOpen} onClose={() => setIsFormOpen(false)} initialDate={selectedDate} />

      <Modal isOpen={isClearModalOpen} onClose={() => setIsClearModalOpen(false)} title="Clear Calendar">
        <div className="space-y-4">
          <p className="text-sm text-zinc-600 dark:text-zinc-400">
            This will unschedule <strong>{scheduledTasks.length}</strong> tasks.
          </p>
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="secondary" onClick={() => setIsClearModalOpen(false)}>Cancel</Button>
            <Button onClick={() => clearCalendar.mutate(scheduledTasks)} className="bg-red-600 hover:bg-red-700 text-white">Unschedule All</Button>
          </div>
        </div>
      </Modal>
    </DndContext>
  );
};