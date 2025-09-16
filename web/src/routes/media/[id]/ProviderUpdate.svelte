<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import type { ProviderValue } from "$lib/api/types";
  import ProviderUpdateModal from "./ProviderUpdateModal.svelte";

  type Props = {
    mediaId: string;
    provider: ProviderValue;
  };

  const { mediaId, provider }: Props = $props();
  let updateModalOpen = $state(false);
</script>

<Button
  onclick={() => {
    updateModalOpen = true;
  }}
>
  {provider.displayName}
</Button>

<ProviderUpdateModal bind:open={updateModalOpen} {mediaId} {provider} />
