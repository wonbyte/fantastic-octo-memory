FROM node:22-alpine

WORKDIR /app

# Install Expo CLI globally
RUN npm install -g expo-cli@latest

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
