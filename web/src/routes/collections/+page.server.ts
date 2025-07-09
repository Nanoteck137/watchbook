import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, url }) => {
  const query = getPagedQueryOptions(url.searchParams);

  const collections = await locals.apiClient.getCollections({ query });
  if (!collections.success) {
    throw error(collections.error.code, {
      message: collections.error.message,
    });
  }

  return {
    page: collections.data.page,
    collections: collections.data.collections,
  };
};
