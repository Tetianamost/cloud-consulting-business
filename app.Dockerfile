# Simple Dockerfile for pre-built binaries (similar to quota-service approach)
FROM public.ecr.aws/docker/library/nginx:alpine

# Install supervisor
RUN apk add --no-cache supervisor

# Copy pre-built frontend
COPY frontend/build /usr/share/nginx/html

# Copy pre-built backend binary
COPY .build/server /usr/local/bin/server

# Copy backend templates
COPY backend/templates /usr/local/bin/templates

# Copy configuration files
COPY nginx.conf /etc/nginx/nginx.conf
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

EXPOSE 80

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
