import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const stats = await locals.apiClient.getUserStats(params.id);
  if (!stats.success) {
    throw error(stats.error.code, { message: stats.error.message });
  }

  return {
    stats: stats.data,
  };
};
