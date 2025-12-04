FROM node:22-alpine

# Update npm to latest version
RUN npm install -g npm@11.6.4

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy application code
COPY . .

# Expose ports for Expo
EXPOSE 3000 19000 19001 19002

# Start Expo development server
CMD ["npm", "start"]
