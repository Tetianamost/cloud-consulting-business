#!/usr/bin/env node

/**
 * Bundle Analysis Script
 * Analyzes the webpack bundle to ensure performance optimizations are working
 */

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

console.log('🔍 Analyzing bundle performance...\n');

// Build the project first
console.log('📦 Building project...');
try {
  execSync('npm run build', { stdio: 'inherit' });
} catch (error) {
  console.error('❌ Build failed:', error.message);
  process.exit(1);
}

// Analyze build directory
const buildDir = path.join(__dirname, '../build/static');

function analyzeDirectory(dir, type) {
  if (!fs.existsSync(dir)) {
    console.log(`⚠️  ${type} directory not found: ${dir}`);
    return { files: [], totalSize: 0 };
  }

  const files = fs.readdirSync(dir)
    .filter(file => file.endsWith(type === 'js' ? '.js' : '.css'))
    .map(file => {
      const filePath = path.join(dir, file);
      const stats = fs.statSync(filePath);
      return {
        name: file,
        size: stats.size,
        sizeKB: Math.round(stats.size / 1024),
      };
    })
    .sort((a, b) => b.size - a.size);

  const totalSize = files.reduce((sum, file) => sum + file.size, 0);

  return { files, totalSize };
}

// Analyze JavaScript files
const jsAnalysis = analyzeDirectory(path.join(buildDir, 'js'), 'js');
console.log('📊 JavaScript Bundle Analysis:');
console.log(`   Total JS Size: ${Math.round(jsAnalysis.totalSize / 1024)} KB`);

jsAnalysis.files.forEach(file => {
  let category = '📄 Other';
  if (file.name.includes('admin')) category = '🔐 Admin';
  else if (file.name.includes('vendor') || file.name.includes('chunk')) category = '📚 Vendor';
  else if (file.name.includes('main')) category = '🏠 Main';
  
  console.log(`   ${category}: ${file.name} (${file.sizeKB} KB)`);
});

// Analyze CSS files
const cssAnalysis = analyzeDirectory(path.join(buildDir, 'css'), 'css');
console.log('\n🎨 CSS Bundle Analysis:');
console.log(`   Total CSS Size: ${Math.round(cssAnalysis.totalSize / 1024)} KB`);

cssAnalysis.files.forEach(file => {
  console.log(`   📄 ${file.name} (${file.sizeKB} KB)`);
});

// Performance recommendations
console.log('\n💡 Performance Recommendations:');

const totalBundleSize = (jsAnalysis.totalSize + cssAnalysis.totalSize) / 1024;
console.log(`   Total Bundle Size: ${Math.round(totalBundleSize)} KB`);

if (totalBundleSize > 500) {
  console.log('   ⚠️  Bundle size is large (>500KB). Consider:');
  console.log('      - Further code splitting');
  console.log('      - Tree shaking optimization');
  console.log('      - Lazy loading more components');
} else if (totalBundleSize > 250) {
  console.log('   ✅ Bundle size is acceptable (250-500KB)');
} else {
  console.log('   🎉 Excellent bundle size (<250KB)');
}

// Check for proper splitting
const hasAdminChunk = jsAnalysis.files.some(f => f.name.includes('admin'));
const hasVendorChunk = jsAnalysis.files.some(f => f.name.includes('vendor') || f.name.includes('chunk'));

console.log('\n🔧 Bundle Splitting Analysis:');
console.log(`   Admin chunk separated: ${hasAdminChunk ? '✅' : '❌'}`);
console.log(`   Vendor chunk separated: ${hasVendorChunk ? '✅' : '❌'}`);

if (hasAdminChunk && hasVendorChunk) {
  console.log('   🎉 Bundle splitting is working correctly!');
} else {
  console.log('   ⚠️  Bundle splitting may need optimization');
}

// Tailwind CSS optimization check
const mainCssFile = cssAnalysis.files.find(f => f.name.includes('main'));
if (mainCssFile) {
  if (mainCssFile.sizeKB < 50) {
    console.log('   ✅ Tailwind CSS appears to be properly purged');
  } else {
    console.log('   ⚠️  CSS bundle is large - check Tailwind purging');
  }
}

console.log('\n✨ Analysis complete!');

// Generate performance report
const report = {
  timestamp: new Date().toISOString(),
  totalBundleSize: Math.round(totalBundleSize),
  jsSize: Math.round(jsAnalysis.totalSize / 1024),
  cssSize: Math.round(cssAnalysis.totalSize / 1024),
  hasAdminChunk,
  hasVendorChunk,
  files: {
    js: jsAnalysis.files.map(f => ({ name: f.name, sizeKB: f.sizeKB })),
    css: cssAnalysis.files.map(f => ({ name: f.name, sizeKB: f.sizeKB })),
  },
};

fs.writeFileSync(
  path.join(__dirname, '../build/bundle-analysis.json'),
  JSON.stringify(report, null, 2)
);

console.log('📋 Detailed report saved to build/bundle-analysis.json');