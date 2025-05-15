FROM oven/bun:1.2-debian AS base

WORKDIR /work
COPY ./ ./
RUN bun install
RUN bun run build

FROM nginx:alpine
COPY --from=base /work/build/client /var/www
COPY ./docker/nginx/ /etc/nginx/conf.d/
CMD ["nginx", "-g", "daemon off;"]
