<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import type { Provider } from "$lib/api/types";
  import toast from "svelte-5-french-toast";
  import ProviderSearchModal from "./ProviderSearchModal.svelte";
  import { invalidateAll } from "$app/navigation";
  import { Button } from "@nanoteck137/nano-ui";

  type Props = {
    provider: Provider;
  };

  const { provider }: Props = $props();
  const apiClient = getApiClient();

  let providerSearchModalOpen = $state(false);
</script>

<div class="flex flex-col items-start gap-1">
  <p>{provider.displayName}</p>

  <div class="flex gap-2">
    <Button
      disabled={!provider.supports.searchMedia}
      onclick={() => {
        providerSearchModalOpen = true;
      }}
    >
      Search
    </Button>

    <Button disabled={!provider.supports.searchMedia} onclick={() => {}}>
      Get Media
    </Button>

    <Button disabled={!provider.supports.searchCollection} onclick={() => {}}>
      Search Collection
    </Button>

    <Button disabled={!provider.supports.getCollection} onclick={() => {}}>
      Get Collection
    </Button>
  </div>
</div>

<ProviderSearchModal
  bind:open={providerSearchModalOpen}
  providerName={provider.name}
  providerDisplayName={provider.displayName}
  onResult={async (items) => {
    const res = await apiClient.providerImportMedia(provider.name, {
      ids: items.map((i) => i.providerId),
    });
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully imported media");
    invalidateAll();
  }}
/>
