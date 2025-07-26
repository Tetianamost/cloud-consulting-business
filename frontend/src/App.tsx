import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from 'styled-components';
import { theme } from './styles/theme';
import GlobalStyles from './styles/GlobalStyles';

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
import AdminLayout from './components/admin/AdminLayout';
import Dashboard from './components/admin/Dashboard';
import InquiriesList from './components/admin/InquiriesList';
import MetricsDashboard from './components/admin/MetricsDashboard';
import EmailStatusMonitor from './components/admin/EmailStatusMonitor';

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
      <Router>
        <Routes>
          {/* Public site route */}
          <Route path="/" element={<MainSite />} />
          
          {/* Admin routes */}
          {enableAdmin && (
            <Route path="/admin" element={<AdminLayout />}>
              <Route index element={<Dashboard />} />
              <Route path="inquiries" element={<InquiriesList />} />
              <Route path="metrics" element={<MetricsDashboard />} />
              <Route path="email-status" element={<EmailStatusMonitor />} />
            </Route>
          )}
          
          {/* Redirect any unknown routes to home */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

export default App;