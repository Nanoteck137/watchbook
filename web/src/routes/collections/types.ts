import { CollectionTypeEnum } from "$lib/api-types";
import { z } from "zod";

export const sortTypes = [
  { label: "Name (A-Z)", value: "name-a-z" },
  { label: "Name (Z-A)", value: "name-z-a" },
] as const;
export type SortType = (typeof sortTypes)[number]["value"];
export const SortTypeEnum = z.enum(
  sortTypes.map((f) => f.value) as [SortType, ...SortType[]],
);

export const defaultSort: SortType = "name-a-z";

export const FullFilter = z.object({
  query: z.string(),
  filters: z.object({
    type: z.array(CollectionTypeEnum),
  }),
  excludes: z.object({
    type: z.array(CollectionTypeEnum),
  }),
  sort: SortTypeEnum.default(defaultSort),
});
export type FullFilter = z.infer<typeof FullFilter>;
