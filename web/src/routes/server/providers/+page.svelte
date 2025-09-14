<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import ProviderSearchModal from "./ProviderSearchModal.svelte";
  import { getApiClient, handleApiError } from "$lib";
  import toast from "svelte-5-french-toast";
  import { invalidateAll } from "$app/navigation";

  const { data } = $props();
  const apiClient = getApiClient();
</script>

<div class="flex flex-col gap-2">
  {#each data.providers as provider}
    <div class="flex flex-col items-start gap-1">
      <p>{provider.name}</p>
      <ProviderSearchModal
        providerName={provider.name}
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
      >
        <Button>Search</Button>
      </ProviderSearchModal>
    </div>
  {/each}
</div>
