<script lang="ts">
  import Part from "./Part.svelte";

  const { data } = $props();
</script>

<div class="flex flex-col gap-4">
  <div>
    <h1 class="text-2xl font-bold">Part Settings</h1>
    <p class="text-sm text-muted-foreground">
      Manage, edit, or delete episodes for this media entry.
    </p>
  </div>

  <section class="flex flex-col gap-4">
    <h2 class="text-lg font-semibold">Parts</h2>

    {#each data.parts as part}
      <Part {part} />
    {/each}
  </section>
</div>
