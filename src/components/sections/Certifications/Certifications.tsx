import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { theme } from '../../../styles/theme';
import Section from '../../ui/Section';
import CertificationBadge from './CertificationBadge';

// AWS certification images would normally be imported from your assets folder
// For this example, we'll use placeholders
const certifications = [
  {
    id: 1,
    title: 'AWS Certified Solutions Architect – Professional',
    description: 'Advanced expertise in designing distributed applications and systems on AWS platform.',
    image: 'https://d1.awsstatic.com/training-and-certification/certification-badges/AWS-Certified-Solutions-Architect-Professional_badge.63c39369def4468dab4822b6fad24158909bffc0.png',
    color: '#FF9900'
  },
  {
    id: 2,
    title: 'AWS Certified DevOps Engineer – Professional',
    description: 'Expertise in continuous delivery and automation of AWS infrastructure for security and compliance.',
    image: 'https://d1.awsstatic.com/training-and-certification/certification-badges/AWS-Certified-DevOps-Engineer-Professional_badge.e1e4bcd2ee73d669a6b490c5622cf62e5494253e.png',
    color: '#FF9900'
  },
  {
    id: 3,
    title: 'AWS Certified Data Analytics – Specialty',
    description: 'Expertise in designing and implementing AWS services to derive insights from data.',
    image: 'https://d1.awsstatic.com/training-and-certification/certification-badges/AWS-Certified-Data-Analytics-Specialty_badge.73b622eb78b8f1118689e28d665a6ddff553651a.png',
    color: '#FF9900'
  },
  {
    id: 4,
    title: 'AWS Certified Security – Specialty',
    description: 'Expertise in security best practices for AWS platform and securing complex AWS environments.',
    image: 'https://d1.awsstatic.com/training-and-certification/certification-badges/AWS-Certified-Security-Specialty_badge.b8d70d28dc20dd8e7643b5b1b0ee6d57f54c95f3.png',
    color: '#FF9900'
  },
  {
    id: 5,
    title: 'AWS Certified Advanced Networking – Specialty',
    description: 'Expertise in designing and implementing AWS and hybrid IT network architectures.',
    image: 'https://d1.awsstatic.com/training-and-certification/certification-badges/AWS-Certified-Advanced-Networking-Specialty_badge.c3579756201966ad228869493772e3422d86bcf0.png',
    color: '#FF9900'
  },
  {
    id: 6,
    title: 'AWS Certified Database – Specialty',
    description: 'Expertise in designing, implementing, and managing AWS database solutions.',
    image: 'https://d1.awsstatic.com/training-and-certification/certification-badges/AWS-Certified-Database-Specialty_badge.d4edf2ddd4c63dbee96258a4a48db0a6b8a25580.png',
    color: '#FF9900'
  }
];

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

const CertificationsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: ${theme.space[6]};
  
  @media (max-width: ${theme.breakpoints.lg}) {
    grid-template-columns: repeat(2, 1fr);
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    grid-template-columns: 1fr;
  }
`;

const PartnerWrapper = styled(motion.div)`
  background-color: rgba(255, 255, 255, 0.05);
  border-radius: ${theme.borderRadius.lg};
  padding: ${theme.space[6]};
  margin-top: ${theme.space[10]};
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
`;

const PartnerTitle = styled.h3`
  font-size: ${theme.fontSizes['2xl']};
  color: ${theme.colors.white};
  margin-bottom: ${theme.space[4]};
`;

const PartnerBadge = styled(motion.div)`
  background-color: ${theme.colors.primary};
  border: 2px solid ${theme.colors.secondary};
  border-radius: ${theme.borderRadius.lg};
  padding: ${theme.space[4]};
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: ${theme.space[4]};
  
  img {
    max-width: 160px;
    height: auto;
  }
`;

const PartnerDescription = styled.p`
  color: ${theme.colors.gray300};
  font-size: ${theme.fontSizes.md};
  max-width: 500px;
  margin: 0 auto;
`;

// Animation variants
const titleVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: { 
    opacity: 1, 
    y: 0,
    transition: {
      duration:.5
    }
  }
};

const partnerVariants = {
  hidden: { opacity: 0, y: 30 },
  visible: { 
    opacity: 1, 
    y: 0,
    transition: {
      duration: 0.6,
      delay: 0.4
    }
  }
};

const Certifications: React.FC = () => {
  return (
    <Section id="certifications" background="primary">
      <SectionTitle
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={titleVariants}
      >
        AWS Certified Professionals
      </SectionTitle>
      <SectionSubtitle
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={titleVariants}
      >
        Our team holds industry-leading certifications, ensuring you receive expert-level cloud migration services
      </SectionSubtitle>
      
      <CertificationsGrid>
        {certifications.map((certification, index) => (
          <CertificationBadge 
            key={certification.id}
            certification={certification}
            index={index}
          />
        ))}
      </CertificationsGrid>
      
      <PartnerWrapper
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={partnerVariants}
      >
        <PartnerTitle>AWS Select Consulting Partner</PartnerTitle>
        <PartnerBadge
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.98 }}
        >
          <motion.img 
            src="https://d1.awsstatic.com/partner-network/logo/AWS_Consulting_Partner.d737d1305b4d89e13a8777bf5f304d7765913639.png" 
            alt="AWS Consulting Partner" 
          />
        </PartnerBadge>
        <PartnerDescription>
          As an AWS Select Consulting Partner, we've demonstrated technical proficiency and proven success in helping customers migrate to and optimize their AWS environments. This partnership allows us to deliver enhanced service and support to our clients.
        </PartnerDescription>
      </PartnerWrapper>
    </Section>
  );
};

export default Certifications;