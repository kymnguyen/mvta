import React from 'react';
import { Outlet, useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

export const AppLayout: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="app">
      <header className="app-header">
        <div className="header-content">
          <h1>Management Vehicle Tracking Application</h1>
          <nav className="app-nav">
            <Link to="/" className="nav-link">Tracking (tracking-svc)</Link>
            <Link to="/vehicle-svc" className="nav-link">Vehicles (vehicle-svc)</Link>
          </nav>
          <div className="header-right">
            <span className="user-info">{user?.email || user?.name}</span>
            <button onClick={handleLogout} className="btn-logout">
              Logout
            </button>
          </div>
        </div>
      </header>
      <main className="app-main">
        <Outlet />
      </main>
    </div>
  );
};
