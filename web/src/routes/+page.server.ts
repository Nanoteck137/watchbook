import type { Collection, Media } from "$lib/api/types";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, parent }) => {
  const data = await parent();

  let userInprogressMedia: Media[] = [];
  let userBacklogMedia: Media[] = [];
  let recentlyCreatedMedia: Media[] = [];
  let recentlyCreatedCollections: Collection[] = [];

  if (data.user) {
    {
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

      userInprogressMedia = res.data.media;
    }

    {
      const res = await locals.apiClient.getMedia({
        query: {
          filter: 'userList == "backlog"',
          sort: "sort=title",
          perPage: "10",
        },
      });
      if (!res.success) {
        throw error(res.error.code, { message: res.error.message });
      }

      userBacklogMedia = res.data.media;
    }

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

      recentlyCreatedMedia = res.data.media;
    }

    {
      const res = await locals.apiClient.getCollections({
        query: {
          filter: "",
          sort: "sort=-created",
          perPage: "10",
        },
      });
      if (!res.success) {
        throw error(res.error.code, { message: res.error.message });
      }

      recentlyCreatedCollections = res.data.collections;
    }
  }

  return {
    ...data,
    userInprogressMedia,
    userBacklogMedia,
    recentlyCreatedMedia,
    recentlyCreatedCollections,
  };
};
