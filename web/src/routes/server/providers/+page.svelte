<script lang="ts">
  const { data } = $props();
</script>

{#each data.providers as provider}
  <p>{provider.name}</p>
{/each}
