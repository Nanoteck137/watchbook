<script lang="ts">
  import MediaCard from "$lib/components/MediaCard.svelte";
  import { Button, buttonVariants } from "@nanoteck137/nano-ui";
  import { FileQuestion, Image as ImageIcon, Plus } from "lucide-svelte";
  import ShowLogoModal from "./ShowLogoModal.svelte";
  import { cn, isRoleAdmin } from "$lib/utils";
  import { getApiClient, handleApiError } from "$lib";
  import toast from "svelte-5-french-toast";
  import { invalidateAll } from "$app/navigation";
  import MediaItemDropdown from "./MediaItemDropdown.svelte";
  import AddMediaItem from "./AddMediaItem.svelte";
  import EditImagesModal from "./EditImagesModal.svelte";
  import CollectionDropdown from "./CollectionDropdown.svelte";
  import Spacer from "$lib/components/Spacer.svelte";

  const { data } = $props();
  const apiClient = getApiClient();

  let openAddMediaModal = $state(false);
  let openEditImagesModal = $state(false);
</script>

<div
  class="relative h-48 w-full overflow-hidden rounded-lg shadow-lg sm:h-60 md:h-72"
>
  {#if data.collection.bannerUrl}
    <img
      src={data.collection.bannerUrl}
      alt="Collection Banner"
      class="h-full w-full object-cover"
    />
  {:else}
    <div class="h-full w-full bg-background"></div>
  {/if}

  <div class="absolute inset-0 rounded-lg border bg-black bg-opacity-40"></div>

  {#if data.collection.logoUrl}
    <div
      class="pointer-events-none absolute inset-0 hidden items-center justify-center md:flex"
    >
      <div class="w-40 rounded-lg border bg-black bg-opacity-60 p-4 shadow-lg">
        <img
          src={data.collection.logoUrl}
          alt="Collection Logo"
          class="mx-auto max-h-24 w-auto object-contain"
        />
      </div>
    </div>
  {/if}

  {#if data.collection.logoUrl}
    <ShowLogoModal
      class={cn(
        "absolute right-4 top-4 z-20 md:hidden",
        buttonVariants({ variant: "secondary", size: "icon" }),
      )}
      logoUrl={data.collection.logoUrl}
      onResult={() => {}}
    >
      <ImageIcon />
    </ShowLogoModal>
  {/if}
</div>

<div
  class="relative mx-2 -mt-16 flex flex-col items-center space-y-4 px-4 sm:-mt-20 sm:flex-row sm:items-start sm:space-x-6 sm:space-y-0 sm:px-0"
>
  <div
    class="relative z-10 aspect-[75/106] w-40 min-w-40 flex-shrink-0 overflow-hidden rounded-lg border bg-black shadow-lg sm:w-48 sm:min-w-48"
  >
    {#if data.collection.coverUrl}
      <img
        src={data.collection.coverUrl}
        alt="Collection Cover"
        class="h-full w-full object-cover"
      />
    {:else}
      <div class="flex h-full w-full items-center justify-center">
        <FileQuestion size={52} />
      </div>
    {/if}

    <CollectionDropdown collection={data.collection} />
  </div>

  <div
    class="flex max-w-xl flex-col justify-start pb-4 text-center sm:pt-20 sm:text-left"
  >
    <h1 class="text-2xl font-bold drop-shadow-lg sm:pt-2">
      {data.collection.name}
    </h1>
    <!-- <p class="text-gray-300">
      All seasons and movies of the Fullmetal Alchemist franchise, including
      Brotherhood and the original series.
    </p>
    <p class="text-sm italic text-gray-400">
      Updated regularly with new releases and extras.
    </p> -->
  </div>
</div>

<Spacer size="md" />

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

<EditImagesModal
  bind:open={openEditImagesModal}
  collectionId={data.collection.id}
/>
