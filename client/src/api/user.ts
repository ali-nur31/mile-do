import { api } from './axios';

export interface UserProfile {
  email: string;
  created_at: string;
}

export const userApi = {
  getMe: async () => {
    const response = await api.get<UserProfile>('/users/me');
    return response.data;
  }
};
