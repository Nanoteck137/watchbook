<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import CollectionCard from "$lib/components/CollectionCard.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { Pagination } from "@nanoteck137/nano-ui";

  const { data } = $props();
</script>

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

<!-- Basic Filters Section -->
<div class="hidden" id="basicSection">
  <!-- Search, Sort & Filter Controls -->
  <div
    class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
  >
    <!-- Search Bar -->
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

    <!-- Sort Select -->
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

<!-- Advanced Filters Section (hidden by default) -->
<div id="advancedSection" class="">
  <div class="mb-6">
    <label for="customFilter" class="mb-2 block font-semibold text-gray-300"
      >Advanced Filter</label
    >
    <input
      type="text"
      id="customFilter"
      placeholder="Enter custom filter language..."
      class="w-full rounded-md bg-gray-800 px-4 py-2 text-gray-200 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
  </div>

  <div>
    <label for="customSort" class="mb-2 block font-semibold text-gray-300"
      >Advanced Sort</label
    >
    <input
      type="text"
      id="customSort"
      placeholder="Enter custom sort language..."
      class="w-full rounded-md bg-gray-800 px-4 py-2 text-gray-200 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
  </div>
</div>

<p id="totalCount" class="mb-4 text-gray-400">Total collections: 0</p>

<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
  <div
    class="grid grid-cols-[repeat(auto-fit,minmax(280px,1fr))] items-center justify-items-center gap-6"
  >
    {#each data.collections as collection}
      <CollectionCard
        href="/collections/{collection.id}"
        name={collection.name}
        coverUrl={collection.coverUrl}
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
