import React from 'react';

const Services: React.FC = () => {
  const services = [
    {
      id: 1,
      title: 'Cloud Readiness Assessment',
      description: 'Comprehensive evaluation of your current infrastructure to determine cloud migration feasibility, potential challenges, and expected ROI.',
      icon: '🔍'
    },
    {
      id: 2,
      title: 'Migration Strategy & Planning',
      description: 'Develop a tailored migration roadmap with prioritized workloads, architecture designs, and detailed implementation timelines.',
      icon: '📝'
    },
    {
      id: 3,
      title: 'Cloud Implementation',
      description: 'Expert execution of your migration plan with minimal disruption to operations, including data transfer, application refactoring, and testing.',
      icon: '🚀'
    },
    {
      id: 4,
      title: 'Post-Migration Support',
      description: 'Ongoing optimization, monitoring, and management of your cloud environment to ensure performance, security, and cost efficiency.',
      icon: '🛡️'
    },
    {
      id: 5,
      title: 'Cloud Cost Optimization',
      description: 'Identify and implement strategies to reduce cloud spending while maintaining or improving performance and reliability.',
      icon: '💰'
    },
    {
      id: 6,
      title: 'Cloud Security & Compliance',
      description: 'Implement industry-leading security practices and ensure your cloud environment meets all relevant regulatory requirements.',
      icon: '🔒'
    }
  ];

  return (
    <section id="services" className="section" style={sectionStyle}>
      <div className="container">
        <h2 className="section-title">Our Cloud Migration Services</h2>
        <p style={sectionDescriptionStyle}>
          Comprehensive solutions to streamline your journey to the cloud
        </p>
        
        <div style={servicesGridStyle}>
          {services.map(service => (
            <div key={service.id} className="card" style={serviceCardStyle}>
              <div style={iconContainerStyle}>
                <span style={iconStyle}>{service.icon}</span>
              </div>
              <h3 style={serviceHeadingStyle}>{service.title}</h3>
              <p>{service.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

// Styles
const sectionStyle: React.CSSProperties = {
  background: '#fff',
  padding: '80px 0'
};

const sectionDescriptionStyle: React.CSSProperties = {
  fontSize: '1.2rem',
  maxWidth: '800px',
  margin: '0 auto 40px',
  textAlign: 'center',
  color: '#666'
};

const servicesGridStyle: React.CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
  gap: '30px'
};

const serviceCardStyle: React.CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  textAlign: 'center',
  height: '100%'
};

const iconContainerStyle: React.CSSProperties = {
  width: '80px',
  height: '80px',
  borderRadius: '50%',
  background: 'var(--light-color)',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  marginBottom: '20px'
};

const iconStyle: React.CSSProperties = {
  fontSize: '2.5rem'
};

const serviceHeadingStyle: React.CSSProperties = {
  fontSize: '1.5rem',
  marginBottom: '15px',
  color: 'var(--primary-color)'
};

export default Services;