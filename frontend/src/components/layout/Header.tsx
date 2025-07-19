import React, { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import { Link } from 'react-scroll';
import { FiMenu, FiX } from 'react-icons/fi';
import { theme } from '../../styles/theme';
import Button from '../ui/Button';
import Icon from '../ui/Icon';

const HeaderContainer = styled(motion.header)<{ $isScrolled: boolean; $isHidden: boolean }>`
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: ${theme.zIndices.sticky};
  background: ${props => props.$isScrolled ? 'rgba(255, 255, 255, 0.95)' : 'transparent'};
  backdrop-filter: ${props => props.$isScrolled ? 'blur(10px)' : 'none'};
  border-bottom: ${props => props.$isScrolled ? `1px solid ${theme.colors.gray200}` : 'none'};
  padding: ${props => props.$isScrolled ? theme.space[3] : theme.space[4]} 0;
  transition: ${theme.transitions.normal}, transform 0.3s cubic-bezier(0.4,0,0.2,1);
  box-shadow: ${props => props.$isScrolled ? '0 4px 24px rgba(0,0,0,0.07)' : 'none'};
  transform: translateY(${props => props.$isHidden ? '-110%' : '0'});
  will-change: transform;
`;

const Nav = styled.nav`
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: ${theme.sizes.container.xl};
  margin: 0 auto;
  padding: 0 ${theme.space[4]};
  
  @media (min-width: ${theme.breakpoints.lg}) {
    padding: 0 ${theme.space[6]};
  }
`;

const Logo = styled.div`
  font-family: ${theme.fonts.heading};
  font-weight: ${theme.fontWeights.bold};
  font-size: ${theme.fontSizes['2xl']};
  color: ${props => props.theme.colors.primary};
  display: flex;
  align-items: center;
  
  span {
    color: ${props => props.theme.colors.secondary};
  }
`;

const NavLinks = styled.div<{ $isOpen: boolean }>`
  display: flex;
  align-items: center;
  
  @media (max-width: ${theme.breakpoints.lg}) {
    flex-direction: column;
    background: ${theme.colors.white};
    position: fixed;
    top: 0;
    right: 0;
    width: 80%;
    max-width: 400px;
    height: 100vh;
    padding: ${theme.space[8]} ${theme.space[4]};
    z-index: ${theme.zIndices.modal};
    box-shadow: -10px 0 30px rgba(0, 0, 0, 0.1);
    transform: ${props => props.$isOpen ? 'translateX(0)' : 'translateX(100%)'};
    transition: transform 0.3s ease-in-out;
  }
`;

const NavLink = styled(Link)`
  margin: 0 ${theme.space[4]};
  font-weight: ${theme.fontWeights.medium};
  cursor: pointer;
  position: relative;
  
  &:after {
    content: '';
    position: absolute;
    width: 0;
    height: 2px;
    bottom: -4px;
    left: 0;
    background-color: ${theme.colors.secondary};
    transition: width 0.3s ease;
  }
  
  &:hover:after,
  &.active:after {
    width: 100%;
  }
  
  @media (max-width: ${theme.breakpoints.lg}) {
    margin: ${theme.space[4]} 0;
    font-size: ${theme.fontSizes.xl};
  }
`;

const NavButton = styled.div`
  margin-left: ${theme.space[4]};
  
  @media (max-width: ${theme.breakpoints.lg}) {
    margin: ${theme.space[6]} 0 0;
    width: 100%;
  }
`;

const MenuButton = styled.button`
  display: none;
  background: transparent;
  border: none;
  font-size: ${theme.fontSizes['2xl']};
  cursor: pointer;
  z-index: ${theme.zIndices.overlay};
  
  @media (max-width: ${theme.breakpoints.lg}) {
    display: block;
  }
`;

const Overlay = styled(motion.div)`
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  z-index: ${theme.zIndices.overlay};
`;

const Header: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [isScrolled, setIsScrolled] = useState(false);
  const [isHidden, setIsHidden] = useState(false);
  const lastScrollY = useRef(0);

  // Improved scroll handler: hide on scroll down, show on scroll up
  useEffect(() => {
    let lastY = window.scrollY;
    let ticking = false;
    const handleScroll = () => {
      if (!ticking) {
        window.requestAnimationFrame(() => {
          const currentY = window.scrollY;
          setIsScrolled(currentY > 50);
          if (currentY > lastY && currentY > 100) {
            setIsHidden(true); // Hide header
          } else {
            setIsHidden(false); // Show header
          }
          lastY = currentY;
          ticking = false;
        });
        ticking = true;
      }
    };
    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  // Prevent background scroll when mobile menu is open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);
  
  const toggleMenu = () => {
    setIsOpen(!isOpen);
  };
  
  const closeMenu = () => {
    setIsOpen(false);
  };
  
  // Animation variants
  const navVariants = {
    hidden: { opacity: 0, y: -20 },
    visible: { 
      opacity: 1, 
      y: 0,
      transition: {
        duration: 0.5,
        delay: 0.2
      }
    }
  };
  
  return (
    <HeaderContainer
      $isScrolled={isScrolled}
      $isHidden={isHidden}
      initial="hidden"
      animate="visible"
      variants={navVariants}
    >
      <Nav>
        <Logo>
        <NavLink
            to="home"
            spy={true}
            smooth={true}
            offset={-80}
            duration={500}
            onClick={closeMenu}
          >
          Cloud<span>Partner</span><span style={{ color: theme.colors.accent }}>Pro</span>
          </NavLink>
        </Logo>
        <MenuButton onClick={toggleMenu}>
        {isOpen ? <Icon icon={FiX} size={24} /> : <Icon icon={FiMenu} size={24} />}
      </MenuButton>
        <AnimatePresence>
          {isOpen && (
            <Overlay
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={closeMenu}
            />
          )}
        </AnimatePresence>
        <NavLinks $isOpen={isOpen}>
          <NavLink
            to="home"
            spy={true}
            smooth={true}
            offset={-80}
            duration={500}
            onClick={closeMenu}
          >
            Home
          </NavLink>
          <NavLink
            to="services"
            spy={true}
            smooth={true}
            offset={-80}
            duration={500}
            onClick={closeMenu}
          >
            Services
          </NavLink>
          <NavLink
            to="certifications"
            spy={true}
            smooth={true}
            offset={-80}
            duration={500}
            onClick={closeMenu}
          >
            Certifications
          </NavLink>
          <NavLink
            to="project-insights"
            spy={true}
            smooth={true}
            offset={-80}
            duration={500}
            onClick={closeMenu}
          >
            Project Insights
          </NavLink>
          <NavLink
            to="pricing"
            spy={true}
            smooth={true}
            offset={-80}
            duration={500}
            onClick={closeMenu}
          >
            Pricing
          </NavLink>
          <NavButton>
            <Button
              onClick={() => {
                closeMenu();
                // Scroll to contact section
                const contactSection = document.getElementById('contact');
                if (contactSection) {
                  contactSection.scrollIntoView({ behavior: 'smooth' });
                }
              }}
            >
              Contact Us
            </Button>
          </NavButton>
        </NavLinks>
      </Nav>
    </HeaderContainer>
  );
};

export default Header;