import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { theme } from '../../../styles/theme';
import Section from '../../ui/section';
import PricingCalculator from './PricingCalculator';

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

const PricingContent = styled.div`
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: ${theme.space[6]};
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-template-columns: 1fr;
  }
`;

const PricingInfo = styled(motion.div)`
  grid-column: span 4;
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-column: span 12;
  }
`;

const PricingCalculatorContainer = styled(motion.div)`
  grid-column: span 8;
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-column: span 12;
  }
`;

const InfoBlock = styled(motion.div)`
  background-color: white;
  border-radius: ${theme.borderRadius.lg};
  padding: ${theme.space[6]};
  margin-bottom: ${theme.space[5]};
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
`;

const InfoTitle = styled.h3`
  font-size: ${theme.fontSizes.xl};
  color: ${theme.colors.primary};
  margin-bottom: ${theme.space[4]};
  display: flex;
  align-items: center;
  
  svg {
    margin-right: ${theme.space[3]};
    color: ${theme.colors.secondary};
  }
`;

const InfoText = styled.p`
  color: ${theme.colors.gray700};
  margin-bottom: ${theme.space[3]};
  line-height: 1.6;
`;

const PricingPoint = styled.div`
  margin-bottom: ${theme.space[3]};
  display: flex;
  align-items: center;
`;

const PricingIcon = styled.span`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: ${theme.borderRadius.full};
  background-color: ${theme.colors.success + '20'};
  color: ${theme.colors.success};
  margin-right: ${theme.space[3]};
  flex-shrink: 0;
`;

const SatisfactionBadge = styled(motion.div)`
  background: linear-gradient(135deg, ${theme.colors.primary}, ${theme.colors.accent});
  color: white;
  border-radius: ${theme.borderRadius.lg};
  padding: ${theme.space[5]};
  text-align: center;
`;

const BadgeHeading = styled.h4`
  font-size: ${theme.fontSizes.xl};
  margin-bottom: ${theme.space[2]};
`;

const BadgeText = styled.p`
  opacity: 0.9;
  margin-bottom: 0;
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

const infoBlockVariants = {
  hidden: { opacity: 0, x: -30 },
  visible: (index: number) => ({
    opacity: 1,
    x: 0,
    transition: {
      duration: 0.5,
      delay: 0.2 * index
    }
  })
};

const calculatorVariants = {
  hidden: { opacity: 0, x: 30 },
  visible: {
    opacity: 1,
    x: 0,
    transition: {
      duration: 0.5,
      delay: 0.4
    }
  }
};

const badgeVariants = {
  hidden: { opacity: 0, scale: 0.9 },
  visible: {
    opacity: 1,
    scale: 1,
    transition: {
      duration: 0.5,
      delay: 0.6
    }
  },
  hover: {
    scale: 1.05,
    transition: {
      duration: 0.3
    }
  }
};

const Pricing: React.FC = () => {
  return (
    <Section id="pricing" background="white">
      <SectionTitle
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={titleVariants}
      >
        Transparent Service Pricing
      </SectionTitle>
      <SectionSubtitle
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={titleVariants}
      >
        Fair, transparent pricing for our cloud consulting services with realistic delivery schedules
      </SectionSubtitle>
      
      <PricingContent>
        <PricingInfo>
          <InfoBlock
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
            variants={infoBlockVariants}
            custom={0}
          >
            <InfoTitle>
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>
              Flexible Scheduling
            </InfoTitle>
            <InfoText>
              We offer flexible scheduling that works around your business needs, with availability 
              during weekends and weekday.
            </InfoText>
            <InfoText>
              Our calculator provides transparent estimates based on our flexible schedule, with clear 
              project timelines and milestones. We're committed to delivering quality results on a schedule 
              that works for everyone.
            </InfoText>
          </InfoBlock>
          
          <InfoBlock
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
            variants={infoBlockVariants}
            custom={1}
          >
            <InfoTitle>
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
              What's Included
            </InfoTitle>
            <PricingPoint>
              <PricingIcon>✓</PricingIcon>
              <div>Focused, personalized attention to your needs</div>
            </PricingPoint>
            <PricingPoint>
              <PricingIcon>✓</PricingIcon>
              <div>Clear scope that fits project timeframes</div>
            </PricingPoint>
            <PricingPoint>
              <PricingIcon>✓</PricingIcon>
              <div>Professional expertise</div>
            </PricingPoint>
            <PricingPoint>
              <PricingIcon>✓</PricingIcon>
              <div>Transparent delivery timeline expectations</div>
            </PricingPoint>
            <PricingPoint>
              <PricingIcon>✓</PricingIcon>
              <div>Documentation and knowledge transfer</div>
            </PricingPoint>
            <PricingPoint>
              <PricingIcon>✓</PricingIcon>
              <div>Two weeks of follow-up email support</div>
            </PricingPoint>
          </InfoBlock>
          
          <SatisfactionBadge
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
            variants={badgeVariants}
            whileHover="hover"
          >
            <BadgeHeading>Dedicated Attention Guarantee</BadgeHeading>
            <BadgeText>
              Our flexible schedule means your projects get our undivided attention when we're working together,
              resulting in efficient, high-quality deliverables.
            </BadgeText>
          </SatisfactionBadge>
        </PricingInfo>
        
        <PricingCalculatorContainer
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true }}
          variants={calculatorVariants}
        >
          <PricingCalculator />
        </PricingCalculatorContainer>
      </PricingContent>
    </Section>
  );
};

const MemoizedPricing = React.memo(Pricing);
export default MemoizedPricing;