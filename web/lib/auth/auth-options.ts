import type { NextAuthOptions } from "next-auth";
import { decode as defaultDecode } from "next-auth/jwt";
import KeycloakProvider from "next-auth/providers/keycloak";

export const authOptions: NextAuthOptions = {
  providers: [
    KeycloakProvider({
      clientId: process.env.KEYCLOAK_CLIENT_ID!,
      clientSecret: process.env.KEYCLOAK_CLIENT_SECRET!,
      issuer: process.env.KEYCLOAK_ISSUER!,
      authorization: {
        params: {
          prompt: "login", // บังคับให้เข้าสู่ระบบใหม่ที่หน้าต่างของ Auth0 เสมอ
        }
      }
    }),

  ],

};
