import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { theme } from '../../../styles/theme';
import Button from '../../ui/Button';
import { FiInfo, FiCheckCircle } from 'react-icons/fi';
import Icon from '../../ui/Icon';

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
  
  svg {
    margin-right: ${theme.space[3]};
    font-size: ${theme.fontSizes.xl};
  }
`;

// Define the pricing calculator parameters
const migrationTypes = [
  {
    id: 'rehost',
    title: 'Rehost (Lift & Shift)',
    description: 'Move applications to the cloud with minimal changes',
    basePrice: 5000,
    pricePerServer: 500,
  },
  {
    id: 'replatform',
    title: 'Replatform (Lift & Reshape)',
    description: 'Optimize applications for cloud while retaining core architecture',
    basePrice: 10000,
    pricePerServer: 800,
  },
  {
    id: 'refactor',
    title: 'Refactor/Rearchitect',
    description: 'Modify applications to take full advantage of cloud-native features',
    basePrice: 15000,
    pricePerServer: 1200,
  },
];

const complexityLevels = [
  {
    id: 'simple',
    title: 'Simple',
    description: 'Basic applications with minimal dependencies',
    multiplier: 1,
  },
  {
    id: 'moderate',
    title: 'Moderate',
    description: 'Multiple integrations and medium complexity',
    multiplier: 1.5,
  },
  {
    id: 'complex',
    title: 'Complex',
    description: 'Legacy systems, many dependencies, strict requirements',
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

const PricingCalculator: React.FC = () => {
  const [migrationType, setMigrationType] = useState('rehost');
  const [complexity, setComplexity] = useState('moderate');
  const [serverCount, setServerCount] = useState(10);
  const [success, setSuccess] = useState(false);
  
  const selectedMigrationType = migrationTypes.find(type => type.id === migrationType)!;
  const selectedComplexity = complexityLevels.find(level => level.id === complexity)!;
  
  const basePrice = selectedMigrationType.basePrice;
  const serverPrice = selectedMigrationType.pricePerServer * serverCount;
  const complexityMultiplier = selectedComplexity.multiplier;
  
  const totalEstimate = Math.round((basePrice + serverPrice) * complexityMultiplier);
  
  const handleSubmit = () => {
    // Simulate form submission
    setSuccess(true);
    setTimeout(() => {
      setSuccess(false);
    }, 5000);
  };
  
  return (
    <CalculatorContainer
      initial="hidden"
      animate="visible"
      variants={calculatorVariants}
    >
      <CalculatorTitle>Estimate Your Migration Cost</CalculatorTitle>
      
      <FormGroup>
        <Label>
          Migration Type
          <InfoIcon 
            data-tooltip="Choose the migration approach that best fits your goals"
          >
            <Icon icon={FiInfo} size={16} />
          </InfoIcon>
        </Label>
        <OptionGrid>
          {migrationTypes.map(type => (
            <OptionCard 
              key={type.id} 
              selected={migrationType === type.id}
            >
              <RadioInput
                type="radio"
                name="migrationType"
                value={type.id}
                checked={migrationType === type.id}
                onChange={() => setMigrationType(type.id)}
              />
              <OptionTitle selected={migrationType === type.id}>
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
          Number of Servers/Applications
          <InfoIcon 
            data-tooltip="How many servers or applications will be migrated"
          >
            <Icon icon={FiInfo} size={16} />
          </InfoIcon>
        </Label>
        <RangeContainer>
          <RangeValue>{serverCount}</RangeValue>
          <RangeInput
            type="range"
            min="1"
            max="50"
            value={serverCount}
            onChange={(e) => setServerCount(parseInt(e.target.value))}
          />
          <RangeLabels>
            <span>1</span>
            <span>25</span>
            <span>50+</span>
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
          <EstimateLabel>Base Migration Cost:</EstimateLabel>
          <EstimateValue>${basePrice.toLocaleString()}</EstimateValue>
        </EstimateRow>
        <EstimateRow>
          <EstimateLabel>Server/Application Cost:</EstimateLabel>
          <EstimateValue>${serverPrice.toLocaleString()}</EstimateValue>
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
        <Button onClick={handleSubmit} size="lg">
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
          <div>Your request has been submitted. Our team will contact you shortly!</div>
        </SuccessMessage>
      )}
    </CalculatorContainer>
  );
};

export default PricingCalculator;