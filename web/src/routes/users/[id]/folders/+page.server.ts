import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  const folders = await locals.apiClient.getFolders();
  if (!folders.success) {
    throw error(folders.error.code, { message: folders.error.message });
  }

  return {
    folders: folders.data.folders,
  };
};
