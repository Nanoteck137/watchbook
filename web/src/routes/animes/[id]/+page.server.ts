import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const anime = await locals.apiClient.getAnimeById(params.id);
  if (!anime.success) {
    throw error(anime.error.code, { message: anime.error.message });
  }

  return {
    anime: anime.data,
  };
};
