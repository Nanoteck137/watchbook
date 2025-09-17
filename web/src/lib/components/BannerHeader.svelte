<script lang="ts">
  import ShowLogoModal from "$lib/components/ShowLogoModal.svelte";
  import { cn } from "$lib/utils";
  import { Button, buttonVariants } from "@nanoteck137/nano-ui";
  import { FileQuestion, Image as ImageIcon } from "lucide-svelte";
  import type { Snippet } from "svelte";

  type Props = {
    title: string;
    description?: string | null;
    coverUrl: string | null;
    bannerUrl: string | null;
    logoUrl: string | null;

    imageContent?: Snippet<[]>;
  };

  const {
    title,
    description,
    coverUrl,
    bannerUrl,
    logoUrl,
    imageContent,
  }: Props = $props();

  let descriptionShowMore = $state(false);
</script>

<div
  class="relative h-48 w-full overflow-hidden rounded-lg shadow-lg sm:h-60 md:h-72"
>
  {#if bannerUrl}
    <img
      src={bannerUrl}
      alt="Collection Banner"
      class="h-full w-full object-cover"
    />
  {:else}
    <div class="h-full w-full bg-background"></div>
  {/if}

  <div class="absolute inset-0 rounded-lg border bg-black bg-opacity-40"></div>

  {#if logoUrl}
    <div
      class="pointer-events-none absolute inset-0 hidden items-center justify-center md:flex"
    >
      <div class="w-40 rounded-lg border bg-black bg-opacity-60 p-4 shadow-lg">
        <img
          src={logoUrl}
          alt="Collection Logo"
          class="mx-auto max-h-24 w-auto object-contain"
        />
      </div>
    </div>
  {/if}

  {#if logoUrl}
    <ShowLogoModal
      class={cn(
        "absolute right-4 top-4 z-20 md:hidden",
        buttonVariants({ variant: "secondary", size: "icon" }),
      )}
      {logoUrl}
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
    {#if coverUrl}
      <img
        src={coverUrl}
        alt="Collection Cover"
        class="h-full w-full object-cover"
      />
    {:else}
      <div class="flex h-full w-full items-center justify-center">
        <FileQuestion size={52} />
      </div>
    {/if}

    {#if imageContent}
      {@render imageContent()}
    {/if}
    <!-- <CollectionDropdown collection={data.collection} /> -->
  </div>

  <div
    class="flex max-w-xl flex-col justify-start gap-2 pb-4 text-center sm:pt-20 sm:text-left"
  >
    <h1 class="text-2xl font-bold drop-shadow-lg sm:pt-2" {title}>
      {title}
    </h1>
    {#if description}
      <div class="flex flex-col gap-1">
        <p
          class={`text-ellipsis whitespace-pre-line text-sm ${!descriptionShowMore ? "line-clamp-4" : ""}`}
        >
          {description}
        </p>

        <Button
          class="w-fit"
          size="sm"
          variant="outline"
          onclick={() => {
            descriptionShowMore = !descriptionShowMore;
          }}
        >
          Show More
        </Button>
      </div>
    {/if}
    <!-- <p class="text-sm italic text-gray-400">
      Updated regularly with new releases and extras.
    </p> -->
  </div>
</div>
