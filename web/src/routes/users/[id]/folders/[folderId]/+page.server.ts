import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const folder = await locals.apiClient.getFolderById(params.folderId);
  if (!folder.success) {
    throw error(folder.error.code, { message: folder.error.message });
  }

  const items = await locals.apiClient.getFolderItems(params.folderId);
  if (!items.success) {
    throw error(items.error.code, { message: items.error.message });
  }

  return {
    folder: folder.data,
    items: items.data.items,
  };
};
