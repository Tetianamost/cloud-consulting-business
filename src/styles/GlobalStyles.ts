import { createGlobalStyle } from 'styled-components';
import { theme } from './theme';

const GlobalStyles = createGlobalStyle`
  /* Box sizing rules and reset */
  *,
  *::before,
  *::after {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
  }

  /* Remove default margin */
  body,
  h1,
  h2,
  h3,
  h4,
  h5,
  h6,
  p,
  figure,
  blockquote,
  dl,
  dd {
    margin: 0;
  }

  /* Set core body defaults */
  html {
    font-size: 16px;
    scroll-behavior: smooth;
  }

  body {
    min-height: 100vh;
    font-family: ${theme.fonts.primary};
    font-weight: ${theme.fontWeights.regular};
    font-size: ${theme.fontSizes.md};
    line-height: ${theme.lineHeights.normal};
    color: ${theme.colors.gray800};
    background-color: ${theme.colors.light};
    text-rendering: optimizeSpeed;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    overflow-x: hidden;
  }

  /* Better typography */
  h1, h2, h3, h4, h5, h6 {
    font-family: ${theme.fonts.heading};
    font-weight: ${theme.fontWeights.bold};
    line-height: ${theme.lineHeights.tight};
    color: ${theme.colors.primary};
    margin-bottom: ${theme.space[4]};
  }

  h1 {
    font-size: ${theme.fontSizes['5xl']};
    margin-bottom: ${theme.space[6]};
  }

  h2 {
    font-size: ${theme.fontSizes['4xl']};
  }

  h3 {
    font-size: ${theme.fontSizes['3xl']};
  }

  h4 {
    font-size: ${theme.fontSizes['2xl']};
  }

  h5 {
    font-size: ${theme.fontSizes.xl};
  }

  h6 {
    font-size: ${theme.fontSizes.lg};
  }

  p {
    margin-bottom: ${theme.space[4]};
  }

  a {
    color: ${theme.colors.accent};
    text-decoration: none;
    transition: ${theme.transitions.fast};
  }

  a:hover {
    color: ${theme.colors.highlight};
  }

  /* Remove list styles on ul, ol elements */
  ul,
  ol {
    list-style: none;
  }

  /* Remove all animations, transitions and smooth scroll for people that prefer not to see them */
  @media (prefers-reduced-motion: reduce) {
    html {
      scroll-behavior: auto;
    }
    
    *,
    *::before,
    *::after {
      animation-duration: 0.01ms !important;
      animation-iteration-count: 1 !important;
      transition-duration: 0.01ms !important;
      scroll-behavior: auto !important;
    }
  }

  /* Create a root stacking context */
  #root {
    isolation: isolate;
  }
  
  /* Improved responsiveness for images */
  img,
  picture,
  video {
    max-width: 100%;
    display: block;
  }
  
  /* Inherit fonts for inputs and buttons */
  input,
  button,
  textarea,
  select {
    font: inherit;
    color: inherit;
  }
  
  /* Container class */
  .container {
    width: 100%;
    max-width: ${theme.sizes.container.xl};
    margin-left: auto;
    margin-right: auto;
    padding-left: ${theme.space[4]};
    padding-right: ${theme.space[4]};
    
    @media (min-width: ${theme.breakpoints.lg}) {
      padding-left: ${theme.space[6]};
      padding-right: ${theme.space[6]};
    }
  }
  
  /* Section spacing */
  section {
    padding: ${theme.space[10]} 0;
    
    @media (min-width: ${theme.breakpoints.md}) {
      padding: ${theme.space[16]} 0;
    }
  }
`;

export default GlobalStyles;