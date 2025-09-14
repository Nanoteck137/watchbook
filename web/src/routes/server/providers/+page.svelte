<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import ProviderSearchModal from "./ProviderSearchModal.svelte";

  const { data } = $props();
</script>

<div class="flex flex-col gap-2">
  {#each data.providers as provider}
    <div class="flex flex-col items-start gap-1">
      <p>{provider.name}</p>
      <ProviderSearchModal providerName={provider.name} onResult={() => {}}>
        <Button>Search</Button>
      </ProviderSearchModal>
    </div>
  {/each}
</div>
