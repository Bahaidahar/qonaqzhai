# Next.js 16 standalone build. NEXT_PUBLIC_API_URL is baked at build time
# because it ships in the client bundle — point it at the public gateway.
FROM node:23-slim AS deps
WORKDIR /app
RUN corepack enable
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
# --ignore-scripts: skip native build scripts (sharp/unrs-resolver). Tailwind's
# oxide and Next ship prebuilt binaries, so the build still works; image
# optimization gracefully degrades without sharp.
RUN pnpm install --frozen-lockfile --ignore-scripts

FROM node:23-alpine AS build
WORKDIR /app
RUN corepack enable
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_TELEMETRY_DISABLED=1
RUN pnpm build

FROM node:23-alpine AS run
WORKDIR /app
ENV NODE_ENV=production NEXT_TELEMETRY_DISABLED=1
COPY --from=build /app/.next/standalone ./
COPY --from=build /app/.next/static ./.next/static
COPY --from=build /app/public ./public
EXPOSE 3000
CMD ["node", "server.js"]
