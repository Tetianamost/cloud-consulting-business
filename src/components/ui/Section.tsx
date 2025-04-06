import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { theme } from '../../styles/theme';

interface SectionProps {
  id?: string;
  background?: 'light' | 'dark' | 'primary' | 'accent' | 'white' | 'gradient';
  paddingTop?: keyof typeof theme.space;
  paddingBottom?: keyof typeof theme.space;
  children: React.ReactNode;
  fullWidth?: boolean;
}

const backgroundMap = {
  light: theme.colors.light,
  dark: theme.colors.dark,
  primary: theme.colors.primary,
  accent: theme.colors.accent,
  white: theme.colors.white,
  gradient: `linear-gradient(135deg, ${theme.colors.primary} 0%, ${theme.colors.accent} 100%)`,
};

const textColorMap = {
  light: theme.colors.gray800,
  dark: theme.colors.white,
  primary: theme.colors.white,
  accent: theme.colors.white,
  white: theme.colors.gray800,
  gradient: theme.colors.white,
};

const SectionWrapper = styled(motion.section)<SectionProps>`
  background: ${props => backgroundMap[props.background || 'white']};
  color: ${props => textColorMap[props.background || 'white']};
  padding-top: ${props => theme.space[props.paddingTop || 10]};
  padding-bottom: ${props => theme.space[props.paddingBottom || 10]};
  position: relative;
  overflow: hidden;
  width: 100%;
  max-width: 100vw;
  box-sizing: border-box;
  
  @media (min-width: ${theme.breakpoints.md}) {
    padding-top: ${props => theme.space[props.paddingTop ? (Number(props.paddingTop) + 2) as keyof typeof theme.space : 12]};
    padding-bottom: ${props => theme.space[props.paddingBottom ? (Number(props.paddingBottom) + 2) as keyof typeof theme.space : 12]};
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    padding-top: ${props => theme.space[props.paddingTop ? (Number(props.paddingTop) - 2) as keyof typeof theme.space : 8]};
    padding-bottom: ${props => theme.space[props.paddingBottom ? (Number(props.paddingBottom) - 2) as keyof typeof theme.space : 8]};
  }
`;

const SectionContainer = styled.div<{ fullWidth?: boolean }>`
  width: 100%;
  max-width: ${props => (props.fullWidth ? '100%' : theme.sizes.container.xl)};
  margin-left: auto;
  margin-right: auto;
  padding-left: ${props => (props.fullWidth ? '0' : theme.space[4])};
  padding-right: ${props => (props.fullWidth ? '0' : theme.space[4])};
  box-sizing: border-box;
  overflow-x: hidden;
  
  @media (min-width: ${theme.breakpoints.lg}) {
    padding-left: ${props => (props.fullWidth ? '0' : theme.space[6])};
    padding-right: ${props => (props.fullWidth ? '0' : theme.space[6])};
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    padding-left: ${props => (props.fullWidth ? '0' : theme.space[2])};
    padding-right: ${props => (props.fullWidth ? '0' : theme.space[2])};
  }
`;

// Animation variants
const sectionVariants = {
  hidden: { opacity: 0 },
  visible: { 
    opacity: 1,
    transition: {
      duration: 0.5,
      when: "beforeChildren",
      staggerChildren: 0.2
    }
  }
};

const Section: React.FC<SectionProps> = ({
  id,
  background = 'white',
  paddingTop,
  paddingBottom,
  children,
  fullWidth = false,
}) => {
  return (
    <SectionWrapper
      id={id}
      background={background}
      paddingTop={paddingTop}
      paddingBottom={paddingBottom}
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true, amount: 0.1 }}
      variants={sectionVariants}
    >
      <SectionContainer fullWidth={fullWidth}>{children}</SectionContainer>
    </SectionWrapper>
  );
};

export default Section;