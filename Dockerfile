# Stage 1: Build frontend
FROM node:22-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend with embedded frontend
FROM golang:1.26-alpine AS backend
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
COPY --from=frontend /app/frontend/dist ./cmd/almaz/dist
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /almaz ./cmd/almaz

# Stage 3: Minimal production image
FROM scratch
COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend /almaz /almaz
COPY migrations/ /migrations/
EXPOSE 8080
ENTRYPOINT ["/almaz"]
