import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { theme } from '../../../styles/theme';
import Section from '../../ui/section';
import ServiceCard from './ServiceCard';
import Icon from '../../ui/icon';
import { FiCloudOff, FiCloudSnow, FiCpu, FiDollarSign, FiShield, FiMonitor } from 'react-icons/fi';


const SectionTitle = styled(motion.h2)`
  text-align: center;
  margin-bottom: ${theme.space[3]};
  color: ${theme.colors.primary};
`;

const SectionSubtitle = styled(motion.p)`
  text-align: center;
  max-width: 600px;
  margin: 0 auto ${theme.space[8]};
  color: ${theme.colors.gray600};
  font-size: ${theme.fontSizes.lg};
`;

const ServicesGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: ${theme.space[6]};
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-template-columns: repeat(2, 1fr);
  }
  
  @media (max-width: ${theme.breakpoints.md}) {
    grid-template-columns: 1fr;
  }
`;

const services = [
  {
    id: 1,
    title: "Cloud Assessment & Roadmap",
    description: "Evaluate your current infrastructure and develop a strategic roadmap for cloud adoption or optimization tailored to your specific business needs.",
    icon: "assessment",
    color: theme.colors.info,
    features: [
      "Infrastructure evaluation & analysis",
      "Cloud readiness assessment",
      "Cost estimation & comparison",
      "Risk identification & planning",
      "Personalized recommendations report"
    ]
  },
  {
    id: 2,
    title: "Small-Scale Migrations",
    description: "Assist with focused, manageable cloud migrations for specific applications or servers that can be completed within our flexible availability.",
    icon: "strategy",
    color: theme.colors.accent,
    features: [
      "Single-application migration planning",
      "Small database migrations",
      "Simple lift-and-shift operations",
      "Step-by-step implementation guidance",
      "Migration verification & testing"
    ]
  },
  {
    id: 3,
    title: "Cloud Optimization Consulting",
    description: "Review your existing cloud environment to identify cost-saving opportunities and performance improvements without extensive resource commitment.",
    icon: "optimization",
    color: theme.colors.success,
    features: [
      "Cost analysis & reduction planning",
      "Resource right-sizing recommendations",
      "Reserved instance & savings plan advice",
      "Architecture review & improvement suggestions",
      "Performance optimization guidance"
    ]
  },
  {
    id: 4,
    title: "Cloud Architecture Review",
    description: "Expert examination of your cloud architecture with actionable recommendations to improve security, reliability, performance, and cost-efficiency.",
    icon: "security",
    color: theme.colors.warning,
    features: [
      "Architecture best practices assessment",
      "Security configuration review",
      "Scalability & resilience evaluation",
      "Compliance check & recommendations",
      "Documentation & knowledge transfer"
    ]
  }
];

// Animation variants
const titleVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: { 
    opacity: 1, 
    y: 0,
    transition: {
      duration:.5
    }
  }
};

const Services: React.FC = () => {
  const renderIcon = (iconName: string) => {
    switch (iconName) {
      case 'assessment':
        return <Icon icon={FiCloudOff} size={24} />;
      case 'strategy':
        return <Icon icon={FiCloudSnow} size={24} />;
      case 'implementation':
        return <Icon icon={FiCpu} size={24} />;
      case 'optimization':
        return <Icon icon={FiDollarSign} size={24} />;
      case 'security':
        return <Icon icon={FiShield} size={24} />;
      case 'managed':
        return <Icon icon={FiMonitor} size={24} />;
      default:
        return null;
    }
  };

  return (
    <Section id="services" background="white">
      <SectionTitle
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={titleVariants}
      >
        Our Cloud Services
      </SectionTitle>
      <SectionSubtitle
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={titleVariants}
      >
        Focused cloud consulting services with flexible scheduling that adapts to your business needs and project timelines
      </SectionSubtitle>
      
      <ServicesGrid>
        {services.map((service, index) => (
          <ServiceCard 
            key={service.id}
            service={{...service, icon: renderIcon(service.icon)}}
            index={index}
          />
        ))}
      </ServicesGrid>
    </Section>
  );
};

const MemoizedServices = React.memo(Services);
export default MemoizedServices;