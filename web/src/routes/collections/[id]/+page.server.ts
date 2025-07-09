import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const collection = await locals.apiClient.getCollectionById(params.id);
  if (!collection.success) {
    throw error(collection.error.code, { message: collection.error.message });
  }

  const items = await locals.apiClient.getCollectionItems(params.id);
  if (!items.success) {
    throw error(items.error.code, { message: items.error.message });
  }

  return {
    collection: collection.data,
    items: items.data.items,
  };
};
