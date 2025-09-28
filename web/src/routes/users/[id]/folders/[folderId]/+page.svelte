<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import AddMediaItem from "./AddMediaItem.svelte";
  import Item from "./Item.svelte";
  import { type MediaType } from "$lib/api-types";
  import { isRoleAdmin } from "$lib/utils";
  import { Plus, Settings } from "lucide-svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import EditFolderModal from "./EditFolderModal.svelte";

  const { data } = $props();

  let openAddItemModal = $state(false);
  let openEditFolderModal = $state(false);
</script>

<Spacer size="md" />

<div class="flex items-center gap-2">
  <p class="text-2xl font-bold">{data.folder.name}</p>
  <Button
    variant="ghost"
    size="icon"
    onclick={() => {
      openEditFolderModal = true;
    }}
  >
    <Settings />
  </Button>
</div>

<Spacer size="md" />

<div class="flex items-center justify-between">
  <h2 class="text-bold text-xl">
    Media
    {#if isRoleAdmin(data.user?.role)}
      <Button
        variant="ghost"
        size="icon"
        onclick={() => {
          openAddItemModal = true;
        }}
      >
        <Plus />
      </Button>
    {/if}
  </h2>
  <p class="text-sm">{data.items.length} item(s)</p>
</div>

<Spacer size="md" />

{#each data.items as item, i}
  <Item
    folderId={data.folder.id}
    mediaId={item.mediaId}
    position={item.position}
    index={i}
    numItems={data.items.length}
    title={item.title}
    type={item.mediaType as MediaType}
    coverUrl={item.coverUrl}
    startDate={item.startDate ?? undefined}
  />
{/each}

<AddMediaItem
  bind:open={openAddItemModal}
  folderId={data.folder.id}
  itemIds={data.items.map((i) => i.mediaId)}
/>

<EditFolderModal bind:open={openEditFolderModal} folder={data.folder} />
