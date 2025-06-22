import React, { useState } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import { theme } from '../../../styles/theme';
import { FiChevronDown } from 'react-icons/fi';
import Icon from '../../ui/Icon';

interface ServiceItem {
  id: number;
  title: string;
  description: string;
  icon: React.ReactNode;
  color: string;
  features: string[];
}

interface ServiceCardProps {
  service: ServiceItem;
  index: number;
}

const Card = styled(motion.div)`
  background-color: white;
  border-radius: ${theme.borderRadius.lg};
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  padding: ${theme.space[6]};
  height: 100%;
  position: relative;
  overflow: hidden;
  cursor: pointer;
  transition: ${theme.transitions.normal};
  
  &:hover {
    transform: translateY(-8px);
    box-shadow: 0 12px 30px rgba(0, 0, 0, 0.12);
    border-top: 3px solid ${theme.colors.primary};
  }
`;

const IconWrapper = styled.div<{ color: string }>`
  width: 60px;
  height: 60px;
  border-radius: ${theme.borderRadius.full};
  background-color: ${props => props.color + '15'}; // 15% opacity
  color: ${props => props.color};
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: ${theme.fontSizes['2xl']};
  margin-bottom: ${theme.space[4]};
`;

const Title = styled.h3`
  font-size: ${theme.fontSizes.xl};
  font-weight: ${theme.fontWeights.bold};
  margin-bottom: ${theme.space[3]};
  color: ${theme.colors.primary};
`;

const Description = styled.p`
  color: ${theme.colors.gray600};
  margin-bottom: ${theme.space[4]};
  font-size: ${theme.fontSizes.md};
`;

const ExpandButton = styled(motion.button)<{ $isExpanded: boolean }>`
  background: transparent;
  border: none;
  color: ${theme.colors.accent};
  font-size: ${theme.fontSizes.md};
  font-weight: ${theme.fontWeights.medium};
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0;
  margin-top: ${theme.space[3]};
  
  svg {
    margin-left: ${theme.space[2]};
    transform: ${props => props.$isExpanded ? 'rotate(180deg)' : 'rotate(0)'};
    transition: transform 0.3s ease;
  }
`;

const FeatureList = styled(motion.ul)`
  margin-top: ${theme.space[4]};
  padding-top: ${theme.space[4]};
  border-top: 1px solid ${theme.colors.gray200};
`;

const FeatureItem = styled(motion.li)`
  display: flex;
  align-items: flex-start;
  margin-bottom: ${theme.space[3]};
  color: ${theme.colors.gray700};
`;

const Bullet = styled.span<{ color: string }>`
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: ${theme.borderRadius.full};
  background-color: ${props => props.color};
  margin-right: ${theme.space[3]};
  margin-top: 8px;
  flex-shrink: 0;
`;

// Animation variants
const cardVariants = {
  hidden: { opacity: 0, y: 30, scale: 0.96 },
  visible: (i: number) => ({
    opacity: 1,
    y: 0,
    scale: 1,
    transition: {
      duration: 0.6,
      delay: i * 0.13,
      ease: 'easeOut',
    }
  })
};

const ServiceCard: React.FC<ServiceCardProps> = ({ service, index }) => {
  const [expanded, setExpanded] = useState(false);

  return (
    <Card
      as={motion.div}
      initial="hidden"
      animate="visible"
      variants={cardVariants}
      custom={index}
      layoutId={`service-card-${service.id}`}
      id={`service-${service.id}`}
    >
      <IconWrapper color={service.color}>
        {service.icon}
      </IconWrapper>
      <Title>{service.title}</Title>
      <Description>{service.description}</Description>
      <ExpandButton 
        onClick={e => { e.stopPropagation(); setExpanded(v => !v); }}
        $isExpanded={expanded}
        as={motion.button}
        aria-expanded={expanded}
      >
        Learn more <Icon icon={FiChevronDown} size={16} />
      </ExpandButton>
      <AnimatePresence>
        {expanded && (
          <FeatureList
            initial="hidden"
            animate="visible"
            exit="hidden"
            variants={{ 
              hidden: { opacity: 0, height: 0 },
              visible: { 
                opacity: 1, 
                height: 'auto',
                transition: {
                  duration: 0.3,
                  when: "beforeChildren",
                  staggerChildren: 0.08
                }
              }
            }}
          >
            {service.features.map((feature, idx) => (
              <FeatureItem key={idx} variants={{ 
                hidden: { opacity: 0, x: -20 },
                visible: { 
                  opacity: 1, 
                  x: 0,
                  transition: {
                    duration: 0.3
                  }
                }
              }}>
                <Bullet color={service.color} />
                {feature}
              </FeatureItem>
            ))}
          </FeatureList>
        )}
      </AnimatePresence>
    </Card>
  );
};

const MemoizedServiceCard = React.memo(ServiceCard);
export default MemoizedServiceCard;