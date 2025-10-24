<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import { Plus } from "lucide-svelte";
  import { isRoleAdmin } from "$lib/utils";
  import { getApiClient } from "$lib";
  import Spacer from "$lib/components/Spacer.svelte";
  import AddSeasonModal from "./AddSeasonModal.svelte";
  import ShowSeasonItem from "./ShowSeasonItem.svelte";

  const { data } = $props();
  const apiClient = getApiClient();

  let openAddSeasonModal = $state(false);
</script>

<div class="flex items-center justify-between">
  <h2 class="text-bold text-xl">
    Collection Items
    {#if isRoleAdmin(data.user?.role)}
      <Button
        variant="ghost"
        size="icon"
        onclick={() => {
          openAddSeasonModal = true;
        }}
      >
        <Plus />
      </Button>
    {/if}
  </h2>
  <!-- <p class="text-sm">{data.items.length} item(s)</p> -->
</div>

<Spacer size="md" />

<div class="flex flex-col">
  {#each data.seasons as season}
    <ShowSeasonItem userRole={data.user?.role} {season} />
  {/each}
</div>

<!-- <div
  class="grid grid-cols-[repeat(auto-fit,minmax(240px,1fr))] items-center justify-items-center gap-6"
>
  {#each data.seasons as item}
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
</div> -->

<AddSeasonModal bind:open={openAddSeasonModal} showId={data.show.id} />
