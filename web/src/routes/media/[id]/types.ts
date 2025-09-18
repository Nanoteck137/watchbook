import { z } from "zod";

export const mediaUserLists = [
  { label: "In Progress", value: "in-progress" },
  { label: "Completed", value: "completed" },
  { label: "On Hold", value: "on-hold" },
  { label: "Dropped", value: "dropped" },
  { label: "Backlog", value: "backlog" },
] as const;
export type MediaUserList = (typeof mediaUserLists)[number]["value"];
export const MediaUserListEnum = z.enum(
  mediaUserLists.map((f) => f.value) as [MediaUserList, ...MediaUserList[]],
);
