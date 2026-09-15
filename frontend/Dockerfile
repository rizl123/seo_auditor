FROM node:25.9.0-alpine AS builder
WORKDIR /app

RUN npm install -g pnpm@11.24.0

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN --mount=type=cache,target=/root/.local/share/pnpm/store pnpm ci

COPY . .
ENV NEXT_TELEMETRY_DISABLED=1
RUN pnpm run build

FROM node:25.9.0-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV HOSTNAME="0.0.0.0"

COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public* ./public

EXPOSE 3000
CMD ["node", "server.js"]
