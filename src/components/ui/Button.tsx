import React from 'react';
import styled, { css } from 'styled-components';
import { motion } from 'framer-motion';
import { theme } from '../../styles/theme';

type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost';
type ButtonSize = 'sm' | 'md' | 'lg';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  fullWidth?: boolean;
  isLoading?: boolean;
  as?: any;
  to?: string;
  icon?: React.ReactElement;
  iconPosition?: 'left' | 'right';
}

// Size variations
const sizeStyles = {
  sm: css`
    font-size: ${theme.fontSizes.sm};
    padding: ${theme.space[2]} ${theme.space[3]};
    border-radius: ${theme.borderRadius.md};
  `,
  md: css`
    font-size: ${theme.fontSizes.md};
    padding: ${theme.space[3]} ${theme.space[5]};
    border-radius: ${theme.borderRadius.md};
  `,
  lg: css`
    font-size: ${theme.fontSizes.lg};
    padding: ${theme.space[4]} ${theme.space[6]};
    border-radius: ${theme.borderRadius.lg};
  `,
};

// Variant styles
const variantStyles = {
  primary: css`
    background-color: ${theme.colors.accent};
    color: ${theme.colors.white};
    border: 2px solid ${theme.colors.accent};
    
    &:hover:not(:disabled) {
      background-color: ${theme.colors.highlight};
      border-color: ${theme.colors.highlight};
    }
  `,
  secondary: css`
    background-color: ${theme.colors.secondary};
    color: ${theme.colors.primary};
    border: 2px solid ${theme.colors.secondary};
    
    &:hover:not(:disabled) {
      background-color: ${theme.colors.warning};
      border-color: ${theme.colors.warning};
    }
  `,
  outline: css`
    background-color: transparent;
    color: ${theme.colors.accent};
    border: 2px solid ${theme.colors.accent};
    
    &:hover:not(:disabled) {
      background-color: ${theme.colors.accent};
      color: ${theme.colors.white};
    }
  `,
  ghost: css`
    background-color: transparent;
    color: ${theme.colors.accent};
    border: 2px solid transparent;
    
    &:hover:not(:disabled) {
      background-color: ${theme.colors.gray100};
    }
  `,
};

const StyledButton = styled(motion.button)<ButtonProps>`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: ${theme.fontWeights.semibold};
  cursor: pointer;
  transition: ${theme.transitions.normal};
  white-space: nowrap;
  position: relative;
  width: ${props => (props.fullWidth ? '100%' : 'auto')};

  // Apply variant styling
  ${props => variantStyles[props.variant || 'primary']}
  
  // Apply size styling
  ${props => sizeStyles[props.size || 'md']}
  
  // Loading state
  &:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }
  
  // Icon spacing
  & > svg {
    ${props => props.iconPosition === 'left' && `margin-right: ${theme.space[2]};`}
    ${props => props.iconPosition === 'right' && `margin-left: ${theme.space[2]};`}
  }
`;

const LoadingSpinner = styled.div`
  border: 2px solid rgba(255, 255, 255, 0.2);
  border-top-color: white;
  border-radius: 50%;
  width: 1em;
  height: 1em;
  animation: spin 0.8s linear infinite;
  margin-right: ${theme.space[2]};

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
`;

const buttonVariants = {
  rest: { scale: 1 },
  hover: { scale: 1.05 },
  tap: { scale: 0.95 },
};

const Button: React.FC<ButtonProps> = ({
  children,
  variant = 'primary',
  size = 'md',
  fullWidth = false,
  icon,
  iconPosition = 'left',
  isLoading = false,
  disabled,
  type = 'button',
  ...rest
}) => {
  return (
    <StyledButton
      type={type}
      variant={variant}
      size={size}
      fullWidth={fullWidth}
      iconPosition={iconPosition}
      disabled={disabled || isLoading}
      whileHover="hover"
      whileTap="tap"
      initial="rest"
      variants={buttonVariants}
      {...rest}
    >
      {isLoading && <LoadingSpinner />}
      {!isLoading && icon && iconPosition === 'left' && icon}
      {children}
      {!isLoading && icon && iconPosition === 'right' && icon}
    </StyledButton>
  );
};

export default Button;