import { getPageOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { FullFilter } from "./types";

function constructFilterSort(
  filter: FullFilter,
  query: Record<string, string>,
) {
  let f = "";

  if (filter.query !== "") {
    f += `name % "%${filter.query}%"`;
  }

  if (filter.excludes.length > 0) {
    if (f !== "") {
      f += "&&";
    }

    const e = filter.excludes.map((i) => `"${i}"`).join(",");
    f += `!hasCollectionType(${e})`;
  }

  if (filter.filter !== "") {
    if (f !== "") {
      f += "&&";
    }

    f += `hasCollectionType("${filter.filter}")`;
  }

  query["filter"] = f;

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
    excludes: url.searchParams.get("excludes")?.split(",") ?? [],
    filter: url.searchParams.get("filter") ?? "",
  });

  constructFilterSort(filter, query);

  console.log(query);

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
