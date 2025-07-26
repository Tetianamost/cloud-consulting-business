import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from 'styled-components';
import { theme } from './styles/theme';
import GlobalStyles from './styles/GlobalStyles';
import './styles/admin.css';
import { AuthProvider } from './contexts/AuthContext';

// Public site components
import Header from './components/layout/Header';
import Footer from './components/layout/Footer';
import Hero from './components/sections/Hero/Hero';
import Services from './components/sections/Services/Services';
import Certifications from './components/sections/Certifications/Certifications';
import ProjectHighlights from './components/sections/ProjectInsights/ProjectInsights';
import Pricing from './components/sections/Pricing/Pricing';
import Contact from './components/sections/Contact/Contact';

// Admin components
import Login from './components/admin/Login';
import ProtectedRoute from './components/admin/ProtectedRoute';
import V0DashboardNew from './components/admin/V0DashboardNew';
import V0Dashboard from './components/admin/V0Dashboard';

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

            {/* Admin routes */}
            {enableAdmin && (
              <>
                <Route path="/admin/login" element={<Login />} />
                <Route path="/admin" element={
                  <ProtectedRoute>
                    <Navigate to="/admin/dashboard" replace />
                  </ProtectedRoute>
                } />
                <Route path="/admin/dashboard" element={
                  <ProtectedRoute>
                    <V0DashboardNew />
                  </ProtectedRoute>
                } />
                <Route path="/admin/inquiries" element={
                  <ProtectedRoute>
                    <V0DashboardNew />
                  </ProtectedRoute>
                } />
                <Route path="/admin/metrics" element={
                  <ProtectedRoute>
                    <V0DashboardNew />
                  </ProtectedRoute>
                } />
                <Route path="/admin/email-status" element={
                  <ProtectedRoute>
                    <V0DashboardNew />
                  </ProtectedRoute>
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