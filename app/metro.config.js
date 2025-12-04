const { getDefaultConfig } = require('expo/metro-config');

const config = getDefaultConfig(__dirname);

// Configure Metro for Docker environment
config.server = {
  ...config.server,
  // Allow connections from any IP for Docker
  enhanceMiddleware: (middleware) => {
    return (req, res, next) => {
      // Add CORS headers for development
      res.setHeader('Access-Control-Allow-Origin', '*');
      res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
      res.setHeader('Access-Control-Allow-Headers', 'Content-Type');
      
      if (req.method === 'OPTIONS') {
        res.statusCode = 200;
        res.end();
        return;
      }
      
      return middleware(req, res, next);
    };
  },
};

// Configure watchman for Docker (optional, may need to be disabled)
config.watchFolders = [];

module.exports = config;
