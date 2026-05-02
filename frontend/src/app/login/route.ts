import { NextResponse } from "next/server";
import { processLogin } from "@/lib/auth/process";

export async function GET() {
  const authorizationUrl = await processLogin();
  return NextResponse.redirect(authorizationUrl);
}
