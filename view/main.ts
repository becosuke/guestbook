import { App, staticFiles } from "fresh";
import type { State } from "./utils.ts";

const API_UPSTREAM = "http://localhost:50080";

export const app = new App<State>();

app.use(staticFiles());

// Proxy /api requests to the backend REST server
app.use(async (ctx) => {
  const url = new URL(ctx.req.url);
  if (url.pathname.startsWith("/api/")) {
    const upstream = `${API_UPSTREAM}${url.pathname}${url.search}`;
    const headers = new Headers(ctx.req.headers);
    headers.delete("host");
    const res = await fetch(upstream, {
      method: ctx.req.method,
      headers,
      body: ctx.req.method !== "GET" && ctx.req.method !== "HEAD"
        ? ctx.req.body
        : undefined,
    });
    return new Response(res.body, {
      status: res.status,
      headers: res.headers,
    });
  }
  return ctx.next();
});

app.fsRoutes();
