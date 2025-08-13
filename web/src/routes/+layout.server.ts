import type { GetMe } from "$lib/api/types";
import { error } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals, cookies }) => {
  let user: GetMe | null = null;
  if (locals.token) {
    const res = await locals.apiClient.getMe();
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    user = res.data;
  }

  return {
    apiAddress: locals.apiAddress,
    userToken: locals.token,
    user,
  };
};
