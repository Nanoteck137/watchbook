<script lang="ts">
  import Spacer from "$lib/components/Spacer.svelte";
  import MediaCard from "$lib/components/MediaCard.svelte";
  import Filter from "./Filter.svelte";
  import StandardPagination from "$lib/components/StandardPagination.svelte";
  import { Button } from "@nanoteck137/nano-ui";
  import { isRoleAdmin } from "$lib/utils";
  import { Plus } from "lucide-svelte";
  import NewMediaModal from "./NewMediaModal.svelte";

  const { data } = $props();

  let openNewMediaModal = $state(false);

  // NOTE(patrik):
  //  - Large image size: 225x318 (h-80 w-56)
  //  - Medium image size: 160x220 (h-56 w-40)
  //  - Small image size: 50x70 (h-20 w-14)
</script>

<Filter fullFilter={data.filter} />

<Spacer size="md" />

<div class="flex items-center justify-between">
  <h2 class="text-bold text-xl">
    Media
    {#if isRoleAdmin(data.user?.role)}
      <Button
        variant="ghost"
        size="icon"
        onclick={() => {
          openNewMediaModal = true;
        }}
      >
        <Plus />
      </Button>
    {/if}
  </h2>
  <p class="text-sm">{data.page.totalItems} item(s)</p>
</div>

<Spacer size="md" />

<div
  class="grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] items-center justify-items-center gap-6"
>
  {#each data.media as media}
    <MediaCard
      href="/media/{media.id}"
      title={media.title}
      coverUrl={media.coverUrl}
      startDate={media.startDate}
      partCount={media.partCount}
      score={media.score}
      userList={media.user?.list ?? null}
    />
  {/each}
</div>

<Spacer size="sm" />

<StandardPagination pageData={data.page} />

<NewMediaModal bind:open={openNewMediaModal} />
