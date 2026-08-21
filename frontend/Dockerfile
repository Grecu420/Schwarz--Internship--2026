
FROM node:22-alpine AS build-stage

WORKDIR /app


ENV CI=true
ENV PNPM_CONFIRM_MODULES_PURGE=false


RUN corepack enable && corepack prepare pnpm@latest --activate


COPY package.json pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile --ignore-scripts


COPY . .


RUN pnpm build


FROM nginx:alpine AS production-stage

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build-stage /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]