import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, url }) => {
  const query = getPagedQueryOptions(url.searchParams);

  const releases = await locals.apiClient.getReleases({ query });
  if (!releases.success) {
    throw error(releases.error.code, { message: releases.error.message });
  }

  return {
    page: releases.data.page,
    releases: releases.data.releases,
  };
};
