<script lang="ts">
  import ReleaseMediaCard from "$lib/components/ReleaseMediaCard.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import StandardPagination from "$lib/components/StandardPagination.svelte";
  import Filter from "./Filter.svelte";

  const { data } = $props();
</script>

<Filter fullFilter={data.filter} />

<Spacer size="md" />

<div class="flex items-center justify-between">
  <h2 class="text-bold text-xl">Media</h2>
  <p class="text-sm">{data.page.totalItems} item(s)</p>
</div>

<Spacer size="md" />

<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
  <div
    class="grid grid-cols-[repeat(auto-fit,minmax(240px,1fr))] items-center justify-items-center gap-6"
  >
    {#each data.media as media}
      <ReleaseMediaCard
        href="/media/{media.id}"
        title={media.title}
        coverUrl={media.coverUrl}
        userList={media.user?.list ?? null}
        release={media.release!}
      />
    {/each}
  </div>
</div>

<Spacer size="md" />

<StandardPagination pageData={data.page} />
