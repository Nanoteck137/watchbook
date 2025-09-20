import { getPageOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { FullFilter } from "./types";

function constructFilterSort(
  filter: FullFilter,
  query: Record<string, string>,
) {
  const filters = [];

  // let list = queryParams.get("list");
  // if (!list || list === "") {
  //   list = "all";
  // }

  if (filter.list !== "") {
    filters.push(`userList == "${filter.list}"`);
  } else {
    filters.push("userList != null");
  }

  // if (filter.query !== "") {
  //   filters.push(`title % "%${filter.query}%"`);
  // }

  if (filter.types.length > 0) {
    const s = filter.types.map((i) => `"${i}"`).join(",");
    filters.push(`hasType(${s})`);
  }

  if (filter.status.length > 0) {
    const s = filter.status.map((i) => `"${i}"`).join(",");
    filters.push(`hasStatus(${s})`);
  }

  // if (filter.filters.rating.length > 0) {
  //   const s = filter.filters.rating.map((i) => `"${i}"`).join(",");
  //   filters.push(`hasRating(${s})`);
  // }

  // if (filter.excludes.type.length > 0) {
  //   const s = filter.excludes.type.map((i) => `"${i}"`).join(",");
  //   filters.push(`!hasType(${s})`);
  // }

  // if (filter.excludes.status.length > 0) {
  //   const s = filter.excludes.status.map((i) => `"${i}"`).join(",");
  //   filters.push(`!hasStatus(${s})`);
  // }

  // if (filter.excludes.rating.length > 0) {
  //   const s = filter.excludes.rating.map((i) => `"${i}"`).join(",");
  //   filters.push(`!hasRating(${s})`);
  // }

  query["filter"] = filters.join(" && ");

  switch (filter.sort) {
    case "title-a-z":
      query["sort"] = "sort=+title";
      break;
    case "title-z-a":
      query["sort"] = "sort=-title";
      break;
    case "score-high":
      query["sort"] = "sort=-score";
      break;
    case "score-low":
      query["sort"] = "sort=+score";
      break;
    case "user-score-high":
      query["sort"] = "sort=-userScore";
      break;
    case "user-score-low":
      query["sort"] = "sort=+userScore";
      break;
  }
}

export const load: PageServerLoad = async ({ locals, params, url }) => {
  const query = getPageOptions(url.searchParams);

  // query["userId"] = params.id;
  // query["filter"] = "userList != null";

  // let list = queryParams.get("list");
  // if (!list || list === "") {
  //   list = "all";
  // }

  // if (list !== "all") {
  //   query["filter"] = `userList == "${list}"`;
  // }

  query["userId"] = params.id;

  const filter = FullFilter.parse({
    list: url.searchParams.get("list") ?? "",
    types:
      url.searchParams
        .get("types")
        ?.split(",")
        .filter((t) => t !== "") ?? [],
    status:
      url.searchParams
        .get("status")
        ?.split(",")
        .filter((t) => t !== "") ?? [],
    sort: url.searchParams.get("sort") ?? undefined,
  });
  console.log(filter);

  constructFilterSort(filter, query);

  console.log(query);

  const media = await locals.apiClient.getMedia({
    query,
  });
  if (!media.success) {
    throw error(media.error.code, { message: media.error.message });
  }

  return {
    page: media.data.page,
    media: media.data.media,
    userId: params.id,
    filter,
  };
};
