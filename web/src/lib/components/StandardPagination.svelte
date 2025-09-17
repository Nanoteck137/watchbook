<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import type { Page } from "$lib/api/types";
  import { Pagination } from "@nanoteck137/nano-ui";

  type Props = {
    pageData: Page;
  };

  const { pageData }: Props = $props();
</script>

<Pagination.Root
  page={pageData.page + 1}
  count={pageData.totalItems}
  perPage={pageData.perPage}
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
        <Pagination.PrevButton class="w-28" />
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
        <Pagination.NextButton class="w-28" />
      </Pagination.Item>
    </Pagination.Content>
  {/snippet}
</Pagination.Root>
