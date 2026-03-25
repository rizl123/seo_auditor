# СТАДИЯ 1: Сборщик (Builder)
# Содержит зависимости (deps + devDeps) и исходный код для сборки
FROM node:22.12.0-alpine AS builder
WORKDIR /app

# Копируем конфиги зависимостей для кэширования
COPY package.json package-lock.json ./
RUN npm ci

# Копируем остальной код и собираем
COPY . .
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

# СТАДИЯ 2: Рантайм (Runner)
# Содержит только скомпилированный сервис (standalone)
FROM node:22.12.0-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV HOSTNAME="0.0.0.0"

# Копируем только артефакты сборки
# Папка .next/standalone — это самодостаточный сервер Node.js
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

EXPOSE 3000
CMD ["node", "server.js"]
