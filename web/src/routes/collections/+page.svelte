<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import CollectionCard from "$lib/components/CollectionCard.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import {
    Button,
    Checkbox,
    Input,
    Label,
    Pagination,
    Select,
  } from "@nanoteck137/nano-ui";
  import NewCollectionModal from "./NewCollectionModal.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import Filter from "./Filter.svelte";

  const { data } = $props();

  let openNewCollectionModal = $state(false);
</script>

<Filter fullFilter={data.filter} />

<p>Total collections: {data.page.totalItems}</p>

<Button
  onclick={() => {
    openNewCollectionModal = true;
  }}
>
  New Collection
</Button>

<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
  <div
    class="grid grid-cols-[repeat(auto-fit,minmax(280px,1fr))] items-center justify-items-center gap-6"
  >
    {#each data.collections as collection}
      <CollectionCard
        href="/collections/{collection.id}"
        name={collection.name}
        coverUrl={collection.coverUrl}
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

<NewCollectionModal bind:open={openNewCollectionModal} />
