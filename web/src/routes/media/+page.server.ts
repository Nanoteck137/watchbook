import { getPageOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { FullFilter } from "./types";

function constructFilterSort(
  filter: FullFilter,
  query: Record<string, string>,
) {
  const filters = [];

  if (filter.query !== "") {
    filters.push(`title % "%${filter.query}%"`);
  }

  if (filter.filters.type.length > 0) {
    const s = filter.filters.type.map((i) => `"${i}"`).join(",");
    filters.push(`hasType(${s})`);
  }

  if (filter.filters.status.length > 0) {
    const s = filter.filters.status.map((i) => `"${i}"`).join(",");
    filters.push(`hasStatus(${s})`);
  }

  if (filter.filters.rating.length > 0) {
    const s = filter.filters.rating.map((i) => `"${i}"`).join(",");
    filters.push(`hasRating(${s})`);
  }

  if (filter.excludes.type.length > 0) {
    const s = filter.excludes.type.map((i) => `"${i}"`).join(",");
    filters.push(`!hasType(${s})`);
  }

  if (filter.excludes.status.length > 0) {
    const s = filter.excludes.status.map((i) => `"${i}"`).join(",");
    filters.push(`!hasStatus(${s})`);
  }

  if (filter.excludes.rating.length > 0) {
    const s = filter.excludes.rating.map((i) => `"${i}"`).join(",");
    filters.push(`!hasRating(${s})`);
  }

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
  }
}

export const load: PageServerLoad = async ({ locals, url }) => {
  const query = getPageOptions(url.searchParams);

  const filter = FullFilter.parse({
    query: url.searchParams.get("query") ?? "",
    sort: url.searchParams.get("sort") ?? undefined,
    filters: {
      type: url.searchParams.get("filterType")?.split(",") ?? [],
      status: url.searchParams.get("filterStatus")?.split(",") ?? [],
      rating: url.searchParams.get("filterRating")?.split(",") ?? [],
    },
    excludes: {
      type: url.searchParams.get("excludeType")?.split(",") ?? [],
      status: url.searchParams.get("excludeStatus")?.split(",") ?? [],
      rating: url.searchParams.get("excludeRating")?.split(",") ?? [],
    },
  });

  constructFilterSort(filter, query);

  const media = await locals.apiClient.getMedia({ query });
  if (!media.success) {
    throw error(media.error.code, { message: media.error.message });
  }

  return {
    page: media.data.page,
    media: media.data.media,
    filter,
  };
};
