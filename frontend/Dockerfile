# =============================================================================
# Stage: base
#
# Build the base image with the necessary tools so that it can be used in the 
# build stage without having to install them again.
# =============================================================================
FROM node:18-alpine AS base

WORKDIR /usr/src/app

# Install necessary tools and dependencies
RUN apk update && apk upgrade && apk add --no-cache git

# =============================================================================
# Stage: build
#
# Download project dependencies and build the Vue.js application.
# =============================================================================
FROM base AS build

WORKDIR /usr/src/app

# Pre-copy/cache package.json and package-lock.json for pre-downloading dependencies
# and only redownloading them in subsequent builds if they change.
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the source code and build the project
COPY . .

# Build the application for production
RUN npm run build

# =============================================================================
# Stage: release
#
# Use a minimal Nginx image to serve the built static files.
# =============================================================================
FROM nginx:stable-alpine AS release

# Copy the built files from the previous stage to the Nginx HTML directory
COPY --from=build /usr/src/app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

# Optional: Copy custom Nginx configuration file if needed
# COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80

# Start Nginx when the container starts
CMD ["nginx", "-g", "daemon off;"]