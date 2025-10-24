<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import { cn } from "$lib/utils";
  import { buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import { EllipsisVertical, Pencil, Trash } from "lucide-svelte";
  import ConfirmBox from "$lib/components/ConfirmBox.svelte";
  import type { ShowSeasonItem } from "$lib/api/types";
  import toast from "svelte-5-french-toast";
  import { invalidateAll } from "$app/navigation";
  import EditMediaItem from "./EditMediaItem.svelte";

  type Props = {
    item: ShowSeasonItem;
  };

  const { item }: Props = $props();

  const apiClient = getApiClient();

  let editModalOpen = $state(false);
  let removeModalOpen = $state(false);
</script>

<DropdownMenu.Root>
  <DropdownMenu.Trigger
    class={cn(
      "absolute right-4 top-4 z-20",
      buttonVariants({ variant: "secondary", size: "icon" }),
    )}
  >
    <EllipsisVertical />
  </DropdownMenu.Trigger>
  <DropdownMenu.Content class="w-40" align="end">
    <DropdownMenu.Group>
      <DropdownMenu.Item
        onclick={() => {
          editModalOpen = true;
        }}
      >
        <Pencil />
        Edit
      </DropdownMenu.Item>

      <DropdownMenu.Item
        onclick={() => {
          removeModalOpen = true;
        }}
      >
        <Trash />
        Remove
      </DropdownMenu.Item>
    </DropdownMenu.Group>
  </DropdownMenu.Content>
</DropdownMenu.Root>

<EditMediaItem
  bind:open={editModalOpen}
  position={item.position}
  onResult={async (data) => {
    const res = await apiClient.editShowSeasonItem(
      item.showId,
      item.showSeasonNum.toString(),
      item.mediaId,
      data,
    );
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully updated collection item");
    invalidateAll();
  }}
/>

<ConfirmBox
  bind:open={removeModalOpen}
  title="Remove Item?"
  description="Are you sure you want to remove this media item? This action cannot be undone."
  confirmText="Remove Item"
  onResult={async () => {
    const res = await apiClient.removeShowSeasonItem(
      item.showId,
      item.showSeasonNum.toString(),
      item.mediaId,
    );
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully removed collection item");
    invalidateAll();
  }}
/>
