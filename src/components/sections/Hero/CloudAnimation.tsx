import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { theme } from '../../../styles/theme';

const AnimationContainer = styled.div`
  position: relative;
  width: 100%;
  height: 400px;
  max-width: 500px;
  margin: 0 auto;
`;

const Cloud = styled(motion.div)`
  position: absolute;
  background-color: #fff;
  border-radius: 20px;
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
`;

const MainCloud = styled(Cloud)`
  width: 70%;
  height: 180px;
  top: 30%;
  left: 15%;
  z-index: 2;
`;

const SecondaryCloud1 = styled(Cloud)`
  width: 40%;
  height: 100px;
  top: 15%;
  left: 5%;
  z-index: 1;
`;

const SecondaryCloud2 = styled(Cloud)`
  width: 30%;
  height: 80px;
  bottom: 20%;
  right: 10%;
  z-index: 1;
`;

const CloudLogo = styled.div`
  color: ${theme.colors.primary};
  font-size: 40px;
  font-weight: bold;
  display: flex;
  align-items: center;
  
  span {
    color: ${theme.colors.secondary};
  }
`;

const Server = styled(motion.div)`
  width: 60px;
  height: 90px;
  background: linear-gradient(145deg, #e6e6e6, #ffffff);
  border-radius: 5px;
  position: absolute;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 10px;
`;

const ServerLight = styled(motion.div)`
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background-color: ${props => props.color || theme.colors.success};
  margin-bottom: 5px;
`;

const ServerLine = styled.div`
  width: 40px;
  height: 2px;
  background-color: ${theme.colors.gray300};
  margin: 5px 0;
`;

const Connection = styled(motion.path)`
  stroke: ${theme.colors.secondary};
  stroke-width: 2;
  stroke-dasharray: 5;
  fill: none;
`;

const DataPacket = styled(motion.circle)`
  fill: ${theme.colors.accent};
  r: 6;
`;

const CloudIcon = styled(motion.div)`
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background-color: #fff;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  color: ${theme.colors.primary};
  font-size: 20px;
  z-index: 10;
`;

// Server positions
const servers = [
  { top: '60%', left: '20%', delay: 0 },
  { top: '50%', left: '60%', delay: 0.1 },
  { top: '70%', left: '40%', delay: 0.2 },
];

// Cloud icons
const cloudIcons = [
  { 
    top: '25%', 
    left: '70%', 
    content: '⚙️', 
    animate: { 
      rotate: [0, 360], 
      transition: { repeat: Infinity, duration: 10, ease: "linear" } 
    } 
  },
  { 
    top: '60%', 
    left: '75%', 
    content: '🔒',
    animate: {
      scale: [1, 1.1, 1],
      transition: { repeat: Infinity, duration: 2, ease: "easeInOut" }
    }
  },
  { 
    top: '35%', 
    left: '25%', 
    content: '📊',
    animate: {
      y: [0, -5, 0],
      transition: { repeat: Infinity, duration: 3, ease: "easeInOut" }
    }
  },
];

const CloudAnimation: React.FC = () => {
  return (
    <AnimationContainer>
      <motion.svg
        width="100%"
        height="100%"
        viewBox="0 0 500 400"
        style={{ position: 'absolute', top: 0, left: 0 }}
      >
        {/* Animated connection paths */}
        <Connection
          d="M 150,200 Q 250,150 350,200"
          animate={{
            strokeDashoffset: [0, -100],
            transition: { repeat: Infinity, duration: 3, ease: "linear" }
          }}
        />
        
        <Connection
          d="M 120,280 Q 250,320 380,280"
          animate={{
            strokeDashoffset: [0, 100],
            transition: { repeat: Infinity, duration: 4, ease: "linear" }
          }}
        />
        
        {/* Data packet animations */}
        <DataPacket
          animate={{
            cx: [150, 250, 350],
            cy: [200, 150, 200],
            transition: { repeat: Infinity, duration: 3, ease: "easeInOut" }
          }}
        />
        
        <DataPacket
          animate={{
            cx: [380, 250, 120],
            cy: [280, 320, 280],
            transition: { repeat: Infinity, duration: 4, ease: "easeInOut", delay: 0.5 }
          }}
        />
      </motion.svg>
      
      <SecondaryCloud1
        animate={{ 
          y: [0, -10, 0],
          transition: { 
            repeat: Infinity, 
            duration: 4,
            ease: "easeInOut"
          }
        }}
      >
        <CloudLogo style={{ fontSize: '20px' }}>DB</CloudLogo>
      </SecondaryCloud1>
      
      <SecondaryCloud2
        animate={{ 
          y: [0, -8, 0],
          transition: { 
            repeat: Infinity, 
            duration: 5,
            ease: "easeInOut"
          }
        }}
      >
        <CloudLogo style={{ fontSize: '18px' }}>S3</CloudLogo>
      </SecondaryCloud2>
      
      <MainCloud
        animate={{ 
          y: [0, -15, 0],
          transition: { 
            repeat: Infinity, 
            duration: 6,
            ease: "easeInOut"
          }
        }}
      >
        <CloudLogo>
          AWS<span>Cloud</span>
        </CloudLogo>
      </MainCloud>
      
      {servers.map((server, index) => (
        <Server
          key={index}
          style={{ top: server.top, left: server.left }}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: server.delay + 0.5 }}
        >
          <ServerLight 
            color={theme.colors.success}
            animate={{ 
              opacity: [1, 0.5, 1],
              transition: { repeat: Infinity, duration: 1.5, ease: "easeInOut" }
            }}
          />
          <ServerLine />
          <ServerLine />
          <ServerLine />
        </Server>
      ))}
      
      {cloudIcons.map((icon, index) => (
        <CloudIcon
          key={index}
          style={{ top: icon.top, left: icon.left }}
          initial={{ opacity: 0, scale: 0 }}
          animate={{ 
            opacity: 1, 
            scale: 1,
            ...icon.animate
          }}
          transition={{ duration: 0.5, delay: index * 0.2 + 1 }}
        >
          {icon.content}
        </CloudIcon>
      ))}
    </AnimationContainer>
  );
};

export default CloudAnimation;