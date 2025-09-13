#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

// Change to the directory where this script is located
process.chdir(__dirname);

console.log('🔍 Verifying frontend build process...');

// Check if required files exist
const requiredFiles = [
  'package.json',
  'tsconfig.json',
  'public/index.html',
  'src/index.tsx',
  'src/App.tsx'
];

console.log('📁 Checking required files...');
for (const file of requiredFiles) {
  if (!fs.existsSync(file)) {
    console.error(`❌ Missing required file: ${file}`);
    process.exit(1);
  }
  console.log(`✅ Found: ${file}`);
}

// Check if dependencies are installed
console.log('📦 Checking dependencies...');
if (!fs.existsSync('node_modules')) {
  console.error('❌ node_modules directory not found. Run npm install first.');
  process.exit(1);
}

if (!fs.existsSync('node_modules/.bin/react-scripts')) {
  console.error('❌ react-scripts not found. Run npm install first.');
  process.exit(1);
}
console.log('✅ Dependencies are installed');

// Check TypeScript compilation
console.log('🔧 Checking TypeScript compilation...');
try {
  execSync('node_modules/.bin/tsc --noEmit', { stdio: 'pipe' });
  console.log('✅ TypeScript compilation successful');
} catch (error) {
  console.error('❌ TypeScript compilation failed');
  console.error(error.stdout?.toString() || error.message);
  process.exit(1);
}

// Verify package.json scripts
console.log('📜 Checking package.json scripts...');
const packageJson = JSON.parse(fs.readFileSync('package.json', 'utf8'));
const requiredScripts = ['start', 'build', 'test'];

for (const script of requiredScripts) {
  if (!packageJson.scripts || !packageJson.scripts[script]) {
    console.error(`❌ Missing script: ${script}`);
    process.exit(1);
  }
  console.log(`✅ Found script: ${script}`);
}

console.log('🎉 Frontend build verification completed successfully!');
console.log('');
console.log('📋 Summary:');
console.log('  ✅ All required files present');
console.log('  ✅ Dependencies installed');
console.log('  ✅ TypeScript compilation works');
console.log('  ✅ Package.json scripts configured');
console.log('');
console.log('🚀 The frontend is ready for development and production builds.');