FROM node:22-slim AS build
WORKDIR /app

# Enable Corepack to get Yarn
RUN corepack enable

# Copy lockfiles first for better cache hits
COPY package.json yarn.lock ./

# Install deps (production+dev for building)
RUN yarn install --frozen-lockfile

# Copy source & build
COPY . .
ENV NEXT_TELEMETRY_DISABLED=1
RUN yarn build

# ---------- Runtime ----------
FROM node:22-slim
WORKDIR /app

# Create non-root user
RUN useradd -m nextjs
USER nextjs

# Copy standalone output
COPY --from=build /app/.next/standalone ./
# Static + public
COPY --from=build /app/.next/static ./.next/static
COPY --from=build /app/public ./public

# Runtime env
ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME=0.0.0.0

EXPOSE 3000
CMD ["node", "server.js"]
