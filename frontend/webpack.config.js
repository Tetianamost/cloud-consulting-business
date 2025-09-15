const path = require('path');

module.exports = function override(config, env) {
  // Only apply optimizations in production
  if (env === 'production') {
    // Optimize bundle splitting for admin vs public components
    config.optimization = {
      ...config.optimization,
      splitChunks: {
        ...config.optimization.splitChunks,
        cacheGroups: {
          // Separate chunk for admin components
          admin: {
            test: /[\\/]src[\\/]components[\\/]admin[\\/]/,
            name: 'admin',
            chunks: 'all',
            priority: 20,
          },
          // Separate chunk for UI components (shared between admin and public)
          ui: {
            test: /[\\/]src[\\/]components[\\/]ui[\\/]/,
            name: 'ui',
            chunks: 'all',
            priority: 15,
          },
          // Separate chunk for public site components
          public: {
            test: /[\\/]src[\\/]components[\\/](layout|sections)[\\/]/,
            name: 'public',
            chunks: 'all',
            priority: 10,
          },
          // Vendor chunk for third-party libraries
          vendor: {
            test: /[\\/]node_modules[\\/]/,
            name: 'vendors',
            chunks: 'all',
            priority: 5,
          },
          // Default chunk
          default: {
            minChunks: 2,
            priority: -20,
            reuseExistingChunk: true,
          },
        },
      },
    };

    // Add performance hints
    config.performance = {
      ...config.performance,
      hints: 'warning',
      maxEntrypointSize: 512000, // 500KB
      maxAssetSize: 512000, // 500KB
    };
  }

  return config;
};

module.exports.override = module.exports;