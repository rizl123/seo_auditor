import type { JWTPayload } from "jose";

export type AuthUserInfo = JWTPayload & {
  at_hash: string;
  created_at: number;
  email: string;
  email_verified: boolean;
  name: string;
  picture: string;
  updated_at: number;
  username: string;
};
