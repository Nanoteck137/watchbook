<script lang="ts">
  import Spacer from "$lib/components/Spacer.svelte";
  import MediaCard from "$lib/components/MediaCard.svelte";
  import Filter from "./Filter.svelte";
  import StandardPagination from "$lib/components/StandardPagination.svelte";

  const { data } = $props();

  // NOTE(patrik):
  //  - Large image size: 225x318 (h-80 w-56)
  //  - Medium image size: 160x220 (h-56 w-40)
  //  - Small image size: 50x70 (h-20 w-14)
</script>

<Spacer size="sm" />

<Filter fullFilter={data.filter} />

<Spacer />

<p>Total: {data.page.totalItems}</p>

<Spacer />

<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
  <div
    class="grid grid-cols-[repeat(auto-fit,minmax(240px,1fr))] items-center justify-items-center gap-6"
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
</div>

<Spacer size="sm" />

<StandardPagination pageData={data.page} />
