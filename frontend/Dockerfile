# Use alpine base image and install what we need
FROM alpine:latest

# Install Node.js, npm, nginx, and other dependencies
RUN apk add --no-cache nodejs npm nginx curl

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Copy built app to nginx directory
RUN cp -r build/* /usr/share/nginx/html/

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Create nginx directories
RUN mkdir -p /var/log/nginx /var/lib/nginx/tmp

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]