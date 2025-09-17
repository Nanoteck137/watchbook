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
    filters.push(`name % "%${filter.query}%"`);
  }

  if (filter.filters.type.length > 0) {
    const s = filter.filters.type.map((i) => `"${i}"`).join(",");
    filters.push(`hasType(${s})`);
  }

  if (filter.excludes.type.length > 0) {
    const s = filter.excludes.type.map((i) => `"${i}"`).join(",");
    filters.push(`!hasType(${s})`);
  }

  query["filter"] = filters.join(" && ");

  switch (filter.sort) {
    case "name-a-z":
      query["sort"] = "sort=+name";
      break;
    case "name-z-a":
      query["sort"] = "sort=-name";
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
    },
    excludes: {
      type: url.searchParams.get("excludeType")?.split(",") ?? [],
    },
  });

  constructFilterSort(filter, query);

  const res = await locals.apiClient.getCollections({ query });
  if (!res.success) {
    throw error(res.error.code, {
      message: res.error.message,
    });
  }

  return {
    page: res.data.page,
    collections: res.data.collections,
    filter,
  };
};
