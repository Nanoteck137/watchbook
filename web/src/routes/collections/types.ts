import { z } from "zod";

export const collectionTypes = [
  { label: "Unknown", value: "unknown" },
  { label: "Series", value: "series" },
  { label: "Anime", value: "anime" },
] as const;
export type CollectionType = (typeof collectionTypes)[number]["value"];

export const sortTypes = [
  { label: "Name (A-Z)", value: "name-a-z" },
  { label: "Name (Z-A)", value: "name-z-a" },
] as const;
export type SortType = (typeof sortTypes)[number]["value"];

export const FullFilter = z.object({
  query: z.string(),
  filter: z
    .enum(
      collectionTypes.map((f) => f.value) as [
        CollectionType,
        ...CollectionType[],
      ],
    )
    .or(z.literal(""))
    .default(""),
  excludes: z.array(
    z.enum(
      collectionTypes.map((f) => f.value) as [
        CollectionType,
        ...CollectionType[],
      ],
    ),
  ),
  sort: z
    .enum(sortTypes.map((f) => f.value) as [SortType, ...SortType[]])
    .default("name-a-z"),
});
export type FullFilter = z.infer<typeof FullFilter>;
