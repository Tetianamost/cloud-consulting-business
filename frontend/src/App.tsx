import React, { Suspense } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from 'styled-components';
import { theme } from './styles/theme';
import GlobalStyles from './styles/GlobalStyles';
import './styles/admin.css';
import { AuthProvider } from './contexts/AuthContext';

// Public site components (eagerly loaded for better initial page performance)
import Header from './components/layout/Header';
import Footer from './components/layout/Footer';
import Hero from './components/sections/Hero/Hero';
import Services from './components/sections/Services/Services';
import Certifications from './components/sections/Certifications/Certifications';
import ProjectHighlights from './components/sections/ProjectInsights/ProjectInsights';
import Pricing from './components/sections/Pricing/Pricing';
import Contact from './components/sections/Contact/Contact';

// Lazy-loaded admin components for better performance
const Login = React.lazy(() => import('./components/admin/Login'));
const ProtectedRoute = React.lazy(() => import('./components/admin/ProtectedRoute'));
const AdminLayoutWrapper = React.lazy(() => import('./components/admin/AdminLayoutWrapper'));
const V0DashboardNew = React.lazy(() => import('./components/admin/V0DashboardNew'));
const AIReportsPage = React.lazy(() => import('./components/admin/AIReportsPage'));
const ReportPage = React.lazy(() => import('./components/admin/ReportPage'));

// Loading component for lazy-loaded routes
const AdminLoadingFallback: React.FC = () => (
  <div className="min-h-screen bg-gray-50 flex items-center justify-center">
    <div className="flex flex-col items-center space-y-4">
      <div className="animate-spin rounded-full h-8 w-8 border-2 border-blue-600 border-t-transparent"></div>
      <div className="text-gray-600 text-sm">Loading admin dashboard...</div>
    </div>
  </div>
);

// Main site layout component
const MainSite: React.FC = () => (
  <>
    <Header />
    <main>
      <Hero />
      <Services />
      <Certifications />
      <ProjectHighlights />
      <Pricing />
      <Contact />
    </main>
    <Footer />
  </>
);

function App() {
  // Check if admin dashboard is enabled
  const enableAdmin = process.env.REACT_APP_ENABLE_ADMIN !== 'false';

  return (
    <ThemeProvider theme={theme}>
      <GlobalStyles />
      <AuthProvider>
        <Router>
          <Routes>
            {/* Public site route */}
            <Route path="/" element={<MainSite />} />

            {/* Admin routes with lazy loading and suspense */}
            {enableAdmin && (
              <>
                <Route path="/admin/login" element={
                  <Suspense fallback={<AdminLoadingFallback />}>
                    <Login />
                  </Suspense>
                } />
                <Route path="/admin" element={
                  <Suspense fallback={<AdminLoadingFallback />}>
                    <ProtectedRoute>
                      <Navigate to="/admin/dashboard" replace />
                    </ProtectedRoute>
                  </Suspense>
                } />
                <Route path="/admin/dashboard" element={
                  <Suspense fallback={<AdminLoadingFallback />}>
                    <ProtectedRoute>
                      <AdminLayoutWrapper>
                        <V0DashboardNew />
                      </AdminLayoutWrapper>
                    </ProtectedRoute>
                  </Suspense>
                } />
                <Route path="/admin/inquiries" element={
                  <Suspense fallback={<AdminLoadingFallback />}>
                    <ProtectedRoute>
                      <AdminLayoutWrapper>
                        <V0DashboardNew />
                      </AdminLayoutWrapper>
                    </ProtectedRoute>
                  </Suspense>
                } />
                <Route path="/admin/metrics" element={
                  <Suspense fallback={<AdminLoadingFallback />}>
                    <ProtectedRoute>
                      <AdminLayoutWrapper>
                        <V0DashboardNew />
                      </AdminLayoutWrapper>
                    </ProtectedRoute>
                  </Suspense>
                } />
                <Route path="/admin/email-status" element={
                  <Suspense fallback={<AdminLoadingFallback />}>
                    <ProtectedRoute>
                      <AdminLayoutWrapper>
                        <V0DashboardNew />
                      </AdminLayoutWrapper>
                    </ProtectedRoute>
                  </Suspense>
                } />
                <Route path="/admin/reports" element={
                  <Suspense fallback={<AdminLoadingFallback />}>
                    <ProtectedRoute>
                      <AdminLayoutWrapper>
                        <AIReportsPage />
                      </AdminLayoutWrapper>
                    </ProtectedRoute>
                  </Suspense>
                } />
                <Route path="/admin/reports/:id" element={
                  <Suspense fallback={<AdminLoadingFallback />}>
                    <ProtectedRoute>
                      <AdminLayoutWrapper>
                        <ReportPage />
                      </AdminLayoutWrapper>
                    </ProtectedRoute>
                  </Suspense>
                } />
              </>
            )}

            {/* Redirect any unknown routes to home */}
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;