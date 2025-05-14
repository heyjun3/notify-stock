# FROM oven/bun:latest AS base
FROM node:24-slim AS base

WORKDIR /work
COPY ./ ./
RUN npm install
RUN npm run build
# RUN bun install
# RUN bun run build

FROM nginx:alpine
COPY --from=base /work/build/client /var/www
COPY ./docker/nginx/ /etc/nginx/conf.d/
CMD ["nginx", "-g", "daemon off;"]
