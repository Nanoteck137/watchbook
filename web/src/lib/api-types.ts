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

export const collectionTypes = [
  { label: "Unknown", value: "unknown" },
  { label: "Series", value: "series" },
  { label: "Anime", value: "anime" },
] as const;
export type CollectionType = (typeof collectionTypes)[number]["value"];
export const CollectionTypeEnum = z.enum(
  collectionTypes.map((f) => f.value) as [CollectionType, ...CollectionType[]],
);
