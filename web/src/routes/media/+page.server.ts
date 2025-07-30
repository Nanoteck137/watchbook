import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, url }) => {
  const query = getPagedQueryOptions(url.searchParams);

  const media = await locals.apiClient.getMedia({ query });
  if (!media.success) {
    throw error(media.error.code, { message: media.error.message });
  }

  return {
    page: media.data.page,
    media: media.data.media,
  };
};
