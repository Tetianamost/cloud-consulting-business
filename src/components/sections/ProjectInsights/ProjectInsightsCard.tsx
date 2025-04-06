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

interface TestimonialData {
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
  testimonial: TestimonialData;
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
  max-width: 100%;
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-template-columns: 1fr;
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    border-radius: ${theme.borderRadius.lg};
  }
`;

const ImageContainer = styled.div`
  position: relative;
  min-height: 400px;
  width: 100%;
  box-sizing: border-box;
  
  @media (max-width: ${theme.breakpoints.lg}) {
    min-height: 350px;
  }
  
  @media (max-width: ${theme.breakpoints.md}) {
    min-height: 300px;
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    min-height: 250px;
    max-height: 400px;
    overflow-y: auto;
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
  overflow: auto;

  background: linear-gradient(
    to bottom,
    rgba(35, 47, 62, 0.8),
    rgba(22, 30, 45, 0.95)
  );
  
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: ${theme.space[8]};
  color: white;
  box-sizing: border-box;
  
  @media (max-width: ${theme.breakpoints.md}) {
    padding: ${theme.space[5]};
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    padding: ${theme.space[4]};
    justify-content: flex-start;
    overflow-y: auto;
  }
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
  font-size: clamp(${theme.fontSizes.lg}, 5vw, ${theme.fontSizes['3xl']});
  font-weight: ${theme.fontWeights.bold};
  margin-bottom: ${theme.space[4]};
  color: ${theme.colors.white};
  text-shadow: 0 1px 3px rgba(0,0,0,0.3);
  word-wrap: break-word;
  
  @media (max-width: ${theme.breakpoints.md}) {
    font-size: ${theme.fontSizes.xl};
    margin-bottom: ${theme.space[3]};
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    font-size: ${theme.fontSizes.lg};
    margin-bottom: ${theme.space[2]};
    margin-top: ${theme.space[3]};
  }
`;

const ProjectDescription = styled.p`
  font-size: clamp(${theme.fontSizes.sm}, 4vw, ${theme.fontSizes.lg});
  margin-bottom: ${theme.space[4]};
  opacity: 0.95;
  text-shadow: 1px 1px 2px rgba(0,0,0,0.5);
  line-height: 1.6;
  word-wrap: break-word;
  
  @media (max-width: ${theme.breakpoints.md}) {
    font-size: ${theme.fontSizes.md};
    margin-bottom: ${theme.space[3]};
    line-height: 1.5;
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    font-size: ${theme.fontSizes.sm};
    margin-bottom: ${theme.space[3]};
    line-height: 1.4;
  }
`;

const Testimonial = styled(motion.div)`
  background-color: rgba(0, 0, 0, 0.5);
  padding: ${theme.space[4]};
  border-radius: ${theme.borderRadius.lg};
  margin-top: auto;
  max-width: 100%;
  box-sizing: border-box;
  word-wrap: break-word;
  
  @media (max-width: ${theme.breakpoints.md}) {
    padding: ${theme.space[3]};
    margin-top: ${theme.space[3]};
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    padding: ${theme.space[3]};
    margin-top: ${theme.space[2]};
    border-radius: ${theme.borderRadius.md};
  }
`;

const Quote = styled.blockquote`
  font-style: italic;
  margin-bottom: ${theme.space[3]};
  font-size: clamp(${theme.fontSizes.sm}, 3vw, ${theme.fontSizes.md});
  line-height: 1.5;
  word-wrap: break-word;
  margin: 0 0 ${theme.space[3]} 0;
  
  @media (max-width: ${theme.breakpoints.md}) {
    margin-bottom: ${theme.space[2]};
    font-size: ${theme.fontSizes.sm};
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    font-size: ${theme.fontSizes.xs};
    line-height: 1.4;
  }
`;

const Author = styled.p`
  font-weight: ${theme.fontWeights.medium};
  font-size: ${theme.fontSizes.sm};
  margin: 0;
  word-wrap: break-word;
  
  span {
    color: ${theme.colors.secondary};
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    font-size: ${theme.fontSizes.xs};
  }
`;

const ContentContainer = styled.div`
  padding: ${theme.space[8]};
  display: flex;
  flex-direction: column;
  
  @media (max-width: ${theme.breakpoints.md}) {
    padding: ${theme.space[5]};
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    padding: ${theme.space[4]};
  }
`;

const SectionTitle = styled.h4`
  font-size: clamp(${theme.fontSizes.md}, 4vw, ${theme.fontSizes.lg});
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.accent};
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
  
  @media (max-width: ${theme.breakpoints.md}) {
    font-size: ${theme.fontSizes.md};
    margin-bottom: ${theme.space[2]};
    
    &::before {
      height: 16px;
      margin-right: ${theme.space[2]};
    }
  }
`;

const ChallengesList = styled.ul`
  margin-bottom: ${theme.space[5]};
  padding-left: ${theme.space[5]};
  
  @media (max-width: ${theme.breakpoints.md}) {
    margin-bottom: ${theme.space[4]};
    padding-left: ${theme.space[4]};
  }
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
  font-size: clamp(${theme.fontSizes.sm}, 3vw, ${theme.fontSizes.md});
  
  @media (max-width: ${theme.breakpoints.md}) {
    margin-bottom: ${theme.space[4]};
  }
`;

const ResultsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: ${theme.space[4]};
  margin-top: ${theme.space[4]};
  width: 100%;
  box-sizing: border-box;
  
  @media (max-width: ${theme.breakpoints.sm}) {
    gap: ${theme.space[2]};
    margin-top: ${theme.space[3]};
    grid-template-columns: 1fr;
  }
`;

const ResultItem = styled(motion.div)`
  background-color: ${theme.colors.info}15;
  padding: ${theme.space[4]};
  border-radius: ${theme.borderRadius.lg};
  text-align: center;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
  width: 100%;
  box-sizing: border-box;
  
  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.1);
  }
  
  @media (max-width: ${theme.breakpoints.md}) {
    padding: ${theme.space[3]};
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    padding: ${theme.space[3]};
    margin-bottom: ${theme.space[2]};
    
    &:hover {
      transform: translateY(-3px);
    }
  }
`;

const ResultValue = styled.div`
  font-size: clamp(${theme.fontSizes.xl}, 5vw, ${theme.fontSizes['2xl']});
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.accent};
  margin-bottom: ${theme.space[1]};
  word-break: break-word;
  
  @media (max-width: ${theme.breakpoints.md}) {
    font-size: ${theme.fontSizes.lg};
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    font-size: ${theme.fontSizes.lg};
  }
`;

const ResultLabel = styled.div`
  font-size: ${theme.fontSizes.sm};
  color: ${theme.colors.gray700};
  font-weight: ${theme.fontWeights.medium};
  line-height: 1.4;
  word-break: break-word;
  
  @media (max-width: ${theme.breakpoints.sm}) {
    font-size: ${theme.fontSizes.xs};
  }
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
          <div style={{ marginTop: 'auto', width: '100%', maxWidth: '100%' }}>
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
          </div>
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