import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { theme } from '../../../styles/theme';

interface ProjectResult {
  certifications?: string;
  professionalExperience?: string;
  technologiesMastered?: string;
  continuousLearning?: string;
  technicalSkills?: string;
  infrastructureKnowledge?: string;
  certificationLevel?: string;
  problemSolvingApproach?: string;
  [key: string]: string | undefined;
}

interface Testimonial {
  quote: string;
  author: string;
  position: string;
}

interface ProjectHighlight {
  id: number;
  title: string;
  industry: string;
  description: string;
  challenges: string[];
  solution: string;
  results: ProjectResult;
  testimonial: Testimonial;
  image: string;
}

interface ProjectHighlightCardProps {
  projectHighlight: ProjectHighlight;
}

const CardContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: ${theme.space[6]};
  background-color: white;
  border-radius: ${theme.borderRadius.xl};
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  width: 100%;
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-template-columns: 1fr;
  }
`;

const ImageContainer = styled.div`
  position: relative;
  min-height: 400px;
  
  @media (max-width: ${theme.breakpoints.lg}) {
    min-height: 300px;
  }
`;

const Image = styled.img`
  width: 100%;
  height: 100%;
  object-fit: cover;
  position: absolute;
  top: 0;
  left: 0;
`;

const Overlay = styled.div`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    to bottom,
    rgba(${props => props.theme.colors.primary.replace('#', '')}, 0.7),
    rgba(${props => props.theme.colors.dark.replace('#', '')}, 0.9)
  );
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: ${theme.space[8]};
  color: white;
`;

const IndustryTag = styled(motion.span)`
  background-color: ${theme.colors.secondary};
  color: ${theme.colors.primary};
  padding: ${theme.space[1]} ${theme.space[3]};
  border-radius: ${theme.borderRadius.full};
  font-size: ${theme.fontSizes.sm};
  font-weight: ${theme.fontWeights.medium};
  display: inline-block;
  margin-bottom: ${theme.space[3]};
`;

const ProjectTitle = styled.h3`
  font-size: ${theme.fontSizes['3xl']};
  font-weight: ${theme.fontWeights.bold};
  margin-bottom: ${theme.space[4]};
`;

const ProjectDescription = styled.p`
  font-size: ${theme.fontSizes.lg};
  margin-bottom: ${theme.space[4]};
  opacity: 0.9;
`;

const Testimonial = styled(motion.div)`
  background-color: rgba(255, 255, 255, 0.1);
  padding: ${theme.space[4]};
  border-radius: ${theme.borderRadius.lg};
  margin-top: auto;
`;

const Quote = styled.blockquote`
  font-style: italic;
  margin-bottom: ${theme.space[3]};
  font-size: ${theme.fontSizes.md};
`;

const Author = styled.p`
  font-weight: ${theme.fontWeights.medium};
  font-size: ${theme.fontSizes.sm};
  
  span {
    color: ${theme.colors.secondary};
  }
`;

const ContentContainer = styled.div`
  padding: ${theme.space[8]};
  display: flex;
  flex-direction: column;
  
  @media (max-width: ${theme.breakpoints.md}) {
    padding: ${theme.space[5]};
  }
`;

const SectionTitle = styled.h4`
  font-size: ${theme.fontSizes.lg};
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.primary};
  margin-bottom: ${theme.space[3]};
  display: flex;
  align-items: center;
  
  &::before {
    content: '';
    display: inline-block;
    width: 4px;
    height: 20px;
    background-color: ${theme.colors.secondary};
    margin-right: ${theme.space[3]};
    border-radius: ${theme.borderRadius.md};
  }
`;

const ChallengesList = styled.ul`
  margin-bottom: ${theme.space[5]};
  padding-left: ${theme.space[5]};
`;

const ChallengeItem = styled(motion.li)`
  margin-bottom: ${theme.space[2]};
  position: relative;
  color: ${theme.colors.gray700};
  
  &::before {
    content: '•';
    color: ${theme.colors.secondary};
    position: absolute;
    left: -${theme.space[5]};
    font-size: ${theme.fontSizes.xl};
  }
`;

const Solution = styled.p`
  color: ${theme.colors.gray700};
  margin-bottom: ${theme.space[5]};
  line-height: 1.6;
`;

const ResultsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: ${theme.space[4]};
  margin-top: ${theme.space[4]};
`;

const ResultItem = styled(motion.div)`
  background-color: ${theme.colors.gray100};
  padding: ${theme.space[4]};
  border-radius: ${theme.borderRadius.lg};
  text-align: center;
`;

const ResultValue = styled.div`
  font-size: ${theme.fontSizes['2xl']};
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.accent};
  margin-bottom: ${theme.space[1]};
`;

const ResultLabel = styled.div`
  font-size: ${theme.fontSizes.sm};
  color: ${theme.colors.gray600};
`;

// Animation variants
const testimonialVariants = {
  initial: { opacity: 0, y: 20 },
  animate: { 
    opacity: 1, 
    y: 0,
    transition: {
      duration: 0.5,
      delay: 0.5
    }
  }
};

const resultVariants = {
  initial: { scale: 0.9, opacity: 0 },
  animate: (index: number) => ({
    scale: 1,
    opacity: 1,
    transition: {
      duration: 0.3,
      delay: 0.1 * index
    }
  })
};

const challengeVariants = {
  initial: { opacity: 0, x: -10 },
  animate: (index: number) => ({
    opacity: 1,
    x: 0,
    transition: {
      duration: 0.3,
      delay: 0.1 * index
    }
  })
};

const ProjectHighlightCard: React.FC<ProjectHighlightCardProps> = ({ projectHighlight }) => {
  const resultsArray = Object.entries(projectHighlight.results).map(([key, value], index) => ({
    key,
    value,
    label: key
      .replace(/([A-Z])/g, ' $1')
      .replace(/^./, str => str.toUpperCase())
      .replace(/Time/g, 'Time Reduction')
  }));
  
  return (
    <CardContainer>
      <ImageContainer>
        <Image src={projectHighlight.image} alt={projectHighlight.title} />
        <Overlay>
          <IndustryTag
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5 }}
          >
            {projectHighlight.industry}
          </IndustryTag>
          <ProjectTitle>{projectHighlight.title}</ProjectTitle>
          <ProjectDescription>{projectHighlight.description}</ProjectDescription>
          <Testimonial
            variants={testimonialVariants}
            initial="initial"
            animate="animate"
          >
            <Quote>"{projectHighlight.testimonial.quote}"</Quote>
            <Author>
              {projectHighlight.testimonial.author} • <span>{projectHighlight.testimonial.position}</span>
            </Author>
          </Testimonial>
        </Overlay>
      </ImageContainer>
      
      <ContentContainer>
        <SectionTitle>Professional Challenges</SectionTitle>
        <ChallengesList>
          {projectHighlight.challenges.map((challenge, index) => (
            <ChallengeItem
              key={index}
              variants={challengeVariants}
              initial="initial"
              animate="animate"
              custom={index}
            >
              {challenge}
            </ChallengeItem>
          ))}
        </ChallengesList>
        
        <SectionTitle>Professional Growth Strategy</SectionTitle>
        <Solution>{projectHighlight.solution}</Solution>
        
        <SectionTitle>Professional Achievements</SectionTitle>
        <ResultsGrid>
          {resultsArray.map((result, index) => (
            <ResultItem
              key={result.key}
              variants={resultVariants}
              initial="initial"
              animate="animate"
              custom={index}
            >
              <ResultValue>{result.value}</ResultValue>
              <ResultLabel>{result.label}</ResultLabel>
            </ResultItem>
          ))}
        </ResultsGrid>
      </ContentContainer>
    </CardContainer>
  );
};

export default ProjectHighlightCard;