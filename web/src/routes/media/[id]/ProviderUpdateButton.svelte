<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import type { ProviderValue } from "$lib/api/types";
  import ProviderUpdateSettingsModal from "./ProviderUpdateSettingsModal.svelte";

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

<ProviderUpdateSettingsModal
  bind:open={updateModalOpen}
  {mediaId}
  {provider}
/>
