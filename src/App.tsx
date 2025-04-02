import React from 'react';
import { ThemeProvider } from 'styled-components';
import { theme } from './styles/theme';
import GlobalStyles from './styles/GlobalStyles';
import Header from './components/layout/Header';
import Footer from './components/layout/Footer';
import Hero from './components/sections/Hero/Hero';
import Services from './components/sections/Services/Services';
import Certifications from './components/sections/Certifications/Certifications';
import CaseStudies from './components/sections/CaseStudies/CaseStudies';
import Pricing from './components/sections/Pricing/Pricing';
import Contact from './components/sections/Contact/Contact';

function App() {
  return (
    <ThemeProvider theme={theme}>
      <GlobalStyles />
      <Header />
      <main>
        <Hero />
        <Services />
        <Certifications />
        <CaseStudies />
        <Pricing />
        <Contact />
      </main>
      <Footer />
    </ThemeProvider>
  );
}

export default App;