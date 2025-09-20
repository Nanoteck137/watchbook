import { z } from "zod";
import { MediaUserListEnum } from "../../../media/[id]/types";
import { MediaStatusEnum, MediaTypeEnum } from "../../../media/types";

export const sortTypes = [
  { label: "Title (A-Z)", value: "title-a-z" },
  { label: "Title (Z-A)", value: "title-z-a" },
  { label: "Score (High–Low)", value: "score-high" },
  { label: "Score (Low–High)", value: "score-low" },
  { label: "User Score (High–Low)", value: "user-score-high" },
  { label: "User Score (Low–High)", value: "user-score-low" },
] as const;
export type SortType = (typeof sortTypes)[number]["value"];
export const SortTypeEnum = z.enum(
  sortTypes.map((f) => f.value) as [SortType, ...SortType[]],
);

export const defaultSort: SortType = "title-a-z";

export const FullFilter = z.object({
  list: MediaUserListEnum.or(z.literal("")),
  types: z.array(MediaTypeEnum),
  status: z.array(MediaStatusEnum),
  sort: SortTypeEnum.default(defaultSort),
});
export type FullFilter = z.infer<typeof FullFilter>;
