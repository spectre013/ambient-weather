FROM node:lts-alpine
RUN apk add --no-cache tzdata
ENV TZ America/Denver

# Create app directory
WORKDIR /usr/local/app

# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# where available (npm@5+)
COPY package*.json ./

RUN npm install
# If you are building your code for production
# RUN npm ci --only=production
RUN npm install pm2 -g
# Bundle app source
COPY . .
ENV NODE_ENV production
CMD ["pm2-runtime", "/usr/local/app/process.yml"]