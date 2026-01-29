import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { VehicleList } from './features/vehicles/VehicleList';
import { VehicleDetail } from './features/vehicles/VehicleDetail';
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
      <BrowserRouter>
        <div className="app">
          <header className="app-header">
            <h1>MVTA Fleet Management</h1>
          </header>
          <main className="app-main">
            <Routes>
              <Route path="/" element={<VehicleList />} />
              <Route path="/vehicles/:id" element={<VehicleDetail />} />
            </Routes>
          </main>
        </div>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;
