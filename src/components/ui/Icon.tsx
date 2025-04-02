import React from 'react';
import { IconType } from 'react-icons';

interface IconProps {
  icon: IconType;
  size?: number | string;
  color?: string;
  className?: string;
}

// Use type assertion to make TypeScript happy
const Icon: React.FC<IconProps> = ({ icon, ...props }) => {
  // Force the type to be a valid React component type
  const IconComponent = icon as React.ComponentType<any>;
  return <IconComponent {...props} />;
};

export default Icon;