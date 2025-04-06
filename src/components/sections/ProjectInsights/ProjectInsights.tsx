import React, { useState } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import { theme } from '../../../styles/theme';
import Section from '../../ui/Section';
import { FiArrowRight, FiArrowLeft } from 'react-icons/fi';
import ProjectHighlightCard from './ProjectInsightsCard';
import Icon from '../../ui/Icon';

const projectHighlights = [
  {
    id: 1,
    title: 'Our Professional Background',
    industry: 'Cloud Expertise',
    description: 'Learn about our professional experience and AWS certifications that form the foundation of our cloud consulting services.',
    challenges: [
      'Providing enterprise-level expertise to small businesses',
      'Delivering personalized, focused cloud solutions',
      'Creating flexible project timelines that work for clients',
      'Bringing professional cloud experience to each project'
    ],
    solution: 'We bring our 5+ combined AWS certifications and years of professional cloud experience to help businesses with targeted cloud solutions. Our flexible scheduling means we can adapt to your specific timeline needs.',
    results: {
      certifications: '5+ AWS Certifications',
      professionalExperience: 'Day Job Expertise',
      approach: 'Personalized Service',
      availability: 'Flexible Schedule'
    },
    testimonial: {
      quote: "We're passionate about cloud technology and love helping businesses leverage it effectively with our flexible schedule that adapts to your needs.",
      author: "Cloud Partners Team",
      position: "Cloud Partners Founders"
    },
    image: 'https://images.unsplash.com/photo-1519389950473-47ba0277781c?ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80'
  },
  {
    id: 2,
    title: 'Cloud Architecture Expertise',
    industry: 'Technical Specialization',
    description: "Our deep knowledge of AWS architecture and infrastructure design helps businesses build scalable, secure, and cost-effective cloud environments.",
    challenges: [
      "Complex infrastructure requirements",
      "Cost optimization across multiple services",
      "Security and compliance implementation",
      "Performance and scalability planning"
    ],
    solution: "We leverage our extensive AWS certification knowledge and practical experience to design robust cloud architectures tailored to your specific business requirements and growth plans.",
    results: {
      architectureReview: "30% Cost Savings",
      securityAssessment: "Compliance Ready",
      performanceTuning: "2x Faster Apps",
      cloudStrategy: "Future-Proof Design"
    },
    testimonial: {
      quote: "Our approach combines technical expertise with a practical understanding of business needs, resulting in cloud solutions that deliver real value for your organization.",
      author: "Cloud Partners Team",
      position: "AWS Certified Professionals"
    },
    image: "https://images.unsplash.com/photo-1488590528505-98d2b5aba04b?ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80"
  }
];

const SectionTitle = styled(motion.h2)`
  text-align: center;
  margin-bottom: ${theme.space[4]};
  font-size: clamp(${theme.fontSizes['2xl']}, 5vw, ${theme.fontSizes['4xl']});
  color: ${theme.colors.primary};
  
  @media (max-width: ${theme.breakpoints.md}) {
    margin-bottom: ${theme.space[3]};
  }
`;

const SliderContainer = styled.div`
  position: relative;
  overflow: hidden;
  padding: ${theme.space[4]} 0;
  max-width: 1200px;
  margin: 0 auto;
  
  @media (max-width: ${theme.breakpoints.lg}) {
    padding: ${theme.space[3]} ${theme.space[3]};
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    padding: ${theme.space[2]};
    overflow-x: hidden;
    width: 100%;
    box-sizing: border-box;
  }
`;

const SliderWrapper = styled(motion.div)`
  display: flex;
  width: 100%;
  box-sizing: border-box;
  overflow: hidden;
  
  @media (max-width: ${theme.breakpoints.sm}) {
    padding: 0;
  }
`;

const SliderControls = styled.div`
  display: flex;
  justify-content: center;
  margin-top: ${theme.space[6]};
  gap: ${theme.space[3]};
  
  @media (max-width: ${theme.breakpoints.md}) {
    margin-top: ${theme.space[4]};
  }
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
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  
  &:hover {
    background-color: ${theme.colors.accent};
    color: white;
    border-color: ${theme.colors.accent};
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  }
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    
    &:hover {
      background-color: transparent;
      color: ${theme.colors.primary};
      border-color: ${theme.colors.gray300};
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    }
  }
  
  @media (max-width: ${theme.breakpoints.md}) {
    width: 40px;
    height: 40px;
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
  background-color: ${props => props.active ? theme.colors.accent : theme.colors.gray300};
  border: none;
  cursor: pointer;
  transition: ${theme.transitions.normal};
  box-shadow: ${props => props.active ? '0 2px 4px rgba(0, 115, 187, 0.3)' : 'none'};
  
  &:hover {
    background-color: ${props => props.active ? theme.colors.accent : theme.colors.primary};
    transform: ${props => props.active ? 'scale(1.2)' : 'scale(1.1)'};
  }
`;

// Animation variants
const titleVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: { 
    opacity: 1, 
    y: 0,
    transition: {
      duration: 0.5,
      ease: "easeOut"
    }
  }
};

const buttonVariants = {
  initial: { scale: 1 },
  hover: { scale: 1.1 },
  tap: { scale: 0.95 }
};
const ProjectInsights: React.FC = () => {
  const [currentSlide, setCurrentSlide] = useState(0);
  
  const nextSlide = () => {
    setCurrentSlide((prev) => (prev + 1) % projectHighlights.length);
  };
  
  const prevSlide = () => {
    setCurrentSlide((prev) => (prev - 1 + projectHighlights.length) % projectHighlights.length);
  };
  
  const goToSlide = (index: number) => {
    setCurrentSlide(index);
  };
  
  return (
    <Section id="project-insights" background="light">
      <SectionTitle
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={titleVariants}
      >
        Our Professional Expertise
      </SectionTitle>
      
      <SliderContainer>
        <AnimatePresence mode="wait">
          <SliderWrapper
            key={currentSlide}
            initial={{ opacity: 0, x: 50 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -50 }}
            transition={{ duration: 0.4, ease: "easeInOut" }}
          >
            <ProjectHighlightCard projectHighlight={projectHighlights[currentSlide]} />
          </SliderWrapper>
        </AnimatePresence>
        
        <SliderControls>
          <SliderButton
            onClick={prevSlide}
            disabled={projectHighlights.length <= 1}
            variants={buttonVariants}
            initial="initial"
            whileHover="hover"
            whileTap="tap"
          >
            <Icon icon={FiArrowLeft} size={18} />
          </SliderButton>
          <SliderButton
            onClick={nextSlide}
            disabled={projectHighlights.length <= 1}
            variants={buttonVariants}
            initial="initial"
            whileHover="hover"
            whileTap="tap"
          >
            <Icon icon={FiArrowRight} size={18} />
          </SliderButton>
        </SliderControls>
        
        <SliderDots>
          {projectHighlights.map((highlight, index) => (
            <SliderDot
              key={highlight.id}
              active={index === currentSlide}
              onClick={() => goToSlide(index)}
            />
          ))}
        </SliderDots>
      </SliderContainer>
    </Section>
  );
};

export default ProjectInsights;