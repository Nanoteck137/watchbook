<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import Image from "$lib/components/Image.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import UserMediaCard from "$lib/components/UserMediaCard.svelte";
  import {
    Button,
    buttonVariants,
    DropdownMenu,
    Pagination,
    ScrollArea,
  } from "@nanoteck137/nano-ui";
  import { ArrowUpDown, ChevronDown, Star } from "lucide-svelte";
  import HeaderButton from "./HeaderButton.svelte";
  import { cn } from "$lib/utils";

  const { data } = $props();

  const buttons = [
    "all",
    "in-progress",
    "completed",
    "on-hold",
    "dropped",
    "backlog",
  ];

  function gotoSort(sort: string) {
    const params = $page.url.searchParams;
    params.set("sort", sort);
    goto("?" + params.toString(), { invalidateAll: true });
  }

  function getName() {
    if (data.userData.displayName) return data.userData.displayName;
    return data.userData.username;
  }
</script>

<!-- Header Section -->
<div class="mb-6">
  <!-- Title -->
  <h1 class="mb-4 text-3xl font-bold text-white">
    {getName()} Watchlist
  </h1>

  <!-- Status Buttons with Strips -->
  <div class="grid grid-cols-6 gap-2 text-center">
    <HeaderButton
      href="?list=all"
      name="All"
      stripClass="bg-white"
      active={data.list === "all"}
    />
    <HeaderButton
      href="?list=completed"
      name="Completed"
      stripClass="bg-green-600"
      active={data.list === "completed"}
    />
    <HeaderButton
      href="?list=in-progress"
      name="In Progress"
      stripClass="bg-blue-500"
      active={data.list === "in-progress"}
    />
    <HeaderButton
      href="?list=on-hold"
      name="On Hold"
      stripClass="bg-yellow-500"
      active={data.list === "on-hold"}
    />
    <HeaderButton
      href="?list=dropped"
      name="Dropped"
      stripClass="bg-red-600"
      active={data.list === "dropped"}
    />
    <HeaderButton
      href="?list=backlog"
      name="Backlog"
      stripClass="bg-gray-600"
      active={data.list === "backlog"}
    />
  </div>

  <Spacer />

  <DropdownMenu.Root>
    <DropdownMenu.Trigger class={cn(buttonVariants({ variant: "outline" }))}>
      <ArrowUpDown />
      Sort
    </DropdownMenu.Trigger>
    <DropdownMenu.Content align="start">
      <DropdownMenu.Group>
        <DropdownMenu.GroupHeading>Sorting</DropdownMenu.GroupHeading>
        <DropdownMenu.Separator />
        <DropdownMenu.Item onclick={() => gotoSort("titleAsc")}>
          Title (A-Z)
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => gotoSort("titleDesc")}>
          Title (Z-A)
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => gotoSort("userScoreDesc")}>
          User Score (+)
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => gotoSort("userScoreAsc")}>
          User Score (-)
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => gotoSort("scoreDesc")}>
          Score (+)
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => gotoSort("scoreAsc")}>
          Score (-)
        </DropdownMenu.Item>
      </DropdownMenu.Group>
    </DropdownMenu.Content>
  </DropdownMenu.Root>
</div>

<!-- Media Grid -->

<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
  <div
    class="grid grid-cols-[repeat(auto-fit,minmax(240px,1fr))] items-center justify-items-center gap-6"
  >
    <!-- <div
    class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6"
  > -->

    {#each data.media as media}
      <UserMediaCard
        href="/media/{media.id}"
        title={media.title}
        coverUrl={media.coverUrl}
        score={media.user?.score}
        list={media.user?.list}
        currentPart={media.user?.currentPart}
        partCount={media.partCount}
      />
    {/each}
  </div>
</div>

{#if data.user?.id === data.userId}
  <Button href="/account">Account</Button>
{/if}

<ScrollArea orientation="horizontal">
  <div class="flex gap-2 py-2">
    {#each buttons as button}
      <Button variant="link" data-sveltekit-replacestate href="?list={button}">
        {button}
      </Button>
    {/each}
  </div>
</ScrollArea>

<div class="flex flex-col gap-4">
  {#each data.media as media}
    <div class="flex justify-between border-b py-2">
      <div class="flex">
        <Image class="h-20 w-14" src={media.coverUrl} alt="cover" />
        <div class="flex flex-col gap-2 px-4 py-1">
          <a
            class="line-clamp-2 text-ellipsis text-sm font-semibold hover:cursor-pointer hover:underline"
            href="/media/{media.id}"
            title={media.title}
          >
            {media.title}
          </a>
          <p class="text-xs">{media.user?.list}</p>
        </div>
      </div>
      <div class="flex min-w-24 max-w-24 items-center justify-center border-l">
        <Star size={18} class="fill-foreground" />
        <Spacer horizontal size="xs" />
        <p class="font-mono text-xs">
          {media.user?.score ?? "0"}
        </p>
      </div>
    </div>
  {/each}
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
