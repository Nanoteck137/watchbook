<script lang="ts">
  import { goto } from "$app/navigation";
  import { getApiClient } from "$lib";
  import Image from "$lib/components/Image.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { pickTitle } from "$lib/utils.js";
  import { Pagination } from "@nanoteck137/nano-ui";
  import { Star } from "lucide-svelte";
  import ImportMalListDialog from "./ImportMalListDialog.svelte";
  import ImportAnime from "./ImportAnime.svelte";
  import { page } from "$app/stores";

  const { data } = $props();
  const apiClient = getApiClient();

  // NOTE(patrik):
  //  - Large image size: 225x318 (h-80 w-56)
  //  - Medium image size: 160x220 (h-56 w-40)
  //  - Small image size: 50x70 (h-20 w-14)
</script>

<ImportMalListDialog />
<ImportAnime />

<Spacer size="sm" />

<div class="flex flex-col gap-4">
  {#each data.animes as anime}
    {@const title = pickTitle(anime)}

    <div class="flex justify-between border-b py-2">
      <div class="flex">
        <Image class="h-20 w-14" src={anime.coverUrl} alt="cover" />
        <div class="px-4 py-1">
          <a
            class="line-clamp-2 text-ellipsis text-sm font-semibold hover:cursor-pointer hover:underline"
            href="/animes/{anime.id}"
            {title}
          >
            {title}
          </a>
        </div>
      </div>
      <div class="flex min-w-24 max-w-24 items-center justify-center border-l">
        <Star size={18} class="fill-foreground" />
        <Spacer horizontal size="xs" />
        <p class="font-mono text-xs">{anime.score?.toFixed(2) ?? "N/A"}</p>
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
