import { api } from './axios';
import type { LoginRequest, RegisterRequest, AuthResponse } from '../types';

export const authApi = {
  login: async (data: LoginRequest) => {
    const response = await api.post<AuthResponse>('/auth/login', data);
    return response.data;
  },

  register: async (data: RegisterRequest) => {
    const response = await api.post<AuthResponse>('/auth/register', data);
    return response.data;
  },

  logout: async () => {
    await api.delete('/auth/logout');
  },

  refreshToken: async (token: string) => {
    const response = await api.post<AuthResponse>('/auth/refresh', { refresh_token: token });
    return response.data;
  }
};
