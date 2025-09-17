import { z } from "zod";

export const mediaTypes = [
  { label: "Unknown", value: "unknown" },
  { label: "TV", value: "tv" },
  { label: "Movie", value: "movie" },
  { label: "Anime Season", value: "anime-season" },
  { label: "Anime Movie", value: "anime-movie" },
  { label: "Game", value: "game" },
  { label: "Manga", value: "manga" },
  { label: "Comic", value: "comic" },
] as const;
export type MediaType = (typeof mediaTypes)[number]["value"];
export const MediaTypeEnum = z.enum(
  mediaTypes.map((f) => f.value) as [MediaType, ...MediaType[]],
);

export const mediaStatus = [
  { label: "Unknown", value: "unknown" },
  { label: "On-Going", value: "ongoing" },
  { label: "Completed", value: "completed" },
  { label: "Upcoming", value: "upcoming" },
] as const;
export type MediaStatus = (typeof mediaStatus)[number]["value"];
export const MediaStatusEnum = z.enum(
  mediaStatus.map((f) => f.value) as [MediaStatus, ...MediaStatus[]],
);

export const mediaRatings = [
  { label: "Unknown", value: "unknown" },
  { label: "All Ages", value: "all-ages" },
  { label: "PG", value: "pg" },
  { label: "PG-13", value: "pg-13" },
  { label: "R-17", value: "r-17" },
  { label: "R-Mild-Nudity", value: "r-mild-nudity" },
  { label: "R-Hentai", value: "r-hentai" },
] as const;
export type MediaRating = (typeof mediaRatings)[number]["value"];
export const MediaRatingEnum = z.enum(
  mediaRatings.map((f) => f.value) as [MediaRating, ...MediaRating[]],
);

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
  filter: z.object({
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
