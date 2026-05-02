import { NextResponse } from "next/server";
import { BASE_URL } from "@/config/urls";
import { processLogout } from "@/lib/auth/process";

export async function GET() {
  const logoutUrl = await processLogout();
  return NextResponse.redirect(logoutUrl || new URL("/", BASE_URL));
}
