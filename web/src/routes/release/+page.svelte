<script lang="ts">
  import Badge from "$lib/components/Badge.svelte";
  import Image from "$lib/components/Image.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { parseUserList } from "$lib/types";
  import { clamp, formatTimeDiff, getTimeDifference } from "$lib/utils";

  const { data } = $props();

  const now = new Date();
</script>

<Spacer />

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
      <!-- TODO(patrik): Fix -->
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
      <!-- TODO(patrik): Fix -->
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
      {@const percent =
        media.release!.currentPart / media.release!.numExpectedParts}
      {@const time = getTimeDifference(
        now,
        new Date(media.release!.nextAiring ?? ""),
      )}
      <!-- svelte-ignore a11y_missing_attribute -->
      <!-- svelte-ignore a11y_invalid_attribute -->
      <a
        class="group relative aspect-[75/106] max-w-[240px] transform cursor-pointer overflow-hidden rounded-3xl bg-gray-800 shadow-md transition-transform duration-300 hover:-translate-y-1 hover:shadow-lg"
        href="/media/{media.id}"
      >
        <!-- Badge -->
        {#if !!media.user?.list}
          {@const list = parseUserList(media.user.list)}
          {#if list}
            <Badge class="absolute left-2 top-2 z-10" {list} />
          {/if}
        {/if}

        <div
          class="flex min-w-[240px] items-center justify-center bg-gray-800 text-gray-400"
        >
          <Image
            src={media.coverUrl}
            alt="Cover image"
            class="aspect-[75/106] h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
          />
        </div>

        <div
          class="absolute bottom-0 left-0 right-0 rounded-t-lg bg-black/80 p-3 text-center backdrop-blur-md"
        >
          <h2
            class="line-clamp-2 text-ellipsis text-base font-semibold"
            title={media.title}
          >
            {media.title}
          </h2>

          <Spacer size="xs" />

          <div class="flex flex-col gap-2">
            <div class="flex flex-col gap-1">
              <p class="text-xs text-gray-400">
                {media.release!.currentPart} / {media.release!
                  .numExpectedParts} eps
              </p>

              <div class="h-1 w-full rounded-full bg-gray-700">
                <div
                  class="h-1 rounded-full bg-blue-500"
                  style={`width: ${clamp(percent * 100, 0, 100).toFixed(0)}%;`}
                ></div>
              </div>
            </div>

            {#if media.release!.status !== "completed"}
              <div class="flex items-center gap-1">
                <p class="text-xs text-gray-300">Next in:</p>
                <p class="text-sm font-bold text-blue-400">
                  {formatTimeDiff(time)}
                </p>
              </div>
            {/if}

            <p>{media.release!.status}</p>
          </div>
        </div>
      </a>
    {/each}
  </div>
</div>
