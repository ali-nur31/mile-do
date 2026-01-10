import { CheckCircle2, Menu, LogOut, User, Moon, Sun } from 'lucide-react';
import { useStore } from '../../store/useUIStore';
import { api } from '../../api/axios';

export const Header = () => {
  const { toggleSidebar, logout, toggleTheme, theme } = useStore();

  const handleLogout = async () => {
    try {
      await api.delete('/auth/logout');
    } catch (error) {
      console.error(error);
    } finally {
      logout();
    }
  };

  return (
    <header className="fixed top-0 left-0 right-0 h-16 bg-white/80 dark:bg-zinc-900/90 backdrop-blur-md border-b border-zinc-200 dark:border-zinc-800 z-50 px-4 md:px-6 flex items-center justify-between transition-colors duration-200">
      <div className="flex items-center gap-4">
        <button 
          onClick={toggleSidebar}
          className="p-2 -ml-2 text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-lg transition-colors"
        >
          <Menu size={20} />
        </button>
        
        <div className="flex items-center gap-2">
          <CheckCircle2 className="text-blue-600 dark:text-blue-500" size={24} strokeWidth={2.5} />
          <span className="text-xl font-bold tracking-tight text-zinc-900 dark:text-zinc-100 hidden md:block">Mile-Do</span>
        </div>
      </div>

      <div className="flex items-center gap-2 md:gap-4">
        <button
          onClick={toggleTheme}
          className="p-2 text-zinc-500 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-lg transition-colors"
        >
          {theme === 'dark' ? <Sun size={20} /> : <Moon size={20} />}
        </button>

        <div className="h-8 w-8 rounded-full bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center text-zinc-400 dark:text-zinc-500 border border-zinc-200 dark:border-zinc-700">
          <User size={16} />
        </div>
        <button 
          onClick={handleLogout}
          className="flex items-center gap-2 text-sm font-medium text-zinc-500 dark:text-zinc-400 hover:text-red-600 dark:hover:text-red-400 transition-colors"
        >
          <LogOut size={16} />
          <span className="hidden md:inline">Sign Out</span>
        </button>
      </div>
    </header>
  );
};
