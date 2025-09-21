<script>
  import { getApiClient, handleApiError } from "$lib";
  import ConfirmBox from "$lib/components/ConfirmBox.svelte";
  import toast from "svelte-5-french-toast";
  import ProviderUpdateButton from "../ProviderUpdateButton.svelte";
  import { goto } from "$app/navigation";
  import { Button } from "@nanoteck137/nano-ui";
  import { Trash } from "lucide-svelte";

  const { data } = $props();
  const apiClient = getApiClient();

  let openDeleteModal = $state(false);
</script>

<Button
  variant="destructive"
  onclick={() => {
    openDeleteModal = true;
  }}
>
  <Trash />
  Delete Media
</Button>

{#each data.media.providers as provider}
  <ProviderUpdateButton mediaId={data.media.id} {provider} />
{/each}

<ConfirmBox
  bind:open={openDeleteModal}
  title="Delete Media?"
  description="Are you sure you want to delete this media? This action cannot be undone."
  confirmText="Delete"
  onResult={async () => {
    const res = await apiClient.deleteMedia(data.media.id);
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully deleted media");
    goto(`/media`, { invalidateAll: true });
  }}
/>
