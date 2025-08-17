import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params, url }) => {
  const queryParams = url.searchParams;
  const query = getPagedQueryOptions(queryParams);
  query["userId"] = params.id;
  query["filter"] = "userList != null";

  let list = queryParams.get("list");
  if (!list || list === "") {
    list = "all";
  }

  if (list !== "all") {
    query["filter"] = `userList == "${list}"`;
  }

  delete query["sort"];
  const sort = queryParams.get("sort");
  if (sort) {
    switch (sort) {
      case "titleAsc":
        query["sort"] = "sort=+title";
        break;
      case "titleDesc":
        query["sort"] = "sort=-title";
        break;
      case "userScoreAsc":
        query["sort"] = "sort=+userScore,+title";
        break;
      case "userScoreDesc":
        query["sort"] = "sort=-userScore,+title";
        break;
      case "scoreAsc":
        query["sort"] = "sort=+score,+title";
        break;
      case "scoreDesc":
        query["sort"] = "sort=-score,+title";
        break;
    }
  }

  const animes = await locals.apiClient.getMedia({
    query,
  });
  if (!animes.success) {
    throw error(animes.error.code, { message: animes.error.message });
  }

  return {
    page: animes.data.page,
    media: animes.data.media,
    userId: params.id,
    list,
  };
};
