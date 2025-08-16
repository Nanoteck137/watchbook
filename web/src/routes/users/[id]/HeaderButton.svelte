<script lang="ts">
  import { cn } from "$lib/utils";

  type Props = {
    href: string;
    name: string;
    active?: boolean;
    class?: string;
    stripClass?: string;
  };

  const { href, name, active, class: className, stripClass }: Props = $props();
</script>

<a {href} class={cn("group flex flex-col items-center", className)}>
  <p class={cn("mb-1 text-xs", active ? "font-semibold" : "")}>{name}</p>
  <div
    class={cn(
      "h-1 w-full rounded bg-white transition-all duration-200 group-hover:h-2 group-hover:brightness-110",
      active ? "h-2 brightness-125" : "",
      stripClass,
    )}
  ></div>
</a>
