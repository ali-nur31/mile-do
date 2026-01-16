import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api/axios';
import { useStore } from '../store/useUIStore';
import type { AuthResponse, LoginRequest, RegisterRequest } from '../types';
import { CheckCircle2, ArrowRight } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { motion } from 'framer-motion';

export const Login = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();
  
  const { setAuthenticated, theme } = useStore();

  useEffect(() => {
    if (theme === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [theme]);

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      if (isLogin) {
        const payload: LoginRequest = { email, password };
        const { data } = await api.post<AuthResponse>('/auth/login', payload);
        handleAuthSuccess(data);
      } else {
        if (password !== confirmPassword) {
          setError("Passwords do not match");
          setIsLoading(false);
          return;
        }
        const payload: RegisterRequest = { email, password, confirm_password: confirmPassword };
        const { data } = await api.post<AuthResponse>('/auth/register', payload);
        handleAuthSuccess(data);
      }
    } catch (err: any) {
      console.error(err);
      setError(err.response?.data?.message || 'Authentication failed. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleAuthSuccess = (data: AuthResponse) => {
    localStorage.setItem('access_token', data.access_token);
    localStorage.setItem('refresh_token', data.refresh_token);
    setAuthenticated(true);
    navigate('/');
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-zinc-50 dark:bg-zinc-950 px-4 transition-colors duration-200">
      <motion.div 
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="w-full max-w-md"
      >
        <div className="bg-white dark:bg-zinc-900 rounded-2xl shadow-xl border border-zinc-100 dark:border-zinc-800 overflow-hidden transition-colors duration-200">
          <div className="px-8 pt-8 pb-6 text-center">
            <div className="flex justify-center mb-4">
              <CheckCircle2 className="w-12 h-12 text-blue-600 dark:text-blue-500" strokeWidth={2} />
            </div>
            <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100 tracking-tight">
              {isLogin ? 'Welcome back' : 'Create an account'}
            </h1>
            <p className="text-zinc-500 dark:text-zinc-400 mt-2 text-sm">
              {isLogin ? 'Enter your credentials to access Mile-Do' : 'Start organizing your life today'}
            </p>
          </div>

          <form onSubmit={handleSubmit} className="px-8 pb-8 space-y-4">
            {error && (
              <div className="p-3 rounded-lg bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 text-sm border border-red-100 dark:border-red-900/50">
                {error}
              </div>
            )}

            <Input 
              label="Email" 
              type="email" 
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="name@example.com"
              required 
            />
            
            <Input 
              label="Password" 
              type="password" 
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required 
            />

            {!isLogin && (
              <Input 
                label="Confirm Password" 
                type="password" 
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="••••••••"
                required 
              />
            )}

            <div className="pt-2">
              <Button type="submit" className="w-full" isLoading={isLoading}>
                {isLogin ? 'Sign In' : 'Create Account'} <ArrowRight size={16} className="ml-2" />
              </Button>
            </div>
          </form>

          <div className="px-8 py-4 bg-zinc-50 dark:bg-zinc-900/50 border-t border-zinc-100 dark:border-zinc-800 text-center transition-colors duration-200">
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              {isLogin ? "Don't have an account? " : "Already have an account? "}
              <button 
                onClick={() => { setIsLogin(!isLogin); setError(''); }}
                className="font-medium text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 transition-colors"
              >
                {isLogin ? 'Sign up' : 'Log in'}
              </button>
            </p>
          </div>
        </div>
        
        <p className="text-center text-xs text-zinc-400 dark:text-zinc-600 mt-8">
          © 2026 Mile-Do.
        </p>
      </motion.div>
    </div>
  );
};
