import React, { useState } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import { theme } from '../../../styles/theme';
import Section from '../../ui/Section';
import { FiArrowRight, FiArrowLeft } from 'react-icons/fi';
import CaseStudyCard from './CaseStudyCard';
import Icon from '../../ui/Icon';

const caseStudies = [
  {
    id: 1,
    title: 'E-commerce Platform Migration',
    industry: 'Retail',
    description: 'Migrated a high-traffic e-commerce platform from on-premises infrastructure to AWS, resulting in improved scalability and performance.',
    challenges: [
      'Legacy monolithic architecture',
      'Seasonal traffic spikes requiring 5x capacity',
      'Payment card industry (PCI) compliance requirements',
      'Zero downtime migration requirement'
    ],
    solution: 'We implemented a phased migration approach, starting with non-critical components. The application was refactored into microservices using containers, with static content moved to S3 and CloudFront. We used AWS Auto Scaling for handling traffic fluctuations and implemented a blue-green deployment strategy to ensure zero downtime.',
    results: {
      costReduction: '35%',
      performanceImprovement: '60%',
      deploymentTime: '90%',
      scalability: 'Unlimited',
    },
    testimonial: {
      quote: "The migration transformed our business. We're now able to handle Black Friday traffic with ease, and our development team can ship features faster than ever.",
      author: "Sarah Johnson",
      position: "CTO, RetailGiant Inc."
    },
    image: 'https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80'
  },
  {
    id: 2,
    title: 'Healthcare Data Platform Modernization',
    industry: 'Healthcare',
    description: 'Redesigned and migrated a critical healthcare data platform to AWS, enhancing security, compliance, and analytics capabilities.',
    challenges: [
      'Strict HIPAA compliance requirements',
      'Massive historical data (10+ years, 50TB+)',
      'Real-time data processing needs',
      'Integration with legacy systems'
    ],
    solution: 'We developed a secure landing zone with enhanced security controls and implemented a data lake architecture using S3, Glue, and Athena. Sensitive data was encrypted and access strictly controlled with IAM policies. Lambda functions were used for real-time processing, and API Gateway provided secure access for legacy systems.',
    results: {
      costReduction: '28%',
      dataProcessingSpeed: '75%',
      complianceAutomation: '100%',
      insightsAccess: 'Real-time',
    },
    testimonial: {
      quote: "Our migration to AWS allowed us to not only meet compliance requirements with greater confidence but also unlock insights from our data that were previously inaccessible.",
      author: "Dr. Michael Chen",
      position: "Director of Informatics, HealthFirst"
    },
    image: 'https://images.unsplash.com/photo-1576091160550-2173dba999ef?ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80'
  },
  {
    id: 3,
    title: 'Financial Services Application Portfolio',
    industry: 'Finance',
    description: 'Migrated a portfolio of 15+ financial applications from a traditional data center to AWS, improving reliability and disaster recovery capabilities.',
    challenges: [
      'Stringent financial regulatory requirements',
      'Complex interdependent applications',
      'Legacy mainframe systems',
      'Aggressive timelines due to data center contract expiration'
    ],
    solution: 'We established a Cloud Center of Excellence and implemented a factory model for application assessment and migration. Critical applications were refactored, while others were rehosted. We implemented comprehensive monitoring, automated disaster recovery, and enhanced security controls throughout the stack.',
    results: {
      costReduction: '42%',
      systemAvailability: '99.99%',
      recoveryTimeObjective: '15 minutes',
      deploymentFrequency: '3x increase',
    },
    testimonial: {
      quote: "The migration enabled us to exceed our regulatory requirements while significantly reducing costs. More importantly, we've transformed our ability to innovate and respond to market changes.",
      author: "James Wilson",
      position: "CIO, Global Financial Group"
    },
    image: 'https://images.unsplash.com/photo-1460925895917-afdab827c52f?ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80'
  }
];

const SectionTitle = styled(motion.h2)`
  text-align: center;
  margin-bottom: ${theme.space[3]};
`;

const SectionSubtitle = styled(motion.p)`
  text-align: center;
  max-width: 600px;
  margin: 0 auto ${theme.space[8]};
  color: ${theme.colors.gray600};
  font-size: ${theme.fontSizes.lg};
`;

const SliderContainer = styled.div`
  position: relative;
  overflow: hidden;
  padding: ${theme.space[4]} 0;
`;

const SliderWrapper = styled(motion.div)`
  display: flex;
  width: 100%;
`;

const SliderControls = styled.div`
  display: flex;
  justify-content: center;
  margin-top: ${theme.space[6]};
  gap: ${theme.space[3]};
`;

const SliderButton = styled(motion.button)`
  background-color: transparent;
  border: 1px solid ${theme.colors.gray300};
  color: ${theme.colors.primary};
  width: 50px;
  height: 50px;
  border-radius: ${theme.borderRadius.full};
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: ${theme.transitions.normal};
  
  &:hover {
    background-color: ${theme.colors.primary};
    color: white;
    border-color: ${theme.colors.primary};
  }
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    
    &:hover {
      background-color: transparent;
      color: ${theme.colors.primary};
      border-color: ${theme.colors.gray300};
    }
  }
`;

const SliderDots = styled.div`
  display: flex;
  justify-content: center;
  gap: ${theme.space[2]};
  margin-top: ${theme.space[4]};
`;

const SliderDot = styled.button<{ active: boolean }>`
  width: 12px;
  height: 12px;
  border-radius: ${theme.borderRadius.full};
  background-color: ${props => props.active ? theme.colors.primary : theme.colors.gray300};
  border: none;
  cursor: pointer;
  transition: ${theme.transitions.normal};
  
  &:hover {
    background-color: ${props => props.active ? theme.colors.primary : theme.colors.gray500};
  }
`;

// Animation variants
const titleVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: { 
    opacity: 1, 
    y: 0,
    transition: {
      duration: 0.5
    }
  }
};

const buttonVariants = {
  initial: { scale: 1 },
  hover: { scale: 1.1 },
  tap: { scale: 0.95 }
};

const CaseStudies: React.FC = () => {
  const [currentSlide, setCurrentSlide] = useState(0);
  
  const nextSlide = () => {
    setCurrentSlide((prev) => (prev + 1) % caseStudies.length);
  };
  
  const prevSlide = () => {
    setCurrentSlide((prev) => (prev - 1 + caseStudies.length) % caseStudies.length);
  };
  
  const goToSlide = (index: number) => {
    setCurrentSlide(index);
  };
  
  return (
    <Section id="case-studies" background="light">
      <SectionTitle
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={titleVariants}
      >
        Success Stories
      </SectionTitle>
      <SectionSubtitle
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={titleVariants}
      >
        Real results from our cloud migration and modernization projects
      </SectionSubtitle>
      
      <SliderContainer>
        <AnimatePresence mode="wait">
          <SliderWrapper
            key={currentSlide}
            initial={{ opacity: 0, x: 100 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -100 }}
            transition={{ duration: 0.5 }}
          >
            <CaseStudyCard caseStudy={caseStudies[currentSlide]} />
          </SliderWrapper>
        </AnimatePresence>
        
        <SliderControls>
          <SliderButton
            onClick={prevSlide}
            disabled={caseStudies.length <= 1}
            variants={buttonVariants}
            initial="initial"
            whileHover="hover"
            whileTap="tap"
          >
            <Icon icon={FiArrowLeft} size={18} />
          </SliderButton>
          <SliderButton
            onClick={nextSlide}
            disabled={caseStudies.length <= 1}
            variants={buttonVariants}
            initial="initial"
            whileHover="hover"
            whileTap="tap"
          >
            <Icon icon={FiArrowRight} size={18} />
          </SliderButton>
        </SliderControls>
        
        <SliderDots>
          {caseStudies.map((study, index) => (
            <SliderDot
              key={study.id}
              active={index === currentSlide}
              onClick={() => goToSlide(index)}
            />
          ))}
        </SliderDots>
      </SliderContainer>
    </Section>
  );
};

export default CaseStudies;