import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, url }) => {
  const query = getPagedQueryOptions(url.searchParams);

  const animes = await locals.apiClient.getAnimes({ query });
  if (!animes.success) {
    throw error(animes.error.code, { message: animes.error.message });
  }

  return {
    page: animes.data.page,
    animes: animes.data.animes,
  };
};
