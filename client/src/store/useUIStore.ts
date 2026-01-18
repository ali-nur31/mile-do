import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface ContextMenuState {
  isOpen: boolean;
  x: number;
  y: number;
  type: 'task' | 'list' | null;
  targetId: number | null;
  data?: any;
}

interface AppState {
  isSidebarOpen: boolean;
  theme: 'light' | 'dark';
  selectedTaskId: number | null;
  contextMenu: ContextMenuState;
  isAuthenticated: boolean;
  calendarView: 'month' | 'week';

  toggleSidebar: () => void;
  toggleTheme: () => void;
  selectTask: (id: number | null) => void;
  openContextMenu: (e: React.MouseEvent | React.TouchEvent, type: 'task' | 'list', id: number, data?: any) => void;
  closeContextMenu: () => void;
  setAuthenticated: (status: boolean) => void;
  logout: () => void;
  setCalendarView: (view: 'month' | 'week') => void;
}

export const useStore = create<AppState>()(
  persist(
    (set) => ({
      isSidebarOpen: true,
      theme: 'light',
      selectedTaskId: null,
      isAuthenticated: !!localStorage.getItem('access_token'),
      contextMenu: { isOpen: false, x: 0, y: 0, type: null, targetId: null },
      calendarView: 'month',

      toggleSidebar: () => set((state) => ({ isSidebarOpen: !state.isSidebarOpen })),

      toggleTheme: () => set((state) => {
        const newTheme = state.theme === 'light' ? 'dark' : 'light';
        if (newTheme === 'dark') {
          document.documentElement.classList.add('dark');
        } else {
          document.documentElement.classList.remove('dark');
        }
        return { theme: newTheme };
      }),

      selectTask: (id) => set({ selectedTaskId: id }),

      openContextMenu: (e, type, targetId, data) => {
        e.preventDefault();
        let clientX, clientY;

        if ('touches' in e) {
           clientX = e.touches[0].clientX;
           clientY = e.touches[0].clientY;
        } else {
           clientX = (e as React.MouseEvent).clientX;
           clientY = (e as React.MouseEvent).clientY;
        }
        const menuWidth = 220;
        const menuHeight = 200;
        if (clientX + menuWidth > window.innerWidth) clientX = window.innerWidth - menuWidth - 10;
        if (clientY + menuHeight > window.innerHeight) clientY = window.innerHeight - menuHeight - 10;

        set({ contextMenu: { isOpen: true, x: clientX, y: clientY, type, targetId, data } });
      },

      closeContextMenu: () => set({ contextMenu: { isOpen: false, x: 0, y: 0, type: null, targetId: null } }),

      setAuthenticated: (status) => set({ isAuthenticated: status }),

      logout: () => {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        set({ isAuthenticated: false, selectedTaskId: null });
      },
      
      setCalendarView: (view) => set({ calendarView: view }),
    }),
    {
      name: 'mile-do-storage',
      partialize: (state) => ({ 
        isSidebarOpen: state.isSidebarOpen, 
        theme: state.theme,
        calendarView: state.calendarView
      }),
    }
  )
);
