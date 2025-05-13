import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  const animes = await locals.apiClient.getAnimes();
  if (!animes.success) {
    throw error(animes.error.code, { message: animes.error.message });
  }

  return {
    page: animes.data.page,
    animes: animes.data.animes,
  };
};
