import type { Config } from "jest"
import nextJest from "next/jest.js"

const createJestConfig = nextJest({
  // Path to the Next.js app, used to load next.config.ts and .env files in tests
  dir: "./",
})

const config: Config = {
  coverageProvider: "v8",
  testEnvironment: "node",
  roots: ["<rootDir>/tests"],
  moduleNameMapper: {
    "^@/(.*)$": "<rootDir>/$1",
  },
}

// Exported this way so next/jest can load the async Next.js config
export default createJestConfig(config)
