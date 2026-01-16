import { useEffect } from 'react';
import { Routes, Route, Navigate, useNavigate } from 'react-router-dom';
import { AppLayout } from './components/layout/AppLayout';
import { Dashboard } from './pages/Dashboard';
import { Today } from './pages/Today';
import { Goals } from './pages/Goals';
import { CalendarPage } from './pages/CalendarPage';
import { Login } from './pages/Login';
import { AllTasks } from './pages/AllTasks';
import { useStore } from './store/useUIStore';
import { GoalDetails } from './pages/GoalDetails';

const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const { isAuthenticated } = useStore();
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
};

function App() {
  const { logout } = useStore();
  const navigate = useNavigate();

  useEffect(() => {
    const handleLogoutEvent = () => {
      logout();
      navigate('/login');
    };
    window.addEventListener('auth:logout', handleLogoutEvent);
    return () => window.removeEventListener('auth:logout', handleLogoutEvent);
  }, [logout, navigate]);

  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/" element={
        <ProtectedRoute>
          <AppLayout />
        </ProtectedRoute>
      }>
        <Route index element={<Dashboard />} />
        <Route path="today" element={<Today />} />
        <Route path="calendar" element={<CalendarPage />} />
        <Route path="goals" element={<Goals />} />
        <Route path="goals/:id" element={<GoalDetails />} />
        <Route path="all" element={<AllTasks />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  );
}

export default App;
