<script lang="ts">
  import type { UserList } from "$lib/types";
  import { cn, userListClass } from "$lib/utils";
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

  const name = $derived(toName());
  const extraClass = $derived(userListClass(list));
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
