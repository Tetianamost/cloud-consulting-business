import * as React from "react"
import styled, { css } from "styled-components"
import { theme } from "../../styles/theme"

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'default' | 'primary' | 'outline' | 'secondary' | 'ghost' | 'link' | 'destructive'
  size?: 'default' | 'sm' | 'lg' | 'icon'
  asChild?: boolean
  children?: React.ReactNode
}

const getVariantStyles = (variant: string) => {
  switch (variant) {
    case 'default':
      return css`
        background-color: ${theme.colors.secondary};
        color: ${theme.colors.white};
        border: none;
        
        &:hover:not(:disabled) {
          background-color: #e6890a;
        }
        
        &:active:not(:disabled) {
          background-color: #cc7700;
        }
        
        &:focus-visible {
          outline: 2px solid ${theme.colors.secondary};
          outline-offset: 2px;
        }
      `
    case 'primary':
      return css`
        background-color: ${theme.colors.primary};
        color: ${theme.colors.white};
        border: none;
        
        &:hover:not(:disabled) {
          background-color: #1a252f;
        }
        
        &:active:not(:disabled) {
          background-color: #0f1419;
        }
        
        &:focus-visible {
          outline: 2px solid ${theme.colors.primary};
          outline-offset: 2px;
        }
      `
    case 'outline':
      return css`
        background-color: transparent;
        color: ${theme.colors.secondary};
        border: 2px solid ${theme.colors.secondary};
        
        &:hover:not(:disabled) {
          background-color: ${theme.colors.secondary};
          color: ${theme.colors.white};
        }
        
        &:active:not(:disabled) {
          background-color: #e6890a;
          color: ${theme.colors.white};
        }
        
        &:focus-visible {
          outline: 2px solid ${theme.colors.secondary};
          outline-offset: 2px;
        }
      `
    case 'secondary':
      return css`
        background-color: ${theme.colors.accent};
        color: ${theme.colors.white};
        border: none;
        
        &:hover:not(:disabled) {
          background-color: #005a94;
        }
        
        &:active:not(:disabled) {
          background-color: #004173;
        }
        
        &:focus-visible {
          outline: 2px solid ${theme.colors.accent};
          outline-offset: 2px;
        }
      `
    case 'ghost':
      return css`
        background-color: transparent;
        color: ${theme.colors.primary};
        border: none;
        
        &:hover:not(:disabled) {
          background-color: rgba(255, 153, 0, 0.1);
          color: ${theme.colors.secondary};
        }
        
        &:focus-visible {
          outline: 2px solid ${theme.colors.secondary};
          outline-offset: 2px;
        }
      `
    case 'link':
      return css`
        background-color: transparent;
        color: ${theme.colors.accent};
        border: none;
        text-decoration: underline;
        text-underline-offset: 4px;
        
        &:hover:not(:disabled) {
          color: #005a94;
        }
        
        &:focus-visible {
          outline: 2px solid ${theme.colors.accent};
          outline-offset: 2px;
        }
      `
    case 'destructive':
      return css`
        background-color: ${theme.colors.danger};
        color: ${theme.colors.white};
        border: none;
        
        &:hover:not(:disabled) {
          background-color: #c82333;
        }
        
        &:active:not(:disabled) {
          background-color: #bd2130;
        }
        
        &:focus-visible {
          outline: 2px solid ${theme.colors.danger};
          outline-offset: 2px;
        }
      `
    default:
      return css`
        background-color: ${theme.colors.secondary};
        color: ${theme.colors.white};
        border: none;
      `
  }
}

const getSizeStyles = (size: string) => {
  switch (size) {
    case 'sm':
      return css`
        height: 32px;
        padding: 0 12px;
        font-size: ${theme.fontSizes.xs};
        border-radius: ${theme.borderRadius.md};
      `
    case 'lg':
      return css`
        height: 48px;
        padding: 0 32px;
        font-size: ${theme.fontSizes.lg};
        border-radius: ${theme.borderRadius.md};
      `
    case 'icon':
      return css`
        height: 40px;
        width: 40px;
        padding: 0;
      `
    default:
      return css`
        height: 40px;
        padding: 0 16px;
        font-size: ${theme.fontSizes.sm};
      `
  }
}

const StyledButton = styled.button<{ $variant: string; $size: string }>`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  white-space: nowrap;
  border-radius: ${theme.borderRadius.md};
  font-weight: ${theme.fontWeights.medium};
  transition: all 0.2s ease;
  cursor: pointer;
  font-family: ${theme.fonts.primary};
  
  ${props => getVariantStyles(props.$variant)}
  ${props => getSizeStyles(props.$size)}
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  svg {
    flex-shrink: 0;
  }
`

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ variant = 'default', size = 'default', asChild = false, children, ...props }, ref) => {
    if (asChild && React.isValidElement(children)) {
      // For asChild, we'll use a simpler approach that works with styled-components
      const childProps = children.props as any;
      return React.cloneElement(children, {
        ...props,
        ...childProps,
        style: {
          display: 'inline-flex',
          alignItems: 'center',
          justifyContent: 'center',
          gap: '8px',
          whiteSpace: 'nowrap',
          borderRadius: theme.borderRadius.md,
          fontWeight: theme.fontWeights.medium,
          transition: 'all 0.2s ease',
          cursor: 'pointer',
          fontFamily: theme.fonts.primary,
          textDecoration: 'none',
          border: 'none',
          ...(childProps.style || {}),
        },
        className: `button-${variant} button-${size} ${childProps.className || ''}`.trim()
      })
    }

    return (
      <StyledButton
        ref={ref}
        $variant={variant}
        $size={size}
        {...props}
      >
        {children}
      </StyledButton>
    )
  }
)

Button.displayName = "Button"

export { Button }