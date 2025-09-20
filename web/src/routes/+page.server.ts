import type { Media } from "$lib/api/types";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, parent }) => {
  const data = await parent();

  let userMedia: Media[] = [];

  if (data.user) {
    const res = await locals.apiClient.getMedia({ query: { perPage: "10" } });
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    userMedia = res.data.media;
  }

  return {
    ...data,
    userMedia,
  };
};
