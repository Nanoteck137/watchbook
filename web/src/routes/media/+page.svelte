<script lang="ts">
  import { goto } from "$app/navigation";
  import { getApiClient } from "$lib";
  import Spacer from "$lib/components/Spacer.svelte";
  import { Pagination } from "@nanoteck137/nano-ui";
  import { page } from "$app/stores";
  import MediaCard from "$lib/components/MediaCard.svelte";

  const { data } = $props();
  const apiClient = getApiClient();

  // NOTE(patrik):
  //  - Large image size: 225x318 (h-80 w-56)
  //  - Medium image size: 160x220 (h-56 w-40)
  //  - Small image size: 50x70 (h-20 w-14)
</script>

<Spacer size="sm" />

<div class="mb-6 flex space-x-4 border-b border-gray-700 text-gray-400">
  <button
    id="tabBasic"
    class="border-b-2 border-blue-500 px-4 py-2 font-semibold text-blue-400 focus:outline-none"
    type="button"
  >
    Basic
  </button>
  <button
    id="tabAdvanced"
    class="border-b-2 border-transparent px-4 py-2 hover:text-blue-400 focus:outline-none"
    type="button"
  >
    Advanced
  </button>
</div>

<div class="" id="basicSection">
  <div
    class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
  >
    <div class="flex-grow">
      <label for="search" class="sr-only">Search Collections</label>
      <input
        type="search"
        id="search"
        name="search"
        placeholder="Search collections..."
        class="w-full rounded-md bg-gray-800 px-4 py-2 text-gray-200 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 sm:w-96"
        autocomplete="off"
      />
    </div>

    <div>
      <label for="sort" class="sr-only">Sort Collections</label>
      <select
        id="sort"
        name="sort"
        class="w-full rounded-md bg-gray-800 px-4 py-2 text-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-500 sm:w-auto"
      >
        <option value="title-asc">Title: A to Z</option>
        <option value="title-desc">Title: Z to A</option>
      </select>
    </div>
  </div>

  <!-- Type Filters -->
  <div class="mb-8 flex flex-wrap gap-4 text-gray-300">
    <label class="inline-flex cursor-pointer items-center">
      <input type="checkbox" class="type-filter" value="anime" checked />
      <span class="ml-2 select-none">Anime</span>
    </label>

    <label class="inline-flex cursor-pointer items-center">
      <input type="checkbox" class="type-filter" value="tv" checked />
      <span class="ml-2 select-none">TV</span>
    </label>

    <label class="inline-flex cursor-pointer items-center">
      <input type="checkbox" class="type-filter" value="game" checked />
      <span class="ml-2 select-none">Game</span>
    </label>

    <label class="inline-flex cursor-pointer items-center">
      <input type="checkbox" class="type-filter" value="more" checked />
      <span class="ml-2 select-none">More</span>
    </label>
  </div>
</div>

<div id="advancedSection" class="hidden">
  <div class="mb-6">
    <label for="customFilter" class="mb-2 block font-semibold text-gray-300">
      Advanced Filter
    </label>

    <input
      type="text"
      id="customFilter"
      placeholder="Enter custom filter language..."
      class="w-full rounded-md bg-gray-800 px-4 py-2 text-gray-200 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
  </div>

  <div>
    <label for="customSort" class="mb-2 block font-semibold text-gray-300">
      Advanced Sort
    </label>

    <input
      type="text"
      id="customSort"
      placeholder="Enter custom sort language..."
      class="w-full rounded-md bg-gray-800 px-4 py-2 text-gray-200 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
  </div>
</div>

<Spacer />

<p class="text-gray-400">Total: {data.page.totalItems}</p>

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

<Pagination.Root
  page={data.page.page + 1}
  count={data.page.totalItems}
  perPage={data.page.perPage}
  siblingCount={1}
  onPageChange={(p) => {
    const query = $page.url.searchParams;
    query.set("page", (p - 1).toString());

    goto(`?${query.toString()}`, { invalidateAll: true, keepFocus: true });
  }}
>
  {#snippet children({ pages, currentPage })}
    <Pagination.Content>
      <Pagination.Item>
        <Pagination.PrevButton />
      </Pagination.Item>
      {#each pages as page (page.key)}
        {#if page.type === "ellipsis"}
          <Pagination.Item>
            <Pagination.Ellipsis />
          </Pagination.Item>
        {:else}
          <Pagination.Item>
            <Pagination.Link
              href="?page={page.value}"
              {page}
              isActive={currentPage === page.value}
            >
              {page.value}
            </Pagination.Link>
          </Pagination.Item>
        {/if}
      {/each}
      <Pagination.Item>
        <Pagination.NextButton />
      </Pagination.Item>
    </Pagination.Content>
  {/snippet}
</Pagination.Root>
