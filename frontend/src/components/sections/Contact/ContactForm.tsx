import React, { useState, useEffect } from 'react';
import styled, { keyframes, css } from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import { Formik, Form, Field, ErrorMessage, FormikHelpers } from 'formik';
import * as Yup from 'yup';
import { theme } from '../../../styles/theme';
import { Button } from '../../ui/button';
import { FiCheckCircle, FiAlertCircle, FiLoader, FiSend, FiMail, FiClock } from 'react-icons/fi';
import Icon from '../../ui/Icon';
import { apiService, CreateInquiryRequest } from '../../../services/api';

// Fixed icon components
const CheckCircleIcon = () => <Icon icon={FiCheckCircle} size={16} />;
const AlertCircleIcon = () => <Icon icon={FiAlertCircle} size={16} />;
const LoaderIcon = () => <Icon icon={FiLoader} size={16} />;
const SendIcon = () => <Icon icon={FiSend} size={16} />;
const MailIcon = () => <Icon icon={FiMail} size={24} />;
const ClockIcon = () => <Icon icon={FiClock} size={16} />;

interface FormValues {
  name: string;
  email: string;
  company: string;
  phone: string;
  services: string[];
  message: string;
}

interface SubmissionState {
  status: 'idle' | 'submitting' | 'success' | 'error';
  progress: number;
  message?: string;
  inquiryId?: string;
}

// Animations
const spin = keyframes`
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
`;

const pulse = keyframes`
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
`;

const slideInUp = keyframes`
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
`;

const FormContainer = styled.div`
  background-color: rgba(255, 255, 255, 0.05);
  border-radius: ${theme.borderRadius.lg};
  padding: ${theme.space[6]};
  position: relative;
  
  @media (max-width: ${theme.breakpoints.md}) {
    padding: ${theme.space[4]};
  }
`;

const FormTitle = styled.h3`
  font-size: ${theme.fontSizes.xl};
  color: ${theme.colors.white};
  margin-bottom: ${theme.space[4]};
`;

const FieldRow = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: ${theme.space[4]};
  
  @media (max-width: ${theme.breakpoints.sm}) {
    grid-template-columns: 1fr;
  }
`;

const FormGroup = styled.div`
  margin-bottom: ${theme.space[4]};
`;

const Label = styled.label`
  display: block;
  font-weight: ${theme.fontWeights.medium};
  margin-bottom: ${theme.space[2]};
  color: ${theme.colors.gray200};
`;

const StyledField = styled(Field)<{ $error?: boolean; $success?: boolean }>`
  width: 100%;
  background-color: rgba(255, 255, 255, 0.08);
  border: 1px solid ${props => {
    if (props.$error) return theme.colors.danger;
    if (props.$success) return theme.colors.success;
    return 'rgba(255, 255, 255, 0.2)';
  }};
  border-radius: ${theme.borderRadius.md};
  padding: ${theme.space[3]};
  color: ${theme.colors.white};
  font-size: ${theme.fontSizes.md};
  transition: ${theme.transitions.fast};
  position: relative;
  
  &:focus {
    outline: none;
    border-color: ${theme.colors.secondary};
    box-shadow: 0 0 0 2px ${theme.colors.secondary + '40'};
  }
  
  &::placeholder {
    color: rgba(255, 255, 255, 0.5);
  }
  
  ${props => props.$success && css`
    &::after {
      content: '✓';
      position: absolute;
      right: 12px;
      top: 50%;
      transform: translateY(-50%);
      color: ${theme.colors.success};
      font-weight: bold;
    }
  `}
`;

const TextArea = styled(Field)<{ error?: boolean; $success?: boolean }>`
  width: 100%;
  background-color: rgba(255, 255, 255, 0.08);
  border: 1px solid ${props => {
    if (props.error) return theme.colors.danger;
    if (props.$success) return theme.colors.success;
    return 'rgba(255, 255, 255, 0.2)';
  }};
  border-radius: ${theme.borderRadius.md};
  padding: ${theme.space[3]};
  color: ${theme.colors.white};
  font-size: ${theme.fontSizes.md};
  transition: ${theme.transitions.fast};
  min-height: 150px;
  resize: vertical;
  position: relative;
  
  &:focus {
    outline: none;
    border-color: ${theme.colors.secondary};
    box-shadow: 0 0 0 2px ${theme.colors.secondary + '40'};
  }
  
  &::placeholder {
    color: rgba(255, 255, 255, 0.5);
  }
`;

const ErrorText = styled.div`
  color: ${theme.colors.danger};
  font-size: ${theme.fontSizes.sm};
  margin-top: ${theme.space[1]};
  display: flex;
  align-items: center;
  
  svg {
    margin-right: ${theme.space[1]};
  }
`;

const CheckboxGroup = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: ${theme.space[3]};
  margin-top: ${theme.space[2]};
  
  @media (max-width: ${theme.breakpoints.sm}) {
    grid-template-columns: 1fr;
  }
`;

const CheckboxLabel = styled.label`
  display: flex;
  align-items: center;
  cursor: pointer;
  color: ${theme.colors.gray200};
  transition: ${theme.transitions.fast};
  
  &:hover {
    color: ${theme.colors.white};
  }
`;

const CheckboxInput = styled(Field)`
  appearance: none;
  width: 18px;
  height: 18px;
  background-color: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: ${theme.borderRadius.sm};
  margin-right: ${theme.space[2]};
  position: relative;
  cursor: pointer;
  flex-shrink: 0;
  
  &:checked {
    background-color: ${theme.colors.secondary};
    border-color: ${theme.colors.secondary};
    
    &::after {
      content: '';
      position: absolute;
      top: 3px;
      left: 6px;
      width: 5px;
      height: 10px;
      border: solid white;
      border-width: 0 2px 2px 0;
      transform: rotate(45deg);
    }
  }
  
  &:focus {
    outline: none;
    box-shadow: 0 0 0 2px ${theme.colors.secondary + '40'};
  }
`;

const SubmitButton = styled(Button)<{ $isSubmitting?: boolean }>`
  width: 100%;
  margin-top: ${theme.space[4]};
  position: relative;
  overflow: hidden;
  
  ${props => props.$isSubmitting && css`
    pointer-events: none;
    
    .spinner {
      animation: ${spin} 1s linear infinite;
    }
  `}
`;

const ProgressBar = styled(motion.div)`
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background: linear-gradient(90deg, 
    ${theme.colors.secondary}80, 
    ${theme.colors.secondary}
  );
  border-radius: ${theme.borderRadius.md};
  z-index: 1;
`;

const ButtonContent = styled.div`
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: ${theme.space[2]};
`;

const FormMessage = styled(motion.div)<{ success?: boolean }>`
  background-color: ${props => props.success ? theme.colors.success + '20' : theme.colors.danger + '20'};
  color: ${props => props.success ? theme.colors.success : theme.colors.danger};
  padding: ${theme.space[4]};
  border-radius: ${theme.borderRadius.md};
  display: flex;
  align-items: flex-start;
  margin-top: ${theme.space[5]};
  border: 1px solid ${props => props.success ? theme.colors.success + '40' : theme.colors.danger + '40'};
  
  svg {
    margin-right: ${theme.space[3]};
    font-size: ${theme.fontSizes.xl};
    flex-shrink: 0;
    margin-top: 2px;
  }
`;

const SuccessMessage = styled.div`
  animation: ${slideInUp} 0.5s ease-out;
`;

const SuccessTitle = styled.h4`
  margin: 0 0 ${theme.space[2]} 0;
  font-size: ${theme.fontSizes.lg};
  font-weight: ${theme.fontWeights.semibold};
  color: ${theme.colors.success};
`;

const SuccessContent = styled.div`
  font-size: ${theme.fontSizes.md};
  line-height: 1.5;
  
  p {
    margin: 0 0 ${theme.space[2]} 0;
    
    &:last-child {
      margin-bottom: 0;
    }
  }
`;

const NextSteps = styled.div`
  margin-top: ${theme.space[3]};
  padding-top: ${theme.space[3]};
  border-top: 1px solid ${theme.colors.success + '30'};
`;

const StepList = styled.ul`
  margin: ${theme.space[2]} 0 0 0;
  padding-left: ${theme.space[4]};
  
  li {
    margin-bottom: ${theme.space[1]};
    display: flex;
    align-items: center;
    
    svg {
      margin-right: ${theme.space[2]};
      margin-left: -${theme.space[4]};
      color: ${theme.colors.success};
    }
  }
`;

const InquiryId = styled.div`
  background-color: rgba(255, 255, 255, 0.1);
  padding: ${theme.space[2]} ${theme.space[3]};
  border-radius: ${theme.borderRadius.sm};
  font-family: monospace;
  font-size: ${theme.fontSizes.sm};
  margin-top: ${theme.space[2]};
  border: 1px solid rgba(255, 255, 255, 0.2);
`;

// Form validation schema
const validationSchema = Yup.object().shape({
  name: Yup.string()
    .min(2, 'Name is too short')
    .max(50, 'Name is too long')
    .required('Name is required'),
  email: Yup.string()
    .email('Invalid email format')
    .required('Email is required'),
  company: Yup.string()
    .optional(),
  phone: Yup.string()
    .matches(/^[0-9+-\s()]*$/, 'Invalid phone number format')
    .optional(),
  services: Yup.array()
    .min(1, 'Please select at least one service'),
  message: Yup.string()
    .min(1, 'Message is required')
    .required('Message is required'),
});

const serviceOptions = [
  { id: 'assessment', label: 'Cloud Assessment' },
  { id: 'migration', label: 'Cloud Migration' },
  { id: 'optimization', label: 'Cloud Optimization' },
  { id: 'architecture_review', label: 'Architecture Review' },
];

// Animation variants
const messageVariants = {
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

const ContactForm: React.FC = () => {
  const [submissionState, setSubmissionState] = useState<SubmissionState>({
    status: 'idle',
    progress: 0
  });
  const [validatedFields, setValidatedFields] = useState<Set<string>>(new Set());
  
  const initialValues: FormValues = {
    name: '',
    email: '',
    company: '',
    phone: '',
    services: [],
    message: '',
  };
  
  // Progress simulation for better UX
  const simulateProgress = (callback: () => void) => {
    setSubmissionState(prev => ({ ...prev, status: 'submitting', progress: 0 }));
    
    const steps = [
      { progress: 20, delay: 200, message: 'Validating information...' },
      { progress: 40, delay: 300, message: 'Connecting to server...' },
      { progress: 60, delay: 400, message: 'Processing inquiry...' },
      { progress: 80, delay: 300, message: 'Generating report...' },
      { progress: 100, delay: 200, message: 'Finalizing...' }
    ];
    
    let currentStep = 0;
    const progressInterval = setInterval(() => {
      if (currentStep < steps.length) {
        const step = steps[currentStep];
        setSubmissionState(prev => ({ 
          ...prev, 
          progress: step.progress, 
          message: step.message 
        }));
        currentStep++;
      } else {
        clearInterval(progressInterval);
        callback();
      }
    }, steps[currentStep]?.delay || 300);
  };

  const handleSubmit = async (
    values: FormValues,
    { setSubmitting, resetForm }: FormikHelpers<FormValues>
  ) => {
    simulateProgress(async () => {
      try {
        // Prepare data for backend API
        const inquiryData: CreateInquiryRequest = {
          name: values.name.trim(),
          email: values.email.trim(),
          company: values.company.trim() || undefined,
          phone: values.phone.trim() || undefined,
          services: values.services, // Send as array of service IDs
          message: values.message.trim(),
          source: 'contact_form'
        };
        
        // Submit to backend API
        const response = await apiService.createInquiry(inquiryData);
        
        if (response.success) {
          console.log('Contact form submitted successfully!', response.data);
          setSubmissionState({
            status: 'success',
            progress: 100,
            message: 'Inquiry submitted successfully!',
            inquiryId: response.data?.id
          });
          resetForm();
          setValidatedFields(new Set());
          
          // Reset success message after 10 seconds
          setTimeout(() => {
            setSubmissionState({ status: 'idle', progress: 0 });
          }, 10000);
        } else {
          throw new Error(response.error || 'Form submission failed');
        }
      } catch (error) {
        console.error('Failed to submit form:', error);
        setSubmissionState({
          status: 'error',
          progress: 0,
          message: error instanceof Error ? error.message : 'Form submission failed'
        });
        
        // Reset error message after 8 seconds
        setTimeout(() => {
          setSubmissionState({ status: 'idle', progress: 0 });
        }, 8000);
      } finally {
        setSubmitting(false);
      }
    });
  };

  // Real-time field validation
  const handleFieldValidation = (fieldName: string, value: any, errors: any) => {
    const isValid = !errors[fieldName] && value && value.toString().trim() !== '';
    
    if (isValid && !validatedFields.has(fieldName)) {
      setValidatedFields(prev => new Set(Array.from(prev).concat(fieldName)));
    } else if (!isValid && validatedFields.has(fieldName)) {
      setValidatedFields(prev => {
        const newSet = new Set(Array.from(prev));
        newSet.delete(fieldName);
        return newSet;
      });
    }
  };
  
  return (
    <FormContainer>
      <FormTitle>Send Us a Message</FormTitle>
      
      <Formik
        initialValues={initialValues}
        validationSchema={validationSchema}
        onSubmit={handleSubmit}
        validate={(values) => {
          // Real-time validation for instant feedback
          Object.keys(values).forEach(fieldName => {
            handleFieldValidation(fieldName, values[fieldName as keyof FormValues], {});
          });
        }}
      >
        {({ isSubmitting, errors, touched, values }) => {
          // Update field validation status in real-time
          React.useEffect(() => {
            Object.keys(values).forEach(fieldName => {
              handleFieldValidation(fieldName, values[fieldName as keyof FormValues], errors);
            });
          }, [values, errors]);

          return (
          <Form>
            <FieldRow>
              <FormGroup>
                <Label htmlFor="name">Full Name *</Label>
                <StyledField
                  type="text"
                  id="name"
                  name="name"
                  placeholder="John Doe"
                  $error={Boolean(errors.name && touched.name)}
                  $success={validatedFields.has('name') && !errors.name}
                />
                <ErrorMessage name="name">
                  {(msg) => (
                    <ErrorText>
                      <AlertCircleIcon /> {msg}
                    </ErrorText>
                  )}
                </ErrorMessage>
              </FormGroup>
              
              <FormGroup>
                <Label htmlFor="email">Email Address *</Label>
                <StyledField
                  type="email"
                  id="email"
                  name="email"
                  placeholder="john@company.com"
                  $error={Boolean(errors.email && touched.email)}
                  $success={validatedFields.has('email') && !errors.email}
                />
                <ErrorMessage name="email">
                  {(msg) => (
                    <ErrorText>
                      <AlertCircleIcon /> {msg}
                    </ErrorText>
                  )}
                </ErrorMessage>
              </FormGroup>
            </FieldRow>
            
            <FieldRow>
              <FormGroup>
                <Label htmlFor="company">Company Name</Label>
                <StyledField
                  type="text"
                  id="company"
                  name="company"
                  placeholder="Your Company (Optional)"
                  $error={Boolean(errors.company && touched.company)}
                  $success={validatedFields.has('company') && !errors.company}
                />
                <ErrorMessage name="company">
                  {(msg) => (
                    <ErrorText>
                      <AlertCircleIcon /> {msg}
                    </ErrorText>
                  )}
                </ErrorMessage>
              </FormGroup>
              
              <FormGroup>
                <Label htmlFor="phone">Phone Number</Label>
                <StyledField
                  type="text"
                  id="phone"
                  name="phone"
                  placeholder="+1 (555) 123-4567 (Optional)"
                  $error={Boolean(errors.phone && touched.phone)}
                  $success={validatedFields.has('phone') && !errors.phone}
                />
                <ErrorMessage name="phone">
                  {(msg) => (
                    <ErrorText>
                      <AlertCircleIcon /> {msg}
                    </ErrorText>
                  )}
                </ErrorMessage>
              </FormGroup>
            </FieldRow>
            
            <FormGroup>
              <Label id="services-group">Services Required *</Label>
              <CheckboxGroup role="group" aria-labelledby="services-group">
                {serviceOptions.map(option => (
                  <CheckboxLabel key={option.id}>
                    <CheckboxInput
                      type="checkbox"
                      name="services"
                      value={option.id}
                    />
                    {option.label}
                  </CheckboxLabel>
                ))}
              </CheckboxGroup>
              <ErrorMessage name="services">
                {(msg) => (
                  <ErrorText>
                    <AlertCircleIcon /> {msg}
                  </ErrorText>
                )}
              </ErrorMessage>
            </FormGroup>
            
            <FormGroup>
              <Label htmlFor="message">Message *</Label>
              <TextArea
                component="textarea"
                id="message"
                name="message"
                placeholder="Tell us about your project and requirements..."
                error={errors.message && touched.message ? true : undefined}
                $success={validatedFields.has('message') && !errors.message}
              />
              <ErrorMessage name="message">
                {(msg) => (
                  <ErrorText>
                    <AlertCircleIcon /> {msg}
                  </ErrorText>
                )}
              </ErrorMessage>
            </FormGroup>
            
            <SubmitButton 
              type="submit" 
              disabled={isSubmitting || submissionState.status === 'submitting'}
              $isSubmitting={submissionState.status === 'submitting'}
            >
              {submissionState.status === 'submitting' && (
                <ProgressBar
                  initial={{ width: '0%' }}
                  animate={{ width: `${submissionState.progress}%` }}
                  transition={{ duration: 0.3 }}
                />
              )}
              <ButtonContent>
                {submissionState.status === 'submitting' ? (
                  <>
                    <span className="spinner">
                      <LoaderIcon />
                    </span>
                    {submissionState.message || 'Processing...'}
                  </>
                ) : (
                  <>
                    <SendIcon />
                    Send Message
                  </>
                )}
              </ButtonContent>
            </SubmitButton>
          </Form>
          );
        }}
      </Formik>
      
      <AnimatePresence>
        {submissionState.status === 'success' && (
          <FormMessage
            success
            initial="hidden"
            animate="visible"
            exit="exit"
            variants={messageVariants}
            viewport={{ once: true }}
          >
            <MailIcon />
            <SuccessMessage>
              <SuccessTitle>🎉 Inquiry Submitted Successfully!</SuccessTitle>
              <SuccessContent>
                <p>Thank you for reaching out! Your inquiry has been received and is being processed.</p>
                {submissionState.inquiryId && (
                  <InquiryId>
                    <strong>Reference ID:</strong> {submissionState.inquiryId}
                  </InquiryId>
                )}
                <NextSteps>
                  <strong>What happens next:</strong>
                  <StepList>
                    <li>
                      <ClockIcon />
                      Our AI system is generating a preliminary assessment report
                    </li>
                    <li>
                      <MailIcon />
                      You'll receive a confirmation email within 30 seconds
                    </li>
                    <li>
                      <CheckCircleIcon />
                      Our consultant will review and respond within 24 hours
                    </li>
                  </StepList>
                </NextSteps>
                <p style={{ marginTop: '16px', fontSize: '14px', opacity: 0.9 }}>
                  <strong>Need immediate assistance?</strong> Email us directly at{' '}
                  <a href="mailto:info@cloudpartner.pro" style={{ color: 'inherit', textDecoration: 'underline' }}>
                    info@cloudpartner.pro
                  </a>
                </p>
              </SuccessContent>
            </SuccessMessage>
          </FormMessage>
        )}
        
        {submissionState.status === 'error' && (
          <FormMessage
            initial="hidden"
            animate="visible"
            exit="exit"
            variants={messageVariants}
            viewport={{ once: true }}
          >
            <Icon icon={FiAlertCircle} size={24} />
            <div>
              <h4 style={{ margin: '0 0 8px 0', fontSize: '18px' }}>Submission Failed</h4>
              <p style={{ margin: '0 0 12px 0' }}>
                {submissionState.message || 'There was an error submitting your inquiry.'}
              </p>
              <p style={{ margin: 0, fontSize: '14px', opacity: 0.9 }}>
                Please try again or email us directly at{' '}
                <a href="mailto:info@cloudpartner.pro" style={{ color: 'inherit', textDecoration: 'underline' }}>
                  info@cloudpartner.pro
                </a>
              </p>
            </div>
          </FormMessage>
        )}
      </AnimatePresence>
    </FormContainer>
  );
};

// Memoize the component to prevent unnecessary re-renders
const MemoizedContactForm = React.memo(ContactForm);
export default MemoizedContactForm;