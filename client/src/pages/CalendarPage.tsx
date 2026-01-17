import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  DndContext,
  useSensor,
  useSensors,
  MouseSensor,
  TouchSensor,
  DragOverlay,
  defaultDropAnimationSideEffects
} from '@dnd-kit/core';
import type { DragEndEvent, DragStartEvent, DropAnimation } from '@dnd-kit/core';
import { format } from 'date-fns';
import { api } from '../api/axios';
import type { Task } from '../types';
import { TaskForm } from '../features/tasks/TaskForm';
import { Calendar as CalendarIcon, LayoutList, Loader2, Menu, X } from 'lucide-react';
import { safeDate, combineToBackend } from '../utils/date';
import { useStore } from '../store/useUIStore';
import { showToast } from '../utils/toast';
import { Button } from '../components/ui/Button';
import { Modal } from '../components/ui/Modal';
import { useTasks } from '../hooks/useTasks';
import { MonthView } from '../features/calendar/MonthView';
import { WeekView } from '../features/calendar/WeekView';
import { DraggableTask } from '../features/calendar/DraggableTask';
import { clsx } from 'clsx';

export const CalendarPage = () => {
  const { selectTask, calendarView, setCalendarView } = useStore();
  const queryClient = useQueryClient();
  const { updateTask } = useTasks();

  const [currentDate] = useState(new Date());
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isClearModalOpen, setIsClearModalOpen] = useState(false);
  const [selectedDate, setSelectedDate] = useState<Date | undefined>();
  const [activeDragTask, setActiveDragTask] = useState<Task | null>(null);
  const [isBacklogOpen, setIsBacklogOpen] = useState(false);

  const sensors = useSensors(
    useSensor(MouseSensor, { activationConstraint: { distance: 5 } }),
    useSensor(TouchSensor, { activationConstraint: { delay: 100, tolerance: 5 } })
  );

  const { data, isLoading } = useQuery({
    queryKey: ['tasks', 'all'],
    queryFn: async () => {
      const res = await api.get('/tasks/');
      return res.data;
    }
  });

  const clearCalendar = useMutation({
    mutationFn: async (tasksToClear: Task[]) => {
      for (const t of tasksToClear) {
        await updateTask.mutateAsync({
          id: t.id,
          data: { scheduled_date_time: "0001-01-01 00:00:00" }
        });
      }
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
      if (scheduledDateTime) {
        updateTask.mutate({
          id: task.id,
          data: { scheduled_date_time: scheduledDateTime }
        });
      }
    }
  };

  const openCreateForm = (date: Date) => {
    setSelectedDate(date);
    setIsFormOpen(true);
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="animate-spin text-zinc-400" />
      </div>
    );
  }

  const allTasks = data?.task_data || [];
  const backlogTasks = allTasks.filter((t: Task) => !safeDate(t.scheduled_date) && !t.is_done);
  const scheduledTasks = allTasks.filter((t: Task) => safeDate(t.scheduled_date));

  const dropAnimation: DropAnimation = {
    sideEffects: defaultDropAnimationSideEffects({ styles: { active: { opacity: '0.5' } } }),
  };

  return (
    <DndContext
      sensors={sensors}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
    >
      <div className="flex h-full bg-white dark:bg-zinc-950 relative">
        <div className={clsx(
          "w-64 border-r border-zinc-200 dark:border-zinc-800 flex flex-col bg-zinc-50 dark:bg-zinc-900/50 transition-all duration-300",
          "md:relative md:translate-x-0",
          "absolute inset-y-0 left-0 z-30",
          isBacklogOpen ? "translate-x-0" : "-translate-x-full md:translate-x-0"
        )}>
          <div className="p-4 border-b border-zinc-200 dark:border-zinc-800 flex items-center justify-between">
            <h2 className="font-semibold text-sm text-zinc-700 dark:text-zinc-200">Backlog</h2>
            <div className="flex items-center gap-2">
              <span className="text-xs bg-zinc-200 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 px-2 py-0.5 rounded-full">
                {backlogTasks.length}
              </span>
              <button
                onClick={() => setIsBacklogOpen(false)}
                className="md:hidden p-1 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded transition-colors"
              >
                <X size={16} />
              </button>
            </div>
          </div>
          <div className="flex-1 overflow-y-auto p-3 custom-scrollbar">
            {backlogTasks.map((task: Task) => (
              <DraggableTask key={task.id} task={task} />
            ))}
            {backlogTasks.length === 0 && (
              <div className="text-center text-xs text-zinc-400 mt-10">All unscheduled tasks done</div>
            )}
          </div>
        </div>

        {isBacklogOpen && (
          <div
            className="fixed inset-0 bg-black/20 z-20 md:hidden"
            onClick={() => setIsBacklogOpen(false)}
          />
        )}

        <div className="flex-1 flex flex-col min-w-0">
          <div className="h-14 border-b border-zinc-200 dark:border-zinc-800 flex items-center justify-between px-3 md:px-6 bg-white dark:bg-zinc-950 gap-2">
            <div className="flex items-center gap-2 md:gap-4 min-w-0 flex-1">
              <button
                onClick={() => setIsBacklogOpen(true)}
                className="md:hidden p-2 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-lg transition-colors flex-shrink-0"
              >
                <Menu size={18} />
              </button>
              <h1 className="text-sm md:text-lg font-bold text-zinc-800 dark:text-zinc-100 truncate">
                {format(currentDate, 'MMMM yyyy')}
              </h1>
              {scheduledTasks.length > 0 && (
                <button
                  onClick={() => setIsClearModalOpen(true)}
                  className="text-xs text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 px-2 py-1.5 rounded flex items-center gap-1 flex-shrink-0"
                >
                  Clear
                </button>
              )}
            </div>

            <div className="flex items-center bg-zinc-100 dark:bg-zinc-800 rounded-lg p-1 flex-shrink-0">
              <button
                onClick={() => setCalendarView('month')}
                className={clsx(
                  "px-2 md:px-3 py-1.5 rounded-md text-xs font-medium transition-all flex items-center gap-1 md:gap-2",
                  calendarView === 'month'
                    ? "bg-white dark:bg-zinc-700 text-zinc-900 dark:text-zinc-100 shadow-sm"
                    : "text-zinc-500 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200"
                )}
              >
                <CalendarIcon size={14} />
                <span className="hidden sm:inline">Month</span>
              </button>
              <button
                onClick={() => setCalendarView('week')}
                className={clsx(
                  "px-2 md:px-3 py-1.5 rounded-md text-xs font-medium transition-all flex items-center gap-1 md:gap-2",
                  calendarView === 'week'
                    ? "bg-white dark:bg-zinc-700 text-zinc-900 dark:text-zinc-100 shadow-sm"
                    : "text-zinc-500 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200"
                )}
              >
                <LayoutList size={14} />
                <span className="hidden sm:inline">Week</span>
              </button>
            </div>
          </div>

          <div className="flex-1 min-h-0 overflow-hidden">
            {calendarView === 'month' ? (
              <MonthView
                currentDate={currentDate}
                tasks={scheduledTasks}
                onAddClick={openCreateForm}
                onTaskClick={(t) => selectTask(t.id)}
              />
            ) : (
              <WeekView
                currentDate={currentDate}
                tasks={scheduledTasks}
                onAddClick={openCreateForm}
                onTaskClick={(t) => selectTask(t.id)}
              />
            )}
          </div>
        </div>
      </div>

      <DragOverlay dropAnimation={dropAnimation}>
        {activeDragTask ? <DraggableTask task={activeDragTask} isOverlay /> : null}
      </DragOverlay>

      <TaskForm isOpen={isFormOpen} onClose={() => setIsFormOpen(false)} initialDate={selectedDate} />

      <Modal isOpen={isClearModalOpen} onClose={() => setIsClearModalOpen(false)} title="Clear Calendar">
        <div className="p-4">
          <p className="text-zinc-600 dark:text-zinc-400 mb-6">
            This will unschedule {scheduledTasks.length} tasks.
          </p>
          <div className="flex justify-end gap-2">
            <Button variant="secondary" onClick={() => setIsClearModalOpen(false)}>Cancel</Button>
            <Button onClick={() => clearCalendar.mutate(scheduledTasks)} className="bg-red-600 hover:bg-red-700 text-white">
              Unschedule All
            </Button>
          </div>
        </div>
      </Modal>
    </DndContext>
  );
};