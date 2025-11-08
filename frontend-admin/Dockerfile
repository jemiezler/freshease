FROM node:22-slim AS build
WORKDIR /app

# Install deps
COPY package*.json ./
RUN npm ci

# Copy source & build (App Router ready; produces standalone server)
COPY . .
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

# --- Runtime stage ---
FROM node:22-slim
WORKDIR /app

# Non-root user
RUN useradd -m nextjs
USER nextjs

# Copy the standalone output (fast, minimal runtime)
# Next.js produces .next/standalone with a server.js entrypoint
COPY --from=build /app/.next/standalone ./
# Static assets
COPY --from=build /app/.next/static ./.next/static
# Public assets (if you use them)
COPY --from=build /app/public ./public

# Env & port for Next server
ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME=0.0.0.0

EXPOSE 3000
# The standalone output contains server.js at repo root (copied above)
CMD ["node", "server.js"]
