<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import AddMediaItem from "./AddMediaItem.svelte";
  import { getApiClient, handleApiError } from "$lib";
  import toast from "svelte-5-french-toast";
  import { invalidateAll } from "$app/navigation";

  const { data } = $props();
  const apiClient = getApiClient();

  let openAddItemModal = $state(false);
</script>

<p>Folder: {data.folder.name}</p>

<Button
  onclick={() => {
    openAddItemModal = true;
  }}>Add Item</Button
>

{#each data.items as item}
  <div>
    <p>{item.position} - {item.title}</p>
    <Button
      onclick={async () => {
        const res = await apiClient.moveFolderItem(
          data.folder.id,
          item.mediaId,
          (item.position - 1).toString(),
        );
        if (!res.success) {
          return handleApiError(res.error);
        }

        toast.success("Successfully removed folder item");
        invalidateAll();
      }}>Move Up</Button
    >
    <Button
      onclick={async () => {
        const res = await apiClient.moveFolderItem(
          data.folder.id,
          item.mediaId,
          (item.position + 1).toString(),
        );
        if (!res.success) {
          return handleApiError(res.error);
        }

        toast.success("Successfully removed folder item");
        invalidateAll();
      }}>Move Down</Button
    >
    <Button
      variant="destructive"
      onclick={async () => {
        const res = await apiClient.removeFolderItem(
          data.folder.id,
          item.mediaId,
        );
        if (!res.success) {
          return handleApiError(res.error);
        }

        toast.success("Successfully removed folder item");
        invalidateAll();
      }}>Delete</Button
    >
  </div>
{/each}

<AddMediaItem
  bind:open={openAddItemModal}
  folderId={data.folder.id}
  itemIds={data.items.map((i) => i.mediaId)}
/>
