# Stage 1: Build Quasar Web Application
FROM node:20-alpine AS builder

WORKDIR /app

# Salin seluruh source code (node_modules dan dist sudah diabaikan lewat .dockerignore)
COPY . .

# Install dependencies dan jalankan quasar prepare
RUN npm ci

# Argumen build untuk base URL API
ARG API_BASE_URL=/api/v1
ENV API_BASE_URL=${API_BASE_URL}

# Build bundle Quasar SPA (menghasilkan dist/spa)
RUN npm run build

# Stage 2: Nginx Web Server
FROM nginx:alpine

# Salin file hasil build SPA ke direktori Nginx
COPY --from=builder /app/dist/spa /usr/share/nginx/html

# Salin konfigurasi custom Nginx (SPA history routing & API proxy)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
