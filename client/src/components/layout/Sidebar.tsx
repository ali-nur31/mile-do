import { motion } from 'framer-motion';
import { Inbox, CalendarDays, Target, Hash, X, Plus, MoreVertical } from 'lucide-react';
import { NavLink, useNavigate } from 'react-router-dom';
import { useStore } from '../../store/useUIStore';
import { useGoals } from '../../hooks/useGoals';

const navItems = [
  { icon: Inbox, label: 'Inbox', path: '/' },
  { icon: CalendarDays, label: 'Today', path: '/today' },
  { icon: Target, label: 'Goals Overview', path: '/goals' },
];

export const Sidebar = () => {
  const { isSidebarOpen, toggleSidebar, openContextMenu } = useStore();
  const { goals } = useGoals();
  const navigate = useNavigate();

  const sidebarContent = (
    <div className="flex flex-col h-full bg-zinc-50/95 dark:bg-zinc-900/95 backdrop-blur-xl md:bg-zinc-50 dark:md:bg-zinc-900 border-r border-zinc-200 dark:border-zinc-800 transition-colors duration-200">
      <div className="flex items-center justify-between p-4 md:hidden">
        <span className="font-bold text-lg text-zinc-900 dark:text-zinc-100">Menu</span>
        <button onClick={toggleSidebar} className="p-2 text-zinc-500 dark:text-zinc-400">
          <X size={20} />
        </button>
      </div>

      <div className="flex-1 overflow-y-auto pt-2 md:pt-6 custom-scrollbar">
        <div className="px-3 mb-6">
          <nav className="space-y-0.5">
            {navItems.map((item) => (
              <NavLink
                key={item.path}
                to={item.path}
                onClick={() => window.innerWidth < 768 && toggleSidebar()}
                className={({ isActive }) => `
                  flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors
                  ${isActive 
                    ? 'bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400' 
                    : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 hover:text-zinc-900 dark:hover:text-zinc-200'}
                `}
              >
                <item.icon size={18} strokeWidth={2} />
                {item.label}
              </NavLink>
            ))}
          </nav>
        </div>

        <div className="px-3">
          <div className="flex items-center justify-between px-3 mb-2 group">
            <h3 className="text-xs font-semibold text-zinc-400 dark:text-zinc-500 uppercase tracking-wider">
              Lists
            </h3>
            <button 
              onClick={() => navigate('/goals')}
              className="text-zinc-400 dark:text-zinc-500 hover:text-blue-600 dark:hover:text-blue-400 opacity-0 group-hover:opacity-100 transition-all"
            >
              <Plus size={14} />
            </button>
          </div>
          
          <div className="space-y-0.5">
            {goals?.map((goal) => (
              <NavLink
                key={goal.id}
                to={`/goals/${goal.id}`}
                onContextMenu={(e) => openContextMenu(e, 'list', goal.id, { title: goal.title })}
                onClick={() => window.innerWidth < 768 && toggleSidebar()}
                className={({ isActive }) => `
                  flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors group
                  ${isActive 
                    ? 'bg-zinc-100 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-100' 
                    : 'text-zinc-500 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 hover:text-zinc-900 dark:hover:text-zinc-200'}
                `}
              >
                <Hash size={16} style={{ color: goal.color || '#9ca3af' }} />
                <span className="truncate flex-1">{goal.title}</span>
                <button 
                  className="md:hidden ml-auto p-1 text-zinc-400 dark:text-zinc-500"
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    openContextMenu(e, 'list', goal.id, { title: goal.title });
                  }}
                >
                  <MoreVertical size={14} />
                </button>
              </NavLink>
            ))}
            {goals?.length === 0 && (
              <div className="px-3 py-2 text-sm text-zinc-400 dark:text-zinc-600 italic">No lists yet</div>
            )}
          </div>
        </div>
      </div>
    </div>
  );

  return (
    <>
      <motion.aside
        initial={false}
        animate={{ width: isSidebarOpen ? 260 : 0, opacity: isSidebarOpen ? 1 : 0 }}
        transition={{ duration: 0.3, ease: [0.4, 0, 0.2, 1] }}
        className="hidden md:block h-full overflow-hidden whitespace-nowrap"
      >
        {sidebarContent}
      </motion.aside>

      <div className="md:hidden">
        {isSidebarOpen && (
          <div 
            className="fixed inset-0 bg-black/40 dark:bg-black/60 z-40 backdrop-blur-sm"
            onClick={toggleSidebar}
          />
        )}
        <motion.aside
          initial={{ x: '-100%' }}
          animate={{ x: isSidebarOpen ? 0 : '-100%' }}
          transition={{ type: "spring", stiffness: 300, damping: 30 }}
          className="fixed inset-y-0 left-0 w-[280px] z-50 shadow-2xl"
        >
          {sidebarContent}
        </motion.aside>
      </div>
    </>
  );
};
