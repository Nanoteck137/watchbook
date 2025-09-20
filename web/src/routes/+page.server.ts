import type { Media } from "$lib/api/types";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, parent }) => {
  const data = await parent();

  let userMedia: Media[] = [];
  let recentlyReleasedMedia: Media[] = [];

  if (data.user) {
    const res = await locals.apiClient.getMedia({
      query: {
        filter: 'userList == "in-progress"',
        sort: "sort=-userUpdated",
        perPage: "10",
      },
    });
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    userMedia = res.data.media;

    {
      const res = await locals.apiClient.getMedia({
        query: {
          filter: "",
          sort: "sort=-created",
          perPage: "10",
        },
      });
      if (!res.success) {
        throw error(res.error.code, { message: res.error.message });
      }

      recentlyReleasedMedia = res.data.media;
    }
  }

  return {
    ...data,
    userMedia,
    recentlyReleasedMedia,
  };
};
