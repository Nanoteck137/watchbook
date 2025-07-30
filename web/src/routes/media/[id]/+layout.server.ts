import { error } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals, params }) => {
  const media = await locals.apiClient.getMediaById(params.id);
  if (!media.success) {
    throw error(media.error.code, { message: media.error.message });
  }

  return {
    media: media.data,
  };
};
