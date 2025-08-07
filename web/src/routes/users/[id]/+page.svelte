<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import Image from "$lib/components/Image.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { Button, Pagination, ScrollArea } from "@nanoteck137/nano-ui";
  import { Star } from "lucide-svelte";

  const { data } = $props();

  const buttons = [
    "all",
    "in-progress",
    "completed",
    "on-hold",
    "dropped",
    "backlog",
  ];
</script>

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
