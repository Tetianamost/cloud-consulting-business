import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { theme } from '../../../styles/theme';
import Section from '../../ui/Section';
import ContactForm from './ContactForm';
import { FiMapPin, FiPhone, FiMail, FiClock } from 'react-icons/fi';
import Icon from '../../ui/Icon';

const SectionTitle = styled(motion.h2)`
  text-align: center;
  margin-bottom: ${theme.space[3]};
  color: ${theme.colors.white};
`;

const SectionSubtitle = styled(motion.p)`
  text-align: center;
  max-width: 600px;
  margin: 0 auto ${theme.space[8]};
  color: ${theme.colors.gray300};
  font-size: ${theme.fontSizes.lg};
`;

const ContactContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: ${theme.space[6]};
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-template-columns: 1fr;
  }
`;

const ContactInfo = styled(motion.div)`
  background-color: rgba(255, 255, 255, 0.05);
  border-radius: ${theme.borderRadius.lg};
  padding: ${theme.space[6]};
`;

const InfoBlock = styled.div`
  margin-bottom: ${theme.space[5]};
`;

const InfoTitle = styled.h3`
  font-size: ${theme.fontSizes.xl};
  color: ${theme.colors.white};
  margin-bottom: ${theme.space[4]};
`;

const ContactItem = styled.div`
  display: flex;
  align-items: flex-start;
  margin-bottom: ${theme.space[4]};
  color: ${theme.colors.gray300};
`;

const ContactIcon = styled.div`
  margin-right: ${theme.space[3]};
  color: ${theme.colors.secondary};
  font-size: ${theme.fontSizes.xl};
  min-width: 24px;
`;

const ContactText = styled.div`
  line-height: 1.5;
`;

const ContactLink = styled.a`
  color: ${theme.colors.gray300};
  transition: ${theme.transitions.fast};
  
  &:hover {
    color: ${theme.colors.secondary};
  }
`;

const SocialLinks = styled.div`
  display: flex;
  gap: ${theme.space[3]};
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
  font-size: ${theme.fontSizes.lg};
  transition: ${theme.transitions.normal};
  
  &:hover {
    background-color: ${theme.colors.secondary};
    transform: translateY(-3px);
  }
`;

const GoogleMap = styled.div`
  margin-top: ${theme.space[5]};
  border-radius: ${theme.borderRadius.md};
  overflow: hidden;
  height: 200px;
  
  iframe {
    width: 100%;
    height: 100%;
    border: none;
  }
`;

const FormContainer = styled(motion.div)``;

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

const infoVariants = {
  hidden: { opacity: 0, x: -30 },
  visible: {
    opacity: 1,
    x: 0,
    transition: {
      duration: 0.5,
      delay: 0.2
    }
  }
};

const formVariants = {
  hidden: { opacity: 0, x: 30 },
  visible: {
    opacity: 1,
    x: 0,
    transition: {
      duration: 0.5,
      delay: 0.2
    }
  }
};

const Contact: React.FC = () => {
  return (
    <Section id="contact" background="primary">
      <SectionTitle
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={titleVariants}
      >
        Get In Touch
      </SectionTitle>
      <SectionSubtitle
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={titleVariants}
      >
        Ready to start your cloud journey? Contact us for a free consultation
      </SectionSubtitle>
      
      <ContactContainer>
        <ContactInfo
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true }}
          variants={infoVariants}
        >
          <InfoBlock>
            <InfoTitle>Contact Information</InfoTitle>
            <ContactItem>
              <ContactIcon>
              <Icon icon={FiMapPin} size={20} />
              </ContactIcon>
              <ContactText>
                123 Cloud Street, Suite 456<br />
                San Francisco, CA 94107<br />
                United States
              </ContactText>
            </ContactItem>
            <ContactItem>
              <ContactIcon>
              <Icon icon={FiPhone} size={20} />
              </ContactIcon>
              <ContactText>
                <ContactLink href="tel:+15551234567">+1 (555) 123-4567</ContactLink>
              </ContactText>
            </ContactItem>
            <ContactItem>
              <ContactIcon>
              <Icon icon={FiMail} size={20} />
              </ContactIcon>
              <ContactText>
                <ContactLink href="mailto:info@cloudmigrate.com">info@cloudmigrate.com</ContactLink>
              </ContactText>
            </ContactItem>
            <ContactItem>
              <ContactIcon>
              <Icon icon={FiClock} size={20} />
              </ContactIcon>
              <ContactText>
                Monday - Friday: 9:00 AM - 6:00 PM PST<br />
                Saturday - Sunday: Closed
              </ContactText>
            </ContactItem>
          </InfoBlock>
          
          <InfoBlock>
            <InfoTitle>Connect With Us</InfoTitle>
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
              <SocialLink 
                href="https://instagram.com"
                target="_blank"
                rel="noopener noreferrer"
                whileHover={{ scale: 1.1 }}
                whileTap={{ scale: 0.95 }}
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <rect x="2" y="2" width="20" height="20" rx="5" ry="5"></rect>
                  <path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z"></path>
                  <line x1="17.5" y1="6.5" x2="17.51" y2="6.5"></line>
                </svg>
              </SocialLink>
            </SocialLinks>
          </InfoBlock>
          
          <GoogleMap>
            {/* This would typically be a real Google Maps embed */}
            <iframe
              title="Office Location"
              src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3153.1034021349707!2d-122.40092368468204!3d37.78425507975649!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x8085807ded297e89%3A0xd9e8mineralsed93ca0!2s123%20Main%20St%2C%20San%20Francisco%2C%20CA%2094105!5e0!3m2!1sen!2sus!4v1650000000000!5m2!1sen!2sus"
              allowFullScreen
              loading="lazy"
              referrerPolicy="no-referrer-when-downgrade"
            ></iframe>
          </GoogleMap>
        </ContactInfo>
        
        <FormContainer
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true }}
          variants={formVariants}
        >
          <ContactForm />
        </FormContainer>
      </ContactContainer>
    </Section>
  );
};

export default Contact;