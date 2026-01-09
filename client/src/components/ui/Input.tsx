import React from 'react';

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export const Input: React.FC<InputProps> = ({ label, error, className = '', ...props }) => {
  return (
    <div className="w-full">
      {label && <label className="block text-sm font-medium text-zinc-700 mb-1.5">{label}</label>}
      <input
        className={`
          w-full px-3 py-2.5 bg-white border rounded-lg text-sm transition-all outline-none
          placeholder:text-zinc-400
          focus:ring-2 focus:ring-blue-100 focus:border-blue-500
          ${error ? 'border-red-500 focus:ring-red-100 focus:border-red-500' : 'border-zinc-300'}
          ${className}
        `}
        {...props}
      />
      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}
    </div>
  );
};
