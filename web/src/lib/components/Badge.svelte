<script lang="ts">
  import type { UserList } from "$lib/types";
  import { cn } from "$lib/utils";
  import type { ClassNameValue } from "tailwind-merge";

  type Props = {
    list: UserList;
    class?: string;
  };

  const { list, class: className }: Props = $props();

  function toName(): string {
    switch (list) {
      case "in-progress":
        return "In Progress";
      case "completed":
        return "Completed";
      case "dropped":
        return "Dropped";
      case "on-hold":
        return "On Hold";
      case "backlog":
        return "Backlog";
    }
  }

  function toClass(): ClassNameValue {
    switch (list) {
      case "in-progress":
        return "bg-purple-600 text-purple-100";
      case "completed":
        return "bg-green-600 text-green-100";
      case "dropped":
        return "bg-red-600 text-red-100";
      case "on-hold":
        return "bg-gray-600 text-gray-100";
      case "backlog":
        return "bg-blue-600 text-blue-100";
    }
  }

  const name = $derived(toName());
  const extraClass = $derived(toClass());
</script>

<span
  class={cn(
    "inline-block select-none rounded-full px-3 py-1 text-xs font-semibold",
    extraClass,
    className,
  )}
>
  {name}
</span>
