import { useMutation } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { authApi } from '../api/auth';
import { useStore } from '../store/useUIStore';
import type { LoginRequest, RegisterRequest } from '../types';

export const useAuth = () => {
  const navigate = useNavigate();
  const setAuthenticated = useStore((state) => state.setAuthenticated);

  const loginMutation = useMutation({
    mutationFn: (data: LoginRequest) => authApi.login(data),
    onSuccess: (data) => {
      localStorage.setItem('access_token', data.access_token);
      localStorage.setItem('refresh_token', data.refresh_token);
      setAuthenticated(true);
      navigate('/');
    }
  });

  const registerMutation = useMutation({
    mutationFn: (data: RegisterRequest) => authApi.register(data),
    onSuccess: (data) => {
      localStorage.setItem('access_token', data.access_token);
      localStorage.setItem('refresh_token', data.refresh_token);
      setAuthenticated(true);
      navigate('/');
    }
  });

  return {
    login: loginMutation,
    register: registerMutation,
  };
};
