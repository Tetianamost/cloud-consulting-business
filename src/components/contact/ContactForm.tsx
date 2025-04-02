import React from 'react';
import { Formik, Form, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';

// Validation schema
const ContactSchema = Yup.object().shape({
  name: Yup.string()
    .min(2, 'Name is too short')
    .max(50, 'Name is too long')
    .required('Name is required'),
  email: Yup.string()
    .email('Invalid email')
    .required('Email is required'),
  company: Yup.string()
    .required('Company name is required'),
  phone: Yup.string()
    .matches(/^[0-9+-\s()]*$/, 'Invalid phone number'),
  message: Yup.string()
    .min(10, 'Message is too short')
    .required('Message is required'),
});

const ContactForm: React.FC = () => {
  return (
    <section id="contact" className="section" style={sectionStyle}>
      <div className="container">
        <h2 className="section-title">Get In Touch</h2>
        <p style={sectionDescriptionStyle}>
          Ready to start your cloud journey? Contact us for a free consultation
        </p>
        
        <div style={formContainerStyle}>
          <div style={contactInfoStyle}>
            <div style={contactBlockStyle}>
              <h3 style={contactHeadingStyle}>Contact Information</h3>
              <p style={contactItemStyle}>
                <span style={contactIconStyle}>📞</span> (555) 123-4567
              </p>
              <p style={contactItemStyle}>
                <span style={contactIconStyle}>✉️</span> info@cloudmigrationpro.com
              </p>
              <p style={contactItemStyle}>
                <span style={contactIconStyle}>🏢</span> 123 Business St, Suite 500<br />
                San Francisco, CA 94103
              </p>
            </div>
            
            <div style={contactBlockStyle}>
              <h3 style={contactHeadingStyle}>Business Hours</h3>
              <p style={contactItemStyle}>Monday - Friday: 9am - 6pm PST</p>
              <p style={contactItemStyle}>Saturday - Sunday: Closed</p>
            </div>
          </div>
          
          <div style={formStyle}>
            <Formik
              initialValues={{
                name: '',
                email: '',
                company: '',
                phone: '',
                message: '',
              }}
              validationSchema={ContactSchema}
              onSubmit={(values, { setSubmitting, resetForm }) => {
                // Here you would typically send the form data to your backend
                console.log(values);
                alert('Thank you for your message! We will get back to you soon.');
                resetForm();
                setSubmitting(false);
              }}
            >
              {({ isSubmitting }) => (
                <Form>
                  <div className="form-group">
                    <label htmlFor="name" style={labelStyle}>Name</label>
                    <Field 
                      type="text" 
                      name="name" 
                      className="form-control" 
                      style={inputStyle} 
                    />
                    <ErrorMessage name="name" component="div" className="error-message" />
                  </div>
                  
                  <div className="form-group">
                    <label htmlFor="email" style={labelStyle}>Email</label>
                    <Field 
                      type="email" 
                      name="email" 
                      className="form-control" 
                      style={inputStyle} 
                    />
                    <ErrorMessage name="email" component="div" className="error-message" />
                  </div>
                  
                  <div className="form-group">
                    <label htmlFor="company" style={labelStyle}>Company</label>
                    <Field 
                      type="text" 
                      name="company" 
                      className="form-control" 
                      style={inputStyle} 
                    />
                    <ErrorMessage name="company" component="div" className="error-message" />
                  </div>
                  
                  <div className="form-group">
                    <label htmlFor="phone" style={labelStyle}>Phone (optional)</label>
                    <Field 
                      type="text" 
                      name="phone" 
                      className="form-control" 
                      style={inputStyle} 
                    />
                    <ErrorMessage name="phone" component="div" className="error-message" />
                  </div>
                  
                  <div className="form-group">
                    <label htmlFor="message" style={labelStyle}>Message</label>
                    <Field 
                      as="textarea" 
                      name="message" 
                      className="form-control" 
                      style={{ ...inputStyle, height: '150px' }} 
                    />
                    <ErrorMessage name="message" component="div" className="error-message" />
                  </div>
                  
                  <button 
                    type="submit" 
                    className="btn btn-primary" 
                    style={submitButtonStyle} 
                    disabled={isSubmitting}
                  >
                    {isSubmitting ? 'Sending...' : 'Send Message'}
                  </button>
                </Form>
              )}
            </Formik>
          </div>
        </div>
      </div>
    </section>
  );
};

// Styles
const sectionStyle: React.CSSProperties = {
  background: '#fff',
  padding: '80px 0'
};

const sectionDescriptionStyle: React.CSSProperties = {
  fontSize: '1.2rem',
  maxWidth: '800px',
  margin: '0 auto 40px',
  textAlign: 'center',
  color: '#666'
};

const formContainerStyle: React.CSSProperties = {
  display: 'grid',
  gridTemplateColumns: '1fr 2fr',
  gap: '40px',
  maxWidth: '1100px',
  margin: '0 auto'
};

const contactInfoStyle: React.CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: '30px'
};

const contactBlockStyle: React.CSSProperties = {
  backgroundColor: 'var(--light-color)',
  padding: '25px',
  borderRadius: '10px',
  boxShadow: '0 3px 10px rgba(0, 0, 0, 0.1)'
};

const contactHeadingStyle: React.CSSProperties = {
  fontSize: '1.4rem',
  marginBottom: '15px',
  color: 'var(--primary-color)'
};

const contactItemStyle: React.CSSProperties = {
  marginBottom: '10px',
  display: 'flex',
  alignItems: 'flex-start',
  gap: '10px'
};

const contactIconStyle: React.CSSProperties = {
  fontSize: '1.2rem'
};

const formStyle: React.CSSProperties = {
  background: 'var(--light-color)',
  padding: '30px',
  borderRadius: '10px',
  boxShadow: '0 3px 10px rgba(0, 0, 0, 0.1)'
};

const labelStyle: React.CSSProperties = {
  display: 'block',
  marginBottom: '5px',
  fontWeight: 'bold'
};

const inputStyle: React.CSSProperties = {
  width: '100%',
  padding: '12px',
  border: '1px solid #ddd',
  borderRadius: '5px',
  fontSize: '1rem',
  transition: 'border-color 0.3s'
};

const errorStyle: React.CSSProperties = {
  color: 'var(--danger-color)',
  fontSize: '0.9rem',
  marginTop: '5px'
};

const submitButtonStyle: React.CSSProperties = {
  width: '100%',
  padding: '15px',
  fontSize: '1.1rem',
  marginTop: '10px'
};

export default ContactForm;