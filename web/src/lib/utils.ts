import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

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
