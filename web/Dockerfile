# --- build stage: compile the SPA to static files -------------------------
# node:22-alpine = current LTS on a small base. Vite 6 needs Node 18+.
FROM node:22-alpine AS build
WORKDIR /app

# Copy manifests first so the `npm ci` layer caches unless deps change.
COPY package.json package-lock.json ./
# npm ci = clean, reproducible install straight from the lockfile.
RUN npm ci

# Then the source (node_modules + dist excluded via .dockerignore).
COPY . .

# VITE_API_BASE_URL is inlined into the bundle at BUILD time. Default "" =
# same-origin relative calls (`/players`); the ALB routes API paths to the
# backend, so the image is environment-agnostic. Override for cross-origin/dev:
#   docker build --build-arg VITE_API_BASE_URL=https://api.example.com
ARG VITE_API_BASE_URL=""
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
RUN npm run build            # tsc -b && vite build  ->  /app/dist

# --- runtime stage: serve the built files with nginx ----------------------
# Pinned minor; multi-arch base so buildx can emit arm64 for the t4g nodes.
FROM nginx:1.27-alpine
# Swap nginx's default server block for our SPA config (fallback + probe).
COPY nginx.conf /etc/nginx/conf.d/default.conf
# Ship only the compiled assets, nothing from the build toolchain.
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
# The nginx base image already runs `nginx -g 'daemon off;'` as its CMD.
