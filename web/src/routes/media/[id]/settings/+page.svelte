<script>
  import ProviderUpdateButton from "../ProviderUpdateButton.svelte";

  const { data } = $props();
</script>

{#each data.media.providers as provider}
  <ProviderUpdateButton mediaId={data.media.id} {provider} />
{/each}
