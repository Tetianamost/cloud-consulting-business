#!/usr/bin/env node

/**
 * Backward Compatibility Test Suite
 * Tests that the dual styling system (styled-components + Tailwind) works correctly
 * and that public site components remain unaffected by admin dashboard changes
 */

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

console.log('🧪 Running Backward Compatibility Tests...\n');

// Test 1: Verify public site components still use styled-components
console.log('1️⃣ Testing public site styled-components usage...');

const publicComponents = [
  'src/components/layout/Header.tsx',
  'src/components/layout/Footer.tsx',
  'src/components/sections/Hero/Hero.tsx',
  'src/components/sections/Services/Services.tsx',
  'src/components/sections/Contact/Contact.tsx'
];

let styledComponentsFound = 0;
let tailwindInPublic = 0;

publicComponents.forEach(componentPath => {
  const fullPath = path.join(__dirname, componentPath);
  if (fs.existsSync(fullPath)) {
    const content = fs.readFileSync(fullPath, 'utf8');
    
    // Check for styled-components usage
    if (content.includes('styled.') || content.includes('styled(') || content.includes('from \'styled-components\'')) {
      styledComponentsFound++;
      console.log(`   ✅ ${componentPath} uses styled-components`);
    } else {
      console.log(`   ⚠️  ${componentPath} may not use styled-components`);
    }
    
    // Check for Tailwind classes (should be minimal in public components)
    const tailwindMatches = content.match(/className="[^"]*(?:bg-|text-|p-|m-|flex|grid|w-|h-)/g);
    if (tailwindMatches && tailwindMatches.length > 5) {
      tailwindInPublic++;
      console.log(`   ⚠️  ${componentPath} has significant Tailwind usage: ${tailwindMatches.length} classes`);
    }
  }
});

console.log(`   📊 Styled-components found in ${styledComponentsFound}/${publicComponents.length} public components`);
console.log(`   📊 Components with heavy Tailwind usage: ${tailwindInPublic}\n`);

// Test 2: Verify admin components use Tailwind
console.log('2️⃣ Testing admin components Tailwind usage...');

const adminComponents = [
  'src/components/admin/V0AdminLayout.tsx',
  'src/components/admin/V0Sidebar.tsx',
  'src/components/admin/V0MetricsCards.tsx',
  'src/components/admin/V0DashboardNew.tsx'
];

let tailwindInAdmin = 0;
let styledComponentsInAdmin = 0;

adminComponents.forEach(componentPath => {
  const fullPath = path.join(__dirname, componentPath);
  if (fs.existsSync(fullPath)) {
    const content = fs.readFileSync(fullPath, 'utf8');
    
    // Check for Tailwind usage
    const tailwindMatches = content.match(/className="[^"]*(?:bg-|text-|p-|m-|flex|grid|w-|h-)/g);
    if (tailwindMatches && tailwindMatches.length > 10) {
      tailwindInAdmin++;
      console.log(`   ✅ ${componentPath} uses Tailwind extensively: ${tailwindMatches.length} classes`);
    } else {
      console.log(`   ⚠️  ${componentPath} has limited Tailwind usage`);
    }
    
    // Check for styled-components (should be minimal in admin components)
    if (content.includes('styled.') || content.includes('styled(')) {
      styledComponentsInAdmin++;
      console.log(`   ⚠️  ${componentPath} still uses styled-components`);
    }
  }
});

console.log(`   📊 Admin components with extensive Tailwind: ${tailwindInAdmin}/${adminComponents.length}`);
console.log(`   📊 Admin components with styled-components: ${styledComponentsInAdmin}\n`);

// Test 3: Check CSS file structure
console.log('3️⃣ Testing CSS file structure...');

const cssFiles = [
  'src/styles/admin.css',
  'src/styles/admin-tailwind.css',
  'src/index.css'
];

cssFiles.forEach(cssPath => {
  const fullPath = path.join(__dirname, cssPath);
  if (fs.existsSync(fullPath)) {
    const content = fs.readFileSync(fullPath, 'utf8');
    const size = content.length;
    console.log(`   ✅ ${cssPath} exists (${size} bytes)`);
    
    // Check for Tailwind directives
    if (content.includes('@tailwind')) {
      console.log(`   📝 ${cssPath} contains Tailwind directives`);
    }
  } else {
    console.log(`   ❌ ${cssPath} not found`);
  }
});

// Test 4: Build process verification
console.log('\n4️⃣ Testing build process...');

try {
  console.log('   🔨 Running production build...');
  const buildOutput = execSync('npm run build', { 
    cwd: __dirname,
    encoding: 'utf8',
    timeout: 120000 // 2 minutes timeout
  });
  
  console.log('   ✅ Build completed successfully');
  
  // Check build output for bundle sizes
  const bundleInfo = buildOutput.match(/(\d+(?:\.\d+)?\s*(?:kB|MB))/g);
  if (bundleInfo) {
    console.log(`   📦 Bundle sizes: ${bundleInfo.slice(0, 3).join(', ')}`);
  }
  
  // Check if build folder exists
  const buildPath = path.join(__dirname, 'build');
  if (fs.existsSync(buildPath)) {
    console.log('   ✅ Build folder created successfully');
    
    // Check for CSS files in build
    const buildCssPath = path.join(buildPath, 'static/css');
    if (fs.existsSync(buildCssPath)) {
      const cssFiles = fs.readdirSync(buildCssPath);
      console.log(`   📄 CSS files in build: ${cssFiles.length}`);
    }
  }
  
} catch (error) {
  console.log('   ❌ Build failed:', error.message);
}

// Test 5: Check for potential CSS conflicts
console.log('\n5️⃣ Checking for potential CSS conflicts...');

const appTsxPath = path.join(__dirname, 'src/App.tsx');
if (fs.existsSync(appTsxPath)) {
  const appContent = fs.readFileSync(appTsxPath, 'utf8');
  
  // Check import order
  const styledImportIndex = appContent.indexOf('styled-components');
  const tailwindImportIndex = appContent.indexOf('admin.css');
  
  if (styledImportIndex !== -1 && tailwindImportIndex !== -1) {
    console.log('   ✅ Both styling systems imported in App.tsx');
    if (tailwindImportIndex > styledImportIndex) {
      console.log('   ✅ Tailwind CSS imported after styled-components (good for specificity)');
    } else {
      console.log('   ⚠️  Tailwind CSS imported before styled-components');
    }
  }
  
  // Check for route separation
  if (appContent.includes('AdminLayoutWrapper') && appContent.includes('MainSite')) {
    console.log('   ✅ Admin and public routes properly separated');
  }
}

// Test 6: Performance impact check
console.log('\n6️⃣ Checking performance impact...');

const packageJsonPath = path.join(__dirname, 'package.json');
if (fs.existsSync(packageJsonPath)) {
  const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
  
  // Check dependencies
  const styledComponentsVersion = packageJson.dependencies?.['styled-components'];
  const tailwindVersion = packageJson.devDependencies?.['tailwindcss'];
  
  if (styledComponentsVersion && tailwindVersion) {
    console.log(`   📦 styled-components: ${styledComponentsVersion}`);
    console.log(`   📦 tailwindcss: ${tailwindVersion}`);
    console.log('   ✅ Both styling systems present in dependencies');
  }
}

// Summary
console.log('\n📋 Test Summary:');
console.log('================');
console.log(`✅ Public components using styled-components: ${styledComponentsFound}/${publicComponents.length}`);
console.log(`✅ Admin components using Tailwind: ${tailwindInAdmin}/${adminComponents.length}`);
console.log(`⚠️  Admin components with styled-components: ${styledComponentsInAdmin}`);
console.log(`⚠️  Public components with heavy Tailwind: ${tailwindInPublic}`);

if (styledComponentsFound >= 3 && tailwindInAdmin >= 3 && styledComponentsInAdmin <= 1 && tailwindInPublic === 0) {
  console.log('\n🎉 Backward compatibility tests PASSED!');
  console.log('   - Public site maintains styled-components');
  console.log('   - Admin dashboard uses Tailwind CSS');
  console.log('   - Minimal cross-contamination between systems');
  process.exit(0);
} else {
  console.log('\n⚠️  Some backward compatibility issues detected');
  console.log('   - Review the warnings above');
  console.log('   - Consider refactoring components with mixed styling');
  process.exit(1);
}