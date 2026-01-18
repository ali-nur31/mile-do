import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { goalsApi } from '../api/goals';
import type { CreateGoalRequest, Goal } from '../types';
import { showToast } from '../utils/toast';

export const useGoals = () => {
  const queryClient = useQueryClient();

  const { data: goals, isLoading } = useQuery({
    queryKey: ['goals'],
    queryFn: async () => {
        try {
            return await goalsApi.getAll();
        } catch (error) {
            return [];
        }
    },
    retry: 1,
  });

  const createGoal = useMutation({
    mutationFn: (data: CreateGoalRequest) => goalsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
    },
    onError: () => {
        showToast('delete', 'Failed to create list');
    }
  });

  const deleteGoal = useMutation({
    mutationFn: (id: number) => goalsApi.delete(id),
    onMutate: async (id) => {
        await queryClient.cancelQueries({ queryKey: ['goals'] });
        const previousGoals = queryClient.getQueryData(['goals']);
        
        queryClient.setQueryData(['goals'], (old: Goal[] | undefined) => 
            old ? old.filter(g => g.id !== id) : []
        );
        
        return { previousGoals };
    },
    onError: (_err, _id, context: any) => {
        if (context?.previousGoals) {
            queryClient.setQueryData(['goals'], context.previousGoals);
        }
        showToast('delete', 'Failed to delete list');
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    }
  });

  return {
    goals,
    isLoading,
    createGoal,
    deleteGoal
  };
};
