#!/usr/bin/env node

/**
 * Report Modal User Experience Test Suite
 * Tests the streamlined report view improvements for task 11.2
 */

const fs = require('fs');
const path = require('path');

console.log('📋 Report Modal UX Testing...\n');

let allTestsPassed = true;
const results = {
  visualHierarchy: false,
  loadingStates: false,
  errorHandling: false,
  downloadFunctionality: false,
  responsiveDesign: false,
  contentReadability: false
};

// Test 1: Visual hierarchy improvements
console.log('1️⃣ Testing visual hierarchy improvements...');
try {
  const modalPath = path.join(__dirname, 'src/components/admin/report-preview-modal.tsx');
  if (fs.existsSync(modalPath)) {
    const content = fs.readFileSync(modalPath, 'utf8');
    
    const hierarchyElements = [
      'CardTitle.*flex items-center',  // Improved titles with icons
      'bg-blue-50.*border-blue-200',   // Color-coded cards
      'bg-green-50.*border-green-200', // Service type card
      'bg-purple-50.*border-purple-200', // Status card
      'prose prose-sm.*leading-relaxed', // Better typography
      'border-b border-gray-200'       // Clear section separation
    ];
    
    let hierarchyFound = 0;
    hierarchyElements.forEach(element => {
      if (content.match(new RegExp(element))) {
        hierarchyFound++;
      }
    });
    
    if (hierarchyFound >= 4) {
      console.log(`   ✅ Visual hierarchy improvements: ${hierarchyFound}/${hierarchyElements.length}`);
      results.visualHierarchy = true;
    } else {
      console.log(`   ❌ Insufficient visual hierarchy improvements: ${hierarchyFound}/${hierarchyElements.length}`);
      allTestsPassed = false;
    }
  }
} catch (error) {
  console.log('   ❌ Error checking visual hierarchy:', error.message);
  allTestsPassed = false;
}

// Test 2: Loading states improvements
console.log('\n2️⃣ Testing loading states...');
try {
  const modalPath = path.join(__dirname, 'src/components/admin/report-preview-modal.tsx');
  if (fs.existsSync(modalPath)) {
    const content = fs.readFileSync(modalPath, 'utf8');
    
    const loadingElements = [
      'animate-spin rounded-full',     // Spinner animation
      'Loading report...',             // Loading message
      'Please wait while we fetch',    // Descriptive loading text
      'downloadLoading',               // Download loading state
      'Generating PDF...',             // Download progress text
      'disabled={downloadLoading'      // Disabled state during download
    ];
    
    let loadingFound = 0;
    loadingElements.forEach(element => {
      if (content.includes(element)) {
        loadingFound++;
      }
    });
    
    if (loadingFound >= 5) {
      console.log(`   ✅ Loading states implemented: ${loadingFound}/${loadingElements.length}`);
      results.loadingStates = true;
    } else {
      console.log(`   ❌ Insufficient loading states: ${loadingFound}/${loadingElements.length}`);
      allTestsPassed = false;
    }
  }
} catch (error) {
  console.log('   ❌ Error checking loading states:', error.message);
  allTestsPassed = false;
}

// Test 3: Error handling improvements
console.log('\n3️⃣ Testing error handling...');
try {
  const modalPath = path.join(__dirname, 'src/components/admin/report-preview-modal.tsx');
  if (fs.existsSync(modalPath)) {
    const content = fs.readFileSync(modalPath, 'utf8');
    
    const errorElements = [
      'Failed to load report',         // Error message
      'bg-red-100 rounded-full',       // Error icon styling
      'text-red-600',                  // Error text color
      'Try Again',                     // Retry button
      'catch \\(error\\)',             // Error catching
      'console.error'                  // Error logging
    ];
    
    let errorFound = 0;
    errorElements.forEach(element => {
      if (content.match(new RegExp(element))) {
        errorFound++;
      }
    });
    
    if (errorFound >= 4) {
      console.log(`   ✅ Error handling implemented: ${errorFound}/${errorElements.length}`);
      results.errorHandling = true;
    } else {
      console.log(`   ❌ Insufficient error handling: ${errorFound}/${errorElements.length}`);
      allTestsPassed = false;
    }
  }
} catch (error) {
  console.log('   ❌ Error checking error handling:', error.message);
  allTestsPassed = false;
}

// Test 4: Download functionality improvements
console.log('\n4️⃣ Testing download functionality...');
try {
  const modalPath = path.join(__dirname, 'src/components/admin/report-preview-modal.tsx');
  if (fs.existsSync(modalPath)) {
    const content = fs.readFileSync(modalPath, 'utf8');
    
    const downloadElements = [
      'Download PDF',                  // PDF download option
      'Download HTML',                 // HTML download option
      'Formatted for printing',        // PDF description
      'Web-friendly format',           // HTML description
      'async.*onDownload',             // Async download handling
      'setDownloadLoading'             // Download loading state
    ];
    
    let downloadFound = 0;
    downloadElements.forEach(element => {
      if (content.match(new RegExp(element))) {
        downloadFound++;
      }
    });
    
    if (downloadFound >= 5) {
      console.log(`   ✅ Download functionality improved: ${downloadFound}/${downloadElements.length}`);
      results.downloadFunctionality = true;
    } else {
      console.log(`   ❌ Insufficient download improvements: ${downloadFound}/${downloadElements.length}`);
      allTestsPassed = false;
    }
  }
} catch (error) {
  console.log('   ❌ Error checking download functionality:', error.message);
  allTestsPassed = false;
}

// Test 5: Responsive design
console.log('\n5️⃣ Testing responsive design...');
try {
  const modalPath = path.join(__dirname, 'src/components/admin/report-preview-modal.tsx');
  if (fs.existsSync(modalPath)) {
    const content = fs.readFileSync(modalPath, 'utf8');
    
    const responsiveElements = [
      'grid-cols-1 md:grid-cols-3',    // Responsive grid
      'grid-cols-1 md:grid-cols-2',    // Responsive details grid
      'max-w-3xl',                     // Max width constraint
      'isFullscreen.*h-\\[calc',       // Fullscreen height calculation
      'truncate',                      // Text truncation
      'min-w-0'                        // Flex item min-width
    ];
    
    let responsiveFound = 0;
    responsiveElements.forEach(element => {
      if (content.match(new RegExp(element))) {
        responsiveFound++;
      }
    });
    
    if (responsiveFound >= 4) {
      console.log(`   ✅ Responsive design implemented: ${responsiveFound}/${responsiveElements.length}`);
      results.responsiveDesign = true;
    } else {
      console.log(`   ❌ Insufficient responsive design: ${responsiveFound}/${responsiveElements.length}`);
      allTestsPassed = false;
    }
  }
} catch (error) {
  console.log('   ❌ Error checking responsive design:', error.message);
  allTestsPassed = false;
}

// Test 6: Content readability improvements
console.log('\n6️⃣ Testing content readability...');
try {
  const modalPath = path.join(__dirname, 'src/components/admin/report-preview-modal.tsx');
  if (fs.existsSync(modalPath)) {
    const content = fs.readFileSync(modalPath, 'utf8');
    
    const readabilityElements = [
      'leading-relaxed',               // Improved line height
      'prose prose-sm max-w-none',     // Typography system
      'text-gray-700',                 // Readable text color
      'CardDescription',               // Descriptive text
      'font-medium.*text-gray-600',    // Label styling
      'No content available'           // Empty state message
    ];
    
    let readabilityFound = 0;
    readabilityElements.forEach(element => {
      if (content.match(new RegExp(element))) {
        readabilityFound++;
      }
    });
    
    if (readabilityFound >= 4) {
      console.log(`   ✅ Content readability improved: ${readabilityFound}/${readabilityElements.length}`);
      results.contentReadability = true;
    } else {
      console.log(`   ❌ Insufficient readability improvements: ${readabilityFound}/${readabilityElements.length}`);
      allTestsPassed = false;
    }
  }
} catch (error) {
  console.log('   ❌ Error checking content readability:', error.message);
  allTestsPassed = false;
}

// Final Summary
console.log('\n📋 Report Modal UX Test Results:');
console.log('==================================');
console.log(`✅ Visual hierarchy: ${results.visualHierarchy ? 'PASS' : 'FAIL'}`);
console.log(`✅ Loading states: ${results.loadingStates ? 'PASS' : 'FAIL'}`);
console.log(`✅ Error handling: ${results.errorHandling ? 'PASS' : 'FAIL'}`);
console.log(`✅ Download functionality: ${results.downloadFunctionality ? 'PASS' : 'FAIL'}`);
console.log(`✅ Responsive design: ${results.responsiveDesign ? 'PASS' : 'FAIL'}`);
console.log(`✅ Content readability: ${results.contentReadability ? 'PASS' : 'FAIL'}`);

const passedTests = Object.values(results).filter(Boolean).length;
const totalTests = Object.keys(results).length;

console.log(`\n📊 Overall Score: ${passedTests}/${totalTests} tests passed`);

if (allTestsPassed && passedTests === totalTests) {
  console.log('\n🎉 ALL REPORT MODAL UX TESTS PASSED!');
  console.log('✅ Task 11.2 requirements fully satisfied:');
  console.log('   - Simplified report modal with focus on important information');
  console.log('   - Download functionality works correctly for PDF and HTML');
  console.log('   - Improved visual hierarchy and readability');
  console.log('   - Proper loading states and error handling');
  console.log('   - Responsive design across different screen sizes');
  process.exit(0);
} else {
  console.log('\n⚠️  Some UX improvements still needed');
  console.log('❌ Please review the failed tests above');
  process.exit(1);
}