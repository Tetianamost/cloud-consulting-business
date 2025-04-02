import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { Link } from 'react-scroll';
import { theme } from '../../styles/theme';

const FooterContainer = styled.footer`
  background-color: ${theme.colors.primary};
  color: ${theme.colors.white};
  padding: ${theme.space[10]} 0 ${theme.space[6]};
`;

const FooterContent = styled.div`
  max-width: ${theme.sizes.container.xl};
  margin: 0 auto;
  padding: 0 ${theme.space[4]};
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: ${theme.space[6]};
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-template-columns: repeat(6, 1fr);
  }
  
  @media (max-width: ${theme.breakpoints.md}) {
    grid-template-columns: 1fr;
  }
`;

const LogoColumn = styled.div`
  grid-column: span 4;
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-column: span 6;
  }
`;

const LinksColumn = styled.div`
  grid-column: span 2;
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-column: span 3;
  }
`;

const ContactColumn = styled.div`
  grid-column: span 4;
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-column: span 6;
  }
`;

const Logo = styled.div`
  font-family: ${theme.fonts.heading};
  font-weight: ${theme.fontWeights.bold};
  font-size: ${theme.fontSizes['3xl']};
  margin-bottom: ${theme.space[4]};
  
  span {
    color: ${theme.colors.secondary};
  }
`;

const Description = styled.p`
  font-size: ${theme.fontSizes.md};
  color: ${theme.colors.gray300};
  margin-bottom: ${theme.space[5]};
`;

const SocialLinks = styled.div`
  display: flex;
  gap: ${theme.space[4]};
  margin-top: ${theme.space[4]};
`;

const SocialLink = styled(motion.a)`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: ${theme.borderRadius.full};
  background-color: rgba(255, 255, 255, 0.1);
  color: ${theme.colors.white};
  font-size: ${theme.fontSizes.xl};
  transition: ${theme.transitions.normal};
  
  &:hover {
    background-color: ${theme.colors.secondary};
    transform: translateY(-3px);
  }
`;

const FooterHeading = styled.h4`
  font-size: ${theme.fontSizes.lg};
  font-weight: ${theme.fontWeights.bold};
  margin-bottom: ${theme.space[4]};
  color: ${theme.colors.white};
`;

const FooterLink = styled(Link)`
  display: block;
  color: ${theme.colors.gray300};
  font-size: ${theme.fontSizes.md};
  margin-bottom: ${theme.space[3]};
  cursor: pointer;
  transition: ${theme.transitions.fast};
  
  &:hover {
    color: ${theme.colors.secondary};
    transform: translateX(5px);
  }
`;

const FooterRegularLink = styled.a`
  display: block;
  color: ${theme.colors.gray300};
  font-size: ${theme.fontSizes.md};
  margin-bottom: ${theme.space[3]};
  cursor: pointer;
  transition: ${theme.transitions.fast};
  
  &:hover {
    color: ${theme.colors.secondary};
    transform: translateX(5px);
  }
`;

const ContactInfo = styled.div`
  margin-bottom: ${theme.space[4]};
`;

const ContactItem = styled.div`
  display: flex;
  align-items: flex-start;
  margin-bottom: ${theme.space[3]};
  color: ${theme.colors.gray300};
`;

const ContactIcon = styled.div`
  margin-right: ${theme.space[3]};
  color: ${theme.colors.secondary};
  font-size: ${theme.fontSizes.lg};
`;

const Divider = styled.hr`
  border: none;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  margin: ${theme.space[6]} 0 ${theme.space[4]};
`;

const Copyright = styled.div`
  max-width: ${theme.sizes.container.xl};
  margin: 0 auto;
  padding: 0 ${theme.space[4]};
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: ${theme.colors.gray400};
  font-size: ${theme.fontSizes.sm};
  
  @media (max-width: ${theme.breakpoints.md}) {
    flex-direction: column;
    text-align: center;
    gap: ${theme.space[3]};
  }
`;

const FooterLinks = styled.div`
  display: flex;
  gap: ${theme.space[4]};
  
  a {
    color: ${theme.colors.gray400};
    transition: ${theme.transitions.fast};
    
    &:hover {
      color: ${theme.colors.secondary};
    }
  }
`;

const currentYear = new Date().getFullYear();

const Footer: React.FC = () => {
  return (
    <FooterContainer>
      <FooterContent>
        <LogoColumn>
          <Logo>
            Cloud<span>Migrate</span>
          </Logo>
          <Description>
            Expert cloud migration services to help your business seamlessly transition to the cloud with confidence. AWS certified professionals with proven results.
          </Description>
          <SocialLinks>
            <SocialLink 
              href="https://linkedin.com" 
              target="_blank" 
              rel="noopener noreferrer"
              whileHover={{ scale: 1.1 }}
              whileTap={{ scale: 0.95 }}
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"></path>
                <rect x="2" y="9" width="4" height="12"></rect>
                <circle cx="4" cy="4" r="2"></circle>
              </svg>
            </SocialLink>
            <SocialLink 
              href="https://twitter.com" 
              target="_blank" 
              rel="noopener noreferrer"
              whileHover={{ scale: 1.1 }}
              whileTap={{ scale: 0.95 }}
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M23 3a10.9 10.9 0 0 1-3.14 1.53 4.48 4.48 0 0 0-7.86 3v1A10.66 10.66 0 0 1 3 4s-4 9 5 13a11.64 11.64 0 0 1-7 2c9 5 20 0 20-11.5a4.5 4.5 0 0 0-.08-.83A7.72 7.72 0 0 0 23 3z"></path>
              </svg>
            </SocialLink>
            <SocialLink 
              href="https://github.com" 
              target="_blank" 
              rel="noopener noreferrer"
              whileHover={{ scale: 1.1 }}
              whileTap={{ scale: 0.95 }}
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"></path>
              </svg>
            </SocialLink>
            <SocialLink 
              href="https://facebook.com" 
              target="_blank" 
              rel="noopener noreferrer"
              whileHover={{ scale: 1.1 }}
              whileTap={{ scale: 0.95 }}
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z"></path>
              </svg>
            </SocialLink>
          </SocialLinks>
        </LogoColumn>
        
        <LinksColumn>
          <FooterHeading>Quick Links</FooterHeading>
          <FooterLink to="home" smooth={true} offset={-80} duration={500}>Home</FooterLink>
          <FooterLink to="services" smooth={true} offset={-80} duration={500}>Services</FooterLink>
          <FooterLink to="certifications" smooth={true} offset={-80} duration={500}>Certifications</FooterLink>
          <FooterLink to="case-studies" smooth={true} offset={-80} duration={500}>Case Studies</FooterLink>
          <FooterLink to="pricing" smooth={true} offset={-80} duration={500}>Pricing</FooterLink>
          <FooterLink to="contact" smooth={true} offset={-80} duration={500}>Contact Us</FooterLink>
        </LinksColumn>
        
        <LinksColumn>
          <FooterHeading>Services</FooterHeading>
          <FooterRegularLink href="#services">Cloud Assessment</FooterRegularLink>
          <FooterRegularLink href="#services">Migration Strategy</FooterRegularLink>
          <FooterRegularLink href="#services">Cloud Implementation</FooterRegularLink>
          <FooterRegularLink href="#services">Cloud Optimization</FooterRegularLink>
          <FooterRegularLink href="#services">Security & Compliance</FooterRegularLink>
          <FooterRegularLink href="#services">Managed Services</FooterRegularLink>
        </LinksColumn>
        
        <ContactColumn>
          <FooterHeading>Contact Us</FooterHeading>
          <ContactInfo>
            <ContactItem>
              <ContactIcon>
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
                  <circle cx="12" cy="10" r="3"></circle>
                </svg>
              </ContactIcon>
              <div>
                123 Tech Street, Suite 456<br />
                San Francisco, CA 94107
              </div>
            </ContactItem>
            <ContactItem>
              <ContactIcon>
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path>
                </svg>
              </ContactIcon>
              <div>+1 (555) 123-4567</div>
            </ContactItem>
            <ContactItem>
              <ContactIcon>
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
                  <polyline points="22,6 12,13 2,6"></polyline>
                </svg>
              </ContactIcon>
              <div>info@cloudmigrate.com</div>
            </ContactItem>
          </ContactInfo>
        </ContactColumn>
      </FooterContent>
      
      <Divider />
      
      <Copyright>
        <div>&copy; {currentYear} CloudMigrate. All rights reserved.</div>
        <FooterLinks>
          <a href="/privacy-policy">Privacy Policy</a>
          <a href="/terms-of-service">Terms of Service</a>
        </FooterLinks>
      </Copyright>
    </FooterContainer>
  );
};

export default Footer;