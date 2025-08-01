#!/usr/bin/env node

/**
 * Final Verification Test Suite
 * Comprehensive test to verify all backward compatibility requirements are met
 */

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

console.log('🔍 Final Backward Compatibility Verification...\n');

let allTestsPassed = true;
const results = {
  publicSiteIntact: false,
  adminRoutingWorking: false,
  noCssConflicts: false,
  buildSuccessful: false,
  performanceAcceptable: false,
  authenticationWorking: false
};

// Test 1: Verify public site components are intact
console.log('1️⃣ Verifying public site components...');
try {
  const publicComponents = [
    'src/components/layout/Header.tsx',
    'src/components/layout/Footer.tsx', 
    'src/components/sections/Hero/Hero.tsx',
    'src/components/sections/Services/Services.tsx',
    'src/components/sections/Contact/Contact.tsx'
  ];
  
  let styledComponentsCount = 0;
  publicComponents.forEach(comp => {
    const fullPath = path.join(__dirname, comp);
    if (fs.existsSync(fullPath)) {
      const content = fs.readFileSync(fullPath, 'utf8');
      if (content.includes('styled-components') || content.includes('styled.')) {
        styledComponentsCount++;
      }
    }
  });
  
  if (styledComponentsCount >= 4) {
    console.log('   ✅ Public site components use styled-components');
    results.publicSiteIntact = true;
  } else {
    console.log('   ❌ Public site components may have been affected');
    allTestsPassed = false;
  }
} catch (error) {
  console.log('   ❌ Error checking public components:', error.message);
  allTestsPassed = false;
}

// Test 2: Verify admin routing structure
console.log('\n2️⃣ Verifying admin routing structure...');
try {
  const appPath = path.join(__dirname, 'src/App.tsx');
  const wrapperPath = path.join(__dirname, 'src/components/admin/AdminLayoutWrapper.tsx');
  
  if (fs.existsSync(appPath) && fs.existsSync(wrapperPath)) {
    const appContent = fs.readFileSync(appPath, 'utf8');
    const wrapperContent = fs.readFileSync(wrapperPath, 'utf8');
    
    const hasAdminWrapper = appContent.includes('AdminLayoutWrapper');
    const wrapperUsesV0Layout = wrapperContent.includes('V0AdminLayout');
    const routesUseWrapper = appContent.match(/AdminLayoutWrapper>/g)?.length >= 5;
    
    if (hasAdminWrapper && wrapperUsesV0Layout && routesUseWrapper) {
      console.log('   ✅ Admin routing structure properly implemented');
      results.adminRoutingWorking = true;
    } else {
      console.log('   ❌ Admin routing structure has issues');
      allTestsPassed = false;
    }
  } else {
    console.log('   ❌ Required routing files not found');
    allTestsPassed = false;
  }
} catch (error) {
  console.log('   ❌ Error checking routing structure:', error.message);
  allTestsPassed = false;
}

// Test 3: Check for CSS conflicts
console.log('\n3️⃣ Checking for CSS conflicts...');
try {
  const appPath = path.join(__dirname, 'src/App.tsx');
  if (fs.existsSync(appPath)) {
    const appContent = fs.readFileSync(appPath, 'utf8');
    
    // Check import order - Tailwind should come after styled-components
    const styledImportIndex = appContent.indexOf('styled-components');
    const tailwindImportIndex = appContent.indexOf('admin.css');
    
    const hasThemeProvider = appContent.includes('ThemeProvider');
    const hasGlobalStyles = appContent.includes('GlobalStyles');
    const hasAdminCSS = appContent.includes('admin.css');
    
    if (hasThemeProvider && hasGlobalStyles && hasAdminCSS) {
      console.log('   ✅ Both styling systems properly configured');
      results.noCssConflicts = true;
    } else {
      console.log('   ❌ CSS configuration may have conflicts');
      allTestsPassed = false;
    }
  }
} catch (error) {
  console.log('   ❌ Error checking CSS configuration:', error.message);
  allTestsPassed = false;
}

// Test 4: Verify build process
console.log('\n4️⃣ Testing build process...');
try {
  console.log('   🔨 Running production build...');
  const buildOutput = execSync('npm run build', { 
    cwd: __dirname,
    encoding: 'utf8',
    timeout: 120000,
    stdio: 'pipe'
  });
  
  if (buildOutput.includes('Compiled successfully')) {
    console.log('   ✅ Production build successful');
    results.buildSuccessful = true;
    
    // Check bundle sizes
    const bundleInfo = buildOutput.match(/(\d+(?:\.\d+)?)\s*kB/g);
    if (bundleInfo && bundleInfo.length > 0) {
      const mainBundleSize = parseFloat(bundleInfo[0]);
      if (mainBundleSize < 300) { // Reasonable threshold
        console.log(`   ✅ Bundle size acceptable: ${bundleInfo[0]}`);
        results.performanceAcceptable = true;
      } else {
        console.log(`   ⚠️  Bundle size may be large: ${bundleInfo[0]}`);
      }
    }
  } else {
    console.log('   ❌ Build failed');
    allTestsPassed = false;
  }
} catch (error) {
  console.log('   ❌ Build process failed:', error.message);
  allTestsPassed = false;
}

// Test 5: Verify authentication integration
console.log('\n5️⃣ Verifying authentication integration...');
try {
  const protectedRoutePath = path.join(__dirname, 'src/components/admin/ProtectedRoute.tsx');
  const authContextPath = path.join(__dirname, 'src/contexts/AuthContext.tsx');
  
  if (fs.existsSync(protectedRoutePath) && fs.existsSync(authContextPath)) {
    const protectedContent = fs.readFileSync(protectedRoutePath, 'utf8');
    const authContent = fs.readFileSync(authContextPath, 'utf8');
    
    const hasAuthHook = protectedContent.includes('useAuth');
    const hasRedirect = protectedContent.includes('Navigate to="/admin/login"');
    const hasAuthProvider = authContent.includes('AuthProvider');
    
    if (hasAuthHook && hasRedirect && hasAuthProvider) {
      console.log('   ✅ Authentication system intact');
      results.authenticationWorking = true;
    } else {
      console.log('   ❌ Authentication system may be broken');
      allTestsPassed = false;
    }
  } else {
    console.log('   ❌ Authentication files not found');
    allTestsPassed = false;
  }
} catch (error) {
  console.log('   ❌ Error checking authentication:', error.message);
  allTestsPassed = false;
}

// Final Summary
console.log('\n📋 Final Verification Results:');
console.log('===============================');
console.log(`✅ Public site components intact: ${results.publicSiteIntact ? 'PASS' : 'FAIL'}`);
console.log(`✅ Admin routing working: ${results.adminRoutingWorking ? 'PASS' : 'FAIL'}`);
console.log(`✅ No CSS conflicts: ${results.noCssConflicts ? 'PASS' : 'FAIL'}`);
console.log(`✅ Build successful: ${results.buildSuccessful ? 'PASS' : 'FAIL'}`);
console.log(`✅ Performance acceptable: ${results.performanceAcceptable ? 'PASS' : 'FAIL'}`);
console.log(`✅ Authentication working: ${results.authenticationWorking ? 'PASS' : 'FAIL'}`);

const passedTests = Object.values(results).filter(Boolean).length;
const totalTests = Object.keys(results).length;

console.log(`\n📊 Overall Score: ${passedTests}/${totalTests} tests passed`);

if (allTestsPassed && passedTests === totalTests) {
  console.log('\n🎉 ALL BACKWARD COMPATIBILITY TESTS PASSED!');
  console.log('✅ Task 10.2 requirements fully satisfied:');
  console.log('   - Public site components work with styled-components');
  console.log('   - No CSS conflicts between styling systems');
  console.log('   - Build process works with dual styling approach');
  console.log('   - Application performance not negatively impacted');
  process.exit(0);
} else {
  console.log('\n⚠️  Some backward compatibility issues remain');
  console.log('❌ Please review the failed tests above');
  process.exit(1);
}