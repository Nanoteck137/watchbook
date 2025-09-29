import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const query: Record<string, string> = {};
  query["userId"] = params.id;

  const folders = await locals.apiClient.getFolders({ query });
  if (!folders.success) {
    throw error(folders.error.code, { message: folders.error.message });
  }

  return {
    folders: folders.data.folders,
  };
};
