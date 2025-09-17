<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import { cn } from "$lib/utils";
  import { buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import {
    Download,
    EllipsisVertical,
    Image,
    Pencil,
    Trash,
  } from "lucide-svelte";
  import type { Collection } from "$lib/api/types";
  import ConfirmBox from "$lib/components/ConfirmBox.svelte";
  import EditImagesModal from "./EditImagesModal.svelte";
  import toast from "svelte-5-french-toast";
  import { goto, invalidateAll } from "$app/navigation";
  import EditCollectionModal from "./EditCollectionModal.svelte";
  import ProviderUpdateModal from "./ProviderUpdateModal.svelte";

  type Props = {
    collection: Collection;
  };

  const { collection }: Props = $props();

  const apiClient = getApiClient();

  let openEditModal = $state(false);
  let openEditImagesModal = $state(false);
  let openDeleteModal = $state(false);
  let openProviderUpdateModal = $state(false);
</script>

<DropdownMenu.Root>
  <DropdownMenu.Trigger
    class={cn(
      "absolute right-2 top-2 z-20",
      buttonVariants({ variant: "secondary", size: "icon" }),
    )}
  >
    <EllipsisVertical />
  </DropdownMenu.Trigger>
  <DropdownMenu.Content class="w-40" align="center">
    <DropdownMenu.Group>
      <DropdownMenu.Item
        onclick={() => {
          openProviderUpdateModal = true;
        }}
      >
        <Download />
        Update
      </DropdownMenu.Item>

      <DropdownMenu.Item
        onclick={() => {
          openEditModal = true;
        }}
      >
        <Pencil />
        Edit
      </DropdownMenu.Item>

      <DropdownMenu.Item
        onclick={() => {
          openEditImagesModal = true;
        }}
      >
        <Image />
        Edit Images
      </DropdownMenu.Item>

      <DropdownMenu.Item
        onclick={() => {
          openDeleteModal = true;
        }}
      >
        <Trash />
        Delete
      </DropdownMenu.Item>
    </DropdownMenu.Group>
  </DropdownMenu.Content>
</DropdownMenu.Root>

<EditCollectionModal bind:open={openEditModal} {collection} />

<!-- <EditMediaItem
  bind:open={editModalOpen}
  name={item.collectionName}
  searchSlug={item.searchSlug}
  position={item.position}
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
/>

-->

<EditImagesModal
  bind:open={openEditImagesModal}
  collectionId={collection.id}
/>

<ConfirmBox
  bind:open={openDeleteModal}
  title="Delete Collection?"
  description="Are you sure you want to delete this collection? This action cannot be undone."
  confirmText="Remove"
  onResult={async () => {
    const res = await apiClient.deleteCollection(collection.id);
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully deleted collection");
    goto(`/collections`, { invalidateAll: true });
  }}
/>

<ProviderUpdateModal
  bind:open={openProviderUpdateModal}
  collectionId={collection.id}
  providers={collection.providers}
/>
