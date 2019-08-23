
# build stage
FROM node:lts-alpine as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# production stage
FROM nginx:stable-alpine as production-stage
RUN apk add bash
COPY --from=build-stage /app/dist /usr/share/nginx/html
RUN cat /usr/share/nginx/html/index.html
CMD ["nginx", "-g", "daemon off;"]