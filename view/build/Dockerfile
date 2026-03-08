FROM denoland/deno:2.7.4 AS builder
WORKDIR /app
COPY deno.json deno.lock ./
RUN deno install
COPY . .
RUN deno task build

FROM denoland/deno:2.7.4
WORKDIR /app
COPY --from=builder /app/_fresh/ ./_fresh/
COPY --from=builder /app/static/ ./static/
EXPOSE 8000
CMD ["serve", "-A", "_fresh/server.js"]
