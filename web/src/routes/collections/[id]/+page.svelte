<script lang="ts">
  import MediaCard from "$lib/components/MediaCard.svelte";
  import { buttonVariants } from "@nanoteck137/nano-ui";
  import { Image } from "lucide-svelte";
  import ShowLogoModal from "./ShowLogoModal.svelte";
  import { cn } from "$lib/utils";

  const { data } = $props();
</script>

<div class="max-w-7xl px-4 text-gray-100 sm:px-6 lg:px-8">
  <div
    class="relative h-48 w-full overflow-hidden rounded-lg shadow-lg sm:h-60 md:h-72"
  >
    <img
      src={data.collection.bannerUrl}
      alt="Collection Banner"
      class="h-full w-full object-cover"
    />

    <div class="absolute inset-0 bg-black bg-opacity-40"></div>

    <div
      class="pointer-events-none absolute inset-0 hidden items-center justify-center md:flex"
    >
      <div
        class="w-32 rounded-lg border border-gray-700 bg-gray-900 bg-opacity-70 p-4 shadow-lg sm:w-40 md:w-48"
      >
        <img
          src={data.collection.logoUrl}
          alt="Collection Logo"
          class="mx-auto max-h-24 w-auto object-contain"
        />
      </div>
    </div>

    {#if data.collection.logoUrl}
      <ShowLogoModal
        class={cn(
          "absolute right-4 top-4 z-20 md:hidden",
          buttonVariants({ variant: "secondary", size: "icon" }),
        )}
        logoUrl={data.collection.logoUrl}
        onResult={() => {}}
      >
        <Image />
      </ShowLogoModal>
    {/if}
  </div>

  <div
    class="relative mx-2 -mt-16 flex flex-col items-center space-y-4 px-4 sm:-mt-20 sm:flex-row sm:items-start sm:space-x-6 sm:space-y-0 sm:px-0"
  >
    <div
      class="z-10 aspect-[75/106] w-40 min-w-40 flex-shrink-0 overflow-hidden rounded-lg border-4 border-gray-900 bg-gray-800 shadow-lg sm:w-48 sm:min-w-48"
    >
      <img
        src={data.collection.coverUrl}
        alt="Collection Cover"
        class="h-full w-full object-cover"
      />
    </div>

    <div
      class="flex max-w-xl flex-col justify-start pb-4 text-center sm:pt-20 sm:text-left"
    >
      <h1 class="text-2xl font-bold drop-shadow-lg sm:pt-2">
        {data.collection.name}
      </h1>
      <!-- <p class="mt-2 text-lg text-gray-300">
        All seasons and movies of the Fullmetal Alchemist franchise, including
        Brotherhood and the original series.
      </p>
      <p class="mt-4 text-sm italic text-gray-400">
        Updated regularly with new releases and extras.
      </p> -->
    </div>
  </div>

  <div class="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
    <div
      class="grid grid-cols-[repeat(auto-fit,minmax(240px,1fr))] items-center justify-items-center gap-6"
    >
      {#each data.items as item}
        <MediaCard
          href="/media/{item.mediaId}"
          coverUrl={item.coverUrl}
          title={item.collectionName}
          startDate={item.startDate}
          partCount={item.partCount}
          score={item.score}
          userList={item.user?.list ?? null}
        />
      {/each}
    </div>
  </div>
</div>
