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
    title: 'Cloud Architecture Learning Journey',
    industry: 'Professional Development',
    description: 'Deep dive into AWS cloud infrastructure, focusing on EKS, database design, and scalable solutions through hands-on professional experience.',
    challenges: [
      'Understanding complex kubernetes architectures',
      'Mastering multi-tier database designs',
      'Implementing scalable cloud solutions',
      'Keeping up with rapidly evolving cloud technologies'
    ],
    solution: 'Pursued and obtained 5 AWS certifications, gaining in-depth knowledge through professional work experience. Focused on hands-on learning in EKS environments, relational databases, and cloud infrastructure design.',
    results: {
      certifications: '5 AWS Certifications',
      professionalExperience: '2 Years',
      technologiesMastered: '10+',
      continuousLearning: 'Ongoing',
    },
    testimonial: {
      quote: "Continuous learning and practical experience are the keys to mastering cloud technologies. Each certification and project brings new insights.",
      author: "Tetiana Mostova",
      position: "Cloud Infrastructure Consultant"
    },
    image: 'https://images.unsplash.com/photo-1519389950473-47ba0277781c?ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80'
  },
  {
    id: 2,
    title: 'Professional Cloud Infrastructure Expertise',
    industry: 'Cloud Consulting',
    description: 'Developed comprehensive understanding of AWS cloud infrastructure through professional work and continuous learning.',
    challenges: [
      'Navigating complex cloud environments',
      'Understanding enterprise-level infrastructure needs',
      'Balancing performance and cost-effectiveness',
      'Staying current with cloud innovations'
    ],
    solution: 'Leveraged professional experience in a cloud-focused role to gain hands-on expertise with EKS, database management, and infrastructure design. Continuously expanded knowledge through certifications and practical application.',
    results: {
      technicalSkills: 'Advanced',
      infrastructureKnowledge: 'Comprehensive',
      certificationLevel: 'Professional',
      problemSolvingApproach: 'Strategic',
    },
    testimonial: {
      quote: "True expertise comes from a combination of formal learning and real-world problem-solving.",
      author: "Tetiana Mostova",
      position: "Cloud Developer"
    },
    image: 'https://images.unsplash.com/photo-1488590528505-98d2b5aba04b?ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80'
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
        Professional Journey
      </SectionTitle>
      
      <SliderContainer>
        <AnimatePresence mode="wait">
          <SliderWrapper
            key={currentSlide}
            initial={{ opacity: 0, x: 100 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -100 }}
            transition={{ duration: 0.5 }}
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