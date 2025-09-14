import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  const res = await locals.apiClient.getProviders();
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  return {
    providers: res.data.providers,
  };
};
