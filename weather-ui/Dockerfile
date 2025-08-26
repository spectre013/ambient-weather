
# build stage
FROM node:lts-alpine3.21 as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN ls -l
RUN npm run build

# production stage
FROM nginx:stable-alpine as production-stage
RUN apk add bash
COPY --from=build-stage /app/dist /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]