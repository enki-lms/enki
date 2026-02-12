import { redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = () => {
  if (typeof document !== "undefined") {
    const hasToken = document.cookie.split(";").some((c) => c.trim().startsWith("token="));
    if (hasToken) {
      throw redirect(302, "/home");
    }
  }
};
