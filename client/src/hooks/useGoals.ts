import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { goalsApi } from '../api/goals';
import type { UpdateGoalRequest } from '../api/goals';
import type { CreateGoalRequest } from '../types';

export const useGoals = () => {
  const queryClient = useQueryClient();

  const { data: goals, isLoading } = useQuery({
    queryKey: ['goals'],
    queryFn: goalsApi.getAll
  });

  const createGoal = useMutation({
    mutationFn: (data: CreateGoalRequest) => goalsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
    }
  });

  const updateGoal = useMutation({
    mutationFn: (data: UpdateGoalRequest) => goalsApi.update(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
    }
  });

  const deleteGoal = useMutation({
    mutationFn: (id: number) => goalsApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
    }
  });

  return {
    goals,
    isLoading,
    createGoal,
    updateGoal,
    deleteGoal
  };
};
