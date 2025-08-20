import type { UserList } from "$lib/types";
import { clsx, type ClassValue } from "clsx";
import { twMerge, type ClassNameValue } from "tailwind-merge";

export function capitilize(s: string) {
  if (s.length === 0) return "";
  return s[0].toUpperCase() + s.substring(1);
}

export function formatTime(s: number) {
  const min = Math.floor(s / 60);
  const sec = Math.floor(s % 60);

  return `${min}:${sec.toString().padStart(2, "0")}`;
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function isRoleAdmin(role: string) {
  switch (role) {
    case "super_user":
    case "admin":
      return true;
    default:
      return false;
  }
}

export function getPagedQueryOptions(searchParams: URLSearchParams) {
  const query: Record<string, string> = {};
  const filter = searchParams.get("filter");
  if (filter) {
    query["filter"] = filter;
  }

  const sort = searchParams.get("sort");
  if (sort) {
    query["sort"] = sort;
  }

  const page = searchParams.get("page");
  if (page) {
    query["page"] = page;
  }

  return query;
}

export function pickTitle(entity: {
  title: string;
  titleEnglish: string | null;
}) {
  // return "Villainess Level 99: I May Be the Hidden Boss but I'm Not the Demon Lord's Lead Level Grinds to Success";
  // return "The Strongest Tank's Labyrinth Raids -A Tank with a Rare 9999 Resistance Skill Got Kicked from the Hero's Party-";
  // return "I've Somehow Gotten Stronger When I Improved My Farm-related Skills Is a Fantasy Parody Starring an OP Farm Boy";
  // return "Trapped in a Dating Sim: The World of Otome Games is Tough for Mobs Somehow Combines Otome Games With Mecha";
  // return "The Misfit Of Demon King Academy: History's Strongest Demon King Reincarnates and Goes To School With His Descendants Is an Amusing High School Power Trip The Misfit Of Demon King Academy: History's Strongest Demon King Reincarnates and Goes To School With His Descendants Is an Amusing High School Power Trip";

  if (entity.titleEnglish) return entity.titleEnglish;

  return entity.title;
}

export function userListClass(list: UserList): ClassNameValue {
  switch (list) {
    case "in-progress":
      return "bg-blue-500 text-white";
    case "completed":
      return "bg-green-600 text-white";
    case "dropped":
      return "bg-red-600 text-white";
    case "on-hold":
      return "bg-yellow-500 text-white";
    case "backlog":
      return "bg-gray-600 text-white";
  }
}
