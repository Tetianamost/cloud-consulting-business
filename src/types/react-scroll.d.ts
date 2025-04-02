declare module 'react-scroll' {
    import * as React from 'react';
  
    export interface LinkProps extends React.HTMLProps<HTMLAnchorElement> {
      to: string;
      containerId?: string;
      activeClass?: string;
      spy?: boolean;
      smooth?: boolean | string;
      offset?: number;
      duration?: number;
      delay?: number;
      isDynamic?: boolean;
      onSetActive?: (to: string) => void;
      onSetInactive?: (to: string) => void;
      ignoreCancelEvents?: boolean;
    }
  
    export class Link extends React.Component<LinkProps> {}
    
    export interface ElementProps {
      name: string;
      id?: string;
      className?: string;
    }
    
    export class Element extends React.Component<ElementProps> {}
    
    export interface ScrollerProps {
      scrollTo: (target: string, options?: {
        duration?: number;
        delay?: number;
        smooth?: boolean | string;
        offset?: number;
        containerId?: string;
      }) => void;
    }
    
    export const scroller: {
      scrollTo(target: string, options?: {
        duration?: number;
        delay?: number;
        smooth?: boolean | string;
        offset?: number;
        containerId?: string;
      }): void;
    };
    
    export const scrollSpy: {
      update(): void;
    };
    
    export const Events: {
      scrollEvent: {
        register(eventName: string, callback: (event: any) => void): void;
        remove(eventName: string): void;
      };
    };
  }