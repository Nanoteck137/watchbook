import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params, url }) => {
  const queryParams = url.searchParams;
  const query = getPagedQueryOptions(queryParams);
  query["filter"] = "userList != null";

  let list = queryParams.get("list");
  if (!list || list === "") {
    list = "all";
  }

  if (list !== "all") {
    query["filter"] = `userList == "${list}"`;
  }

  const animes = await locals.apiClient.getUserAnimeList(params.id, {
    query,
  });
  if (!animes.success) {
    throw error(animes.error.code, { message: animes.error.message });
  }

  return {
    page: animes.data.page,
    animes: animes.data.animes,
  };
};
