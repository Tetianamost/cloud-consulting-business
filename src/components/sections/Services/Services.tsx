import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { theme } from '../../../styles/theme';
import Section from '../../ui/Section';
import ServiceCard from './ServiceCard';
import Icon from '../../ui/Icon';
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
    title: "Cloud Assessment & Discovery",
    description: "Comprehensive analysis of your current infrastructure to identify migration candidates, dependencies, and create a detailed roadmap.",
    icon: "assessment",
    color: theme.colors.info,
    features: [
      "Infrastructure audit & dependency mapping",
      "Application portfolio analysis",
      "Total cost of ownership calculation",
      "Migration complexity assessment",
      "Risk identification & mitigation strategies"
    ]
  },
  {
    id: 2,
    title: "Migration Strategy & Planning",
    description: "Develop a tailored migration strategy using proven methodologies to minimize disruption and maximize business value.",
    icon: "strategy",
    color: theme.colors.accent,
    features: [
      "6R assessment (Rehost, Replatform, Repurchase, Refactor, Retire, Retain)",
      "Workload prioritization",
      "Migration wave planning",
      "Resource allocation & timeline development",
      "Governance & compliance planning"
    ]
  },
  {
    id: 3,
    title: "Cloud Implementation",
    description: "End-to-end migration execution with our expert team handling data transfer, application refactoring, and testing.",
    icon: "implementation",
    color: theme.colors.secondary,
    features: [
      "Cloud foundation setup (landing zone)",
      "Infrastructure as Code implementation",
      "Data migration with minimal downtime",
      "Application modernization & refactoring",
      "Comprehensive testing & validation"
    ]
  },
  {
    id: 4,
    title: "Cost Optimization",
    description: "Continuous monitoring and optimization of your cloud environment to reduce costs while maintaining performance.",
    icon: "optimization",
    color: theme.colors.success,
    features: [
      "Right-sizing resources & eliminating waste",
      "Reserved instance & savings plan analysis",
      "Auto-scaling implementation",
      "Strategic workload scheduling",
      "Regular cost optimization reviews"
    ]
  },
  {
    id: 5,
    title: "Security & Compliance",
    description: "Implement robust security controls and ensure compliance with industry regulations in your cloud environment.",
    icon: "security",
    color: theme.colors.warning,
    features: [
      "Security posture assessment",
      "Identity & access management",
      "Network security & encryption",
      "Compliance monitoring & reporting",
      "Security automation & remediation"
    ]
  },
  {
    id: 6,
    title: "Managed Cloud Services",
    description: "Ongoing management and support of your cloud infrastructure to ensure optimal performance and reliability.",
    icon: "managed",
    color: theme.colors.danger,
    features: [
      "24/7 monitoring & incident response",
      "Patch management & updates",
      "Performance optimization",
      "Backup & disaster recovery",
      "Regular health checks & reporting"
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
        Our Cloud Migration Services
      </SectionTitle>
      <SectionSubtitle
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={titleVariants}
      >
        Comprehensive solutions to transform your business with secure, efficient, and cost-effective cloud services
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

export default Services;