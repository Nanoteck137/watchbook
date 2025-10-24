<script>
  import { getApiClient, handleApiError } from "$lib";
  import ConfirmBox from "$lib/components/ConfirmBox.svelte";
  import toast from "svelte-5-french-toast";
  import { goto } from "$app/navigation";
  import { Button } from "@nanoteck137/nano-ui";
  import { Trash } from "lucide-svelte";
  import ProviderUpdateButton from "../ProviderUpdateButton.svelte";

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
  Delete Collection
</Button>

{#each data.collection.providers as provider}
  <ProviderUpdateButton collectionId={data.collection.id} {provider} />
{/each}

<ConfirmBox
  bind:open={openDeleteModal}
  title="Delete Collection?"
  description="Are you sure you want to delete this collection? This action cannot be undone."
  confirmText="Delete"
  onResult={async () => {
    const res = await apiClient.deleteCollection(data.collection.id);
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully deleted collection");
    goto(`/media`, { invalidateAll: true });
  }}
/>
