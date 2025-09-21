<script lang="ts">
  import MediaCard from "$lib/components/MediaCard.svelte";
  import { Button } from "@nanoteck137/nano-ui";
  import { Plus } from "lucide-svelte";
  import { isRoleAdmin } from "$lib/utils";
  import { getApiClient, handleApiError } from "$lib";
  import toast from "svelte-5-french-toast";
  import { invalidateAll } from "$app/navigation";
  import MediaItemDropdown from "./MediaItemDropdown.svelte";
  import AddMediaItem from "./AddMediaItem.svelte";
  import Spacer from "$lib/components/Spacer.svelte";

  const { data } = $props();
  const apiClient = getApiClient();

  let openAddMediaModal = $state(false);
</script>

<div class="flex items-center justify-between">
  <h2 class="text-bold text-xl">
    Collection Items
    {#if isRoleAdmin(data.user?.role)}
      <Button
        variant="ghost"
        size="icon"
        onclick={() => {
          openAddMediaModal = true;
        }}
      >
        <Plus />
      </Button>
    {/if}
  </h2>
  <p class="text-sm">{data.items.length} item(s)</p>
</div>

<Spacer size="md" />

<div
  class="grid grid-cols-[repeat(auto-fit,minmax(240px,1fr))] items-center justify-items-center gap-6"
>
  {#each data.items as item}
    <div class="relative">
      <MediaCard
        href="/media/{item.mediaId}"
        coverUrl={item.coverUrl}
        title={item.collectionName}
        startDate={item.startDate}
        partCount={item.partCount}
        score={item.score}
        userList={item.user?.list ?? null}
      />

      <MediaItemDropdown {item} />
    </div>
  {/each}
</div>

<AddMediaItem
  bind:open={openAddMediaModal}
  itemIds={data.items.map((i) => i.mediaId)}
  onResult={async (results) => {
    for (const item of results) {
      const res = await apiClient.addCollectionItem(data.collection.id, {
        mediaId: item.id,
        name: item.name,
        searchSlug: item.name,
        position: item.position,
      });
      if (!res.success) {
        return handleApiError(res.error);
      }
    }

    toast.success("Successfully added new media items");
    invalidateAll();
  }}
/>
