<script lang="ts">
  import { goto } from "$app/navigation";
  import { getApiClient } from "$lib";
  import Spacer from "$lib/components/Spacer.svelte";
  import { Pagination } from "@nanoteck137/nano-ui";
  import { page } from "$app/stores";
  import MediaCard from "$lib/components/MediaCard.svelte";
  import Filter from "./Filter.svelte";

  const { data } = $props();
  const apiClient = getApiClient();

  // NOTE(patrik):
  //  - Large image size: 225x318 (h-80 w-56)
  //  - Medium image size: 160x220 (h-56 w-40)
  //  - Small image size: 50x70 (h-20 w-14)
</script>

<Spacer size="sm" />

<Filter fullFilter={data.filter} />

<Spacer />

<p>Total: {data.page.totalItems}</p>

<Spacer />

<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
  <div
    class="grid grid-cols-[repeat(auto-fit,minmax(240px,1fr))] items-center justify-items-center gap-6"
  >
    {#each data.media as media}
      <MediaCard
        href="/media/{media.id}"
        title={media.title}
        coverUrl={media.coverUrl}
        startDate={media.startDate}
        partCount={media.partCount}
        score={media.score}
        userList={media.user?.list ?? null}
      />
    {/each}
  </div>
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
