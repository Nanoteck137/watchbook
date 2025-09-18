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
    ScrollArea,
  } from "@nanoteck137/nano-ui";
  import { ArrowUpDown, Star } from "lucide-svelte";
  import HeaderButton from "./HeaderButton.svelte";
  import { cn, userListClass } from "$lib/utils";
  import StandardPagination from "$lib/components/StandardPagination.svelte";

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

{#if data.user?.id === data.userId}
  <Button href="/account">Account</Button>
{/if}

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
      stripClass={userListClass("completed")}
      active={data.list === "completed"}
    />
    <HeaderButton
      href="?list=in-progress"
      name="In Progress"
      stripClass={userListClass("in-progress")}
      active={data.list === "in-progress"}
    />
    <HeaderButton
      href="?list=on-hold"
      name="On Hold"
      stripClass={userListClass("on-hold")}
      active={data.list === "on-hold"}
    />
    <HeaderButton
      href="?list=dropped"
      name="Dropped"
      stripClass={userListClass("dropped")}
      active={data.list === "dropped"}
    />
    <HeaderButton
      href="?list=backlog"
      name="Backlog"
      stripClass={userListClass("backlog")}
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

<Spacer />

<StandardPagination pageData={data.page} />

<Spacer size="md" />

<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
  <div
    class="grid grid-cols-[repeat(auto-fit,minmax(240px,1fr))] items-center justify-items-center gap-6"
  >
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

<Spacer size="md" />

<StandardPagination pageData={data.page} />
