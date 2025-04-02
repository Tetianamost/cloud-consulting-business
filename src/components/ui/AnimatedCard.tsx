import React from 'react';
import styled from 'styled-components';
import { motion, MotionProps } from 'framer-motion';
import { theme } from '../../styles/theme';

interface CardProps extends MotionProps {
  elevation?: 'sm' | 'md' | 'lg';
  interactive?: boolean;
  padding?: keyof typeof theme.space;
  borderRadius?: keyof typeof theme.borderRadius;
  background?: string;
  children: React.ReactNode;
}

const CardContainer = styled(motion.div)<CardProps>`
  background: ${props => props.background || theme.colors.white};
  border-radius: ${props =>
    theme.borderRadius[props.borderRadius || 'lg']};
  padding: ${props => theme.space[props.padding || 5]};
  box-shadow: ${props => theme.shadows[props.elevation || 'md']};
  transition: ${theme.transitions.normal};
  height: 100%;
  overflow: hidden;
  
  ${props =>
    props.interactive &&
    `
    cursor: pointer;
    &:hover {
      box-shadow: ${theme.shadows.lg};
    }
  `}
`;

// Animation variants
const cardVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: { 
    opacity: 1, 
    y: 0,
    transition: {
      duration: 0.5
    }
  },
  hover: { 
    scale: 1.03,
    boxShadow: theme.shadows.lg,
    transition: {
      duration: 0.3
    }
  },
  tap: { 
    scale: 0.98,
    transition: {
      duration: 0.3
    }
  }
};

const AnimatedCard: React.FC<CardProps> = ({
  elevation = 'md',
  interactive = true,
  padding = '5',
  borderRadius = 'lg',
  background,
  children,
  ...rest
}) => {
  return (
    <CardContainer
      elevation={elevation}
      interactive={interactive}
      padding={padding as keyof typeof theme.space}
      borderRadius={borderRadius}
      background={background}
      variants={cardVariants}
      initial="hidden"
      whileInView="visible"
      whileHover={interactive ? 'hover' : undefined}
      whileTap={interactive ? 'tap' : undefined}
      viewport={{ once: true, amount: 0.1 }}
      {...rest}
    >
      {children}
    </CardContainer>
  );
};

export default AnimatedCard;