<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import type { ProviderValue } from "$lib/api/types";
  import ProviderUpdateModal from "./ProviderUpdateModal.svelte";

  type Props = {
    collectionId: string;
    provider: ProviderValue;
  };

  const { collectionId, provider }: Props = $props();
  let updateModalOpen = $state(false);
</script>

<Button
  onclick={() => {
    updateModalOpen = true;
  }}
>
  {provider.displayName}
</Button>

<ProviderUpdateModal bind:open={updateModalOpen} {collectionId} {provider} />
