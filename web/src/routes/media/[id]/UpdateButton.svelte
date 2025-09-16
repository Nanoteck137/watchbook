<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import UpdateModal from "./UpdateModal.svelte";
  import type { MediaProvider } from "$lib/api/types";
  import { getApiClient, handleApiError } from "$lib";
  import toast from "svelte-5-french-toast";
  import { invalidateAll } from "$app/navigation";

  type Props = {
    mediaId: string;
    provider: MediaProvider;
  };

  const { mediaId, provider }: Props = $props();
  const apiClient = getApiClient();

  let updateModalOpen = $state(false);
</script>

<Button
  onclick={() => {
    updateModalOpen = true;
  }}
>
  {provider.displayName}
</Button>

<UpdateModal
  bind:open={updateModalOpen}
  providerDisplayName={provider.displayName}
  onResult={async (data) => {
    const res = await apiClient.providerUpdateMedia(
      provider.name,
      mediaId,
      data,
    );
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully update media");
    invalidateAll();
  }}
/>
