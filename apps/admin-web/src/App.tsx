import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AuthProvider } from './shared/hooks/useAuth';
import { ProtectedRoute } from './shared/components/ProtectedRoute';
import { Login } from './features/auth/Login';
import { VehicleList } from './features/vehicles/VehicleList';
import { VehicleDetail } from './features/vehicles/VehicleDetail';
import { AppLayout } from './shared/components/AppLayout';
import './App.css';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route element={<ProtectedRoute />}>
              <Route element={<AppLayout />}>
                <Route path="/" element={<VehicleList />} />
                <Route path="/vehicles/:id" element={<VehicleDetail />} />
              </Route>
            </Route>
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
