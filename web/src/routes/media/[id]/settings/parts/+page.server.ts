import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const res = await locals.apiClient.getMediaParts(params.id);
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  return {
    parts: res.data.parts,
  };
};
