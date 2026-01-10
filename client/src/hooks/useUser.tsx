import { useQuery } from '@tanstack/react-query';
import { userApi } from '../api/user';

export const useUser = () => {
  return useQuery({
    queryKey: ['user', 'me'],
    queryFn: userApi.getMe,
    retry: false
  });
};
