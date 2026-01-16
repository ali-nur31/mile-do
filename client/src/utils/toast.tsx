import { toast } from 'sonner';
import { CheckCircle2, Trash2, Info } from 'lucide-react';

export const showToast = (type: 'success' | 'delete' | 'info', title: string, message?: string) => {
  if (type === 'success') {
    toast.custom(() => (
      <div className="bg-green-50 dark:bg-green-900/20 border border-green-100 dark:border-green-900/50 text-green-700 dark:text-green-400 px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 w-full">
        <CheckCircle2 size={18} />
        <div>
          <span className="font-bold text-sm block">{title}</span>
          {message && <span className="text-xs opacity-90">{message}</span>}
        </div>
      </div>
    ), { duration: 3000 });
  } else if (type === 'delete') {
    toast.custom(() => (
      <div className="bg-red-50 dark:bg-red-900/20 border border-red-100 dark:border-red-900/50 text-red-600 dark:text-red-400 px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 w-full">
        <Trash2 size={18} />
        <div>
          <span className="font-bold text-sm block">{title}</span>
          {message && <span className="text-xs opacity-90">{message}</span>}
        </div>
      </div>
    ), { duration: 3000 });
  } else {
    toast.custom(() => (
      <div className="bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 text-zinc-700 dark:text-zinc-200 px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 w-full">
        <Info size={18} />
        <div>
          <span className="font-bold text-sm block">{title}</span>
          {message && <span className="text-xs opacity-90">{message}</span>}
        </div>
      </div>
    ), { duration: 3000 });
  }
};
