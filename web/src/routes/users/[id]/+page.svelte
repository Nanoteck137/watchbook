<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import Spacer from "$lib/components/Spacer.svelte";
  import UserMediaCard from "$lib/components/UserMediaCard.svelte";
  import { Card, ScrollArea, Select } from "@nanoteck137/nano-ui";
  import StandardPagination from "$lib/components/StandardPagination.svelte";

  const { data } = $props();

  const lists = [
    { label: "All", value: "" },
    { label: "Completed", value: "completed" },
    { label: "In Progress", value: "in-progress" },
    { label: "On Hold", value: "on-hold" },
    { label: "Dropped", value: "dropped" },
    { label: "Backlog", value: "backlog" },
  ];

  function gotoSort(sort: string) {
    const params = $page.url.searchParams;
    params.set("sort", sort);
    params.set("page", "0");

    goto("?" + params.toString(), { invalidateAll: true });
  }

  function gotoList(list: string) {
    const params = $page.url.searchParams;
    params.set("list", list);
    params.set("page", "0");

    goto("?" + params.toString(), { invalidateAll: true });
  }

  function getName() {
    if (data.userData.displayName) return data.userData.displayName;
    return data.userData.username;
  }
</script>

<Card.Root>
  <div class="relative w-full overflow-hidden rounded-xl shadow-lg">
    <div
      class="h-full w-full bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 p-6"
    >
      <h1 class="text-3xl font-bold">
        {getName()}
      </h1>
    </div>
  </div>

  <div class="flex w-full justify-center gap-4">
    <!-- <div class="flex flex-col items-center gap-2 sm:items-start">
        <h1 class="text-2xl font-bold sm:line-clamp-2 sm:pt-2">
          {getName()}
        </h1>
      </div> -->

    <ScrollArea class="px-2 pb-4 sm:pb-0" orientation="horizontal">
      <div class="flex gap-4 pt-4">
        <!-- <Button variant="secondary" data-sveltekit-noscroll>Overview</Button>
      <Button variant="secondary" data-sveltekit-noscroll>Parts</Button> -->

        <button
          class="border-b-2 px-2 py-3 text-sm font-medium hover:brightness-75"
        >
          Overview
        </button>

        <button
          class="border-b-2 px-2 py-3 text-sm font-medium hover:brightness-75"
        >
          Watchlist
        </button>

        <button
          class="border-b-2 px-2 py-3 text-sm font-medium hover:brightness-75"
        >
          Folders
        </button>

        <button
          class="border-b-2 px-2 py-3 text-sm font-medium hover:brightness-75"
        >
          Settings
        </button>
      </div>
    </ScrollArea>
  </div>

  <!-- <div
    class="relative flex flex-col items-center p-2 sm:flex-row sm:items-stretch"
  ></div> -->
</Card.Root>

<Spacer size="md" />

<Card.Root>
  <div class="flex flex-col">
    <div class="border-b p-4">
      <div class="hidden justify-center gap-2 sm:flex">
        {#each lists as list (list.value)}
          <button
            class="border-b-2 px-2 py-3 text-sm font-medium hover:brightness-75"
            onclick={() => gotoList(list.value)}
          >
            {list.label}
          </button>
        {/each}
      </div>

      <div class="sm:hidden">
        <Select.Root
          type="single"
          value={data.list}
          allowDeselect={false}
          onValueChange={(value) => {
            gotoList(value);
          }}
        >
          <Select.Trigger>
            {lists.find((i) => i.value === data.list)?.label ?? "List"}
          </Select.Trigger>
          <Select.Content>
            {#each lists as list (list.value)}
              <Select.Item value={list.value} label={list.label} />
            {/each}
            <!-- <Select.Item value={""} label={"All"} />
          <Select.Item value={"completed"} label={"Completed"} />
          <Select.Item value={"in-progress"} label={"In Progress"} />
          <Select.Item value={"on-hold"} label={"On Hold"} />
          <Select.Item value={"dropped"} label={"Dropped"} />
          <Select.Item value={"backlog"} label={"Backlog"} /> -->
          </Select.Content>
        </Select.Root>
      </div>

      <!-- <Filter
        fullFilter={{
          sort: "title-a-z",
          query: "",
          filters: { type: [], rating: [], status: [] },
          excludes: { type: [], rating: [], status: [] },
        }}
      /> -->
    </div>

    <div class="flex flex-col">
      <Spacer size="md" />

      <StandardPagination pageData={data.page} />

      <Spacer size="md" />

      <div
        class="grid w-full grid-cols-1 items-center justify-items-center gap-4 px-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4"
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

      <Spacer size="md" />

      <StandardPagination pageData={data.page} />
    </div>
  </div>
</Card.Root>
