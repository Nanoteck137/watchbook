<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import { cn } from "$lib/utils";
  import { buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import { EllipsisVertical, Pencil, Trash } from "lucide-svelte";
  import toast from "svelte-5-french-toast";
  import EditMediaItem from "./EditMediaItem.svelte";
  import type { CollectionItem } from "$lib/api/types";
  import RemoveConfirm from "./RemoveConfirm.svelte";

  type Props = {
    item: CollectionItem;
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
  name={item.collectionName}
  searchSlug={item.searchSlug}
  order={item.order}
  onResult={async (data) => {
    const res = await apiClient.editCollectionItem(
      item.collectionId,
      item.mediaId,
      data,
    );
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully updated collection item");
    invalidateAll();
  }}
></EditMediaItem>

<RemoveConfirm
  bind:open={removeModalOpen}
  onResult={async () => {
    const res = await apiClient.removeCollectionItem(
      item.collectionId,
      item.mediaId,
    );
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully removed collection item");
    invalidateAll();
  }}
/>
