<script lang="ts">
  const { data } = $props();
</script>

{#each data.parts as part}
  <p>{part.index}. {part.name}</p>
{/each}
