import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { Link } from 'react-scroll';
import { theme } from '../../../styles/theme';
import { 
  FiMapPin, FiCloudSnow, FiDollarSign, FiShield, 
  FiDatabase, FiServer, FiCloudOff, FiCloud
} from 'react-icons/fi';
import type { IconType } from 'react-icons';
import type { ComponentType } from 'react';

/**
 * Enhanced Cloud Animation Component
 * 
 * Improvements:
 * 1. Repositioned icons to be visible on all screen sizes
 * 2. Updated icons to match service types
 * 3. Made all elements interactive with appropriate service links to SPECIFIC services
 * 4. Added responsive styling for mobile compatibility
 * 5. Enhanced visual feedback for better user experience
 * 6. Fixed TypeScript errors with icon component types
 */

// Main container with responsive height - optimized for better layout and mobile positioning
const AnimationContainer = styled.div`
  position: relative;
  width: 100%;
  height: 300px;
  max-width: 480px;
  margin: 0 auto;
  overflow: visible;
  
  @media (max-width: ${theme.breakpoints.md}) {
    height: 240px;
    max-height: 30vh; /* Limit height on mobile to prevent overlapping with text */
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    height: 220px;
    max-height: 28vh;
    margin-top: -15px; /* Push up slightly to make more room for text content */
  }
  
  /* Extra small screens get a more compact layout */
  @media (max-width: 375px) {
    height: 200px;
    max-height: 25vh;
    margin-top: -25px;
  }
`;

// Tooltip component for interactive elements
const Tooltip = styled.span`
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
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 20;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  
  &::after {
    content: '';
    position: absolute;
    top: 100%;
    left: 50%;
    margin-left: -5px;
    border-width: 5px;
    border-style: solid;
    border-color: ${theme.colors.primary} transparent transparent transparent;
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    font-size: ${theme.fontSizes.xs};
    padding: ${theme.space[1]} ${theme.space[2]};
    top: -30px;
  }
`;

// Cloud elements with hover effects
const Cloud = styled(motion.div)`
  position: absolute;
  background-color: #fff;
  border-radius: 16px;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  z-index: 3;
  
  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
    
    ${Tooltip} {
      opacity: 1;
    }
  }
`;

// DB Cloud styling
const DbCloud = styled(Cloud)`
  width: 28%;
  height: 60px;
  top: 15%;
  left: 20%;
  
  @media (max-width: ${theme.breakpoints.md}) {
    height: 55px;
    top: 18%;
    left: 15%;
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    height: 50px;
    top: 15%;
    left: 10%;
  }
  
  @media (max-width: 375px) {
    height: 45px;
    top: 12%;
    left: 8%;
  }
`;

// S3 Cloud styling
const S3Cloud = styled(Cloud)`
  width: 22%;
  height: 55px;
  top: 65%;
  right: 15%;
  
  @media (max-width: ${theme.breakpoints.md}) {
    height: 50px;
    top: 68%;
    right: 12%;
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    height: 45px;
    top: 70%;
    right: 10%;
  }
  
  @media (max-width: 375px) {
    height: 40px;
    top: 68%;
    right: 8%;
    /* Ensure S3 doesn't overlap with text on small screens */
    z-index: 2;
  }
`;

// Logo styling for cloud elements
const CloudLogo = styled.div`
  color: ${theme.colors.primary};
  font-size: 24px;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  
  @media (max-width: ${theme.breakpoints.sm}) {
    font-size: 20px;
  }
  
  @media (max-width: 375px) {
    font-size: 18px;
  }
`;

// New central lines container
const FlowingLines = styled(motion.div)`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
`;

// Network nodes styling
const Node = styled(motion.div)`
  width: 46px;
  height: 46px;
  background: linear-gradient(145deg, #ffffff, #f5f5f5);
  border-radius: 10px;
  position: absolute;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 2;
  
  &:hover {
    transform: translateY(-5px) !important;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
    
    ${Tooltip} {
      opacity: 1;
    }
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    width: 40px;
    height: 40px;
  }
  
  @media (max-width: 375px) {
    width: 36px;
    height: 36px;
  }
`;

// Status indicator
const StatusIndicator = styled(motion.div)`
  position: absolute;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: ${props => props.color || theme.colors.success};
  top: 6px;
  right: 6px;
  
  @media (max-width: ${theme.breakpoints.sm}) {
    width: 6px;
    height: 6px;
    top: 5px;
    right: 5px;
  }
`;

// Connection paths between elements
const Connection = styled(motion.path)`
  stroke: ${theme.colors.secondary};
  stroke-width: 2;
  stroke-dasharray: 5;
  fill: none;
  z-index: 1;
`;

// Main connection (central network)
const MainConnection = styled(Connection)`
  stroke-width: 2.5;
  opacity: 0.8;
`;

// Data packet animations
const DataPacket = styled(motion.circle)`
  fill: ${theme.colors.accent};
  r: 5;
  filter: drop-shadow(0 0 2px ${theme.colors.accent});
  
  @media (max-width: ${theme.breakpoints.sm}) {
    r: 4;
  }
  
  @media (max-width: 375px) {
    r: 3;
  }
`;

// Service icon styling
const ServiceIcon = styled(motion.div)`
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border-radius: 50%;
  background-color: #fff;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  color: ${theme.colors.primary};
  font-size: 17px;
  z-index: 3;
  cursor: pointer;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  
  &:hover {
    transform: translateY(-5px) scale(1.1) !important;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
    
    ${Tooltip} {
      opacity: 1;
    }
  }
  
  @media (max-width: ${theme.breakpoints.sm}) {
    width: 34px;
    height: 34px;
    font-size: 15px;
  }
  
  @media (max-width: 375px) {
    width: 30px;
    height: 30px;
    font-size: 13px;
  }
`;

// Interface for service data
interface ServiceData {
  id: number;
  name: string;
  element: string;
  tooltip: string;
}

// Interface for node data with proper typing for IconType
interface NodeInfo {
  top: string;
  left: string;
  delay: number;
  serviceId: string;
  tooltip: string;
  iconComponent: ComponentType<any>;
}

// Interface for service icon data
interface IconInfo {
  top: string;
  left: string;
  serviceId: string;
  tooltip: string;
  iconComponent: ComponentType<any>;
  animate: object;
}

// Define services based on your provided list
const services: ServiceData[] = [
  {
    id: 1,
    name: "Cloud Assessment & Roadmap",
    element: "service-1",
    tooltip: "Cloud Assessment & Roadmap"
  },
  {
    id: 2,
    name: "Small-Scale Migrations",
    element: "service-2",
    tooltip: "Small-Scale Migrations"
  },
  {
    id: 3,
    name: "Cloud Optimization Consulting",
    element: "service-3",
    tooltip: "Cloud Optimization"
  },
  {
    id: 4,
    name: "Cloud Architecture Review",
    element: "service-4",
    tooltip: "Architecture Review"
  }
];

// Network nodes with links to services - arranged around the central connections
const nodes: NodeInfo[] = [
  { 
    top: '28%',
    left: '35%',
    delay: 0,
    serviceId: services[0].element,  // Cloud Assessment
    tooltip: 'Cloud Assessment',
    iconComponent: FiCloudOff as ComponentType<any>
  },
  { 
    top: '25%',
    left: '65%',
    delay: 0.2,
    serviceId: services[1].element,  // Small-Scale Migrations
    tooltip: 'Cloud Migration',
    iconComponent: FiCloud as ComponentType<any>
  },
  { 
    top: '60%',
    left: '48%',
    delay: 0.1,
    serviceId: services[2].element,  // Cloud Optimization
    tooltip: 'Cost Optimization',
    iconComponent: FiServer as ComponentType<any>
  },
];

// Service icons positioned around the central flowing lines
const serviceIcons: IconInfo[] = [
  { 
    top: '20%',
    left: '50%',
    serviceId: services[0].element,  // Cloud Assessment
    tooltip: 'Assessment & Roadmap',
    iconComponent: FiMapPin as ComponentType<any>,
    animate: { 
      rotate: [0, 8, 0, -8, 0], 
      transition: { repeat: Infinity, duration: 4, ease: "easeInOut" } 
    } 
  },
  { 
    top: '40%',
    left: '75%',
    serviceId: services[3].element,  // Cloud Architecture Review
    tooltip: 'Architecture & Security',
    iconComponent: FiShield as ComponentType<any>,
    animate: {
      scale: [1, 1.15, 1],
      transition: { repeat: Infinity, duration: 2, ease: "easeInOut" }
    }
  },
  { 
    top: '70%',
    left: '25%',
    serviceId: services[2].element,  // Cloud Optimization
    tooltip: 'Cost Optimization',
    iconComponent: FiDollarSign as ComponentType<any>,
    animate: {
      y: [0, -5, 0],
      transition: { repeat: Infinity, duration: 3, ease: "easeInOut" }
    }
  },
  { 
    top: '45%',
    left: '15%',
    serviceId: services[1].element,  // Small-Scale Migrations
    tooltip: 'Database Migration',
    iconComponent: FiDatabase as ComponentType<any>,
    animate: {
      scale: [1, 0.85, 1],
      transition: { repeat: Infinity, duration: 2.5, ease: "easeInOut" }
    }
  }
];

// Central connection paths for the flowing network effect
const centralPaths = [
  {
    path: "M 150,140 Q 250,100 350,140",
    direction: -1,
    duration: 4,
    delay: 0
  },
  {
    path: "M 120,160 Q 250,190 380,160",
    direction: 1,
    duration: 5,
    delay: 0.5
  },
  {
    path: "M 140,200 Q 250,150 360,200",
    direction: -1,
    duration: 3.5,
    delay: 1
  },
  {
    path: "M 130,120 Q 240,180 350,120",
    direction: 1,
    duration: 6,
    delay: 0.2
  },
  {
    path: "M 140,220 Q 240,240 340,220",
    direction: -1,
    duration: 4.5,
    delay: 0.7
  }
];

const CloudAnimation: React.FC = () => {
  return (
    <AnimationContainer>
      {/* Central network flowing lines - replaces the previous central cloud element */}
      <FlowingLines>
        <motion.svg
          width="100%"
          height="100%"
          viewBox="0 0 480 300"
          style={{ position: 'absolute', top: 0, left: 0 }}
        >
          {/* Multiple flowing connection paths in the center */}
          {centralPaths.map((path, index) => (
            <MainConnection
              key={`path-${index}`}
              d={path.path}
              animate={{
                strokeDashoffset: [0, path.direction * 80],
                transition: { 
                  repeat: Infinity, 
                  duration: path.duration, 
                  ease: "linear",
                  delay: path.delay
                }
              }}
            />
          ))}
          
          {/* Data packets flowing along the connections */}
          <DataPacket
            animate={{
              cx: [150, 250, 350],
              cy: [140, 100, 140],
              transition: { repeat: Infinity, duration: 4, ease: "easeInOut" }
            }}
          />
          
          <DataPacket
            animate={{
              cx: [380, 250, 120],
              cy: [160, 190, 160],
              transition: { repeat: Infinity, duration: 5, ease: "easeInOut", delay: 0.5 }
            }}
          />
          
          <DataPacket
            animate={{
              cx: [140, 250, 360],
              cy: [200, 150, 200],
              transition: { repeat: Infinity, duration: 3.5, ease: "easeInOut", delay: 1 }
            }}
          />
          
          <DataPacket
            animate={{
              cx: [350, 240, 130],
              cy: [120, 180, 120],
              transition: { repeat: Infinity, duration: 6, ease: "easeInOut", delay: 0.2 }
            }}
          />
          
          <DataPacket
            animate={{
              cx: [140, 240, 340],
              cy: [220, 240, 220],
              transition: { repeat: Infinity, duration: 4.5, ease: "easeInOut", delay: 0.7 }
            }}
          />
        </motion.svg>
      </FlowingLines>
      
      {/* DB Cloud element - link to Small-Scale Migrations */}
      <Link 
        to={services[1].element}  // Link to Small-Scale Migrations
        smooth={true} 
        offset={-70} 
        duration={500} 
        spy={true}
      >
        <DbCloud
          animate={{ 
            y: [0, -8, 0],
            transition: { 
              repeat: Infinity, 
              duration: 4,
              ease: "easeInOut"
            }
          }}
          aria-label="Go to Database Migration Service"
        >
          <Tooltip>Database Migrations</Tooltip>
          <CloudLogo>DB</CloudLogo>
        </DbCloud>
      </Link>
      
      {/* S3 Cloud element - link to Small-Scale Migrations */}
      <Link 
        to={services[1].element}  // Link to Small-Scale Migrations
        smooth={true} 
        offset={-70} 
        duration={500}
        spy={true}
      >
        <S3Cloud
          animate={{ 
            y: [0, -6, 0],
            transition: { 
              repeat: Infinity, 
              duration: 5,
              ease: "easeInOut"
            }
          }}
          aria-label="Go to Storage Migration Service"
        >
          <Tooltip>Storage Migrations</Tooltip>
          <CloudLogo>S3</CloudLogo>
        </S3Cloud>
      </Link>
      
      {/* Network nodes placed around the central connections */}
      {nodes.map((node, index) => (
        <Link 
          key={`node-${index}`}
          to={node.serviceId}  // Link to specific service
          smooth={true} 
          offset={-70} 
          duration={500}
          spy={true}
        >
          <Node
            style={{ top: node.top, left: node.left }}
            initial={{ opacity: 0, y: 15 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: node.delay + 0.3 }}
            aria-label={`Go to ${node.tooltip}`}
            whileHover={{ y: -5, boxShadow: '0 8px 16px rgba(0, 0, 0, 0.15)' }}
          >
            <Tooltip>{node.tooltip}</Tooltip>
            <StatusIndicator 
              color={theme.colors.success}
              animate={{ 
                opacity: [1, 0.5, 1],
                transition: { repeat: Infinity, duration: 1.5, ease: "easeInOut" }
              }}
            />
            <div style={{ color: theme.colors.primary, fontSize: '18px' }}>
              {/* Use JSX element syntax with properly typed component */}
              <node.iconComponent />
            </div>
          </Node>
        </Link>
      ))}
      
      {/* Service icons arranged around the central connections */}
      {serviceIcons.map((icon, index) => (
        <Link 
          key={`icon-${index}`}
          to={icon.serviceId}  // Link to specific service
          smooth={true} 
          offset={-70} 
          duration={500}
          spy={true}
        >
          <ServiceIcon
            style={{ top: icon.top, left: icon.left }}
            initial={{ opacity: 0, scale: 0 }}
            animate={{ 
              opacity: 1, 
              scale: 1,
              ...icon.animate
            }}
            transition={{ duration: 0.5, delay: index * 0.15 + 0.8 }}
            aria-label={`Go to ${icon.tooltip}`}
            whileHover={{ y: -5, scale: 1.1, boxShadow: '0 8px 16px rgba(0, 0, 0, 0.15)' }}
          >
            <Tooltip>{icon.tooltip}</Tooltip>
            {/* Use JSX element syntax with properly typed component */}
            <icon.iconComponent />
          </ServiceIcon>
        </Link>
      ))}
    </AnimationContainer>
  );
};

export default CloudAnimation;