FROM oven/bun:1.2-debian AS base

ARG VITE_BACKEND_URL

WORKDIR /work
COPY ./ ./
RUN bun install
RUN bun run build

FROM nginx:alpine
COPY --from=base /work/build/client /var/www
COPY ./docker/nginx/default.conf.template /etc/nginx/templates/default.conf.template
CMD ["nginx", "-g", "daemon off;"]
