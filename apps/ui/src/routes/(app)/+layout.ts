import { redirect } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";

interface JWTPayload {
  role: string;
  user_id: number;
  email: string;
  full_name: string;
  institution: string;
}

function parseJWT(token: string): JWTPayload | null {
  try {
    const base64Url = token.split(".")[1];
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
        .join(""),
    );
    return JSON.parse(jsonPayload);
  } catch {
    return null;
  }
}

function getTokenFromCookie(): string | null {
  const tokenCookie = document.cookie
    .split(";")
    .find((c) => c.trim().startsWith("token="));
  if (!tokenCookie) return null;
  return tokenCookie.split("=")[1];
}

export const load: LayoutLoad = ({route}) => {
  if (typeof document !== "undefined") {
    const token = getTokenFromCookie();
    if (!token) {
      throw redirect(302, "/login");
    }

    const payload = parseJWT(token);
    console.log(payload?.role)
    if (payload && (payload.role === "teacher" || payload.role === "admin") && !route.id.includes("admin")) {
      throw redirect(302, "/admin");
    }

    return {
      user: payload ? {
        id: payload.user_id,
        email: payload.email,
        fullName: payload.full_name,
        role: payload.role,
        institution: payload.institution,
      } : null,
    };
  }

  return { user: null };
};
