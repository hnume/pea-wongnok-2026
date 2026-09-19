import middleware from "next-auth/middleware";

// Next.js 16 renamed `middleware.ts` -> `proxy.ts` and checks that the
// default export is literally a function at load time. Re-exporting
// next-auth's CJS default (`export { default } from "next-auth/middleware"`)
// doesn't satisfy that check, so wrap it in an explicit function instead.
export default function proxy(...args: Parameters<typeof middleware>) {
  return middleware(...args);
}

export const config = {
  matcher: ["/me/:path*"],
};
