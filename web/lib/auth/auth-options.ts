import type { NextAuthOptions } from "next-auth";
import KeycloakProvider from "next-auth/providers/keycloak";

export const authOptions: NextAuthOptions = {
  providers: [
    KeycloakProvider({
      clientId: process.env.KEYCLOAK_CLIENT_ID!,
      clientSecret: process.env.KEYCLOAK_CLIENT_SECRET!,
      issuer: process.env.KEYCLOAK_ISSUER!,
      // authorization: {
      //   params: {
      //     prompt: "login", // บังคับให้เข้าสู่ระบบใหม่ที่หน้าต่างของ Auth0 เสมอ
      //   }
      // }
    }),

  ],
  callbacks: {
    async jwt({ token, account }) {
      console.log("jwt work")
      // Only present on the initial sign-in request.
      if (account) {
        // token.idToken = account.id_token;
        token.accessToken = account.access_token;
        // token.refreshToken = account.refresh_token;
      }
      return token;
    },
    async session({ session, token }) {
      console.log("session work")

      // session.idToken = token.idToken;
      session.accessToken = token.accessToken;
      // session.refreshToken = token.refreshToken;
      return session;
    },
  },

};
