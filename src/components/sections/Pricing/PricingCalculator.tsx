import React, { useState, useRef } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import { theme } from '../../../styles/theme';
import Button from '../../ui/Button';
import { FiInfo, FiCheckCircle, FiX, FiMail, FiUser, FiPhone, FiBriefcase } from 'react-icons/fi';
import Icon from '../../ui/Icon';
import emailjs from '@emailjs/browser';

const CalculatorContainer = styled(motion.div)`
  background-color: white;
  border-radius: ${theme.borderRadius.lg};
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  padding: ${theme.space[6]};
`;

const CalculatorTitle = styled.h3`
  font-size: ${theme.fontSizes.xl};
  color: ${theme.colors.primary};
  margin-bottom: ${theme.space[5]};
  text-align: center;
`;

const FormGroup = styled.div`
  margin-bottom: ${theme.space[5]};
`;

const Label = styled.label`
  display: block;
  font-weight: ${theme.fontWeights.medium};
  margin-bottom: ${theme.space[2]};
  display: flex;
  align-items: center;
`;

const InfoIcon = styled.span`
  margin-left: ${theme.space[2]};
  color: ${theme.colors.gray500};
  cursor: help;
  position: relative;
  
  &:hover::after {
    content: attr(data-tooltip);
    position: absolute;
    top: -40px;
    left: 50%;
    transform: translateX(-50%);
    background-color: ${theme.colors.primary};
    color: white;
    padding: ${theme.space[2]} ${theme.space[3]};
    border-radius: ${theme.borderRadius.md};
    font-size: ${theme.fontSizes.sm};
    white-space: nowrap;
    z-index: 10;
  }
  
  &:hover::before {
    content: '';
    position: absolute;
    top: -10px;
    left: 50%;
    transform: translateX(-50%);
    border-width: 5px;
    border-style: solid;
    border-color: ${theme.colors.primary} transparent transparent transparent;
    z-index: 10;
  }
`;

const OptionGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: ${theme.space[3]};
  
  @media (max-width: ${theme.breakpoints.sm}) {
    grid-template-columns: 1fr;
  }
`;

const OptionCard = styled.label<{ selected: boolean }>`
  display: block;
  background-color: ${props => props.selected ? theme.colors.primary + '10' : theme.colors.gray100};
  border: 2px solid ${props => props.selected ? theme.colors.primary : 'transparent'};
  border-radius: ${theme.borderRadius.md};
  padding: ${theme.space[3]};
  cursor: pointer;
  transition: ${theme.transitions.fast};
  
  &:hover {
    background-color: ${props => props.selected ? theme.colors.primary + '10' : theme.colors.gray200};
  }
`;

const RadioInput = styled.input`
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
`;

const OptionTitle = styled.div<{ selected: boolean }>`
  font-weight: ${theme.fontWeights.medium};
  color: ${props => props.selected ? theme.colors.primary : theme.colors.gray800};
  margin-bottom: ${theme.space[1]};
`;

const OptionDescription = styled.div`
  font-size: ${theme.fontSizes.sm};
  color: ${theme.colors.gray600};
`;

const RangeContainer = styled.div`
  margin-top: ${theme.space[3]};
`;

const RangeInput = styled.input`
  width: 100%;
  -webkit-appearance: none;
  height: 10px;
  border-radius: ${theme.borderRadius.full};
  background: ${theme.colors.gray200};
  outline: none;
  margin: ${theme.space[3]} 0;
  
  &::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 24px;
    height: 24px;
    border-radius: 50%;
    background: ${theme.colors.primary};
    cursor: pointer;
    transition: ${theme.transitions.fast};
    
    &:hover {
      transform: scale(1.1);
    }
  }
  
  &::-moz-range-thumb {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    background: ${theme.colors.primary};
    cursor: pointer;
    transition: ${theme.transitions.fast};
    
    &:hover {
      transform: scale(1.1);
    }
  }
`;

const RangeLabels = styled.div`
  display: flex;
  justify-content: space-between;
  font-size: ${theme.fontSizes.sm};
  color: ${theme.colors.gray600};
`;

const RangeValue = styled.div`
  text-align: center;
  font-weight: ${theme.fontWeights.medium};
  color: ${theme.colors.primary};
  font-size: ${theme.fontSizes.lg};
  margin-bottom: ${theme.space[2]};
`;

const Divider = styled.hr`
  border: none;
  border-top: 1px solid ${theme.colors.gray200};
  margin: ${theme.space[6]} 0;
`;

const EstimateContainer = styled(motion.div)`
  background-color: ${theme.colors.gray100};
  border-radius: ${theme.borderRadius.lg};
  padding: ${theme.space[5]};
  margin-top: ${theme.space[5]};
`;

const EstimateTitle = styled.h4`
  font-size: ${theme.fontSizes.lg};
  margin-bottom: ${theme.space[3]};
  color: ${theme.colors.primary};
`;

const EstimateRow = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: ${theme.space[2]};
  font-size: ${theme.fontSizes.md};
`;

const EstimateLabel = styled.div`
  color: ${theme.colors.gray700};
`;

const EstimateValue = styled.div`
  font-weight: ${theme.fontWeights.medium};
  color: ${theme.colors.primary};
`;

const TotalRow = styled.div`
  display: flex;
  justify-content: space-between;
  border-top: 1px dashed ${theme.colors.gray300};
  padding-top: ${theme.space[3]};
  margin-top: ${theme.space[3]};
  font-size: ${theme.fontSizes.lg};
  font-weight: ${theme.fontWeights.bold};
`;

const TotalLabel = styled.div`
  color: ${theme.colors.gray800};
`;

const TotalValue = styled.div`
  color: ${theme.colors.accent};
`;

const ButtonContainer = styled.div`
  margin-top: ${theme.space[5]};
  display: flex;
  justify-content: center;
`;

const SuccessMessage = styled(motion.div)`
  background-color: ${theme.colors.success + '20'};
  color: ${theme.colors.success};
  padding: ${theme.space[4]};
  border-radius: ${theme.borderRadius.md};
  display: flex;
  align-items: center;
  margin-top: ${theme.space[5]};
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.05);
  border-left: 4px solid ${theme.colors.success};
  
  svg {
    margin-right: ${theme.space[3]};
    font-size: ${theme.fontSizes.xl};
  }
`;

// Modal components
const ModalOverlay = styled(motion.div)`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: ${theme.space[4]};
`;

const ModalContainer = styled(motion.div)`
  background-color: white;
  border-radius: ${theme.borderRadius.lg};
  max-width: 600px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
`;

const ModalHeader = styled.div`
  padding: ${theme.space[5]};
  border-bottom: 1px solid ${theme.colors.gray200};
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const ModalTitle = styled.h3`
  font-size: ${theme.fontSizes.xl};
  color: ${theme.colors.primary};
  margin: 0;
`;

const CloseButton = styled.button`
  background: transparent;
  border: none;
  color: ${theme.colors.gray500};
  font-size: ${theme.fontSizes.xl};
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: ${theme.space[1]};
  
  &:hover {
    color: ${theme.colors.primary};
  }
`;

const ModalBody = styled.div`
  padding: ${theme.space[5]};
`;

const ModalFooter = styled.div`
  padding: ${theme.space[5]};
  border-top: 1px solid ${theme.colors.gray200};
  display: flex;
  justify-content: flex-end;
  gap: ${theme.space[3]};
`;

const QuoteDetail = styled.div`
  background-color: ${theme.colors.gray100};
  border-radius: ${theme.borderRadius.md};
  padding: ${theme.space[4]};
  margin-bottom: ${theme.space[4]};
`;

const QuoteDetailRow = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: ${theme.space[2]};
  font-size: ${theme.fontSizes.md};
  
  &:last-child {
    margin-bottom: 0;
  }
`;

const QuoteDetailLabel = styled.div`
  color: ${theme.colors.gray700};
  font-weight: ${theme.fontWeights.medium};
`;

const QuoteDetailValue = styled.div`
  color: ${theme.colors.primary};
`;

const QuoteDetailTotal = styled.div`
  display: flex;
  justify-content: space-between;
  border-top: 1px solid ${theme.colors.gray300};
  margin-top: ${theme.space[3]};
  padding-top: ${theme.space[3]};
  font-weight: ${theme.fontWeights.bold};
  font-size: ${theme.fontSizes.lg};
`;

const Input = styled.input`
  width: 100%;
  padding: ${theme.space[3]};
  border: 1px solid ${theme.colors.gray300};
  border-radius: ${theme.borderRadius.md};
  font-size: ${theme.fontSizes.md};
  margin-bottom: ${theme.space[4]};
  transition: ${theme.transitions.fast};
  
  &:focus {
    outline: none;
    border-color: ${theme.colors.primary};
    box-shadow: 0 0 0 3px ${theme.colors.primary}20;
  }
`;

const Textarea = styled.textarea`
  width: 100%;
  padding: ${theme.space[3]};
  border: 1px solid ${theme.colors.gray300};
  border-radius: ${theme.borderRadius.md};
  font-size: ${theme.fontSizes.md};
  margin-bottom: ${theme.space[4]};
  min-height: 100px;
  resize: vertical;
  transition: ${theme.transitions.fast};
  
  &:focus {
    outline: none;
    border-color: ${theme.colors.primary};
    box-shadow: 0 0 0 3px ${theme.colors.primary}20;
  }
`;

const InputGroup = styled.div`
  position: relative;
  margin-bottom: ${theme.space[4]};
`;

const InputIcon = styled.div`
  position: absolute;
  left: ${theme.space[3]};
  top: 50%;
  transform: translateY(-50%);
  color: ${theme.colors.gray500};
`;

const IconInput = styled.input`
  width: 100%;
  padding: ${theme.space[3]};
  padding-left: ${theme.space[8]};
  border: 1px solid ${theme.colors.gray300};
  border-radius: ${theme.borderRadius.md};
  font-size: ${theme.fontSizes.md};
  transition: ${theme.transitions.fast};
  
  &:focus {
    outline: none;
    border-color: ${theme.colors.primary};
    box-shadow: 0 0 0 3px ${theme.colors.primary}20;
  }
`;

const InputLabel = styled.label`
  display: block;
  margin-bottom: ${theme.space[2]};
  font-weight: ${theme.fontWeights.medium};
  color: ${theme.colors.gray700};
`;

const ErrorMessage = styled.div`
  color: ${theme.colors.danger};
  font-size: ${theme.fontSizes.sm};
  margin-top: ${theme.space[1]};
`;

// Define the pricing calculator parameters
const serviceTypes = [
  {
    id: 'assessment',
    title: 'Initial Assessment',
    description: 'Comprehensive cloud readiness evaluation and recommendations',
    basePrice: 750,
    pricePerServer: 50,
  },
  {
    id: 'migration',
    title: 'Migration Planning',
    description: 'Detailed planning for moving specific applications to the cloud',
    basePrice: 500,
    pricePerServer: 150,
  },
  {
    id: 'architecture',
    title: 'Cloud Architecture Review',
    description: 'Expert review of existing or planned cloud architecture',
    basePrice: 600,
    pricePerServer: 75,
  },
  {
    id: 'optimization',
    title: 'Implementation Assistance',
    description: 'Hands-on help implementing cloud solutions (hourly rate)',
    basePrice: 0,
    pricePerServer: 125, // Per hour rate
  },
];

const complexityLevels = [
  {
    id: 'simple',
    title: 'Simple',
    description: 'Basic setup with minimal dependencies',
    multiplier: 1,
  },
  {
    id: 'moderate',
    title: 'Moderate',
    description: 'Multiple services with some integrations',
    multiplier: 1.5,
  },
  {
    id: 'complex',
    title: 'Complex',
    description: 'Multiple environments, many dependencies',
    multiplier: 2,
  },
];

// Animation variants
const calculatorVariants = {
  hidden: { opacity: 0 },
  visible: { 
    opacity: 1,
    transition: {
      duration: 0.5
    }
  }
};

const estimateVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: { 
    opacity: 1, 
    y: 0,
    transition: {
      duration: 0.5
    }
  }
};

const successVariants = {
  hidden: { opacity: 0, height: 0 },
  visible: { 
    opacity: 1, 
    height: 'auto',
    transition: {
      duration: 0.3
    }
  },
  exit: { 
    opacity: 0, 
    height: 0,
    transition: {
      duration: 0.3
    }
  }
};

const modalOverlayVariants = {
  hidden: { opacity: 0 },
  visible: { 
    opacity: 1,
    transition: {
      duration: 0.2
    }
  },
  exit: { 
    opacity: 0,
    transition: {
      duration: 0.2
    }
  }
};

const modalContainerVariants = {
  hidden: { opacity: 0, y: 50, scale: 0.95 },
  visible: { 
    opacity: 1, 
    y: 0,
    scale: 1,
    transition: {
      duration: 0.3,
      ease: 'easeOut'
    }
  },
  exit: { 
    opacity: 0,
    y: 50,
    scale: 0.95,
    transition: {
      duration: 0.2
    }
  }
};

const PricingCalculator: React.FC = () => {
  const [serviceType, setServiceType] = useState('assessment');
  const [complexity, setComplexity] = useState('moderate');
  const [count, setCount] = useState(5);
  const [success, setSuccess] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  
  // Form state
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const [company, setCompany] = useState('');
  const [requirements, setRequirements] = useState('');
  
  // Form validation
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  const selectedService = serviceTypes.find(type => type.id === serviceType)!;
  const selectedComplexity = complexityLevels.find(level => level.id === complexity)!;
  
  const basePrice = selectedService.basePrice;
  const variablePrice = serviceType === 'optimization' 
    ? selectedService.pricePerServer * count // For optimization, count = hours
    : selectedService.pricePerServer * count; // For others, count = servers/apps
  const complexityMultiplier = selectedComplexity.multiplier;
  
  const totalEstimate = Math.round((basePrice + variablePrice) * complexityMultiplier);
  
  const openModal = () => {
    setIsModalOpen(true);
  };
  
  const closeModal = () => {
    setIsModalOpen(false);
  };
  
  const validateForm = () => {
    const newErrors: Record<string, string> = {};
    
    if (!name.trim()) {
      newErrors.name = 'Name is required';
    }
    
    if (!email.trim()) {
      newErrors.email = 'Email is required';
    } else if (!/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i.test(email)) {
      newErrors.email = 'Invalid email address';
    }
    
    if (!phone.trim()) {
      newErrors.phone = 'Phone number is required';
    }
    
    if (!company.trim()) {
      newErrors.company = 'Company name is required';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };
  
  const sendEmail = () => {
    setIsSubmitting(true);
    
    // Get current date for the email template
    const currentDate = new Date().toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
    
    // Prepare email data with all fields matching the template
    const emailData = {
      // Contact Information
      name,
      email,
      phone,
      company,
      
      // Service Details
      serviceType: selectedService.title,
      complexity: selectedComplexity.title,
      count,
      
      // Pricing Breakdown
      basePrice: basePrice.toLocaleString(),
      variablePrice: variablePrice.toLocaleString(),
      complexityMultiplier,
      totalEstimate: totalEstimate.toLocaleString(),
      
      // Additional Fields
      requirements: requirements || "No additional requirements specified.",
      current_date: currentDate,
      
      // These are conditional fields used in the template
      // The actual template will handle the conditional display with Handlebars
      isHourly: serviceType === 'optimization'
    };
    
    // FOR TESTING - Simulate EmailJS send
    // Remove this block and uncomment the EmailJS code below when ready to deploy
    // setTimeout(() => {
    //   console.log('Email data for template:', emailData);
      
    //   // Reset form and close modal
    //   setName('');
    //   setEmail('');
    //   setPhone('');
    //   setCompany('');
    //   setRequirements('');
    //   setErrors({});
    //   setIsSubmitting(false);
    //   setIsModalOpen(false);
      
    //   // Show success message
    //   setSuccess(true);
    //   setTimeout(() => {
    //     setSuccess(false);
    //   }, 5000);
    // }, 1000);
    
    // REAL IMPLEMENTATION WITH EMAILJS
    // Uncomment this block when ready to use EmailJS
    emailjs.send(
      'mosttn18@gmail.com', // Replace with your EmailJS service ID
      'template_9lu3gzb', // Replace with your template ID that uses the HTML you provided
      emailData,
      'hz-jZI5Vs-LNtGM4T' // Replace with your EmailJS public key
    )
    .then((result) => {
      console.log('Email successfully sent!', result.text);
      
      // Reset form and close modal
      setName('');
      setEmail('');
      setPhone('');
      setCompany('');
      setRequirements('');
      setErrors({});
      setIsSubmitting(false);
      setIsModalOpen(false);
      
      // Show success message
      setSuccess(true);
      setTimeout(() => {
        setSuccess(false);
      }, 5000);
    })
    .catch((error) => {
      console.error('Failed to send email:', error);
      setIsSubmitting(false);
    });
  };
  
  const handleSubmit = () => {
    if (validateForm()) {
      sendEmail();
    }
  };
  
  const handleRequestQuote = () => {
    openModal();
  };
  
  return (
    <CalculatorContainer
      initial="hidden"
      animate="visible"
      variants={calculatorVariants}
    >
      <CalculatorTitle>Estimate Your Service Cost</CalculatorTitle>
      
      <FormGroup>
        <Label>
          Service Type
          <InfoIcon 
            data-tooltip="Select the specific service you need"
          >
            <Icon icon={FiInfo} size={16} />
          </InfoIcon>
        </Label>
        <OptionGrid>
          {serviceTypes.map(type => (
            <OptionCard 
              key={type.id} 
              selected={serviceType === type.id}
            >
              <RadioInput
                type="radio"
                name="serviceType"
                value={type.id}
                checked={serviceType === type.id}
                onChange={() => setServiceType(type.id)}
              />
              <OptionTitle selected={serviceType === type.id}>
                {type.title}
              </OptionTitle>
              <OptionDescription>
                {type.description}
              </OptionDescription>
            </OptionCard>
          ))}
        </OptionGrid>
      </FormGroup>
      
      <FormGroup>
        <Label>
          Project Complexity
          <InfoIcon 
            data-tooltip="Assess your application's complexity level"
          >
            <Icon icon={FiInfo} size={16} />
          </InfoIcon>
        </Label>
        <OptionGrid>
          {complexityLevels.map(level => (
            <OptionCard 
              key={level.id} 
              selected={complexity === level.id}
            >
              <RadioInput
                type="radio"
                name="complexity"
                value={level.id}
                checked={complexity === level.id}
                onChange={() => setComplexity(level.id)}
              />
              <OptionTitle selected={complexity === level.id}>
                {level.title}
              </OptionTitle>
              <OptionDescription>
                {level.description}
              </OptionDescription>
            </OptionCard>
          ))}
        </OptionGrid>
      </FormGroup>
      
      <FormGroup>
        <Label>
          {serviceType === 'optimization' ? 'Estimated Hours' : 'Number of Servers/Applications'}
          <InfoIcon 
            data-tooltip={serviceType === 'optimization' ? 
              "Estimated hours of implementation assistance needed" : 
              "How many servers or applications are included in your project"}
          >
            <Icon icon={FiInfo} size={16} />
          </InfoIcon>
        </Label>
        <RangeContainer>
          <RangeValue>{count}</RangeValue>
          <RangeInput
            type="range"
            min="1"
            max={serviceType === 'optimization' ? "20" : "10"}
            value={count}
            onChange={(e) => setCount(parseInt(e.target.value))}
          />
          <RangeLabels>
            <span>1</span>
            <span>{serviceType === 'optimization' ? "10" : "5"}</span>
            <span>{serviceType === 'optimization' ? "20" : "10"}</span>
          </RangeLabels>
        </RangeContainer>
      </FormGroup>
      
      <Divider />
      
      <EstimateContainer
        initial="hidden"
        animate="visible"
        variants={estimateVariants}
      >
        <EstimateTitle>Your Cost Estimate</EstimateTitle>
        <EstimateRow>
          <EstimateLabel>Base Service Fee:</EstimateLabel>
          <EstimateValue>${basePrice.toLocaleString()}</EstimateValue>
        </EstimateRow>
        <EstimateRow>
          <EstimateLabel>{serviceType === 'optimization' ? 'Hourly Rate Cost:' : 'Per-Server/App Cost:'}</EstimateLabel>
          <EstimateValue>${variablePrice.toLocaleString()}</EstimateValue>
        </EstimateRow>
        <EstimateRow>
          <EstimateLabel>Complexity Multiplier:</EstimateLabel>
          <EstimateValue>{complexityMultiplier}x</EstimateValue>
        </EstimateRow>
        <TotalRow>
          <TotalLabel>Estimated Total:</TotalLabel>
          <TotalValue>${totalEstimate.toLocaleString()}</TotalValue>
        </TotalRow>
      </EstimateContainer>
      
      <ButtonContainer>
        <Button onClick={handleRequestQuote} size="lg">
          Request Detailed Quote
        </Button>
      </ButtonContainer>
      
      {success && (
        <SuccessMessage
          initial="hidden"
          animate="visible"
          exit="exit"
          variants={successVariants}
        >
          <Icon icon={FiCheckCircle} size={16} />
          <div>Your detailed quote request has been submitted successfully. We'll review it and get back to you within 24-48 hours.</div>
        </SuccessMessage>
      )}
      
      <AnimatePresence>
        {isModalOpen && (
          <ModalOverlay
            initial="hidden"
            animate="visible"
            exit="exit"
            variants={modalOverlayVariants}
            onClick={closeModal}
          >
            <ModalContainer
              variants={modalContainerVariants}
              onClick={(e) => e.stopPropagation()}
            >
              <ModalHeader>
                <ModalTitle>Request a Detailed Quote</ModalTitle>
                <CloseButton onClick={closeModal}>
                  <Icon icon={FiX} size={24} />
                </CloseButton>
              </ModalHeader>
              
              <ModalBody>
                <QuoteDetail>
                  <QuoteDetailRow>
                    <QuoteDetailLabel>Service:</QuoteDetailLabel>
                    <QuoteDetailValue>{selectedService.title}</QuoteDetailValue>
                  </QuoteDetailRow>
                  <QuoteDetailRow>
                    <QuoteDetailLabel>Complexity:</QuoteDetailLabel>
                    <QuoteDetailValue>{selectedComplexity.title}</QuoteDetailValue>
                  </QuoteDetailRow>
                  <QuoteDetailRow>
                    <QuoteDetailLabel>
                      {serviceType === 'optimization' ? 'Hours:' : 'Servers/Applications:'}
                    </QuoteDetailLabel>
                    <QuoteDetailValue>{count}</QuoteDetailValue>
                  </QuoteDetailRow>
                  <QuoteDetailRow>
                    <QuoteDetailLabel>Base Fee:</QuoteDetailLabel>
                    <QuoteDetailValue>${basePrice.toLocaleString()}</QuoteDetailValue>
                  </QuoteDetailRow>
                  <QuoteDetailRow>
                    <QuoteDetailLabel>
                      {serviceType === 'optimization' ? 'Hourly Rate Cost:' : 'Per-Server/App Cost:'}
                    </QuoteDetailLabel>
                    <QuoteDetailValue>${variablePrice.toLocaleString()}</QuoteDetailValue>
                  </QuoteDetailRow>
                  <QuoteDetailRow>
                    <QuoteDetailLabel>Complexity Multiplier:</QuoteDetailLabel>
                    <QuoteDetailValue>{complexityMultiplier}x</QuoteDetailValue>
                  </QuoteDetailRow>
                  <QuoteDetailTotal>
                    <div>Total Estimate:</div>
                    <div>${totalEstimate.toLocaleString()}</div>
                  </QuoteDetailTotal>
                </QuoteDetail>
                
                <InputGroup>
                  <InputLabel>Your Name</InputLabel>
                  <div style={{ position: 'relative' }}>
                    <InputIcon>
                      <Icon icon={FiUser} size={16} />
                    </InputIcon>
                    <IconInput
                      type="text"
                      placeholder="Enter your full name"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                    />
                  </div>
                  {errors.name && <ErrorMessage>{errors.name}</ErrorMessage>}
                </InputGroup>
                
                <InputGroup>
                  <InputLabel>Email Address</InputLabel>
                  <div style={{ position: 'relative' }}>
                    <InputIcon>
                      <Icon icon={FiMail} size={16} />
                    </InputIcon>
                    <IconInput
                      type="email"
                      placeholder="Enter your email address"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                    />
                  </div>
                  {errors.email && <ErrorMessage>{errors.email}</ErrorMessage>}
                </InputGroup>
                
                <InputGroup>
                  <InputLabel>Phone Number</InputLabel>
                  <div style={{ position: 'relative' }}>
                    <InputIcon>
                      <Icon icon={FiPhone} size={16} />
                    </InputIcon>
                    <IconInput
                      type="tel"
                      placeholder="Enter your phone number"
                      value={phone}
                      onChange={(e) => setPhone(e.target.value)}
                    />
                  </div>
                  {errors.phone && <ErrorMessage>{errors.phone}</ErrorMessage>}
                </InputGroup>
                
                <InputGroup>
                  <InputLabel>Company</InputLabel>
                  <div style={{ position: 'relative' }}>
                    <InputIcon>
                      <Icon icon={FiBriefcase} size={16} />
                    </InputIcon>
                    <IconInput
                      type="text"
                      placeholder="Enter your company name"
                      value={company}
                      onChange={(e) => setCompany(e.target.value)}
                    />
                  </div>
                  {errors.company && <ErrorMessage>{errors.company}</ErrorMessage>}
                </InputGroup>
                
                <InputLabel>Additional Requirements</InputLabel>
                <Textarea
                  placeholder="Please provide any additional requirements or questions..."
                  value={requirements}
                  onChange={(e) => setRequirements(e.target.value)}
                />
              </ModalBody>
              
              <ModalFooter>
                <Button
                  variant="outline"
                  onClick={closeModal}
                  disabled={isSubmitting}
                >
                  Cancel
                </Button>
                <Button
                  onClick={handleSubmit}
                  disabled={isSubmitting}
                >
                  {isSubmitting ? 'Sending...' : 'Submit Quote Request'}
                </Button>
              </ModalFooter>
            </ModalContainer>
          </ModalOverlay>
        )}
      </AnimatePresence>
    </CalculatorContainer>
  );
};

export default PricingCalculator;