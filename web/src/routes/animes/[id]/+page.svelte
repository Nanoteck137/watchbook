<script lang="ts">
  import { Anime } from "$lib/api/types.js";
  import Image from "$lib/components/Image.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { pickTitle } from "$lib/utils.js";
  import {
    Breadcrumb,
    Button,
    DropdownMenu,
    Separator,
  } from "@nanoteck137/nano-ui";
  import { ChevronDown, Eye, Star } from "lucide-svelte";

  const { data } = $props();

  let showMore = $state(false);
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/animes">Animes</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page class="line-clamp-1 max-w-96 text-ellipsis"
          >{pickTitle(data.anime)}</Breadcrumb.Page
        >
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<Spacer size="sm" />

<div class="flex flex-col items-center gap-2">
  <Image class="min-h-80 w-56" src={data.anime.coverUrl} alt="cover" />
  <div class="flex flex-col gap-1">
    <p class="text-center text-base">{data.anime.title}</p>
    {#if data.anime.titleEnglish}
      <Separator />
      <p class="text-center text-sm text-zinc-300">
        {data.anime.titleEnglish}
      </p>
    {/if}
  </div>
</div>

<Spacer size="lg" />

<div class="flex flex-col justify-around gap-2">
  <div class="flex items-center justify-center rounded bg-gray-600 py-1">
    <Star size={18} class="text-yellow-200" />
    <Spacer horizontal size="sm" />
    <p class="text-base">8.82</p>
  </div>

  <DropdownMenu.Root>
    <DropdownMenu.Trigger
      class="relative flex items-center justify-center rounded bg-gray-600 py-1"
    >
      <Star size={18} class="fill-yellow-200 text-yellow-200" />
      <Spacer horizontal size="sm" />
      <p class="text-base">8</p>
      <ChevronDown class="absolute right-4" size={20} />
    </DropdownMenu.Trigger>
    <DropdownMenu.Content class="w-40">
      <DropdownMenu.Group>
        <DropdownMenu.GroupHeading>Select Rating</DropdownMenu.GroupHeading>
        <DropdownMenu.Separator />
        <DropdownMenu.Item>(10) Masterpiece</DropdownMenu.Item>
        <DropdownMenu.Item>(9) Great</DropdownMenu.Item>
        <DropdownMenu.Item>(8) Very Good</DropdownMenu.Item>
        <DropdownMenu.Item>(7) Good</DropdownMenu.Item>
        <DropdownMenu.Item>(6) Fine</DropdownMenu.Item>
        <DropdownMenu.Item>(5) Average</DropdownMenu.Item>
        <DropdownMenu.Item>(4) Bad</DropdownMenu.Item>
        <DropdownMenu.Item>(3) Very Bad</DropdownMenu.Item>
        <DropdownMenu.Item>(2) Horrible</DropdownMenu.Item>
        <DropdownMenu.Item>(1) Appalling</DropdownMenu.Item>
      </DropdownMenu.Group>
    </DropdownMenu.Content>
  </DropdownMenu.Root>

  <!-- <Separator orientation="vertical" /> -->

  <div class="flex items-center justify-center rounded bg-gray-600 py-1">
    <Eye size={18} class="" />
    <Spacer horizontal size="sm" />
    <p class="text-base">1000 / 1000</p>
  </div>
</div>

<Spacer size="lg" />

{#if data.anime.description}
  <div class="flex flex-col gap-1">
    <p
      class={`text-ellipsis whitespace-pre-line text-sm ${!showMore ? "line-clamp-2" : ""}`}
    >
      {data.anime.description}
    </p>
    <Button
      class="w-fit"
      size="sm"
      variant="outline"
      onclick={() => {
        showMore = !showMore;
      }}
    >
      Show More
    </Button>
  </div>
{/if}
