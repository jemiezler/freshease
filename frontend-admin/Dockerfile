# syntax=docker/dockerfile:1.5
# Multi-platform build for Next.js frontend
ARG TARGETPLATFORM=linux/amd64
FROM --platform=$TARGETPLATFORM node:22-slim AS build
WORKDIR /app

# Enable Corepack for Yarn
RUN corepack enable

# Copy lockfiles first for better cache hits
COPY package.json yarn.lock ./

# Use npm registry with retries to avoid network errors
RUN yarn config set registry https://registry.npmjs.org/ && \
    yarn install --frozen-lockfile --network-concurrency 1 --network-timeout 600000

# Copy source & build
COPY . .
ARG NEXT_PUBLIC_API_BASE_URL
ENV NEXT_PUBLIC_API_BASE_URL=${NEXT_PUBLIC_API_BASE_URL:-https://freshease.seezdev.com/api}
RUN yarn run build

# ---------- Runtime ----------
FROM --platform=$TARGETPLATFORM node:22-slim
WORKDIR /app

# Create non-root user
RUN useradd -m nextjs
USER nextjs

# Copy standalone output
COPY --from=build /app/.next/standalone ./
# Copy static and public files
COPY --from=build /app/.next/static ./.next/static
COPY --from=build /app/public ./public

# Runtime environment
ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME=0.0.0.0

EXPOSE 3000
CMD ["node", "server.js"]
