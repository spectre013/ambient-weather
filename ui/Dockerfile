
# build stage
FROM node:lts-alpine3.21 AS build-stage
WORKDIR /app
COPY package*.json ./
COPY nginx.conf ./
RUN npm install
COPY . .
RUN npm run build

# production stage
FROM nginx:stable-alpine AS production-stage
RUN apk add bash
COPY --from=build-stage /app/nginx.conf /etc/nginx/
COPY --from=build-stage /app/dist /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]