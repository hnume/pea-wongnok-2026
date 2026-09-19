import { getServerSession } from "next-auth";
import { NextRequest, NextResponse } from "next/server";

import { authOptions } from "@/lib/auth/auth-options";

const NEXT_AUTH_COOKIE = /^(__Secure-|__Host-)?next-auth/;

export async function GET(request: NextRequest) {
  // Must read the session before the cookies are cleared.
  const session = await getServerSession(authOptions);
  const idToken = session?.idToken;

  // NEXTAUTH_URL first: request.url can reflect 0.0.0.0 when bound with `next dev -H 0.0.0.0`.
  const origin = new URL(process.env.NEXTAUTH_URL ?? request.url).origin;

  let destination = origin;
  if (idToken && process.env.KEYCLOAK_ISSUER) {
    // KEYCLOAK_ISSUER must be the browser-reachable URL, not the docker-internal hostname.
    const issuer = process.env.KEYCLOAK_ISSUER.replace(/\/+$/, "");
    const logoutUrl = new URL(`${issuer}/protocol/openid-connect/logout`);
    logoutUrl.searchParams.set("id_token_hint", idToken);
    logoutUrl.searchParams.set("post_logout_redirect_uri", origin);
    destination = logoutUrl.toString();
  }

  const response = NextResponse.redirect(destination);

  for (const { name } of request.cookies.getAll()) {
    if (!NEXT_AUTH_COOKIE.test(name)) continue;
    if (name.startsWith("__Secure-") || name.startsWith("__Host-")) {
      response.cookies.delete({ name, path: "/", secure: true });
    } else {
      response.cookies.delete(name);
    }
  }

  return response;
}
