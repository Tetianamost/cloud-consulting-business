import React from 'react';

const Hero: React.FC = () => {
  return (
    <section id="home" style={heroStyle}>
      <div className="container" style={heroContentStyle}>
        <h1 style={headingStyle}>Expert Cloud Migration Services</h1>
        <p style={subheadingStyle}>
          Transform your business with secure, efficient, and cost-effective cloud solutions
        </p>
        <div style={ctaContainerStyle}>
          <a href="#services" className="btn btn-accent" style={ctaButtonStyle}>
            Explore Services
          </a>
          <a href="#contact" className="btn" style={ctaButtonStyle}>
            Get in Touch
          </a>
        </div>
      </div>
    </section>
  );
};

// Styles
const heroStyle: React.CSSProperties = {
  height: '100vh',
  background: 'linear-gradient(135deg, rgba(0,102,204,0.9) 0%, rgba(0,153,255,0.9) 100%)',
  backgroundSize: 'cover',
  backgroundPosition: 'center',
  color: '#fff',
  display: 'flex',
  alignItems: 'center',
  textAlign: 'center',
  padding: '0 20px',
  position: 'relative'
};

const heroContentStyle: React.CSSProperties = {
  maxWidth: '800px',
  margin: '0 auto'
};

const headingStyle: React.CSSProperties = {
  fontSize: '3.5rem',
  marginBottom: '20px',
  textShadow: '2px 2px 5px rgba(0,0,0,0.2)'
};

const subheadingStyle: React.CSSProperties = {
  fontSize: '1.5rem',
  marginBottom: '30px',
  fontWeight: '300'
};

const ctaContainerStyle: React.CSSProperties = {
  display: 'flex',
  justifyContent: 'center',
  gap: '20px',
  flexWrap: 'wrap'
};

const ctaButtonStyle: React.CSSProperties = {
  padding: '15px 30px',
  fontSize: '1.1rem',
  boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
  transition: 'all 0.3s ease'
};

export default Hero;