import {
  MediaRatingEnum,
  MediaStatusEnum,
  MediaTypeEnum,
} from "$lib/api-types";
import { z } from "zod";

export const sortTypes = [
  { label: "Title (A-Z)", value: "title-a-z" },
  { label: "Title (Z-A)", value: "title-z-a" },
  { label: "Score (High–Low)", value: "score-high" },
  { label: "Score (Low–High)", value: "score-low" },
] as const;
export type SortType = (typeof sortTypes)[number]["value"];
export const SortTypeEnum = z.enum(
  sortTypes.map((f) => f.value) as [SortType, ...SortType[]],
);

export const defaultSort: SortType = "title-a-z";

export const FullFilter = z.object({
  query: z.string(),
  filters: z.object({
    type: z.array(MediaTypeEnum),
    status: z.array(MediaStatusEnum),
    rating: z.array(MediaRatingEnum),
  }),
  excludes: z.object({
    type: z.array(MediaTypeEnum),
    status: z.array(MediaStatusEnum),
    rating: z.array(MediaRatingEnum),
  }),
  sort: SortTypeEnum.default(defaultSort),
});
export type FullFilter = z.infer<typeof FullFilter>;
