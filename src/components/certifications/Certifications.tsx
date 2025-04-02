import React from 'react';

const Certifications: React.FC = () => {
  const awsCertifications = [
    {
      id: 1,
      title: 'AWS Certified Solutions Architect - Professional',
      description: 'Advanced expertise in designing distributed applications and systems on the AWS platform.',
      icon: '🏆'
    },
    {
      id: 2,
      title: 'AWS Certified DevOps Engineer - Professional',
      description: 'Specialized in implementing and managing continuous delivery systems and methodologies on AWS.',
      icon: '🏆'
    },
    {
      id: 3,
      title: 'AWS Certified Security - Specialty',
      description: 'Deep expertise in security best practices for AWS platform and securing complex AWS environments.',
      icon: '🏆'
    },
    {
      id: 4,
      title: 'AWS Certified Data Analytics - Specialty',
      description: 'Expertise in designing and implementing AWS services to derive valuable insights from data.',
      icon: '🏆'
    }
  ];

  return (
    <section id="certifications" className="section" style={sectionStyle}>
      <div className="container">
        <h2 className="section-title">AWS Certifications</h2>
        <p style={sectionDescriptionStyle}>
          Our team holds industry-leading certifications, ensuring you receive expert-level cloud migration services
        </p>
        
        <div style={certificationsGridStyle}>
          {awsCertifications.map(cert => (
            <div key={cert.id} className="card" style={certCardStyle}>
              <div style={iconContainerStyle}>
                <span style={iconStyle}>{cert.icon}</span>
              </div>
              <h3 style={certHeadingStyle}>{cert.title}</h3>
              <p>{cert.description}</p>
            </div>
          ))}
        </div>
        
        <div style={trustBadgeContainerStyle}>
          <div style={trustBadgeStyle}>
            <h3 style={{marginBottom: '15px'}}>AWS Select Consulting Partner</h3>
            <p>We're part of the Amazon Partner Network, with validated AWS expertise</p>
          </div>
        </div>
      </div>
    </section>
  );
};

// Styles
const sectionStyle: React.CSSProperties = {
  background: 'var(--light-color)',
  padding: '80px 0'
};

const sectionDescriptionStyle: React.CSSProperties = {
  fontSize: '1.2rem',
  maxWidth: '800px',
  margin: '0 auto 40px',
  textAlign: 'center',
  color: '#666'
};

const certificationsGridStyle: React.CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
  gap: '30px',
  marginBottom: '50px'
};

const certCardStyle: React.CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  textAlign: 'center',
  height: '100%'
};

const iconContainerStyle: React.CSSProperties = {
  width: '70px',
  height: '70px',
  borderRadius: '50%',
  background: '#fff',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  marginBottom: '20px',
  boxShadow: '0 3px 10px rgba(0, 0, 0, 0.1)'
};

const iconStyle: React.CSSProperties = {
  fontSize: '2rem'
};

const certHeadingStyle: React.CSSProperties = {
  fontSize: '1.3rem',
  marginBottom: '15px',
  color: 'var(--primary-color)'
};

const trustBadgeContainerStyle: React.CSSProperties = {
  display: 'flex',
  justifyContent: 'center'
};

const trustBadgeStyle: React.CSSProperties = {
  background: 'linear-gradient(135deg, #232F3E 0%, #3B4D61 100%)',
  color: '#fff',
  padding: '30px',
  borderRadius: '10px',
  textAlign: 'center',
  maxWidth: '600px',
  boxShadow: '0 5px 15px rgba(0, 0, 0, 0.2)'
};

export default Certifications;