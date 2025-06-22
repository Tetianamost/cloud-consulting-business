import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { theme } from '../../../styles/theme';
import { IconBaseProps, IconType } from 'react-icons';

interface CertificationItem {
  id: number;
  title: string;
  description: string;
  image: IconType;
  color: string;
}

interface CertificationBadgeProps {
  certification: CertificationItem;
  index: number;
}

const Badge = styled(motion.div)`
  background-color: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border-radius: ${theme.borderRadius.lg};
  padding: ${theme.space[5]};
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  transition: ${theme.transitions.normal};
  cursor: pointer;
  position: relative;
  overflow: hidden;
  height: 100%;
  
  &:hover {
    transform: translateY(-8px);
    background-color: rgba(255, 255, 255, 0.1);
  }
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 4px;
    background: linear-gradient(90deg, ${theme.colors.secondary}, ${theme.colors.highlight});
    opacity: 0;
    transition: opacity 0.3s ease;
  }
  
  &:hover::before {
    opacity: 1;
  }
`;

const BadgeImage = styled(motion.div)`
  width: 120px;
  height: 120px;
  margin-bottom: ${theme.space[4]};
  position: relative;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: contain;
    transition: transform 0.3s ease;
  }
`;

const Glow = styled(motion.div)`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, ${props => props.color + '40'} 0%, transparent 70%);
  opacity: 0.5;
  z-index: -1;
  filter: blur(10px);
`;

const BadgeTitle = styled.h3`
  font-size: ${theme.fontSizes.lg};
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.white};
  margin-bottom: ${theme.space[3]};
`;

const BadgeDescription = styled.p`
  color: ${theme.colors.gray300};
  font-size: ${theme.fontSizes.md};
  line-height: 1.6;
`;

// Animation variants
const badgeVariants = {
  hidden: { opacity: 0, y: 30 },
  visible: (index: number) => ({
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.5,
      delay: 0.1 * index
    }
  })
};

const imageVariants = {
  initial: { rotateY: 0 },
  hover: { 
    rotateY: 360,
    transition: {
      duration: 1.2,
      ease: "easeInOut"
    }
  }
};

const glowVariants = {
  initial: { 
    scale: 0.95, 
    opacity: 0.7 
  },
  hover: { 
    scale: 1.05, 
    opacity: 1, 
    transition: { 
      duration: 0.3, 
      repeat: 0, 
      repeatType: 'loop' as const 
    } 
  }
};

const CertificationBadge: React.FC<CertificationBadgeProps> = ({ certification, index }) => {
  // Explicitly type the icon with IconType and additional props
  const IconComponent = certification.image as React.ComponentType<IconBaseProps>;

  return (
    <Badge
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true}}
      variants={badgeVariants}
      custom={index}
      whileHover="hover"
    >
      <BadgeImage>
        <motion.div variants={imageVariants}>
          <IconComponent size={120} color={certification.color} />
        </motion.div>
        <Glow
          color={certification.color}
          variants={glowVariants}
        />
      </BadgeImage>
      <BadgeTitle>{certification.title}</BadgeTitle>
      <BadgeDescription>{certification.description}</BadgeDescription>
    </Badge>
  );
};

export default CertificationBadge;