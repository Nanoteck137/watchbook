<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import Spacer from "$lib/components/Spacer.svelte";
  import UserMediaCard from "$lib/components/UserMediaCard.svelte";
  import { Button, Card, Select } from "@nanoteck137/nano-ui";
  import StandardPagination from "$lib/components/StandardPagination.svelte";
  import { mediaStatus, mediaTypes } from "../../../media/types.js";
  import { sortTypes } from "./types.js";

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

    goto("?" + params.toString(), { invalidateAll: true, noScroll: true });
  }

  function gotoList(list: string) {
    const params = $page.url.searchParams;
    params.set("list", list);
    params.set("page", "0");

    goto("?" + params.toString(), { invalidateAll: true, noScroll: true });
  }
</script>

<StandardPagination pageData={data.page} />

<Spacer size="md" />

<Card.Root>
  <div class="flex flex-col">
    <div class="flex flex-col gap-4 border-b p-4">
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
          value={data.filter.list}
          allowDeselect={false}
          onValueChange={(value) => {
            gotoList(value);
          }}
        >
          <Select.Trigger>
            {lists.find((i) => i.value === data.filter.list)?.label ?? "List"}
          </Select.Trigger>
          <Select.Content>
            {#each lists as list (list.value)}
              <Select.Item value={list.value} label={list.label} />
            {/each}
          </Select.Content>
        </Select.Root>
      </div>

      <Select.Root
        type="multiple"
        value={data.filter.types}
        onValueChange={(values) => {
          const query = $page.url.searchParams;
          query.set("types", values.join(","));

          goto("?" + query.toString(), {
            invalidateAll: true,
            noScroll: true,
          });
        }}
      >
        <Select.Trigger class="sm:max-w-[120px]">Type</Select.Trigger>
        <Select.Content>
          {#each mediaTypes as ty (ty.value)}
            <Select.Item value={ty.value} label={ty.label} />
          {/each}
        </Select.Content>
      </Select.Root>

      <Select.Root
        type="multiple"
        value={data.filter.status}
        onValueChange={(values) => {
          const query = $page.url.searchParams;
          query.set("status", values.join(","));

          goto("?" + query.toString(), {
            invalidateAll: true,
            noScroll: true,
          });
        }}
      >
        <Select.Trigger class="sm:max-w-[120px]">Status</Select.Trigger>
        <Select.Content>
          {#each mediaStatus as status (status.value)}
            <Select.Item value={status.value} label={status.label} />
          {/each}
        </Select.Content>
      </Select.Root>

      <Select.Root
        type="single"
        allowDeselect={false}
        value={data.filter.sort}
        onValueChange={(sort) => {
          gotoSort(sort);
        }}
      >
        <Select.Trigger class="sm:max-w-[180px]">
          {sortTypes.find((i) => i.value === data.filter.sort)?.label ??
            "Sort"}
        </Select.Trigger>
        <Select.Content>
          {#each sortTypes as ty (ty.value)}
            <Select.Item value={ty.value} label={ty.label} />
          {/each}
        </Select.Content>
      </Select.Root>

      <Button
        variant="outline"
        onclick={() => {
          const query = $page.url.searchParams;
          query.delete("sort");
          query.delete("types");
          query.delete("status");

          goto("?" + query.toString(), {
            invalidateAll: true,
            noScroll: true,
          });
        }}
      >
        Reset Filters
      </Button>

      <!-- <Filter
        fullFilter={{
          sort: "title-a-z",
          query: "",
          filters: { type: [], rating: [], status: [] },
          excludes: { type: [], rating: [], status: [] },
        }}
      /> -->
    </div>

    <div class="flex flex-col py-4">
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
    </div>
  </div>
</Card.Root>

<Spacer size="md" />

<StandardPagination pageData={data.page} />
