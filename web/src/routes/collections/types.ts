import { z } from "zod";

export const collectionTypes = [
  { label: "Unknown", value: "unknown" },
  { label: "Series", value: "series" },
  { label: "Anime", value: "anime" },
] as const;
export type CollectionType = (typeof collectionTypes)[number]["value"];
export const CollectionTypeEnum = z.enum(
  collectionTypes.map((f) => f.value) as [CollectionType, ...CollectionType[]],
);

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
