import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { theme } from '../../../styles/theme';
import Button from '../../ui/Button';
import CloudAnimation from './CloudAnimation';
import { FiArrowRight } from 'react-icons/fi';
import Icon from '../../ui/Icon';

const HeroContainer = styled.section`
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  background: linear-gradient(135deg, ${theme.colors.primary} 0%, ${theme.colors.dark} 100%);
  overflow: hidden;
  padding: ${theme.space[8]} 0;
  
  @media (max-width: ${theme.breakpoints.md}) {
    padding-top: ${theme.space[10]};
  }
`;

const HeroContent = styled.div`
  max-width: ${theme.sizes.container.xl};
  margin: 0 auto;
  padding: 0 ${theme.space[4]};
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: ${theme.space[6]};
  position: relative;
  z-index: 2;
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-template-columns: 1fr;
    text-align: center;
  }
`;

const TextContent = styled(motion.div)`
  display: flex;
  flex-direction: column;
  justify-content: center;
  color: ${theme.colors.white};
`;

const AnimationContent = styled(motion.div)`
  display: flex;
  align-items: center;
  justify-content: center;
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-row: 1;
    margin-bottom: ${theme.space[4]};
  }
`;

const Heading = styled(motion.h1)`
  font-size: clamp(${theme.fontSizes['4xl']}, 5vw, ${theme.fontSizes['6xl']});
  font-weight: ${theme.fontWeights.bold};
  line-height: 1.1;
  margin-bottom: ${theme.space[4]};
  color: ${theme.colors.white};
  
  span {
    color: ${theme.colors.secondary};
    position: relative;
    display: inline-block;
    
    &::after {
      content: '';
      position: absolute;
      width: 100%;
      height: 6px;
      bottom: 5px;
      left: 0;
      background-color: ${theme.colors.secondary};
      opacity: 0.3;
      z-index: -1;
    }
  }
`;

const Subheading = styled(motion.p)`
  font-size: ${theme.fontSizes.xl};
  color: ${theme.colors.gray300};
  margin-bottom: ${theme.space[6]};
  max-width: 600px;
  line-height: ${theme.lineHeights.loose};
  
  @media (max-width: ${theme.breakpoints.lg}) {
    margin-left: auto;
    margin-right: auto;
  }
`;

const CTAButtons = styled(motion.div)`
  display: flex;
  gap: ${theme.space[4]};
  
  @media (max-width: ${theme.breakpoints.lg}) {
    justify-content: center;
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    flex-direction: column;
  }
`;

const BackgroundShapes = styled.div`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 1;
`;

const Shape = styled(motion.div)`
  position: absolute;
  background-color: rgba(255, 255, 255, 0.03);
  border-radius: 50%;
`;

const Shape1 = styled(Shape)`
  width: 400px;
  height: 400px;
  top: -200px;
  left: -100px;
`;

const Shape2 = styled(Shape)`
  width: 300px;
  height: 300px;
  bottom: -150px;
  right: -50px;
`;

const Shape3 = styled(Shape)`
  width: 200px;
  height: 200px;
  top: 40%;
  right: 10%;
`;

const Stats = styled(motion.div)`
  display: flex;
  gap: ${theme.space[8]};
  margin-top: ${theme.space[8]};
  
  @media (max-width: ${theme.breakpoints.lg}) {
    justify-content: center;
  }
  
  @media (max-width: ${theme.breakpoints.md}) {
    flex-direction: column;
    gap: ${theme.space[4]};
  }
`;

const StatItem = styled(motion.div)`
  display: flex;
  flex-direction: column;
`;

const StatNumber = styled.div`
  font-size: ${theme.fontSizes['3xl']};
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.secondary};
  margin-bottom: ${theme.space[1]};
`;

const StatLabel = styled.div`
  font-size: ${theme.fontSizes.sm};
  color: ${theme.colors.gray300};
`;

// Animation variants
const textVariants = {
  hidden: { opacity: 0, y: 30 },
  visible: { 
    opacity: 1, 
    y: 0,
    transition: {
      duration: 0.8,
      ease: "easeOut"
    }
  }
};

const statsVariants = {
  hidden: { opacity: 0 },
  visible: { 
    opacity: 1,
    transition: {
      staggerChildren: 0.2,
      delayChildren: 0.8
    }
  }
};

const statItemVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.5
    }
  }
};

const shapeVariants = {
  animate: {
    y: [0, -10, 0],
    transition: {
      duration: 3,
      repeat: Infinity,
      repeatType: 'loop' as const,
      ease: 'easeInOut'
    }
  }
};

const Hero: React.FC = () => {
  return (
    <HeroContainer id="home">
      <BackgroundShapes>
        <Shape1 variants={shapeVariants} animate="animate" />
        <Shape2 variants={shapeVariants} animate="animate" custom={1} />
        <Shape3 variants={shapeVariants} animate="animate" custom={2} />
      </BackgroundShapes>
      
      <HeroContent>
        <TextContent
          initial="hidden"
          animate="visible"
          variants={textVariants}
        >
          <Heading>
            Seamless <span>Cloud Migration</span> Solutions for Your Business
          </Heading>
          <Subheading>
            Transform your infrastructure with our expert-led cloud migration services.
            AWS certified specialists ensuring security, efficiency, and cost optimization.
          </Subheading>
          
          <CTAButtons>
          <Button 
            size="lg" 
            icon={<Icon icon={FiArrowRight} size={18} />} 
            iconPosition="right"
            onClick={() => {
              console.log("Get Started clicked");
              document.getElementById('contact')?.scrollIntoView({ 
                behavior: 'smooth' 
              });
            }}
          >
            Get Started
          </Button>
          
          <Button 
            size="lg" 
            variant="outline"
            onClick={() => {
              console.log("Learn More clicked");
              // Using the native scrollIntoView method
              document.getElementById('services')?.scrollIntoView({ 
                behavior: 'smooth' 
              });
            }}
          >
            Learn More
          </Button>
        </CTAButtons>
          
          <Stats variants={statsVariants}>
            <StatItem variants={statItemVariants}>
              <StatNumber>100+</StatNumber>
              <StatLabel>Successful Migrations</StatLabel>
            </StatItem>
            <StatItem variants={statItemVariants}>
              <StatNumber>40%</StatNumber>
              <StatLabel>Avg. Cost Reduction</StatLabel>
            </StatItem>
            <StatItem variants={statItemVariants}>
              <StatNumber>24/7</StatNumber>
              <StatLabel>Expert Support</StatLabel>
            </StatItem>
          </Stats>
        </TextContent>
        
        <AnimationContent
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.8, delay: 0.2 }}
        >
          <CloudAnimation />
        </AnimationContent>
      </HeroContent>
    </HeroContainer>
  );
};

export default Hero;