import { useEffect } from 'react';
import { Outlet } from 'react-router-dom';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { RightPanel } from './RightPanel';
import { ContextMenu } from '../ui/ContextMenu';
import { useStore } from '../../store/useUIStore';

export const AppLayout = () => {
  const { isSidebarOpen, selectedTaskId, theme } = useStore();

  useEffect(() => {
    if (theme === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [theme]);

  return (
    <div className="min-h-screen flex flex-col bg-white dark:bg-zinc-950 text-zinc-900 dark:text-zinc-100 font-sans transition-colors duration-200">
      <Header />
      <div className="flex flex-1 pt-16 h-screen overflow-hidden">
        
        {/* sidebar */}
        <div 
          className={`
            flex-shrink-0 transition-all duration-300 ease-[cubic-bezier(0.25,1,0.5,1)] 
            ${isSidebarOpen ? 'w-[260px]' : 'w-0'} hidden md:block border-r border-zinc-200 dark:border-zinc-800 bg-zinc-50/50 dark:bg-zinc-950
          `}
        >
          <div className="w-[260px] h-full overflow-hidden">
             <Sidebar />
          </div>
        </div>
        
        <div className="md:hidden">
           <Sidebar />
        </div>
        
        {/* main */}
        <main className="flex-1 min-w-0 bg-white dark:bg-zinc-950 flex relative transition-colors duration-200">
          <div className="flex-1 overflow-y-auto custom-scrollbar">
            <div className="max-w-4xl mx-auto p-4 md:px-10 md:py-8 pb-32">
              <Outlet />
            </div>
          </div>

          <div className={`
            hidden lg:block h-full border-l border-zinc-200 dark:border-zinc-800 z-10 transition-all duration-300 ease-in-out bg-white dark:bg-zinc-900
            ${selectedTaskId ? 'w-[360px] opacity-100 translate-x-0' : 'w-0 opacity-0 translate-x-10 overflow-hidden'}
          `}>
             {selectedTaskId && <RightPanel />}
          </div>
          
          {selectedTaskId && (
            <div className="lg:hidden fixed inset-0 z-50 flex justify-end">
               <div className="absolute inset-0 bg-black/20 backdrop-blur-sm" />
               <div className="w-[350px] h-full shadow-2xl relative z-10">
                 <RightPanel />
               </div>
            </div>
          )}
        </main>
      </div>
      <ContextMenu />
    </div>
  );
};
