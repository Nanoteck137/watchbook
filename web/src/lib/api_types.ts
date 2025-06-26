import { z } from "zod";

export const MediaType = z.enum([
  "unknown",
  "season",
  "movie",
  "anime-season",
  "anime-movie",
]);
export type MediaType = z.infer<typeof MediaType>;

export function parseMediaType(t: unknown): MediaType {
  const res = MediaType.safeParse(t);
  if (!res.success) {
    return "unknown";
  }

  return res.data;
}

export const MediaStatus = z.enum([
  "unknown",
  "airing",
  "finished",
  "not-aired",
]);
export type MediaStatus = z.infer<typeof MediaStatus>;

export function parseMediaStatus(t: unknown): MediaStatus {
  const res = MediaStatus.safeParse(t);
  if (!res.success) {
    return "unknown";
  }

  return res.data;
}

export const MediaRating = z.enum([
  "unknown",
  "all-ages",
  "pg",
  "pg-13",
  "r-17",
  "r-mild-nudity",
  "r-hentai",
]);
export type MediaRating = z.infer<typeof MediaRating>;

export function parseMediaRating(t: unknown): MediaRating {
  const res = MediaRating.safeParse(t);
  if (!res.success) {
    return "unknown";
  }

  return res.data;
}
