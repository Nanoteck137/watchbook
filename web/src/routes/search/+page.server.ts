import type { Media } from "$lib/api/types";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ request, locals }) => {
  const url = new URL(request.url);
  const query = url.searchParams.get("query") ?? "";

  let media = [] as Media[];
  let mediaError: string | null = null;

  const [mediaQuery] = await Promise.all([
    locals.apiClient.getMedia({
      query: { filter: `title % "%${query}%"`, perPage: "10" },
    }),
  ]);

  if (!mediaQuery.success) {
    mediaError = mediaQuery.error.message;
  } else {
    media = mediaQuery.data.media;
  }

  return {
    query,

    mediaError,
    media,
  };
};
