import { error } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals, params }) => {
  const user = await locals.apiClient.getUser(params.id);
  if (!user.success) {
    return error(user.error.code, { message: user.error.message });
  }

  return {
    userData: user.data,
  };
};
