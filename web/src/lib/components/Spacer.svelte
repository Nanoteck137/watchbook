<script lang="ts">
  export type Type = "vertical" | "horizontal";
  export type Size = "xs" | "sm" | "md" | "lg";

  export type Props = {
    horizontal?: boolean;
    size?: Size;
  };

  const mapping: Record<Type, Record<Size, string>> = {
    horizontal: {
      xs: "w-1",
      sm: "w-2",
      md: "w-4",
      lg: "w-6",
    },
    vertical: {
      xs: "h-1",
      sm: "h-2",
      md: "h-4",
      lg: "h-6",
    },
  };

  const { horizontal, size }: Props = $props();
</script>

<div
  class={mapping[!!horizontal ? "horizontal" : "vertical"][size ?? "sm"]}
></div>
