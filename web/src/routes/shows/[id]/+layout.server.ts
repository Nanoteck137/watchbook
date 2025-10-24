import { error } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals, params }) => {
  const show = await locals.apiClient.getShowById(params.id);
  if (!show.success) {
    throw error(show.error.code, { message: show.error.message });
  }

  const seasons = await locals.apiClient.getShowSeasons(params.id);
  if (!seasons.success) {
    throw error(seasons.error.code, { message: seasons.error.message });
  }

  return {
    show: show.data,
    seasons: seasons.data.seasons,
  };
};
