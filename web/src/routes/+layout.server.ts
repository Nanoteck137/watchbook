import type { GetMe, Notification } from "$lib/api/types";
import { error } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals }) => {
  let user: GetMe | null = null;
  let notifications: Notification[] | null = null;
  if (locals.token) {
    const res = await locals.apiClient.getMe();
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    user = res.data;

    {
      const res = await locals.apiClient.getNotifications({
        query: { filter: "isRead <= 0" },
      });
      if (!res.success) {
        throw error(res.error.code, { message: res.error.message });
      }

      notifications = res.data.notifications;
    }
  }

  return {
    apiAddress: locals.apiAddress,
    userToken: locals.token,
    user,
    notifications,
  };
};
