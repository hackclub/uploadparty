import { config as loadEnv } from 'dotenv';
import { existsSync } from 'fs';
import { resolve } from 'path';

// Centralized env loading: prefer repo-root .env.local, then .env
const repoRoot = resolve(process.cwd(), '..');
const rootEnvLocal = resolve(repoRoot, '.env.local');
const rootEnv = resolve(repoRoot, '.env');
if (existsSync(rootEnvLocal)) {
  loadEnv({ path: rootEnvLocal });
} else if (existsSync(rootEnv)) {
  loadEnv({ path: rootEnv });
}

/** @type {import('next').NextConfig} */
const nextConfig = {
  // No need to manually expose NEXT_PUBLIC_*; Next automatically inlines them.
};

export default nextConfig;
