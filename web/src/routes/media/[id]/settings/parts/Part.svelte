<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import { MediaPart } from "$lib/api/types";
  import ConfirmBox from "$lib/components/ConfirmBox.svelte";
  import { Button } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import EditPartModal from "./EditPartModal.svelte";

  type Props = {
    part: MediaPart;
  };

  const { part }: Props = $props();
  const apiClient = getApiClient();

  let openEditModal = $state(false);
  let openRemoveModal = $state(false);
</script>

<div class="flex items-center justify-between border-b pb-4">
  <div>
    <p class="font-medium">
      Part {part.index} — {part.name}
    </p>
    <p class="text-xs text-muted-foreground">Air Date: 2013-04-07</p>
  </div>

  <div class="flex items-center gap-2">
    <Button
      variant="outline"
      size="sm"
      onclick={() => {
        openEditModal = true;
      }}
    >
      Edit
    </Button>
    <Button
      variant="destructive"
      size="sm"
      onclick={() => {
        openRemoveModal = true;
      }}
    >
      Delete
    </Button>
  </div>
</div>

<EditPartModal bind:open={openEditModal} {part} />

<ConfirmBox
  bind:open={openRemoveModal}
  title="Remove Part?"
  description="Are you sure you want to delete this part? This action cannot be undone."
  confirmText="Remove"
  onResult={async () => {
    const res = await apiClient.removePart(
      part.mediaId,
      part.index.toString(),
    );
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully removed part");
    invalidateAll();
  }}
/>
