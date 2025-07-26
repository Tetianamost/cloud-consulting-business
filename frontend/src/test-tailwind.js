// Simple test to verify Tailwind CSS configuration
console.log('Testing Tailwind CSS configuration...');

// Check if Tailwind directives are properly loaded
const testElement = document.createElement('div');
testElement.className = 'bg-blue-500 text-white p-4 rounded-lg';
testElement.textContent = 'Tailwind Test';
document.body.appendChild(testElement);

// Check computed styles
const computedStyles = window.getComputedStyle(testElement);
console.log('Background color:', computedStyles.backgroundColor);
console.log('Color:', computedStyles.color);
console.log('Padding:', computedStyles.padding);
console.log('Border radius:', computedStyles.borderRadius);

// Clean up
document.body.removeChild(testElement);

console.log('Tailwind CSS test completed. Check the console output above.');