import React, { useState } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import { Formik, Form, Field, ErrorMessage, FormikHelpers } from 'formik';
import * as Yup from 'yup';
import { theme } from '../../../styles/theme';
import Button from '../../ui/Button';
import { FiCheckCircle, FiAlertCircle } from 'react-icons/fi';
import Icon from '../../ui/Icon';
import { apiService, CreateInquiryRequest } from '../../../services/api';

// Fixed icon components
const CheckCircleIcon = () => <Icon icon={FiCheckCircle} size={16} />;
const AlertCircleIcon = () => <Icon icon={FiAlertCircle} size={16} />;

interface FormValues {
  name: string;
  email: string;
  company: string;
  phone: string;
  services: string[];
  message: string;
}

const FormContainer = styled.div`
  background-color: rgba(255, 255, 255, 0.05);
  border-radius: ${theme.borderRadius.lg};
  padding: ${theme.space[6]};
  
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

const StyledField = styled(Field)<{ $error?: boolean }>`
  width: 100%;
  background-color: rgba(255, 255, 255, 0.08);
  border: 1px solid ${props => props.$error ? theme.colors.danger : 'rgba(255, 255, 255, 0.2)'};
  border-radius: ${theme.borderRadius.md};
  padding: ${theme.space[3]};
  color: ${theme.colors.white};
  font-size: ${theme.fontSizes.md};
  transition: ${theme.transitions.fast};
  
  &:focus {
    outline: none;
    border-color: ${theme.colors.secondary};
    box-shadow: 0 0 0 2px ${theme.colors.secondary + '40'};
  }
  
  &::placeholder {
    color: rgba(255, 255, 255, 0.5);
  }
`;

const TextArea = styled(Field)`
  width: 100%;
  background-color: rgba(255, 255, 255, 0.08);
  border: 1px solid ${props => props.error ? theme.colors.danger : 'rgba(255, 255, 255, 0.2)'};
  border-radius: ${theme.borderRadius.md};
  padding: ${theme.space[3]};
  color: ${theme.colors.white};
  font-size: ${theme.fontSizes.md};
  transition: ${theme.transitions.fast};
  min-height: 150px;
  resize: vertical;
  
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

const SubmitButton = styled(Button)`
  width: 100%;
  margin-top: ${theme.space[4]};
`;

const FormMessage = styled(motion.div)<{ success?: boolean }>`
  background-color: ${props => props.success ? theme.colors.success + '20' : theme.colors.danger + '20'};
  color: ${props => props.success ? theme.colors.success : theme.colors.danger};
  padding: ${theme.space[4]};
  border-radius: ${theme.borderRadius.md};
  display: flex;
  align-items: center;
  margin-top: ${theme.space[5]};
  
  svg {
    margin-right: ${theme.space[3]};
    font-size: ${theme.fontSizes.xl};
  }
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
  const [formState, setFormState] = useState<'idle' | 'success' | 'error'>('idle');
  
  const initialValues: FormValues = {
    name: '',
    email: '',
    company: '',
    phone: '',
    services: [],
    message: '',
  };
  
  const handleSubmit = async (
    values: FormValues,
    { setSubmitting, resetForm }: FormikHelpers<FormValues>
  ) => {
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
        setFormState('success');
        resetForm();
        
        // Reset success message after 5 seconds
        setTimeout(() => {
          setFormState('idle');
        }, 5000);
      } else {
        throw new Error(response.error || 'Form submission failed');
      }
    } catch (error) {
      console.error('Failed to submit form:', error);
      setFormState('error');
      
      // Reset error message after 5 seconds
      setTimeout(() => {
        setFormState('idle');
      }, 5000);
    } finally {
      setSubmitting(false);
    }
  };
  
  return (
    <FormContainer>
      <FormTitle>Send Us a Message</FormTitle>
      
      <Formik
        initialValues={initialValues}
        validationSchema={validationSchema}
        onSubmit={handleSubmit}
      >
        {({ isSubmitting, errors, touched }) => (
          <Form>
            <FieldRow>
              <FormGroup>
                <Label htmlFor="name">Full Name *</Label>
                <StyledField
                  type="text"
                  id="name"
                  name="name"
                  placeholder="John Doe"
                  $error={Boolean(errors.name && touched.name) || undefined}
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
                  $error={Boolean(errors.email && touched.email) || undefined}
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
              <Label htmlFor="services">Services Required *</Label>
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
                error={Boolean(errors.message && touched.message)}
              />
              <ErrorMessage name="message">
                {(msg) => (
                  <ErrorText>
                    <AlertCircleIcon /> {msg}
                  </ErrorText>
                )}
              </ErrorMessage>
            </FormGroup>
            
            <SubmitButton type="submit" disabled={isSubmitting}>
              {isSubmitting ? 'Sending...' : 'Send Message'}
            </SubmitButton>
          </Form>
        )}
      </Formik>
      
      <AnimatePresence>
        {formState === 'success' && (
          <FormMessage
            success
            initial="hidden"
            animate="visible"
            exit="exit"
            variants={messageVariants}
            viewport={{ once: true }}
          >
            <CheckCircleIcon />
            <div>Your inquiry has been submitted successfully! We'll get back to you as soon as possible.</div>
          </FormMessage>
        )}
        
        {formState === 'error' && (
          <FormMessage
            initial="hidden"
            animate="visible"
            exit="exit"
            variants={messageVariants}
            viewport={{ once: true }}
          >
            <Icon icon={FiAlertCircle} size={24} />
            <div>There was an error submitting your inquiry. Please try again later or email us directly at info@cloudpartner.pro</div>
          </FormMessage>
        )}
      </AnimatePresence>
    </FormContainer>
  );
};

// Memoize the component to prevent unnecessary re-renders
const MemoizedContactForm = React.memo(ContactForm);
export default MemoizedContactForm;