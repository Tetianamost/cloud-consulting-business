import React from 'react';
import { useLocation } from 'react-router-dom';
import V0AdminLayout from './V0AdminLayout';

interface AdminLayoutWrapperProps {
  children: React.ReactNode;
}

/**
 * AdminLayoutWrapper - Wraps admin components with V0AdminLayout
 * This component provides the consistent admin layout for all admin routes
 */
const AdminLayoutWrapper: React.FC<AdminLayoutWrapperProps> = ({ children }) => {
  const location = useLocation();

  return (
    <V0AdminLayout currentPath={location.pathname}>
      {children}
    </V0AdminLayout>
  );
};

export default AdminLayoutWrapper;